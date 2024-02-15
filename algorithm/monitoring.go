package algorithm

import (
	"../db"
	f "../filtering"
	fb "../firebase"
	v "../variables"
	"log"
	"strconv"
)

func (a *Algorithm) GetViableStocks () {
	log.Println("Getting viable stocks for algorithm: ", a.Name)
	var filteredStocks = f.GetSymbols()
	filteredStocks = a.RunThroughViableFilters(filteredStocks)
	for _, filter := range a.ViableFilter.Checks {
		switch filter.TimePeriod {
		case v.OneWeek:
			filteredStocks = f.Filters(filteredStocks, filter)
		case v.OneDay:
			filteredStocks = f.Filters(filteredStocks, filter)
		case v.FourHour:
			filteredStocks = f.Filters(filteredStocks, filter)
		case v.OneHour:
			filteredStocks = f.Filters(filteredStocks, filter)
		}
	}
	fb.PingUpdate(a.UserName, "Viable filter completed for algorithm: " + a.Name, strconv.Itoa(len(filteredStocks)) + " viable stocks found", a.UniqueID)
	db.UpdateViableStocks(a.UserName, filteredStocks, a.UniqueID)
	log.Println("GetViableStocks for algorithm: ", a.Name, " complete")
}
func (a *Algorithm) MonitorViableStocks () {
	log.Println("Getting purchasable stocks for algorithm: ", a.Name)
	var newPurchasable = db.GetViableStocks(a.UserName, a.UniqueID)
	var oldPurchasable = db.GetPurchasableStocks(a.UserName, a.UniqueID)

	newPurchasable = a.RunThroughViableFilters(newPurchasable)
	newPurchasable = a.RunThroughPurchasableFilters(newPurchasable)
	for _, filter := range a.PurchasableFilter.Checks {
		switch filter.TimePeriod {
		case v.OneWeek:
			newPurchasable = f.Filters(newPurchasable, filter)
		case v.OneDay:
			newPurchasable = f.Filters(newPurchasable, filter)
		case v.FourHour:
			newPurchasable = f.Filters(newPurchasable, filter)
		case v.OneHour:
			newPurchasable = f.Filters(newPurchasable, filter)
		}
	}
	a.CheckIfNewPurchasableStockFound(oldPurchasable, newPurchasable)
	db.UpdatePurchasableStocks(a.UserName, newPurchasable, a.UniqueID)
	log.Println("MonitorViableStocks for algorithm: ", a.Name, " complete")
}
func (a *Algorithm) MonitorUserStocks () {
	log.Println("Getting sellable stocks for algorithm: ", a.Name)
	var toSell []v.Stock
	var oldBought []v.Stock
	oldBought = db.GetUserStocks(a.UserName, a.UniqueID)
	oldBought = a.RunThroughViableFilters(oldBought)
	oldBought = a.RunThroughPurchasableFilters(oldBought)
	oldBought = a.RunThroughSellFilters(oldBought)

	toSell = append(toSell, oldBought...)
	for _, filter := range a.SellFilter.Checks {
		switch filter.TimePeriod {
		case v.OneWeek:
			toSell = f.Filters(toSell, filter)
		case v.OneDay:
			toSell = f.Filters(toSell, filter)
		case v.FourHour:
			toSell = f.Filters(toSell, filter)
		case v.OneHour:
			toSell = f.Filters(toSell, filter)
		}
	}

	for i,_ := range oldBought {
		sell := true
		for _, dontSell := range toSell {
			if dontSell.Symbol == oldBought[i].Symbol {
				sell = false
				break
			}
		}
		if sell == true {
			oldBought[i].Sell = true
		} else {
			oldBought[i].Sell = false
		}
	}

	a.SendSellSignals(oldBought)
	db.UpdateUserStocks(a.UserName, oldBought, a.UniqueID)
	log.Println("MonitorUserStocks for algorithm: ", a.Name, " complete")
}
func (a *Algorithm) SendSellSignals (stocks []v.Stock) {
	var stocksToSell string
	for _, stock := range stocks {
		if stock.Sell == true {
			stocksToSell += stock.Symbol + " | "
		}
	}
	if stocksToSell != "" {
		stocksToSell = stocksToSell[:len(stocksToSell) - 2]
		fb.PingStockToSell(a.UserName, stocksToSell, a.UniqueID)
	} else {
		fb.PingUpdate(a.UserName, "No stocks ready to sell for algorithm: " + a.Name, "", a.UniqueID)
	}
}
func (a *Algorithm) CheckIfNewPurchasableStockFound (old []v.Stock, new []v.Stock) {
	var viableStocksString string
	for _, newStock := range new {
		var isNew = true
		for _, oldStock := range old {
			if oldStock.Symbol == newStock.Symbol {
				isNew = false
				break
			}
		}
		if isNew == true {
			viableStocksString += newStock.Symbol + " | "
		}
	}
	if viableStocksString != "" {
		viableStocksString = viableStocksString[:len(viableStocksString) - 2]
		fb.PingUpdate(a.UserName, "New purchasable Stocks for algorithm: " + a.Name, "", a.UniqueID)
		log.Println("New purchasable stocks found for algorithm: ", a.Name, "\n", viableStocksString)
	} else {
		fb.PingUpdate(a.UserName, "No new purchasable stocks", "Algorithm: " + a.Name , a.UniqueID)
		log.Println("No new purchasable stocks found for algorithm: ", a.Name)
	}
}

func (a *Algorithm) RunThroughViableFilters (stocks []v.Stock) []v.Stock {
	var filteredStocks []v.Stock
	filteredStocks = stocks
	for _, filter := range a.ViableFilter.Checks {
		switch filter.TimePeriod {
		case v.OneWeek:
			filteredStocks = f.CreateAndRunThreadsForGettingOneWeekBars(filteredStocks, filter.Period)
			filteredStocks = f.FilterStocks(filteredStocks, filter)
		case v.OneDay:
			filteredStocks = f.CreateAndRunThreadsForGettingOneDayBars(filteredStocks, filter.Period)
			filteredStocks = f.FilterStocks(filteredStocks, filter)
		case v.FourHour:
			if len(filteredStocks) == 0 || len(filteredStocks[0].HourBars) == 0 {
				filteredStocks = f.CreateAndRunThreadsForGettingOneHourBars(filteredStocks, filter.Period)
			}
			filteredStocks = f.FilterStocks(filteredStocks, filter)
		case v.OneHour:
			filteredStocks = f.CreateAndRunThreadsForGettingOneHourBars(filteredStocks, filter.Period)
			filteredStocks = f.FilterStocks(filteredStocks, filter)
		}
	}
	return filteredStocks
}
func (a *Algorithm) RunThroughPurchasableFilters (stocks []v.Stock) []v.Stock {
	var filteredStocks []v.Stock
	filteredStocks = stocks
	for _, filter := range a.PurchasableFilter.Checks {
		switch filter.TimePeriod {
		case v.OneWeek:
			filteredStocks = f.CreateAndRunThreadsForGettingOneWeekBars(filteredStocks, filter.Period)
			filteredStocks = f.FilterStocks(filteredStocks, filter)
		case v.OneDay:
			filteredStocks = f.CreateAndRunThreadsForGettingOneDayBars(filteredStocks, filter.Period)
			filteredStocks = f.FilterStocks(filteredStocks, filter)
		case v.FourHour:
			if len(filteredStocks) == 0 || len(filteredStocks[0].HourBars) == 0 {
				filteredStocks = f.CreateAndRunThreadsForGettingOneHourBars(filteredStocks, filter.Period)
			}
			filteredStocks = f.FilterStocks(filteredStocks, filter)
		case v.OneHour:
			filteredStocks = f.CreateAndRunThreadsForGettingOneHourBars(filteredStocks, filter.Period)
			filteredStocks = f.FilterStocks(filteredStocks, filter)
		}
	}
	return filteredStocks
}
func (a *Algorithm) RunThroughSellFilters (stocks []v.Stock) []v.Stock {
	var filteredStocks []v.Stock
	filteredStocks = stocks
	for _, filter := range a.SellFilter.Checks {
		switch filter.TimePeriod {
		case v.OneWeek:
			filteredStocks = f.CreateAndRunThreadsForGettingOneWeekBars(filteredStocks, filter.Period)
			filteredStocks = f.FilterStocks(filteredStocks, filter)
		case v.OneDay:
			filteredStocks = f.CreateAndRunThreadsForGettingOneDayBars(filteredStocks, filter.Period)
			filteredStocks = f.FilterStocks(filteredStocks, filter)
		case v.FourHour:
			if len(filteredStocks) == 0 || len(filteredStocks[0].HourBars) == 0 {
				filteredStocks = f.CreateAndRunThreadsForGettingOneHourBars(filteredStocks, filter.Period)
			}
			filteredStocks = f.FilterStocks(filteredStocks, filter)
		case v.OneHour:
			filteredStocks = f.CreateAndRunThreadsForGettingOneHourBars(filteredStocks, filter.Period)
			filteredStocks = f.FilterStocks(filteredStocks, filter)
		}
	}
	return filteredStocks
}


