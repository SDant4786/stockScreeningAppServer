package algorithm

import (
	f "../filtering"
	v "../variables"
)
var BlankAlgorithm = Algorithm{
	UserName:                "",
	Name:                    "blank",
	UniqueID:                0,
	IsRunning:               false,
	RunOnStart:              false,
	shutDown:                make(chan bool),
	StraightToMonitoring:    false,
	DailyAnalysis:           false,
	ViableFilter:            f.Viable{},
	PurchasableFilter:       f.Purchasable{},
	SellFilter:              f.Sell{},
	CheckViableHrStart:      0,
	CheckViableMinStart:     0,
	CheckViableHrIncrement:  0,
	CheckViableMinIncrement: 0,
	CheckSellHrStart:        0,
	CheckSellMinStart:       0,
	CheckSellHrIncrement:    0,
	CheckSellMinIncrement:   0,
}
var DefaultAlgorithm = Algorithm{
	UserName:             "",
	Name:                 "default",
	UniqueID:             1,
	IsRunning:            false,
	RunOnStart:           true,
	shutDown:             make(chan bool),
	StraightToMonitoring: false,
	DailyAnalysis:        false,
	ViableFilter: f.Viable{
		Checks: []f.FilterCheck{
			{
				Period:     9,
				TimePeriod: v.OneWeek,
				ADX: f.ADXCheck{
					Set: true,
					ADXIncreasing: 0,
					PDIIncreasing: 0,
					MDIIncreasing: 0,
				},
				AccumulationDistribution: f.AccumulationDistributionCheck{
					Set: true,
					Increasing: 1,
				},
				EMA: f.EMACheck{
					Set: true,
					Increasing:     1,
					GreaterThan0:   0,
					GreaterThanSMA: 0,
				},
				SMA: f.SMACheck{
					Set: true,
					Increasing:   1,
					GreaterThan0: 0,
				},
				RSI: f.RSICheck{
					Set: true,
					Increasing:  0,
					GreaterThan: 20,
					LessThan:    80,
				},
				AverageVolume: f.AverageVolumeCheck{
					Set: true,
					GreaterThan: 500000,
				},
				LastClose: f.LastCloseCheck{
					Set: true,
					GreaterThan: 1,
					LessThan:    50,
				},
				BB: f.BBCheck{
					Set: true,
					GreaterThan: 1,
					LessThan:    1,
				},
			},
			{
				Period:     9,
				TimePeriod: v.OneDay,
				ADX: f.ADXCheck{
					Set: true,
					ADXIncreasing: 0,
					PDIIncreasing: 1,
					MDIIncreasing: 2,
				},
				AccumulationDistribution: f.AccumulationDistributionCheck{
					Set: true,
					Increasing: 1,
				},
				EMA: f.EMACheck{
					Set: true,
					Increasing:     1,
					GreaterThan0:   0,
					GreaterThanSMA: 0,
				},
				SMA: f.SMACheck{
					Set: true,
					Increasing:   1,
					GreaterThan0: 0,
				},
				RSI: f.RSICheck{
					Set: true,
					Increasing:  0,
					GreaterThan: 20,
					LessThan:    80,
				},
				AverageVolume: f.AverageVolumeCheck{
					Set: true,
					GreaterThan: 500000,
				},
				LastClose: f.LastCloseCheck{
					Set: true,
					GreaterThan: 1,
					LessThan:    50,
				},
				BB: f.BBCheck{
					Set: true,
					GreaterThan: 1,
					LessThan:    1,
				},
			},
		},
	},
	PurchasableFilter: f.Purchasable{
		Checks: []f.FilterCheck{
			{
				Period:     9,
				TimePeriod: v.FourHour,
				ADX: f.ADXCheck{
					Set: true,
					ADXIncreasing: 1,
					PDIIncreasing: 1,
					MDIIncreasing: 2,
				},
				AccumulationDistribution: f.AccumulationDistributionCheck{
					Set: true,
					Increasing: 1,
				},
				EMA: f.EMACheck{
					Set: true,
					Increasing:     1,
					GreaterThan0:   1,
					GreaterThanSMA: 2,
				},
				SMA: f.SMACheck{
					Set: true,
					Increasing:   0,
					GreaterThan0: 0,
				},
				RSI: f.RSICheck{
					Set: true,
					Increasing:  1,
					GreaterThan: 20,
					LessThan:    80,
				},
				AverageVolume: f.AverageVolumeCheck{
					Set: true,
					GreaterThan: 0,
				},
				LastClose: f.LastCloseCheck{
					Set: true,
					GreaterThan: 0,
					LessThan:    0,
				},
				BB: f.BBCheck{
					Set: true,
					GreaterThan: 1,
					LessThan:    1,
				},
			},
			{
				Period:     9,
				TimePeriod: v.OneHour,
				ADX: f.ADXCheck{
					Set: true,
					ADXIncreasing: 0,
					PDIIncreasing: 1,
					MDIIncreasing: 1,
				},
				AccumulationDistribution: f.AccumulationDistributionCheck{
					Set: true,
					Increasing: 0,
				},
				EMA: f.EMACheck{
					Set: true,
					Increasing:     1,
					GreaterThan0:   1,
					GreaterThanSMA: 1,
				},
				SMA: f.SMACheck{
					Set: true,
					Increasing:   0,
					GreaterThan0: 0,
				},
				RSI: f.RSICheck{
					Set: true,
					Increasing:  1,
					GreaterThan: 0,
					LessThan:    0,
				},
				AverageVolume: f.AverageVolumeCheck{
					Set: true,
					GreaterThan: 0,
				},
				LastClose: f.LastCloseCheck{
					Set: true,
					GreaterThan: 0,
					LessThan:    0,
				},
				BB: f.BBCheck{
					Set: true,
					GreaterThan: 1,
					LessThan:    1,
				},
			},
		},
	},
	SellFilter: f.Sell{
		Checks: []f.FilterCheck{
			{
				Period:     9,
				TimePeriod: v.FourHour,
				ADX: f.ADXCheck{
					Set: true,
					ADXIncreasing: 0,
					PDIIncreasing: 2,
					MDIIncreasing: 1,
				},
				AccumulationDistribution: f.AccumulationDistributionCheck{
					Set: true,
					Increasing: 2,
				},
				EMA: f.EMACheck{
					Set: true,
					Increasing:     2,
					GreaterThan0:   0,
					GreaterThanSMA: 2,
				},
				SMA: f.SMACheck{
					Set: true,
					Increasing:   2,
					GreaterThan0: 0,
				},
				RSI: f.RSICheck{
					Set: true,
					Increasing:  2,
					GreaterThan: 0,
					LessThan:    0,
				},
				AverageVolume: f.AverageVolumeCheck{
					Set: true,
					GreaterThan: 0,
				},
				LastClose: f.LastCloseCheck{
					Set: true,
					GreaterThan: 0,
					LessThan:    0,
				},
				BB: f.BBCheck{
					Set: true,
					GreaterThan: 0,
					LessThan:    0,
				},
			},
			{
				Period:     9,
				TimePeriod: v.OneHour,
				ADX: f.ADXCheck{
					Set: true,
					ADXIncreasing: 0,
					PDIIncreasing: 2,
					MDIIncreasing: 1,
				},
				AccumulationDistribution: f.AccumulationDistributionCheck{
					Set: true,
					Increasing: 0,
				},
				EMA: f.EMACheck{
					Set: true,
					Increasing:     2,
					GreaterThan0:   2,
					GreaterThanSMA: 2,
				},
				SMA: f.SMACheck{
					Set: true,
					Increasing:   2,
					GreaterThan0: 0,
				},
				RSI: f.RSICheck{
					Set: true,
					Increasing:  2,
					GreaterThan: 0,
					LessThan:    0,
				},
				AverageVolume: f.AverageVolumeCheck{
					Set: true,
					GreaterThan: 0,
				},
				LastClose: f.LastCloseCheck{
					Set: true,
					GreaterThan: 0,
					LessThan:    0,
				},
				BB: f.BBCheck{
					Set: true,
					GreaterThan: 0,
					LessThan:    0,
				},
			},
		},
	},
	CheckViableHrStart:      9,
	CheckViableMinStart:     30,
	CheckViableHrIncrement:  0,
	CheckViableMinIncrement: 30,
	CheckSellHrStart:        9,
	CheckSellMinStart:       30,
	CheckSellHrIncrement:    0,
	CheckSellMinIncrement:   30,
}