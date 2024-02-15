package db

import (
	v "../variables"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

func CreateUser(username string, password string) bool {
	if checkForUserName(username) == true {
		log.Println("User ", username, " already exists")
		return false
	}
	hashedPassword, err := HashPassword(password)
	if err != nil {
		log.Fatal("Error hashing password in createUser:", err.Error())
	}

	db, err := sql.Open("postgres", v.ConnectionString)
	if err != nil {
		log.Fatal("Error connecting to postgres in createUser:", err.Error())
	}
	defer db.Close()

	sql := "INSERT INTO stock_server.users " +
		"(username, password) " +
		"VALUES ($1, $2)"

	_, err = db.Exec(
		sql,
		username,
		hashedPassword)
	if err != nil {
		log.Fatal("Error inserting to postgres in createUser:", err.Error())
	}

	return true
}
func LoginUser(username string, password string) bool {
	var hashedPassword = getUserPassword(username)
	if checkForUserName(username) == false {
		log.Println("Login failed. User does not exist")
		return false
	}
	if CheckPasswordHash(password, hashedPassword) == false {
		log.Println("Login failed. Incorrect password")
		return false
	}
	return true
}
func GetUser (userName string) v.UserAccount {
	var userAccount v.UserAccount
	var firebaseId sql.NullString

	db, err := sql.Open(v.PGUser, v.ConnectionString)
	if err != nil {
		log.Fatal("Error connecting to postgres in GetUser:", err.Error())
	}
	defer db.Close()

	sql := "SELECT username, firebaseId " +
		"FROM stock_server.users " +
		"WHERE username = $1"

	qr, err := db.Query(
		sql,
		userName)
	if err != nil {
		log.Fatal("Error querying postgres in GetUser:", err.Error())
	}

	if qr.Next() == true {
		err = qr.Scan(&userAccount.UserName, &firebaseId)
		if err != nil {
			log.Fatal("Error scanning query row in GetUser:", err.Error())
		}
	} else {
		return v.UserAccount{}
	}

	userAccount.FirebaseId = firebaseId.String

	return userAccount
}
func GetUsers () []v.UserAccount {
	var userAccounts []v.UserAccount

	var userAccount v.UserAccount
	var firebaseId sql.NullString

	db, err := sql.Open(v.PGUser, v.ConnectionString)
	if err != nil {
		log.Fatal("Error connecting to postgres in GetUsers:", err.Error())
	}
	defer db.Close()

	sql := "SELECT username, firebaseId, userstocks,  algorithm" +
		"FROM stock_server.users u" +
		"WHERE username = $1"

	qr, err := db.Query(sql)
	if err != nil {
		log.Fatal("Error querying postgres in GetUsers:", err.Error())
	}

	for qr.Next() {
		err = qr.Scan(&userAccount.UserName, &firebaseId)
		if err != nil {
			log.Fatal("Error scanning query row in GetUsers:", err.Error())
		} else {
			userAccount.FirebaseId = firebaseId.String
			userAccounts = append(userAccounts, userAccount)
		}
	}
	return userAccounts
}
func UpdateFirebaseId(userAccount v.UserAccount) bool {
	if checkForUserName(userAccount.UserName) == false {
		return false
	}

	db, err := sql.Open(v.PGUser, v.ConnectionString)
	if err != nil {
		log.Fatal("Error connecting to postgres in UpdateFirebaseId:", err.Error())
	}
	defer db.Close()

	sql := "UPDATE stock_server.users " +
		"SET firebaseId = $1 " +
		"WHERE userName = $2"

	_, err = db.Exec(
		sql,
		userAccount.FirebaseId,
		userAccount.UserName)

	if err != nil {
		log.Fatal("Error querying postgres in UpdateFirebaseId:", err.Error())
		return false
	}
	return true

}

func GetUserStocks(username string, algorithm int) []v.Stock {
	var str sql.NullString
	var userStocks []v.Stock

	db, err := sql.Open(v.PGUser, v.ConnectionString)
	if err != nil {
		log.Fatal("Error connecting to postgres in GetUserStocks:", err.Error())
	}
	defer db.Close()

	sql := "SELECT userstocks " +
		"FROM stock_server.algorithms " +
		"WHERE username = $1 AND uniqueid = $2"

	qr, err := db.Query(
		sql,
		username,
		algorithm)
	if err != nil {
		log.Fatal("Error querying postgres in GetUserStocks:", err.Error())
	}

	for qr.Next() {
		err = qr.Scan(&str)
		if err != nil {
			log.Fatal("Error scanning query row in GetUserStocks:", err.Error())
		}
	}
	if str.String == ""{
		return nil
	}
	userStocks = JsonToStock(str.String)

	return userStocks
}
func UpdateUserStocks (username string, stocks []v.Stock, algorithm int) {
	stocks = trimBars(stocks)
	var stringStocks = StocksToJson(stocks)

	db, err := sql.Open(v.PGUser, v.ConnectionString)
	if err != nil {
		log.Fatal("Error connecting to postgres in UpdateUserStocks:", err.Error())
	}
	defer db.Close()

	sql := "UPDATE stock_server.algorithms " +
		"SET userstocks = $1 " +
		"WHERE username = $2 AND uniqueid = $3"

	_, err = db.Exec(
		sql,
		stringStocks,
		username,
		algorithm)
	if err != nil {
		log.Fatal("Error querying postgres in UpdateUserStocks:", err.Error())
	}
}

func GetViableStocks (username string, algorithm int) []v.Stock {
	var str sql.NullString
	var stocksToMonitor []v.Stock

	db, err := sql.Open(v.PGUser, v.ConnectionString)
	if err != nil {
		log.Fatal("Error connecting to postgres in GetViableStocks:", err.Error())
	}
	defer db.Close()

	sql := "SELECT viable " +
		"FROM stock_server.algorithms " +
		"WHERE username = $1 AND uniqueid = $2"

	qr, err := db.Query(
		sql,
		username,
		algorithm)
	if err != nil {
		log.Fatal("Error querying postgres in GetViableStocks:", err.Error())
	}

	for qr.Next() {
		err = qr.Scan(&str)
		if err != nil {
			log.Fatal("Error scanning query row in GetViableStocks:", err.Error())
		}
	}
	if str.String == ""{
		return nil
	}
	stocksToMonitor = JsonToStock(str.String)

	return stocksToMonitor

}
func DeleteViableStocks (username string) {
	db, err := sql.Open(v.PGUser, v.ConnectionString)
	if err != nil {
		log.Fatal("Error connecting to postgres in DeleteViableStocks:", err.Error())
	}
	defer db.Close()

	sql := "DELETE viable " +
		"FROM stock_server.algorithms " +
		"WHERE username = $1"

	_, err = db.Exec(
		sql,
		username)
	if err != nil {
		log.Fatal("Error querying postgres in DeleteViableStocks:", err.Error())
	}
}
func UpdateViableStocks (username string, stocks []v.Stock, algorithm int) {
	stocks = trimBars(stocks)
	var stringStocks = StocksToJson(stocks)

	db, err := sql.Open(v.PGUser, v.ConnectionString)
	if err != nil {
		log.Fatal("Error connecting to postgres in UpdateViableStocks:", err.Error())
	}
	defer db.Close()

	sql := "UPDATE stock_server.algorithms " +
		"SET viable = $1 " +
		"WHERE username = $2 AND uniqueid = $3"

	_, err = db.Exec(
		sql,
		stringStocks,
		username,
		algorithm)
	if err != nil {
		log.Fatal("Error querying postgres in UpdateViableStocks:", err.Error())
	}

}

func GetPurchasableStocks (username string, algorithm int) []v.Stock {
	var str sql.NullString
	var purchasable []v.Stock

	db, err := sql.Open(v.PGUser, v.ConnectionString)
	if err != nil {
		log.Fatal("Error connecting to postgres in GetPurchasableStocks:", err.Error())
	}
	defer db.Close()

	sql := "SELECT purchasable " +
		"FROM stock_server.algorithms " +
		"WHERE username = $1 AND uniqueid = $2"

	qr, err := db.Query(
		sql,
		username,
		algorithm)
	if err != nil {
		log.Fatal("Error querying postgres in GetPurchasableStocks:", err.Error())
	}

	for qr.Next() {
		err = qr.Scan(&str)
		if err != nil {
			log.Fatal("Error scanning query row in GetPurchasableStocks:", err.Error())
		}
	}
	if str.String == ""{
		return nil
	}
	purchasable = JsonToStock(str.String)

	return purchasable
}
func UpdatePurchasableStocks (username string, stocks []v.Stock, algorithm int) {
	stocks = trimBars(stocks)
	var stringStocks = StocksToJson(stocks)

	db, err := sql.Open(v.PGUser, v.ConnectionString)
	if err != nil {
		log.Fatal("Error connecting to postgres in UpdatePurchasableStocks:", err.Error())
	}
	defer db.Close()

	sql := "UPDATE stock_server.algorithms " +
		"SET purchasable = $1 " +
		"WHERE username = $2 AND uniqueid = $3"

	_, err = db.Exec(
		sql,
		stringStocks,
		username,
		algorithm)

	if err != nil {
		log.Fatal("Error querying postgres in UpdatePurchasableStocks:", err.Error())
	}
}

func GetNotifications (username string, algorithm int) []v.Notification {
	var usernamestr sql.NullString
	var datestr sql.NullString
	var titlestr sql.NullString
	var messagestr sql.NullString
	var notifications []v.Notification

	db, err := sql.Open(v.PGUser, v.ConnectionString)
	if err != nil {
		log.Fatal("Error connecting to postgres in GetNotifications:", err.Error())
	}
	defer db.Close()

	sql := "SELECT username, date, title, message " +
		"FROM stock_server.notifications " +
		"WHERE username = $1 AND algorithm = $2 " +
		"ORDER BY date DESC"

	qr, err := db.Query(
		sql,
		username,
		algorithm)
	if err != nil {
		log.Fatal("Error querying postgres in GetNotifications:", err.Error())
	}

	for qr.Next() {
		err = qr.Scan(&usernamestr, &datestr, &titlestr, &messagestr)
		if err != nil {
			log.Fatal("Error scanning query row in GetNotifications:", err.Error())
		}
		notifications = append(notifications, v.Notification{
			UserName: usernamestr.String,
			Date:     datestr.String,
			Title:    titlestr.String,
			Message:  messagestr.String,
		})
	}

	return notifications
}
func DeleteNotifications (notifications []v.Notification, algorithm int) {
	db, err := sql.Open(v.PGUser, v.ConnectionString)
	if err != nil {
		log.Fatal("Error connecting to postgres in DeleteNotifications:", err.Error())
	}
	defer db.Close()

	sql := "DELETE " +
		"FROM stock_server.notifications " +
		"WHERE username = $1 and " +
		"date = $2 and " +
		"title = $3 and " +
		"message = $4 and " +
		"algorithm = $5"

	for _, notification := range notifications {
		_, err = db.Exec(
			sql,
			notification.UserName,
			notification.Date,
			notification.Title,
			notification.Message,
			algorithm)
		if err != nil {
			log.Fatal("Error querying postgres in DeleteNotifications:", err.Error())
		}
	}
}
func StoreNotification (notification v.Notification, algorithm int){
	db, err := sql.Open(v.PGUser, v.ConnectionString)
	if err != nil {
		log.Fatal("Error connecting to postgres in StoreNotification:", err.Error())
	}
	defer db.Close()

	sql := "INSERT INTO stock_server.notifications " +
		"(username, date, title, message, algorithm) " +
		"VALUES($1, $2, $3, $4, $5) "

	_, err = db.Exec(
		sql,
		notification.UserName,
		notification.Date,
		notification.Title,
		notification.Message,
		algorithm)
	if err != nil {
		log.Fatal("Error querying postgres in StoreNotification:", err.Error())
	}
}





