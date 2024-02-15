package algorithm

import (
	f "../filtering"
	"log"
	"time"
)

//Key is the users name
var AlgorithmMap = map[string]map[int]*Algorithm{}

type AlgorithmSendIn struct {
	UserName string
	Algo Algorithm
}
type LoginSendBack struct {
	Successful bool
	Algorithms []UserAlgorithm
}

type UserAlgorithm struct {
	Name string
	UniqueId int32
}

type Algorithm struct {
	UserName 				string
	Name 					string
	UniqueID				int
	IsRunning				bool
	RunOnStart 				bool
	shutDown                chan bool
	StraightToMonitoring    bool
	DailyAnalysis           bool
	ViableFilter			f.Viable
	PurchasableFilter		f.Purchasable
	SellFilter			    f.Sell
	CheckViableHrStart      int
	CheckViableMinStart     int
	CheckViableHrIncrement  int
	CheckViableMinIncrement int
	CheckSellHrStart      int
	CheckSellMinStart     int
	CheckSellHrIncrement  int
	CheckSellMinIncrement int
}
func (a *Algorithm) Run () {
	var boughtCheckHr int = a.CheckSellHrStart
	var boughtCheckMin int = a.CheckSellMinStart
	var viableCheckHr int = a.CheckViableHrStart
	var viableCheckMin int = a.CheckViableMinStart
	var t time.Time
	var day time.Weekday

	log.Println("Algorithm: " + a.Name + " started")
	for {
		select {
		case <-a.shutDown:
			a.ShutDownProcedure()
			log.Println("Algorithm: " + a.Name + " stopped")
			return
		default:
			day = time.Now().Weekday()
			//If market not open and viable stocks not found
			if (time.Now().Hour() == 8 &&
				a.DailyAnalysis == false &&
				day != time.Saturday &&
				day != time.Sunday) &&
				a.StraightToMonitoring == false {

				go a.GetViableStocks()
				a.DailyAnalysis = true
			}

			t = time.Now()
			//If market open
			if  day != time.Saturday &&
				day != time.Sunday &&
				t.Hour() >= 9 && t.Minute() >= 30 &&
				t.Hour() >= boughtCheckHr &&
				t.Minute() >= boughtCheckMin &&
				t.Hour() < 16 {

				go a.MonitorUserStocks()
				boughtCheckHr += a.CheckSellHrIncrement
				boughtCheckMin += a.CheckSellMinIncrement

				if boughtCheckMin == 60 {
					boughtCheckMin = 0
					if a.CheckSellHrIncrement == 0 {
						boughtCheckHr++
					}
				}
			}

			t = time.Now()
			if  day != time.Saturday &&
				day != time.Sunday &&
				t.Hour() >= 9 && t.Minute() >= 30 &&
				t.Hour() >= viableCheckHr &&
				t.Minute() >= viableCheckMin &&
				t.Hour() < 16 {

				go a.MonitorViableStocks()
				viableCheckHr += a.CheckViableHrIncrement
				viableCheckMin += a.CheckViableMinIncrement

				if viableCheckMin == 60 {
					viableCheckMin = 0
					if  a.CheckViableHrIncrement == 0 {
						viableCheckHr++
					}
				}
			}

			//if market close
			t = time.Now()
			if t.Hour() == 16  &&
				day != time.Saturday &&
				day != time.Sunday &&
				a.DailyAnalysis == true {

				a.DailyAnalysis = false
				boughtCheckHr = a.CheckSellHrStart
				boughtCheckMin = a.CheckSellMinStart
				viableCheckHr = a.CheckViableHrStart
				viableCheckMin = a.CheckViableMinStart
			}

			time.Sleep(time.Minute)
		}
	}

}
func (a *Algorithm) CheckAlgorithm() bool {
	return true
}
func (a *Algorithm) ShutDownProcedure () {
	a.IsRunning = false
	UpdateAlgorithm(*a)
}
func AddToAlgorithmMap (a Algorithm) {
	if AlgorithmMap[a.UserName] == nil {
		AlgorithmMap[a.UserName] = map[int]*Algorithm{}
	}
	AlgorithmMap[a.UserName][a.UniqueID] = &a
}
func UpdateAlgorithmMap (a Algorithm) {
	AlgorithmMap[a.UserName][a.UniqueID] = &a
}
func StartAlgorithm (userName string, algorithm int) bool {
	if AlgorithmMap[userName][algorithm].IsRunning == false || AlgorithmMap[userName][algorithm].RunOnStart == true {
		go AlgorithmMap[userName][algorithm].Run()
		AlgorithmMap[userName][algorithm].IsRunning = true
		UpdateAlgorithm(*AlgorithmMap[userName][algorithm])
		return true
	} else {
		return false
	}
}
func ShutDownAlgorithm (userName string, algorithm int) bool {
	if AlgorithmMap[userName][algorithm].IsRunning == true {
		AlgorithmMap[userName][algorithm].shutDown <- true
		return true
	} else {
		return false
	}
}
func RunViable (userName string, algorithm int) {
	go AlgorithmMap[userName][algorithm].GetViableStocks()
}
func RunPurchasable (userName string, algorithm int) {
	go AlgorithmMap[userName][algorithm].MonitorViableStocks()
}
func RunSell (userName string, algorithm int) {
	go AlgorithmMap[userName][algorithm].MonitorUserStocks()
}
func RemovedFromAlgorithmMap (uniqueId int, userName string) {
	delete(AlgorithmMap[userName], uniqueId)
}
func DeleteAndStopAlgorithm (userName string, uniqueId int){
	ShutDownAlgorithm(userName, uniqueId)
	RemovedFromAlgorithmMap(uniqueId, userName)
}