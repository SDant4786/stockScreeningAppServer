package variables

import (
	"fmt"
	"sync"
	"time"

	firebase "firebase.google.com/go"
)

var FireBase *firebase.App
var ConnectionString = fmt.Sprintf("host=%s port=%d user=%s "+
	"password=%s dbname=%s sslmode=disable",
	host, port, PGUser, password, dbname)

type UserAccount struct {
	UserName   string
	Password   string
	FirebaseId string
}
type ThreadSafeSlice struct {
	Mutex *sync.Mutex
	Slice []Stock
}
type Stock struct {
	Symbol       string
	WeekBars     []Bar
	DayBars      []Bar
	FourHourBars []Bar
	HourBars     []Bar
	OneWeek      FilterResults
	OneDay       FilterResults
	OneHour      FilterResults
	FourHour     FilterResults
	PossibleTop  bool
	Sell         bool
}
type Bar struct {
	Time   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}
type FilterResults struct {
	EMA    EMAResult
	SMA    SMAResult
	RSI    RSIResult
	Adx    []AdxResult
	Volume VolumeResult
	Close  CloseResult
	BB     BBResult
}
type VolumeResult struct {
	AverageVolume    float64
	IncreasingVolume int
	VolumeJump       float64
}
type CloseResult struct {
	LastClose       float64
	GreaterThanLast int
}
type EMAResult struct {
	EMA1                    float64
	EMA2                    float64
	EMA3                    float64
	EMASlopeCurrent         float64
	EMASlopePrevious        float64
	EMASlopeIncreasePercent float64
}
type SMAResult struct {
	SMA1                    float64
	SMA2                    float64
	SMA3                    float64
	SMASlopeCurrent         float64
	SMASlopePrevious        float64
	SMASlopeIncreasePercent float64
}
type RSIResult struct {
	RSICurrent         float64
	RSIPrevious        float64
	RSIPercentIncrease float64
}
type BBResult struct {
	UpperBand float64
	LowerBand float64
}
type AdxResult struct {
	Pdi float64
	Mdi float64
	Adx float64
}
type AtrResult struct {
	Tr   float64
	Atr  float64
	Atrp float64
}
type Notification struct {
	UserName string
	Date     string
	Title    string
	Message  string
}

const (
	host     = "localhost"
	port     = 5432
	PGUser   = "postgres"
	password = "password"
	dbname   = "stock_server"
	OneWeek  = "oneWeek"
	OneDay   = "oneDay"
	FourHour = "fourHour"
	OneHour  = "oneHour"
)
