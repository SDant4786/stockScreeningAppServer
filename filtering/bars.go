package filtering

import (
	"../repos/finance-go"
	"../repos/finance-go/chart"
	"../repos/finance-go/datetime"
	v "../variables"
	"context"
	"log"
	"sync"
	"time"
)
/*
	One Week
 */
func CreateAndRunThreadsForGettingOneWeekBars (stocks []v.Stock, periods int) []v.Stock {
	log.Println("Running CreateAndRunThreadsForGettingOneWeekBars")
	var stocksSlice = v.ThreadSafeSlice{
		Mutex: &sync.Mutex{},
		Slice: []v.Stock{},
	}
	var wg sync.WaitGroup

	lenFilters := len(stocks)
	splitDifference := lenFilters / 40
	splitStart := 0
	splitEnd := splitDifference

	for i:=0; i < 41; i++ {
		wg.Add(1)
		if i == 40 {
			go GetOneWeekBars(stocks[splitStart:], &wg, &stocksSlice, periods)
		} else {
			go GetOneWeekBars(stocks[splitStart:splitEnd], &wg, &stocksSlice, periods)
			splitStart += splitDifference
			splitEnd += splitDifference
		}
	}
	wg.Wait()
	log.Println("Completed CreateAndRunThreadsForGettingOneWeekBars")
	return stocksSlice.Slice
}
func GetOneWeekBars(stocks []v.Stock, wg *sync.WaitGroup, stocksSlice *v.ThreadSafeSlice, periods int) {
	var filteredStocks []v.Stock
	var cont = context.TODO()
	var client = chart.GetC()

	var end = time.Now()
	var start = end.AddDate(0, -periods*10, 0)

	var params = chart.Params{
		Params:     finance.Params{Context: &cont},
		Symbol:     "",
		Start:      datetime.New(&start),
		End:        datetime.New(&end),
		Interval:   datetime.FiveDay,
		IncludeExt: false,
	}

	for _, stock := range stocks {
		var bars []v.Bar
		params.Symbol = stock.Symbol
		it := client.Get(&params)
		for it.Next() {
			b := it.Bar()
			close := b.Close
			if b.AdjClose != 0 {
				close = b.AdjClose
			}
			tm := time.Unix(int64(b.Timestamp), 0)

			bars = append(bars, v.Bar{
				Time:   tm,
				Open:   float64(b.Open),
				High:   float64(b.High),
				Low:    float64(b.Low),
				Close:  float64(close),
				Volume: float64(b.Volume),
			})
		}
		//if len(bars) > 112 {
			stock.WeekBars = bars
			filteredStocks = append(filteredStocks, stock)
		//}
	}

	stocksSlice.Mutex.Lock()
	stocksSlice.Slice = append(stocksSlice.Slice, filteredStocks...)
	stocksSlice.Mutex.Unlock()
	wg.Done()
}
/*
	One Day
 */
func CreateAndRunThreadsForGettingOneDayBars (stocks []v.Stock, periods int) []v.Stock {
	log.Println("Running CreateAndRunThreadsForGettingOneDayBars")
	var stocksSlice = v.ThreadSafeSlice{
		Mutex: &sync.Mutex{},
		Slice: []v.Stock{},
	}
	var wg sync.WaitGroup

	lenFilters := len(stocks)
	splitDifference := lenFilters / 40
	splitStart := 0
	splitEnd := splitDifference

	for i:=0; i < 41; i++ {
		wg.Add(1)
		if i == 40 {
			go GetOneDayBars(stocks[splitStart:], &wg, &stocksSlice, periods)
		} else {
			go GetOneDayBars(stocks[splitStart:splitEnd], &wg, &stocksSlice, periods)
			splitStart += splitDifference
			splitEnd += splitDifference
		}
	}
	wg.Wait()
	log.Println("Completed CreateAndRunThreadsForGettingOneDayBars")
	return stocksSlice.Slice
}
func GetOneDayBars(stocks []v.Stock, wg *sync.WaitGroup, stocksSlice *v.ThreadSafeSlice, periods int) {
	var filteredStocks []v.Stock
	var cont = context.TODO()
	var client = chart.GetC()

	var end = time.Now()
	var start = end.AddDate(0, 0, -periods * 15)

	var params = chart.Params{
		Params:     finance.Params{Context: &cont},
		Symbol:     "",
		Start:      datetime.New(&start),
		End:        datetime.New(&end),
		Interval:   datetime.OneDay,
		IncludeExt: false,
	}

	for _, stock := range stocks {
		var bars []v.Bar
		params.Symbol = stock.Symbol
		it := client.Get(&params)
		for it.Next() {
			b := it.Bar()
			close := b.Close
			if b.AdjClose != 0 {
				close = b.AdjClose
			}

			tm := time.Unix(int64(b.Timestamp), 0)

			bars = append(bars, v.Bar {
				Time:   tm,
				Open:   float64(b.Open),
				High:   float64(b.High),
				Low:    float64(b.Low),
				Close:  float64(close),
				Volume: float64(b.Volume),
			})
		}
		//if len(bars) > 50 {
			stock.DayBars = bars
			filteredStocks = append(filteredStocks, stock)
		//}
	}
	stocksSlice.Mutex.Lock()
	stocksSlice.Slice = append(stocksSlice.Slice, filteredStocks...)
	stocksSlice.Mutex.Unlock()
	wg.Done()
}
/*
	One Hour
 */
func CreateAndRunThreadsForGettingOneHourBars (stocks []v.Stock, periods int) []v.Stock {
	log.Println("Running CreateAndRunThreadsForGettingOneHourBars")
	var stocksSlice = v.ThreadSafeSlice{
		Mutex: &sync.Mutex{},
		Slice: []v.Stock{},
	}
	var wg sync.WaitGroup

	lenFilters := len(stocks)
	splitDifference := lenFilters / 40
	splitStart := 0
	splitEnd := splitDifference

	for i:=0; i < 41; i++ {
		wg.Add(1)
		if i == 40 {
			go GetOneHrBars(stocks[splitStart:], &wg, &stocksSlice, periods)
		} else {
			go GetOneHrBars(stocks[splitStart:splitEnd], &wg, &stocksSlice, periods)
			splitStart += splitDifference
			splitEnd += splitDifference
		}
	}
	wg.Wait()
	log.Println("Completed CreateAndRunThreadsForGettingOneHourBars")
	return stocksSlice.Slice
}
func GetOneHrBars(stocks []v.Stock, wg *sync.WaitGroup, stocksSlice *v.ThreadSafeSlice, periods int) {
	var filteredStocks []v.Stock
	var cont = context.TODO()
	var client = chart.GetC()

	var end = time.Now()
	var start = end.AddDate(0, 0, -periods * 10)

	var params = chart.Params{
		Params:     finance.Params{Context: &cont},
		Symbol:     "",
		Start:      datetime.New(&start),
		End:        datetime.New(&end),
		Interval:   datetime.OneHour,
		IncludeExt: false,
	}

	if end.Minute() >= 0 && end.Minute() < 30 {
		var min = end.Minute()
		end.Add(time.Duration(-min) * time.Minute)
	} else {
		var min = end.Minute()
		min = min - 30
		end.Add(time.Duration(-min) * time.Minute)
	}
	var second = end.Second() + 1
	end.Add(time.Duration(-second) * time.Second)

	for _, stock := range stocks {
		var bars []v.Bar
		params.Symbol = stock.Symbol
		it := client.Get(&params)
		for it.Next() {
			b := it.Bar()
			close := b.Close

			tm := time.Unix(int64(b.Timestamp), 0)

			bars = append(bars, v.Bar {
				Time:   tm,
				Open:   float64(b.Open),
				High:   float64(b.High),
				Low:    float64(b.Low),
				Close:  float64(close),
				Volume: float64(b.Volume),
			})

		}
		//if len(bars) > 116 {
			stock.HourBars = bars
			stock.FourHourBars = GetFourHourBars(stock)
			filteredStocks = append(filteredStocks, stock)
		//}
	}


	stocksSlice.Mutex.Lock()
	stocksSlice.Slice = append(stocksSlice.Slice, filteredStocks...)
	stocksSlice.Mutex.Unlock()
	wg.Done()
}