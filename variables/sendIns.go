package variables

type AddSingleStockSendIn struct {
	Algorithm int
	UserName string
	SingleStock string
}
type MultipleStocksSendIn struct {
	Algorithm int
	UserName string
	UserStocks []Stock
}
type BasicUserSendIn struct {
	UserName string
	Algorithm int
}
type DeleteNotificationsSendIn struct {
	UserName      string
	Algorithm     int
	Notifications []Notification
}

