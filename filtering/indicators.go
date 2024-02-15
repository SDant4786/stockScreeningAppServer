package filtering

import (
	v "../variables"
	"log"
	"math"
)

/*
	ADX
 */
func Adx (stocks []v.Stock, timeFrame string) []v.Stock {
	log.Println("Running Adx")
	for i, stock := range stocks {
		switch timeFrame {
		case v.OneWeek:
			stocks[i].OneWeek.Adx = CalcAdxHr(stock.WeekBars)
		case v.OneDay:
			stocks[i].OneDay.Adx = CalcAdxHr(stock.DayBars)
		case v.OneHour:
			stocks[i].OneHour.Adx = CalcAdxHr(stock.HourBars)
		case v.FourHour:
			stocks[i].FourHour.Adx = CalcAdxHr(stock.FourHourBars)
		}
	}
	log.Println("Completed Adx")
	return stocks
}
func CalcAdxHr (bars []v.Bar) []v.AdxResult{
	var start = len(bars) - 30
	if start < 0 {
		return []v.AdxResult{
			v.AdxResult{
				Pdi: 0,
				Mdi: 0,
				Adx: 0,
			},
		}
	}
	bars = bars[start:]
	var prevHigh float64
	var prevLow float64
	var prevTrs float64
	var prevPdm float64
	var prevMdm float64
	var prevAdx float64
	var sumTr float64
	var sumPdm float64
	var sumMdm float64
	var sumDx float64
	var lookBackPeriod = 14
	var results []v.AdxResult
	var atrResults = CalcTrueRangeHr(bars)

	for i := 0; i < len(bars); i++ {
		b := bars[i]
		index := i + 1
		result := v.AdxResult{}

		if index == 1 {
			results = append(results, result)
			prevHigh = b.High
			prevLow = b.Low
			continue
		}

		var tr = atrResults[i].Tr

		var pdm1 float64
		var mdm1 float64
		if (b.High - prevHigh) > (prevLow - b.Low) {
			pdm1 = math.Max(b.High-prevHigh, 0)
		} else {
			pdm1 = 0
		}

		if (prevLow - b.Low) > (b.High - prevHigh) {
			mdm1 = math.Max(prevLow-b.Low, 0)
		} else {
			mdm1 = 0
		}

		prevHigh = b.High
		prevLow = b.Low

		if index <= lookBackPeriod + 1 {
			sumTr += tr
			sumPdm += pdm1
			sumMdm += mdm1
		}

		if index <= lookBackPeriod {
			results = append(results, result)
			continue
		}

		var trs float64
		var pdm float64
		var mdm float64

		if index == lookBackPeriod + 1 {
			trs = sumTr
			pdm = sumPdm
			mdm = sumMdm
		} else {
			trs = prevTrs - (prevTrs / float64(lookBackPeriod)) + tr
			pdm = prevPdm - (prevPdm / float64(lookBackPeriod)) + pdm1
			mdm = prevMdm - (prevMdm / float64(lookBackPeriod)) + mdm1
		}

		prevTrs = trs
		prevPdm = pdm
		prevMdm = mdm

		var pdi = 100 * pdm / trs
		var mdi = 100 * mdm / trs
		var dx = 100 * math.Abs((pdi-mdi) / (pdi+mdi))

		result.Pdi = pdi
		result.Mdi = mdi

		var adx float64

		if index > 2 * lookBackPeriod {
			adx = (prevAdx * (float64(lookBackPeriod)-1) + dx) / float64(lookBackPeriod)
			result.Adx = adx
			prevAdx = adx
		} else if index == 2 * lookBackPeriod {
			sumDx += dx
			adx = sumDx / float64(lookBackPeriod)
			result.Adx = adx
			prevAdx = adx
		} else {
			sumDx += dx
		}
		if math.IsNaN(adx) {
			return []v.AdxResult{
				v.AdxResult{
					Pdi: 0,
					Mdi: 0,
					Adx: 0,
				},
			}
		}
		results = append(results, result)
	}

	return results
}
func CalcTrueRangeHr (bars []v.Bar) []v.AtrResult{
	var prevAtr float64
	var prevClose float64
	var highMinusPrevClose float64
	var lowMinusPrevClose float64
	var sumTr float64
	var lookbackPeriod float64 = 14
	var results []v.AtrResult

	for i := 0; i < len(bars); i++ {
		h := bars[i]
		var index float64 = float64(i) + 1

		var result v.AtrResult

		if index > 1 {
			highMinusPrevClose = math.Abs(h.High - prevClose)
			lowMinusPrevClose = math.Abs(h.Low - prevClose)
		}

		tr := math.Max(h.High - h.Low, math.Max(highMinusPrevClose, lowMinusPrevClose))
		result.Tr = tr

		if index > lookbackPeriod {
			// calculate ATR
			result.Atr = (prevAtr * (lookbackPeriod - 1) + tr) / lookbackPeriod
			if h.Close == 0 {
				result.Atrp = 0
			} else {
				result.Atrp = (result.Atr / h.Close) * 100
			}
			prevAtr = result.Atr
		} else if index == lookbackPeriod {
			// initialize ATR
			sumTr += tr
			result.Atr = sumTr / lookbackPeriod
			if h.Close == 0 {
				result.Atrp = 0
			} else {
				result.Atrp = (result.Atr / h.Close) * 100
			}
			prevAtr = result.Atr
		} else
		{
			// only used for periods before ATR initialization
			sumTr += tr
		}

		results = append(results, result)
		prevClose = h.Close
	}
	return results
}
/*
	Accumulation Distribution
 */
func AccumulationDistribution(stocks []v.Stock, timeFrame string, period int) []v.Stock{
	log.Println("Running AccumulationDistribution")
		for i, stock := range stocks {
			switch timeFrame {
			case v.OneWeek:
				stocks[i].OneWeek.Volume.IncreasingVolume = CalcAccumulationDistribution(stock.WeekBars, period)
			case v.OneDay:
				stocks[i].OneDay.Volume.IncreasingVolume = CalcAccumulationDistribution(stock.DayBars, period)
			case v.OneHour:
				stocks[i].OneHour.Volume.IncreasingVolume = CalcAccumulationDistribution(stock.HourBars, period)
			case v.FourHour:
				stocks[i].FourHour.Volume.IncreasingVolume = CalcAccumulationDistribution(stock.FourHourBars, period)
			}
		}
	log.Println("Completed AccumulationDistribution")
	return stocks
}
func CalcAccumulationDistribution (bars[]v.Bar, period int) int {
	var b []float64
	var missingVolume bool = true
	var start = len(bars) - period - 1
	if start < 0 {
		return 2
	}
	bars = bars[start:]

	for i, bar := range bars {
		ad := (((bar.Close - bar.Low) - (bar.High - bar.Close)) / (bar.High - bar.Low)) * bar.Volume
		if bar.Volume == 0 {
			missingVolume = true
			break
		}
		if i > 0 {
			ad += b[i-1]
		}
		b = append(b, ad)
	}

	if missingVolume == true {
		return 0
	} else if b[len(b)-1] >= b[len(b)-2] {
		return 1
	} else {
		return 2
	}
}
func CalcIncreasingVolume (bars []v.Bar) bool {
	var up float64
	var down float64
	for i := 0; i < len(bars) - 1; i++ {
		if bars[i+1].Volume == 0 || bars[i].Volume == 0 {
			continue
		}
		if bars[i+1].Volume >= bars[i].Volume {
			up++
		} else {
			down++
		}
	}
	if up > down {
		return true
	}

	return false
}
/*
	EMA
 */
func EMA (stocks []v.Stock, timeFrame string, period int) []v.Stock {
	log.Println("Running EMA")
	for i, stock := range stocks {
		switch timeFrame {
		case v.OneWeek:
			stocks[i].OneWeek.EMA = CalcEMA(stock.WeekBars, period)
		case v.OneDay:
			stocks[i].OneDay.EMA = CalcEMA(stock.DayBars, period)
		case v.OneHour:
			stocks[i].OneHour.EMA = CalcEMA(stock.HourBars, period)
		case v.FourHour:
			stocks[i].FourHour.EMA = CalcEMA(stock.FourHourBars, period)
		}
	}
	log.Println("Completed EMA")
	return stocks

}
func CalcEMA (bars[]v.Bar, period int) v.EMAResult {
	var emaResult v.EMAResult
	var ema1 float64
	var ema2 float64
	var ema3 float64
	var emaSlope1 float64
	var emaSlope2 float64
	var multiplier float64
	var emaStart = len(bars) - period - 2
	if emaStart < 0 {
		return v.EMAResult{}
	}

	l := float64(period) + 1
	m := 2 / l
	multiplier = m

	ema1 = bars[emaStart].Close
	ema2 = bars[emaStart+1].Close
	ema3 = bars[emaStart+2].Close

	for k := emaStart + 1; k < len(bars)-2; k++ {
		ema1 = (bars[k].Close * multiplier) + (ema1 * (1 - multiplier))
	}
	for k := emaStart + 2; k < len(bars)-1; k++ {
		ema2 = (bars[k].Close * multiplier) + (ema2 * (1 - multiplier))
	}
	for k := emaStart + 3; k < len(bars); k++ {
		ema3 = (bars[k].Close * multiplier) + (ema3 * (1 - multiplier))
	}

	emaSlope1 = (ema2 - ema1) / ema2
	emaSlope2 = (ema3 - ema2) / ema3

	emaResult.EMA1 = ema1
	emaResult.EMA2 = ema2
	emaResult.EMA3 = ema3
	emaResult.EMASlopeCurrent = emaSlope2
	emaResult.EMASlopePrevious = emaSlope1
	emaResult.EMASlopeIncreasePercent = math.Abs((emaSlope2 / emaSlope1) * 100)
	if emaResult.EMA3 == 0 {
		emaResult.EMASlopeCurrent = 0
	}
	if emaResult.EMA2 == 0 {
		emaResult.EMASlopePrevious = 0
	}
	if emaResult.EMASlopePrevious == 0 {
		emaResult.EMASlopeIncreasePercent = 0
	}

	return emaResult
}
/*
	SMA
 */
func SMA (stocks []v.Stock, timeFrame string, period int) []v.Stock {
	log.Println("Running SMA")
	for i, stock := range stocks {
		switch timeFrame {
		case v.OneWeek:
			stocks[i].OneWeek.SMA = CalcSMA(stock.WeekBars, period)
		case v.OneDay:
			stocks[i].OneDay.SMA = CalcSMA(stock.DayBars, period)
		case v.OneHour:
			stocks[i].OneHour.SMA = CalcSMA(stock.HourBars, period)
		case v.FourHour:
			stocks[i].FourHour.SMA = CalcSMA(stock.FourHourBars, period)
		}
	}
	log.Println("Completed SMA")
	return stocks
}
func CalcSMA (bars[]v.Bar, period int) v.SMAResult {
	var smaResult v.SMAResult
	var sma1 float64
	var sma2 float64
	var sma3 float64
	var smaSlope1 float64
	var smaSlope2 float64
	var smaStart = len(bars) - period - 2
	if smaStart < 0 {
		return v.SMAResult{}
	}

	sma1 = bars[smaStart].Close
	sma2 = bars[smaStart+1].Close
	sma3 = bars[smaStart+2].Close

	for k := smaStart + 1; k < len(bars)-2; k++ {
		sma1 += bars[k].Close
	}
	for k := smaStart + 2; k < len(bars)-1; k++ {
		sma2 += bars[k].Close
	}
	for k := smaStart + 3; k < len(bars); k++ {
		sma3 += bars[k].Close
	}

	sma1 /= float64(period)
	sma2 /= float64(period)
	sma3 /= float64(period)

	smaSlope1 = (sma2 - sma1) / sma2
	smaSlope2 = (sma3 - sma2) / sma3

	smaResult.SMA1 = sma1
	smaResult.SMA2 = sma2
	smaResult.SMA3 = sma3
	smaResult.SMASlopeCurrent = smaSlope2
	smaResult.SMASlopePrevious = smaSlope1
	smaResult.SMASlopeIncreasePercent = math.Abs((smaSlope2 / smaSlope1) * 100)

	if smaResult.SMA3 == 0 {
		smaResult.SMASlopeCurrent = 0
	}
	if smaResult.SMA2 == 0 {
		smaResult.SMASlopePrevious = 0
	}
	if smaResult.SMASlopePrevious == 0 {
		smaResult.SMASlopeIncreasePercent = 0
	}
	return smaResult
}
/*
	RSI
 */
func RSI (stocks []v.Stock, timeFrame string, period int) []v.Stock {
	log.Println("Running RSI")
	for i, stock := range stocks {
		switch timeFrame {
		case v.OneWeek:
			stocks[i].OneWeek.RSI = CalcRSI(stock.WeekBars, period)
		case v.OneDay:
			stocks[i].OneDay.RSI = CalcRSI(stock.DayBars, period)
		case v.OneHour:
			stocks[i].OneHour.RSI = CalcRSI(stock.HourBars, period)
		case v.FourHour:
			stocks[i].FourHour.RSI = CalcRSI(stock.FourHourBars, period)
		}
	}
	log.Println("Completed SMA")
	return stocks
}
func CalcRSI(bars []v.Bar, period int) v.RSIResult {
	var rsiResult v.RSIResult
	var totalGain float64
	var totalLoss float64
	var rs1 float64
	var rs2 float64
	var rsi1 float64
	var rsi2 float64

	var start = len(bars) - period - 1
	if start < 1 {
		return v.RSIResult{}
	}
	for i := start; i < len(bars)-1; i++ {

		prevClose := bars[i-1].Close
		curClose := bars[i].Close

		diff := curClose - prevClose

		if diff >= 0 {
			totalGain = totalGain + diff
		} else {
			totalLoss = totalLoss + diff
		}

	}
	rs1 = (totalGain / float64(period)) / math.Abs(totalLoss/float64(period))
	rsi1 = 100 - (100 / (1 + rs1))

	for i := start + 1; i < len(bars); i++ {

		prevClose := bars[i-1].Close
		curClose := bars[i].Close

		diff := curClose - prevClose

		if diff >= 0 {
			totalGain = totalGain + diff
		} else {
			totalLoss = totalLoss + diff
		}

	}
	rs2 = (totalGain / float64(period)) / math.Abs(totalLoss/float64(period))
	rsi2 = 100 - (100 / (1 + rs2))
	if math.IsNaN(rsi1) {
		rsiResult.RSIPrevious = -1
	} else {
		rsiResult.RSIPrevious = rsi1
	}
	if math.IsNaN(rsi2) {
		rsiResult.RSICurrent = -1
	} else {
		rsiResult.RSICurrent = rsi2
	}
	if rsiResult.RSIPrevious == -1 || rsiResult.RSICurrent == -1 {
		rsiResult.RSIPercentIncrease = -1
	} else {
		rsiResult.RSIPercentIncrease = (rsi1 / rsi2) * 100
	}

	return rsiResult
}
/*
	Average Volume
 */
func AverageVolume(stocks []v.Stock, timeFrame string) []v.Stock {
	log.Println("Running AverageVolume")
	for i, stock := range stocks {
		switch timeFrame {
		case v.OneWeek:
			stocks[i].OneWeek.Volume.AverageVolume = CalcAverageVolume(stock.WeekBars)
		case v.OneDay:
			stocks[i].OneDay.Volume.AverageVolume = CalcAverageVolume(stock.DayBars)
		case v.OneHour:
			stocks[i].OneHour.Volume.AverageVolume = CalcAverageVolume(stock.HourBars)
		case v.FourHour:
			stocks[i].FourHour.Volume.AverageVolume = CalcAverageVolume(stock.FourHourBars)
		}
	}
	log.Println("Completed AverageVolume")
	return stocks
}
func CalcAverageVolume (bars []v.Bar) float64 {
	var totalVolume float64
	for _, bar := range bars {
		if bar.Volume != 0 {
			totalVolume = totalVolume + bar.Volume
		}
	}
	totalVolume = totalVolume / float64(len(bars))
	return totalVolume
}
/*
	Volume Jump
 */
func VolumeJump(stocks []v.Stock, timeFrame string) []v.Stock {
	log.Println("Running VolumeJump")
	for i, stock := range stocks {
		switch timeFrame {
		case v.OneWeek:
			stocks[i].OneWeek.Volume.VolumeJump = CalcVolumeJump(stock.WeekBars)
		case v.OneDay:
			stocks[i].OneDay.Volume.VolumeJump = CalcVolumeJump(stock.DayBars)
		case v.OneHour:
			stocks[i].OneHour.Volume.VolumeJump = CalcVolumeJump(stock.HourBars)
		case v.FourHour:
			stocks[i].FourHour.Volume.VolumeJump = CalcVolumeJump(stock.FourHourBars)
		}
	}
	log.Println("Completed VolumeJump")
	return stocks
}
func CalcVolumeJump (bars []v.Bar) float64 {
	l := len(bars) - 1
	if l < 1 {
		return 0
	}
	curBar := bars[l].Volume
	prevBar := bars[l - 1].Volume

	jump := curBar / prevBar
	if math.IsNaN(jump) {
		return -1
	}
	if math.IsInf(jump,1) {
		return -1
	}
	return jump
}
/*
	Average Close
 */
func LastClose(stocks []v.Stock, timeFrame string) []v.Stock {
	log.Println("Running LastClose")
	for i, stock := range stocks {
		switch timeFrame {
		case v.OneWeek:
			if len(stock.WeekBars) > 0 {
				lastClose := stock.WeekBars[len(stock.WeekBars)-1].Close
				stocks[i].OneWeek.Close.LastClose = lastClose
			}
		case v.OneDay:
			if len(stock.DayBars) > 0 {
				lastClose := stock.DayBars[len(stock.DayBars)-1].Close
				stocks[i].OneDay.Close.LastClose = lastClose
			}
		case v.OneHour:
			if len(stock.HourBars) > 0 {
				lastClose := stock.HourBars[len(stock.HourBars)-1].Close
				stocks[i].OneHour.Close.LastClose = lastClose
			}
		case v.FourHour:
			if len(stock.FourHourBars) > 0 {
				lastClose := stock.FourHourBars[len(stock.FourHourBars)-1].Close
				stocks[i].FourHour.Close.LastClose = lastClose
			}
		}
	}
	log.Println("Completed LastClose")
	return stocks
}
/*
	Close Greater Than Previous
 */
func CloseGreaterThanPrevious(stocks []v.Stock, timeFrame string) []v.Stock {
	log.Println("Running CloseGreaterThanPrevious")
	for i, stock := range stocks {
		switch timeFrame {
		case v.OneWeek:
			l := len(stock.WeekBars) - 1
			if l < 1 {
				stocks[i].OneWeek.Close.GreaterThanLast = 2
				continue
			}
			curBar := stock.WeekBars[l].Close
			prevBar := stock.WeekBars[l - 1].Close
			if prevBar < curBar {
				stocks[i].OneWeek.Close.GreaterThanLast = 1
			} else {
				stocks[i].OneWeek.Close.GreaterThanLast = 2
			}
		case v.OneDay:
			l := len(stock.DayBars) - 1
			if l < 1 {
				stocks[i].OneDay.Close.GreaterThanLast = 2
				continue
			}
			curBar := stock.DayBars[l].Close
			prevBar := stock.DayBars[l - 1].Close
			if prevBar < curBar {
				stocks[i].OneDay.Close.GreaterThanLast = 1
			} else {
				stocks[i].OneDay.Close.GreaterThanLast = 2
			}
		case v.OneHour:
			l := len(stock.HourBars) - 1
			if l < 1 {
				stocks[i].OneHour.Close.GreaterThanLast = 2
				continue
			}
			curBar := stock.HourBars[l].Close
			prevBar := stock.HourBars[l - 1].Close
			if prevBar < curBar {
				stocks[i].OneHour.Close.GreaterThanLast = 1
			} else {
				stocks[i].OneHour.Close.GreaterThanLast = 2
			}
		case v.FourHour:
			l := len(stock.FourHourBars) - 1
			if l < 1 {
				stocks[i].FourHour.Close.GreaterThanLast = 2
				continue
			}
			curBar := stock.FourHourBars[l].Close
			prevBar := stock.FourHourBars[l - 1].Close
			if prevBar < curBar {
				stocks[i].FourHour.Close.GreaterThanLast = 1
			} else {
				stocks[i].FourHour.Close.GreaterThanLast = 2
			}
		}
	}
	log.Println("Completed CloseGreaterThanPrevious")
	return stocks
}
/*
	Bollinger bands
 */
func BB(stocks []v.Stock, timeFrame string, period int) []v.Stock {
	log.Println("Running BB")
	for i, stock := range stocks {
		switch timeFrame {
		case v.OneWeek:
			stocks[i].OneWeek.BB = CalcBB(stock.WeekBars, period)
		case v.OneDay:
			stocks[i].OneDay.BB = CalcBB(stock.DayBars, period)
		case v.OneHour:
			stocks[i].OneHour.BB = CalcBB(stock.HourBars, period)
		case v.FourHour:
			stocks[i].FourHour.BB = CalcBB(stock.FourHourBars, period)
		}
	}
	log.Println("Completed BB")
	return stocks
}
func CalcBB(bars []v.Bar, period int) v.BBResult {
	var bbResult v.BBResult
	var smaStart = len(bars) - period
	var sma float64
	var result float64
	if smaStart < 0 {
		return v.BBResult{}
	}

	for k := smaStart; k < len(bars); k++ {
		sma += bars[k].Close
	}

	sma /= float64(period)

	result = math.Pow(bars[len(bars)-1].Close-sma, 2)
	result = result / 2

	std := math.Sqrt(result)
	upper := sma + (std * 2)
	lower := sma - (std * 2)

	bbResult.LowerBand = lower
	bbResult.UpperBand = upper

	return bbResult
}
