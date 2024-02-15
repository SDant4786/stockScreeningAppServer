package filtering

import (
	v "../variables"
	"log"
)

//todo split up the volume jump and close greater filters
func FilterStocks (stocks []v.Stock, filter FilterCheck) []v.Stock{
	var filteredStocks []v.Stock
	filteredStocks = FilterCalculations(stocks, filter)
	return filteredStocks
}
/*
	Calculates all the values needed for filters
 */
func FilterCalculations (stocks []v.Stock, filter FilterCheck) []v.Stock {
	log.Println("Running FilterCalculations")
	stocks = Adx(stocks, filter.TimePeriod)
	stocks = AccumulationDistribution(stocks, filter.TimePeriod, filter.Period)
	stocks = EMA(stocks, filter.TimePeriod, filter.Period)
	stocks = SMA(stocks, filter.TimePeriod, filter.Period)
	stocks = RSI(stocks, filter.TimePeriod, filter.Period)
	stocks = AverageVolume(stocks, filter.TimePeriod)
	stocks = LastClose(stocks, filter.TimePeriod)
	stocks = BB(stocks, filter.TimePeriod, filter.Period)
	stocks = VolumeJump(stocks, filter.TimePeriod)
	stocks = CloseGreaterThanPrevious(stocks, filter.TimePeriod)
	log.Println("Completed FilterCalculations")
	return stocks
}
/*
	Remove stocks that dont fit the parameters
 */
func Filters (stocks []v.Stock, filter FilterCheck) []v.Stock {
	log.Println("Running Filters")
	if filter.ADX.Set == true {
		stocks = ADXFilter(stocks, filter)
	}
	if filter.AccumulationDistribution.Set == true {
		stocks = AccumulationDistributionFilter(stocks, filter)
	}
	if filter.EMA.Set == true {
		stocks = EMAFilter(stocks, filter)
	}
	if filter.SMA.Set == true {
		stocks = SMAFilter(stocks, filter)
	}
	if filter.RSI.Set == true {
		stocks = RSIFilter(stocks, filter)
	}
	if filter.AverageVolume.Set == true {
		stocks = AverageVolumeFilter(stocks, filter)
	}
	if filter.LastClose.Set == true {
		stocks = LastCloseFilter(stocks, filter)
	}
	if filter.BB.Set == true {
		stocks = BBFilter(stocks, filter)
	}
	if filter.VolumeJump.Set == true {
		stocks = VolumeJumpFilter(stocks, filter)
	}
	if filter.CloseGreaterThanPrevious.Set == true {
		stocks = CloseGreaterThanPreviousFilter(stocks, filter)
	}
	log.Println("Completed Filters")
	return stocks
}
/*
	Filters
 */
func ADXFilter (stocks []v.Stock, filter FilterCheck) []v.Stock {
	log.Println("Running ADXFilter")
	var stocksToRemove []v.Stock
	if filter.ADX.ADXIncreasing == 1 {
		for _, stock := range stocks {
			switch filter.TimePeriod {
			case v.OneWeek:
				l := len(stock.OneWeek.Adx) - 1
				if l < 1 || stock.OneWeek.Adx[l].Adx < stock.OneWeek.Adx[l-1].Adx {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneDay:
				l := len(stock.OneDay.Adx) - 1
				if l < 1 || stock.OneDay.Adx[l].Adx < stock.OneDay.Adx[l-1].Adx {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneHour:
				l := len(stock.OneHour.Adx) - 1
				if l < 1 || stock.OneHour.Adx[l].Adx < stock.OneHour.Adx[l-1].Adx {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.FourHour:
				l := len(stock.FourHour.Adx) - 1
				if l < 1 || stock.FourHour.Adx[l].Adx < stock.FourHour.Adx[l-1].Adx {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			}
		}
	} else if filter.ADX.ADXIncreasing == 2 {
		for _, stock := range stocks {
			switch filter.TimePeriod {
			case v.OneWeek:
				l := len(stock.OneWeek.Adx) - 1
				if l < 1 || stock.OneWeek.Adx[l].Adx > stock.OneWeek.Adx[l-1].Adx {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneDay:
				l := len(stock.OneDay.Adx) - 1
				if l < 1 || stock.OneDay.Adx[l].Adx > stock.OneDay.Adx[l-1].Adx {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneHour:
				l := len(stock.OneHour.Adx) - 1
				if l < 1 || stock.OneHour.Adx[l].Adx > stock.OneHour.Adx[l-1].Adx {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.FourHour:
				l := len(stock.FourHour.Adx) - 1
				if l < 1 || stock.FourHour.Adx[l].Adx > stock.FourHour.Adx[l-1].Adx {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			}
		}
	}
	if filter.ADX.PDIIncreasing == 1 {
		for _, stock := range stocks {
			switch filter.TimePeriod {
			case v.OneWeek:
				l := len(stock.OneWeek.Adx) - 1
				if l < 1 || stock.OneWeek.Adx[l].Pdi < stock.OneWeek.Adx[l-1].Pdi {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneDay:
				l := len(stock.OneDay.Adx) - 1
				if l < 1 || stock.OneDay.Adx[l].Pdi < stock.OneDay.Adx[l-1].Pdi {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneHour:
				l := len(stock.OneHour.Adx) - 1
				if l < 1 || stock.OneHour.Adx[l].Pdi < stock.OneHour.Adx[l-1].Pdi {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.FourHour:
				l := len(stock.FourHour.Adx) - 1
				if l < 1 || stock.FourHour.Adx[l].Pdi < stock.FourHour.Adx[l-1].Pdi {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			}
		}
	} else if filter.ADX.PDIIncreasing == 2 {
		for _, stock := range stocks {
			switch filter.TimePeriod {
			case v.OneWeek:
				l := len(stock.OneWeek.Adx) - 1
				if l < 1 || stock.OneWeek.Adx[l].Pdi > stock.OneWeek.Adx[l-1].Pdi {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneDay:
				l := len(stock.OneDay.Adx) - 1
				if l < 1 || stock.OneDay.Adx[l].Pdi > stock.OneDay.Adx[l-1].Pdi {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneHour:
				l := len(stock.OneHour.Adx) - 1
				if l < 1 || stock.OneHour.Adx[l].Pdi > stock.OneHour.Adx[l-1].Pdi {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.FourHour:
				l := len(stock.FourHour.Adx) - 1
				if l < 1 || stock.FourHour.Adx[l].Pdi > stock.FourHour.Adx[l-1].Pdi {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			}
		}
	}
	if filter.ADX.MDIIncreasing == 1 {
		for _, stock := range stocks {
			switch filter.TimePeriod {
			case v.OneWeek:
				l := len(stock.OneWeek.Adx) - 1
				if l < 1 || stock.OneWeek.Adx[l].Mdi < stock.OneWeek.Adx[l-1].Mdi {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneDay:
				l := len(stock.OneDay.Adx) - 1
				if l < 1 || stock.OneDay.Adx[l].Mdi < stock.OneDay.Adx[l-1].Mdi {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneHour:
				l := len(stock.OneHour.Adx) - 1
				if l < 1 || stock.OneHour.Adx[l].Mdi < stock.OneHour.Adx[l-1].Mdi {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.FourHour:
				l := len(stock.FourHour.Adx) - 1
				if l < 1 || stock.FourHour.Adx[l].Mdi < stock.FourHour.Adx[l-1].Mdi {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			}
		}
	} else if filter.ADX.MDIIncreasing == 2 {
		for _, stock := range stocks {
			switch filter.TimePeriod {
			case v.OneWeek:
				l := len(stock.OneWeek.Adx) - 1
				if l < 1 || stock.OneWeek.Adx[l].Mdi > stock.OneWeek.Adx[l-1].Mdi {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneDay:
				l := len(stock.OneDay.Adx) - 1
				if l < 1 || stock.OneDay.Adx[l].Mdi > stock.OneDay.Adx[l-1].Mdi {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneHour:
				l := len(stock.OneHour.Adx) - 1
				if l < 1 || stock.OneHour.Adx[l].Mdi > stock.OneHour.Adx[l-1].Mdi {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.FourHour:
				l := len(stock.FourHour.Adx) - 1
				if l < 1 || stock.FourHour.Adx[l].Mdi > stock.FourHour.Adx[l-1].Mdi {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			}
		}
	}
	log.Println("Completed ADXFilter")
	return removeStocks(stocks, stocksToRemove)
}
func AccumulationDistributionFilter (stocks []v.Stock, filter FilterCheck) []v.Stock {
	log.Println("Running AccumulationDistributionFilter")
	var stocksToRemove []v.Stock
	if filter.AccumulationDistribution.Increasing == 1 {
		for _, stock := range stocks {
			switch filter.TimePeriod {
			case v.OneWeek:
				if stock.OneWeek.Volume.IncreasingVolume == 2 {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneDay:
				if stock.OneDay.Volume.IncreasingVolume == 2 {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneHour:
				if stock.OneHour.Volume.IncreasingVolume == 2 {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.FourHour:
				if stock.FourHour.Volume.IncreasingVolume == 2 {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			}
		}
	} else if filter.AccumulationDistribution.Increasing == 2 {
		for _, stock := range stocks {
			switch filter.TimePeriod {
			case v.OneWeek:
				if stock.OneWeek.Volume.IncreasingVolume == 1 {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneDay:
				if stock.OneDay.Volume.IncreasingVolume == 1 {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneHour:
				if stock.OneHour.Volume.IncreasingVolume == 1 {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.FourHour:
				if stock.FourHour.Volume.IncreasingVolume == 1 {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			}
		}
	}
	log.Println("Completed AccumulationDistributionFilter")
	return removeStocks(stocks, stocksToRemove)
}
func EMAFilter (stocks []v.Stock, filter FilterCheck) []v.Stock {
	log.Println("Running EMAFilter")
	var stocksToRemove []v.Stock
	if filter.EMA.Increasing == 1 {
		for _, stock := range stocks {
			switch filter.TimePeriod {
			case v.OneWeek:
				if stock.OneWeek.EMA.EMASlopePrevious > stock.OneWeek.EMA.EMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneDay:
				if stock.OneDay.EMA.EMASlopePrevious > stock.OneWeek.EMA.EMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneHour:
				if stock.OneHour.EMA.EMASlopePrevious > stock.OneWeek.EMA.EMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.FourHour:
				if stock.FourHour.EMA.EMASlopePrevious > stock.OneWeek.EMA.EMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			}
		}
	} else if filter.EMA.Increasing == 2 {
		for _, stock := range stocks {
			switch filter.TimePeriod {
			case v.OneWeek:
				if stock.OneWeek.EMA.EMASlopePrevious < stock.OneWeek.EMA.EMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneDay:
				if stock.OneDay.EMA.EMASlopePrevious < stock.OneWeek.EMA.EMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneHour:
				if stock.OneHour.EMA.EMASlopePrevious < stock.OneWeek.EMA.EMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.FourHour:
				if stock.FourHour.EMA.EMASlopePrevious < stock.OneWeek.EMA.EMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			}
		}
	}
	if filter.EMA.GreaterThan0 == 1 {
		for _, stock := range stocks {
			switch filter.TimePeriod {
			case v.OneWeek:
				if 0 > stock.OneWeek.EMA.EMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneDay:
				if 0 > stock.OneDay.EMA.EMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneHour:
				if 0 > stock.OneHour.EMA.EMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.FourHour:
				if 0 > stock.FourHour.EMA.EMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			}
		}
	} else if filter.EMA.GreaterThan0 == 2 {
		for _, stock := range stocks {
			switch filter.TimePeriod {
			case v.OneWeek:
				if 0 < stock.OneWeek.EMA.EMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneDay:
				if 0 < stock.OneDay.EMA.EMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneHour:
				if 0 < stock.OneHour.EMA.EMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.FourHour:
				if 0 < stock.FourHour.EMA.EMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			}
		}
	}
	if filter.EMA.GreaterThanSMA == 1 {
		for _, stock := range stocks {
			switch filter.TimePeriod {
			case v.OneWeek:
				if stock.OneWeek.SMA.SMA3 != 0 {
					if stock.OneWeek.EMA.EMA3 < stock.OneWeek.SMA.SMA3 {
						stocksToRemove = append(stocksToRemove, stock) 
					}
				}
			case v.OneDay:
				if stock.OneDay.SMA.SMA3 != 0 {
					if stock.OneDay.EMA.EMA3 < stock.OneDay.SMA.SMA3 {
						stocksToRemove = append(stocksToRemove, stock) 
					}
				}
			case v.OneHour:
				if stock.OneHour.SMA.SMA3 != 0 {
					if stock.OneHour.EMA.EMA3 < stock.OneHour.SMA.SMA3 {
						stocksToRemove = append(stocksToRemove, stock) 
					}
				}
			case v.FourHour:
				if stock.FourHour.SMA.SMA3 != 0 {
					if stock.FourHour.EMA.EMA3 < stock.FourHour.SMA.SMA3 {
						stocksToRemove = append(stocksToRemove, stock) 
					}
				}
			}
		}
	} else if filter.EMA.GreaterThanSMA == 2 {
		for _, stock := range stocks {
			switch filter.TimePeriod {
			case v.OneWeek:
				if stock.OneWeek.SMA.SMA3 != 0 {
					if stock.OneWeek.EMA.EMA3 > stock.OneWeek.SMA.SMA3 {
						stocksToRemove = append(stocksToRemove, stock) 
					}
				}
			case v.OneDay:
				if stock.OneDay.SMA.SMA3 != 0 {
					if stock.OneDay.EMA.EMA3 > stock.OneDay.SMA.SMA3 {
						stocksToRemove = append(stocksToRemove, stock) 
					}
				}
			case v.OneHour:
				if stock.OneHour.SMA.SMA3 != 0 {
					if stock.OneHour.EMA.EMA3 > stock.OneHour.SMA.SMA3 {
						stocksToRemove = append(stocksToRemove, stock) 
					}
				}
			case v.FourHour:
				if stock.FourHour.SMA.SMA3 != 0 {
					if stock.FourHour.EMA.EMA3 > stock.FourHour.SMA.SMA3 {
						stocksToRemove = append(stocksToRemove, stock) 
					}
				}
			}
		}
	}
	log.Println("Completed EMAFilter")
	return removeStocks(stocks, stocksToRemove)
}
func SMAFilter (stocks []v.Stock, filter FilterCheck) []v.Stock {
	log.Println("Running SMAFilter")
	var stocksToRemove []v.Stock
	if filter.SMA.Increasing == 1 {
		for _, stock := range stocks {
			switch filter.TimePeriod {
			case v.OneWeek:
				if stock.OneWeek.SMA.SMASlopePrevious > stock.OneWeek.SMA.SMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneDay:
				if stock.OneDay.SMA.SMASlopePrevious > stock.OneWeek.SMA.SMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneHour:
				if stock.OneHour.SMA.SMASlopePrevious > stock.OneWeek.SMA.SMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.FourHour:
				if stock.FourHour.SMA.SMASlopePrevious > stock.OneWeek.SMA.SMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			}
		}
	} else if filter.SMA.Increasing == 2 {
		for _, stock := range stocks {
			switch filter.TimePeriod {
			case v.OneWeek:
				if stock.OneWeek.SMA.SMASlopePrevious < stock.OneWeek.SMA.SMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneDay:
				if stock.OneDay.SMA.SMASlopePrevious < stock.OneWeek.SMA.SMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneHour:
				if stock.OneHour.SMA.SMASlopePrevious < stock.OneWeek.SMA.SMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.FourHour:
				if stock.FourHour.SMA.SMASlopePrevious < stock.OneWeek.SMA.SMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			}
		}
	}
	if filter.SMA.GreaterThan0 == 1 {
		for _, stock := range stocks {
			switch filter.TimePeriod {
			case v.OneWeek:
				if 0 > stock.OneWeek.SMA.SMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneDay:
				if 0 > stock.OneDay.SMA.SMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneHour:
				if 0 > stock.OneHour.SMA.SMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.FourHour:
				if 0 > stock.FourHour.SMA.SMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			}
		}
	} else if filter.SMA.GreaterThan0 == 2 {
		for _, stock := range stocks {
			switch filter.TimePeriod {
			case v.OneWeek:
				if 0 < stock.OneWeek.SMA.SMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneDay:
				if 0 < stock.OneDay.SMA.SMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneHour:
				if 0 < stock.OneHour.SMA.SMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.FourHour:
				if 0 < stock.FourHour.SMA.SMASlopeCurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			}
		}
	}
	log.Println("Completed SMAFilter")
	return removeStocks(stocks, stocksToRemove)
}
func RSIFilter (stocks []v.Stock, filter FilterCheck) []v.Stock {
	log.Println("Running RSIFilter")
	var stocksToRemove []v.Stock
	if filter.RSI.Increasing == 1 {
		for _, stock := range stocks {
			switch filter.TimePeriod {
			case v.OneWeek:
				if stock.OneWeek.RSI.RSIPrevious > stock.OneWeek.RSI.RSICurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneDay:
				if stock.OneDay.RSI.RSIPrevious > stock.OneWeek.RSI.RSICurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneHour:
				if stock.OneHour.RSI.RSIPrevious > stock.OneWeek.RSI.RSICurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.FourHour:
				if stock.FourHour.RSI.RSIPrevious > stock.OneWeek.RSI.RSICurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			}
		}
	} else if filter.RSI.Increasing == 2 {
		for _, stock := range stocks {
			switch filter.TimePeriod {
			case v.OneWeek:
				if stock.OneWeek.RSI.RSIPrevious < stock.OneWeek.RSI.RSICurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneDay:
				if stock.OneDay.RSI.RSIPrevious < stock.OneWeek.RSI.RSICurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneHour:
				if stock.OneHour.RSI.RSIPrevious < stock.OneWeek.RSI.RSICurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.FourHour:
				if stock.FourHour.RSI.RSIPrevious < stock.OneWeek.RSI.RSICurrent {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			}
		}
	}
	if filter.RSI.GreaterThan != 0 && filter.RSI.LessThan != 0 {
		for _, stock := range stocks {
			switch filter.TimePeriod {
			case v.OneWeek:
				if filter.RSI.GreaterThan > stock.OneWeek.RSI.RSICurrent || filter.RSI.LessThan < stock.OneWeek.RSI.RSICurrent{
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneDay:
				if filter.RSI.GreaterThan > stock.OneDay.RSI.RSICurrent || filter.RSI.LessThan < stock.OneDay.RSI.RSICurrent{
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneHour:
				if filter.RSI.GreaterThan > stock.OneHour.RSI.RSICurrent || filter.RSI.LessThan < stock.OneHour.RSI.RSICurrent{
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.FourHour:
				if filter.RSI.GreaterThan > stock.FourHour.RSI.RSICurrent || filter.RSI.LessThan < stock.FourHour.RSI.RSICurrent{
					stocksToRemove = append(stocksToRemove, stock) 
				}
			}
		}
	}
	log.Println("Completed RSIFilter")
	return removeStocks(stocks, stocksToRemove)
}
func AverageVolumeFilter (stocks []v.Stock, filter FilterCheck) []v.Stock {
	log.Println("Running AverageVolumeFilter")
	var stocksToRemove []v.Stock
	if filter.AverageVolume.GreaterThan != 0 {
		for _, stock := range stocks {
			switch filter.TimePeriod {
			case v.OneWeek:
				if stock.OneWeek.Volume.AverageVolume < filter.AverageVolume.GreaterThan {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneDay:
				if stock.OneDay.Volume.AverageVolume < filter.AverageVolume.GreaterThan {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneHour:
				if stock.OneHour.Volume.AverageVolume < filter.AverageVolume.GreaterThan {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.FourHour:
				if stock.FourHour.Volume.AverageVolume < filter.AverageVolume.GreaterThan {
					stocksToRemove = append(stocksToRemove, stock) 
				}
			}
		}
	}
	log.Println("Completed AverageVolumeFilter")
	return removeStocks(stocks, stocksToRemove)
}
func LastCloseFilter (stocks []v.Stock, filter FilterCheck) []v.Stock {
	log.Println("Running LastCloseFilter")
	var stocksToRemove []v.Stock
	if filter.LastClose.GreaterThan != 0 && filter.LastClose.LessThan != 0 {
		for _, stock := range stocks {
			switch filter.TimePeriod {
			case v.OneWeek:
				if stock.OneWeek.Close.LastClose < filter.LastClose.GreaterThan || stock.OneWeek.Close.LastClose > filter.LastClose.LessThan{
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneDay:
				if stock.OneDay.Close.LastClose < filter.LastClose.GreaterThan || stock.OneDay.Close.LastClose > filter.LastClose.LessThan{
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.OneHour:
				if stock.OneHour.Close.LastClose < filter.LastClose.GreaterThan || stock.OneHour.Close.LastClose > filter.LastClose.LessThan{
					stocksToRemove = append(stocksToRemove, stock) 
				}
			case v.FourHour:
				if stock.FourHour.Close.LastClose < filter.LastClose.GreaterThan || stock.FourHour.Close.LastClose > filter.LastClose.LessThan{
					stocksToRemove = append(stocksToRemove, stock) 
				}
			}
		}
	}
	log.Println("Completed LastCloseFilter")
	return removeStocks(stocks, stocksToRemove)
}
func BBFilter (stocks []v.Stock, filter FilterCheck) []v.Stock {
	log.Println("Running BBFilter")
	var stocksToRemove []v.Stock
	if filter.BB.GreaterThan != 0 && filter.BB.LessThan != 0 {
		for _, stock := range stocks {
			switch filter.TimePeriod {
			case v.OneWeek:
				if stock.OneWeek.SMA.SMA3 != 0 {
					if stock.OneWeek.BB.LowerBand > stock.OneWeek.SMA.SMA3 || stock.OneWeek.BB.UpperBand < stock.OneWeek.SMA.SMA3 {
						stocksToRemove = append(stocksToRemove, stock)
					}
				}
			case v.OneDay:
				if stock.OneWeek.SMA.SMA3 != 0 {
					if stock.OneDay.BB.LowerBand > stock.OneWeek.SMA.SMA3 || stock.OneDay.BB.UpperBand < stock.OneWeek.SMA.SMA3 {
						stocksToRemove = append(stocksToRemove, stock)
					}
				}
			case v.OneHour:
				if stock.OneWeek.SMA.SMA3 != 0 {
					if stock.OneHour.BB.LowerBand > stock.OneWeek.SMA.SMA3 || stock.OneHour.BB.UpperBand < stock.OneWeek.SMA.SMA3 {
						stocksToRemove = append(stocksToRemove, stock)
					}
				}
			case v.FourHour:
				if stock.OneWeek.SMA.SMA3 != 0 {
					if stock.FourHour.BB.LowerBand > stock.OneWeek.SMA.SMA3 || stock.FourHour.BB.UpperBand < stock.OneWeek.SMA.SMA3 {
						stocksToRemove = append(stocksToRemove, stock)
					}
				}
			}
		}
	}
	log.Println("Completed BBFilter")
	return removeStocks(stocks, stocksToRemove)
}
func VolumeJumpFilter(stocks []v.Stock, filter FilterCheck) []v.Stock {
	log.Println("Running VolumeJumpFilter")
	var stocksToRemove []v.Stock
	for _, stock := range stocks {
		switch filter.TimePeriod {
		case v.OneWeek:
			if stock.OneWeek.Volume.VolumeJump < filter.VolumeJump.Multiplier {
				stocksToRemove = append(stocksToRemove, stock)
			}
		case v.OneDay:
			if stock.OneDay.Volume.VolumeJump < filter.VolumeJump.Multiplier {
				stocksToRemove = append(stocksToRemove, stock)
			}
		case v.OneHour:
			if stock.OneHour.Volume.VolumeJump < filter.VolumeJump.Multiplier {
				stocksToRemove = append(stocksToRemove, stock)
			}
		case v.FourHour:
			if stock.FourHour.Volume.VolumeJump < filter.VolumeJump.Multiplier {
				stocksToRemove = append(stocksToRemove, stock)
			}
		}
	}
	log.Println("Completed VolumeJumpFilter")
	return removeStocks(stocks, stocksToRemove)
}
func CloseGreaterThanPreviousFilter(stocks []v.Stock, filter FilterCheck) []v.Stock {
	log.Println("Running VolumeJumpFilter")
	var stocksToRemove []v.Stock
	if filter.CloseGreaterThanPrevious.Greater == 1 {
		for _, stock := range stocks {
			switch filter.TimePeriod {
			case v.OneWeek:
				if stock.OneWeek.Close.GreaterThanLast == 2 {
					stocksToRemove = append(stocksToRemove, stock)
				}
			case v.OneDay:
				if stock.OneDay.Close.GreaterThanLast == 2 {
					stocksToRemove = append(stocksToRemove, stock)
				}
			case v.OneHour:
				if stock.OneHour.Close.GreaterThanLast == 2 {
					stocksToRemove = append(stocksToRemove, stock)
				}
			case v.FourHour:
				if stock.FourHour.Close.GreaterThanLast == 2 {
					stocksToRemove = append(stocksToRemove, stock)
				}
			}
		}
	} else if filter.CloseGreaterThanPrevious.Greater == 2 {
		for _, stock := range stocks {
			switch filter.TimePeriod {
			case v.OneWeek:
				if stock.OneWeek.Close.GreaterThanLast == 1 {
					stocksToRemove = append(stocksToRemove, stock)
				}
			case v.OneDay:
				if stock.OneDay.Close.GreaterThanLast == 1 {
					stocksToRemove = append(stocksToRemove, stock)
				}
			case v.OneHour:
				if stock.OneHour.Close.GreaterThanLast == 1 {
					stocksToRemove = append(stocksToRemove, stock)
				}
			case v.FourHour:
				if stock.FourHour.Close.GreaterThanLast == 1 {
					stocksToRemove = append(stocksToRemove, stock)
				}
			}
		}

	}
	log.Println("Completed VolumeJumpFilter")
	return removeStocks(stocks, stocksToRemove)
}
/*
	Removes stock from slice
 */
func removeStocks (slice []v.Stock, toRemove []v.Stock) []v.Stock {
	for _, r := range toRemove {
		for i, s := range slice {
			if r.Symbol == s.Symbol {
				slice = remove(slice, i)
				break
			}
		}
	}
	return slice
}
func remove (slice []v.Stock, s int) []v.Stock {
	return append(slice[:s], slice[s+1:]...)
}