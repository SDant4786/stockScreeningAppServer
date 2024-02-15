package filtering

import (
	"../repos/alpaca/alpaca"
	"../repos/alpaca/common"
	v "../variables"
	"log"
	"time"
)

//TODO fix this so it always ends on last bar
func GetFourHourBars(stock v.Stock) []v.Bar {
	var bars []v.Bar
	if len(stock.HourBars) > 0 {
		bars = append(bars, stock.HourBars[0])
		for l := 3; l < len(stock.HourBars); l += 4 {
			bars = append(bars, stock.HourBars[l])
		}
	}
	return bars
}
func GetSymbols() []v.Stock{
	var stocks []v.Stock
	var filteredSymbols []string
	alpaca.SetBaseUrl("https://api.alpaca.markets")
	symbols, err := InitClient().GetAssett()

	if err != nil {
		log.Fatal("Error getting stock list from Alpaca:", err.Error())
	}
	filteredSymbols = filterSymbols(symbols)
	for _, filteredSymbol := range filteredSymbols {
		stocks = append(stocks, v.Stock{
			Symbol:       filteredSymbol,
			WeekBars:     []v.Bar{
				v.Bar{
					Time:   time.Time{},
					Open:   0,
					High:   0,
					Low:    0,
					Close:  0,
					Volume: 0,
				},
			},
			DayBars:      []v.Bar{
				v.Bar{
					Time:   time.Time{},
					Open:   0,
					High:   0,
					Low:    0,
					Close:  0,
					Volume: 0,
				},
			},
			FourHourBars: []v.Bar{
				v.Bar{
					Time:   time.Time{},
					Open:   0,
					High:   0,
					Low:    0,
					Close:  0,
					Volume: 0,
				},
			},
			HourBars:     []v.Bar{
				v.Bar{
					Time:   time.Time{},
					Open:   0,
					High:   0,
					Low:    0,
					Close:  0,
					Volume: 0,
				},
			},
			OneWeek:      v.FilterResults{
				EMA:              v.EMAResult{
					EMA1:                    0,
					EMA2:                    0,
					EMA3:                    0,
					EMASlopeCurrent:         0,
					EMASlopePrevious:        0,
					EMASlopeIncreasePercent: 0,
				},
				SMA:              v.SMAResult{
					SMA1:                    0,
					SMA2:                    0,
					SMA3:                    0,
					SMASlopeCurrent:         0,
					SMASlopePrevious:        0,
					SMASlopeIncreasePercent: 0,
				},
				RSI:              v.RSIResult{
					RSICurrent:         0,
					RSIPrevious:        0,
					RSIPercentIncrease: 0,
				},
				Volume: 		  v.VolumeResult {
					IncreasingVolume: 0,
					AverageVolume:    0,
					VolumeJump:       0,
				},
				Close: 			  v.CloseResult{
					LastClose:        0,
					GreaterThanLast:  0,
				},
				Adx:              []v.AdxResult{
					v.AdxResult{
						Pdi: 0,
						Mdi: 0,
						Adx: 0,
					},
				},
				BB:               v.BBResult{
					UpperBand: 0,
					LowerBand: 0,
				},
			},
			OneDay:        v.FilterResults{
				EMA:              v.EMAResult{
					EMA1:                    0,
					EMA2:                    0,
					EMA3:                    0,
					EMASlopeCurrent:         0,
					EMASlopePrevious:        0,
					EMASlopeIncreasePercent: 0,
				},
				SMA:              v.SMAResult{
					SMA1:                    0,
					SMA2:                    0,
					SMA3:                    0,
					SMASlopeCurrent:         0,
					SMASlopePrevious:        0,
					SMASlopeIncreasePercent: 0,
				},
				RSI:              v.RSIResult{
					RSICurrent:         0,
					RSIPrevious:        0,
					RSIPercentIncrease: 0,
				},
				Volume: 		  v.VolumeResult {
					IncreasingVolume: 0,
					AverageVolume:    0,
					VolumeJump:       0,
				},
				Close: 			  v.CloseResult{
					LastClose:        0,
					GreaterThanLast:  0,
				},
				Adx:               []v.AdxResult{
					v.AdxResult{
						Pdi: 0,
						Mdi: 0,
						Adx: 0,
					},
				},
				BB:               v.BBResult{
					UpperBand: 0,
					LowerBand: 0,
				},
			},
			OneHour:       v.FilterResults{
				EMA:              v.EMAResult{
					EMA1:                    0,
					EMA2:                    0,
					EMA3:                    0,
					EMASlopeCurrent:         0,
					EMASlopePrevious:        0,
					EMASlopeIncreasePercent: 0,
				},
				SMA:              v.SMAResult{
					SMA1:                    0,
					SMA2:                    0,
					SMA3:                    0,
					SMASlopeCurrent:         0,
					SMASlopePrevious:        0,
					SMASlopeIncreasePercent: 0,
				},
				RSI:              v.RSIResult{
					RSICurrent:         0,
					RSIPrevious:        0,
					RSIPercentIncrease: 0,
				},
				Volume: 		  v.VolumeResult {
					IncreasingVolume: 0,
					AverageVolume:    0,
					VolumeJump:       0,
				},
				Close: 			  v.CloseResult{
					LastClose:        0,
					GreaterThanLast:  0,
				},
				Adx:               []v.AdxResult{
					v.AdxResult{
						Pdi: 0,
						Mdi: 0,
						Adx: 0,
					},
				},
				BB:               v.BBResult{
					UpperBand: 0,
					LowerBand: 0,
				},
			},
			FourHour:      v.FilterResults{
				EMA:              v.EMAResult{
					EMA1:                    0,
					EMA2:                    0,
					EMA3:                    0,
					EMASlopeCurrent:         0,
					EMASlopePrevious:        0,
					EMASlopeIncreasePercent: 0,
				},
				SMA:              v.SMAResult{
					SMA1:                    0,
					SMA2:                    0,
					SMA3:                    0,
					SMASlopeCurrent:         0,
					SMASlopePrevious:        0,
					SMASlopeIncreasePercent: 0,
				},
				RSI:              v.RSIResult{
					RSICurrent:         0,
					RSIPrevious:        0,
					RSIPercentIncrease: 0,
				},
				Volume: 		  v.VolumeResult {
					IncreasingVolume: 0,
					AverageVolume:    0,
					VolumeJump:       0,
				},
				Close: 			  v.CloseResult{
					LastClose:        0,
					GreaterThanLast:  0,
				},
				Adx:               []v.AdxResult{
					v.AdxResult{
						Pdi: 0,
						Mdi: 0,
						Adx: 0,
					},
				},
				BB:               v.BBResult{
					UpperBand: 0,
					LowerBand: 0,
				},
			},
			PossibleTop:  false,
			Sell:         false,
		})
	}

	return stocks
}
func filterSymbols(symbols []alpaca.Asset) []string {
	var filteredSymbols []string
	for _, symbol := range symbols {
		//if asset.Tradable == true && (asset.Exchange == "NASDAQ" || asset.Exchange == "NYSE") {
		if symbol.Tradable == true{
			filteredSymbols = append(filteredSymbols, symbol.Symbol)
		}
	}
	return filteredSymbols
}
func InitClient() *alpaca.Client{
	credentials := &common.APIKey{
		ID:           "AK6Z4JA90P2JR0BSR3OG",
		Secret:       "kFPrbHr3YUA9giQ4fLmz4eFfI91od4m8cFj0vCHR",
		OAuth:        "",
		PolygonKeyID: "",
	}
	return alpaca.NewClient(credentials)
}