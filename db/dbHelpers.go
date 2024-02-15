package db

import (
	v "../variables"
	"database/sql"
	"encoding/json"
	"log"
)

const (
	DateFormat  = "01/02/2006"
)

func StocksToJson (stocks []v.Stock) string {
	/*for _, stock := range stocks {
		_, err := json.Marshal(stock)
		if err != nil {
			log.Println(stock)
			log.Fatal("Error turning Stocks into json in StocksToJson: ", err.Error())
		}
	}

	 */
	str, err := json.Marshal(stocks)
	if err != nil {
		log.Fatal("Error turning Stocks into json in StocksToJson: ", err.Error())
	}
	return string(str)
}
func JsonToStock (str string) []v.Stock{
	var filterStocks []v.Stock
	err := json.Unmarshal([]byte(str), &filterStocks)
	if err != nil {
		log.Fatal("Error turning json into Stocks in JsonToStock: ", err.Error())
	}
	return filterStocks
}

func getUserPassword (userName string) string {
	var hashedPassword string

	db, err := sql.Open(v.PGUser, v.ConnectionString)
	if err != nil {
		log.Fatal("Error connecting to postgres in getUserPassword:", err.Error())
	}
	defer db.Close()

	sql := "SELECT password " +
		"FROM stock_server.users " +
		"WHERE userName = $1"

	qr, err := db.Query(
		sql,
		userName)
	if err != nil {
		log.Fatal("Error querying postgres in getUserPassword:", err.Error())
	}

	if qr.Next() == true {
		err = qr.Scan(&hashedPassword)
		if err != nil {
			log.Fatal("Error scanning query row in getUserPassword:", err.Error())
		}
	} else {
		return ""
	}

	return hashedPassword
}
func checkForUserName(userName string) bool {
	db, err := sql.Open(v.PGUser, v.ConnectionString)
	if err != nil {
		log.Fatal("Error connecting to postgres in checkForUserName:", err.Error())
	}
	defer db.Close()

	sql := "SELECT * FROM stock_server.users " +
		"WHERE userName = $1"

	qr, err := db.Query(
		sql,
		userName)
	if err != nil {
		log.Fatal("Error querying postgres in checkForUserName:", err.Error())
		return false
	}

	if qr.Next() == false {
		return false
	}
	return true
}

func trimBars (stocks []v.Stock) []v.Stock {
	for i, _ := range stocks {
		stocks[i].WeekBars = nil
		stocks[i].DayBars = nil
		stocks[i].FourHourBars = nil
		stocks[i].HourBars = nil
		if len(stocks[i].OneWeek.Adx) > 0 {
			stocks[i].OneWeek.Adx = []v.AdxResult{stocks[i].OneWeek.Adx[len(stocks[i].OneWeek.Adx)-1]}
		}
		if len(stocks[i].OneDay.Adx) > 0 {
			stocks[i].OneDay.Adx = []v.AdxResult{stocks[i].OneDay.Adx[len(stocks[i].OneDay.Adx)-1]}
		}
		if len(stocks[i].FourHour.Adx) > 0 {
			stocks[i].FourHour.Adx = []v.AdxResult{stocks[i].FourHour.Adx[len(stocks[i].FourHour.Adx)-1]}
		}
		if len(stocks[i].OneHour.Adx) > 0 {
			stocks[i].OneHour.Adx = []v.AdxResult{stocks[i].OneHour.Adx[len(stocks[i].OneHour.Adx)-1]}
		}
	}
	return stocks
}

