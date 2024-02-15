package filtering

type Viable struct {
	Checks []FilterCheck
}
type Purchasable struct {
	Checks []FilterCheck
}
type Sell struct {
	Checks []FilterCheck
}
type FilterCheck struct {
	Period 					 int
	TimePeriod               string
	ADX                      ADXCheck
	AccumulationDistribution AccumulationDistributionCheck
	EMA                      EMACheck
	SMA                      SMACheck
	RSI                      RSICheck
	AverageVolume            AverageVolumeCheck
	VolumeJump 				 VolumeJumpCheck
	LastClose                LastCloseCheck
	CloseGreaterThanPrevious CloseGreaterThanPreviousCheck
	BB                       BBCheck
}
type VolumeJumpCheck struct {
	Set        bool
	Multiplier float64
}
type CloseGreaterThanPreviousCheck struct {
	Set bool
	Greater int
}
type ADXCheck struct {
	Set bool
	ADXIncreasing int
	PDIIncreasing int
	MDIIncreasing int
}
type AccumulationDistributionCheck struct {
	Set bool
	Increasing int
}
type EMACheck struct {
	Set bool
	Increasing     int
	GreaterThan0   int
	GreaterThanSMA int
}
type SMACheck struct {
	Set bool
	Increasing int
	GreaterThan0 int
}
type RSICheck struct {
	Set bool
	Increasing  int
	GreaterThan float64
	LessThan    float64
}
type AverageVolumeCheck struct {
	Set bool
	GreaterThan float64
}
type LastCloseCheck struct {
	Set bool
	GreaterThan float64
	LessThan    float64
}
type BBCheck struct {
	Set bool
	GreaterThan float64
	LessThan    float64
}