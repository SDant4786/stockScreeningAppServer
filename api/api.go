package api

import (
	a "../algorithm"
	"../db"
	"../firebase"
	v "../variables"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)
var ProgramShutDown = make(chan bool)
//TODO add a verifing user function
func StartServer () {
	var err error

	firebase.InitFirebase()

	http.HandleFunc("/create_user", CreateUser)
	http.HandleFunc("/login", LoginUser)
	http.HandleFunc("/receive_firebase_id", ReceiveFirebaseID)

	http.HandleFunc("/send_purchasable_stocks", SendPurchasableStocks)
	http.HandleFunc("/add_single_stock", AddSingleStock)
	http.HandleFunc("/add_to_user_stocks", AddToUserStocks)
	http.HandleFunc("/send_user_stocks", SendUserStocks)
	http.HandleFunc("/remove_user_stocks", RemoveFromUserStocks)

	http.HandleFunc("/send_algorithms", SendAlgorithms)
	http.HandleFunc("/send_algorithm", SendAlgorithm)
	http.HandleFunc("/create_algorithm", CreateAlgorithm)
	http.HandleFunc("/update_algorithm", UpdateAlgorithm)
	http.HandleFunc("/stop_algorithm", StopAlgorithm)
	http.HandleFunc("/start_algorithm", StartAlgorithm)
	http.HandleFunc("/delete_algorithm", DeleteAlgorithm)
	http.HandleFunc("/run_viable", RunViable)
	http.HandleFunc("/run_purchasable", RunPurchasable)
	http.HandleFunc("/run_sell", RunSell)

	http.HandleFunc("/send_notifications", SendNotifications)
	http.HandleFunc("/delete_notifications", DeleteNotifications)

	err = http.ListenAndServeTLS(
		":443",
		"files/https-server.crt",
		"files/https-server.key",
		nil)
	if err != nil {panic(err)}
}

func CreateUser (w http.ResponseWriter, r * http.Request) {

	var loginSendBack a.LoginSendBack
	var userAccount v.UserAccount
	var decoder *json.Decoder

	log.Println("CreateUser attempted")

	decoder = json.NewDecoder(r.Body)
	err := decoder.Decode(&userAccount)

	if err != nil {
		log.Fatal("Error decoding JSON in CreateUser: ", err.Error())
	}

	if db.CreateUser(userAccount.UserName, userAccount.Password) == false {
		loginSendBack.Successful = false
		send, err := json.Marshal(loginSendBack)
		if err != nil {
			log.Println("Error sending json in CreateUser:", err.Error())
		}
		defer r.Body.Close()
		fmt.Fprint(w, string(send))
		log.Println("CreateUser failed: UserName already in database")
		return
	}
	defaultAlgorithm := a.DefaultAlgorithm
	defaultAlgorithm.UserName = userAccount.UserName

	a.InsertAlgorithm(userAccount.UserName, defaultAlgorithm)

	loginSendBack.Successful = true
	loginSendBack.Algorithms = a.GetUserAlgorithms(userAccount.UserName)

	send, err := json.Marshal(loginSendBack)
	if err != nil {
		log.Fatal("Error sending json in CreateUser:", err.Error())
	}
	defer r.Body.Close()
	fmt.Fprint(w, string(send))
	log.Println("CreateUser succeeded")
}
func LoginUser (w http.ResponseWriter, r * http.Request) {
	var userAccount v.UserAccount
	var loginSendBack a.LoginSendBack
	var decoder *json.Decoder

	log.Println("LoginUser attempted")

	decoder = json.NewDecoder(r.Body)
	err := decoder.Decode(&userAccount)

	if err != nil {
		log.Fatal("Error decoding JSON in LoginUser: ", err.Error())
	}

	if db.LoginUser(userAccount.UserName, userAccount.Password) == false {
		loginSendBack.Successful = false
		send, err := json.Marshal(loginSendBack)
		if err != nil {
			log.Println("Error sending json in LoginUser:", err.Error())
		}
		defer r.Body.Close()
		fmt.Fprint(w, string(send))
		log.Println("LoginUser failed")
		return
	}

	loginSendBack.Successful = true
	loginSendBack.Algorithms = a.GetUserAlgorithms(userAccount.UserName)

	send, err := json.Marshal(loginSendBack)
	if err != nil {
		log.Fatal("Error sending json in LoginUser:", err.Error())
	}
	defer r.Body.Close()
	fmt.Fprint(w, string(send))
	log.Println("LoginUser succeeded")
}
func ReceiveFirebaseID (w http.ResponseWriter, r *http.Request){

	var userAccount v.UserAccount
	var decoder *json.Decoder

	log.Println("Updating FirebaseId attempted")

	decoder = json.NewDecoder(r.Body)
	err := decoder.Decode(&userAccount)

	if err != nil {
		log.Fatal("Error decoding JSON in ReceiveFirebaseID: ", err.Error())
	}

	send, err := json.Marshal(true)
	if err != nil {
		log.Fatal("Error sending json in ReceiveFirebaseID:", err.Error())
	}
	defer r.Body.Close()
	fmt.Fprint(w, string(send))

	db.UpdateFirebaseId(userAccount)

	log.Println("FirebaseId update succeeded")

}

func SendPurchasableStocks (w http.ResponseWriter, r *http.Request) {

	var basicUserSendIn v.BasicUserSendIn
	var decoder *json.Decoder
	var purchasableStocks []v.Stock
	var userStocks []v.Stock
	var stocksToSend []v.Stock

	log.Println("SendPurchasableStocks attempted")

	decoder = json.NewDecoder(r.Body)
	err := decoder.Decode(&basicUserSendIn)

	if err != nil {
		log.Fatal("Error decoding JSON in SendPurchasableStocks: ", err.Error())
	}
	//purchasableStocks = db.GetViableStocks(basicUserSendIn.UserName, basicUserSendIn.Algorithm)
	purchasableStocks = db.GetPurchasableStocks(basicUserSendIn.UserName, basicUserSendIn.Algorithm)
	userStocks = db.GetUserStocks(basicUserSendIn.UserName, basicUserSendIn.Algorithm)
	for _, purchasableStock := range purchasableStocks {
		alreadyBought := false
		for _, userStock := range userStocks {
			if purchasableStock.Symbol == userStock.Symbol {
				alreadyBought = true
				break
			}
		}
		if alreadyBought == false {
			stocksToSend = append(stocksToSend, purchasableStock)
		}
	}

	stocksToSend = trimBars(stocksToSend)

	send, err := json.Marshal(stocksToSend)
	if err != nil {
		log.Fatal("Error sending json in SendPurchasableStocks:", err.Error())
	}
	defer r.Body.Close()
	fmt.Fprint(w, string(send))

	log.Println("SendPurchasableStocks succeeded")

}

func AddSingleStock (w http.ResponseWriter, r *http.Request) {
	var addSingleStockSendIn v.AddSingleStockSendIn
	var decoder *json.Decoder
	var viableStocks []v.Stock

	log.Println("AddSingleStock attempted")
	decoder = json.NewDecoder(r.Body)
	err := decoder.Decode(&addSingleStockSendIn)

	if err != nil {
		log.Fatal("Error decoding JSON in AddToUserStocks: ", err.Error())
	}
	addSingleStockSendIn.SingleStock = strings.ToUpper(addSingleStockSendIn.SingleStock)

	userStocks := db.GetUserStocks(addSingleStockSendIn.UserName, addSingleStockSendIn.Algorithm)
	viableStocks = db.GetViableStocks(addSingleStockSendIn.UserName, addSingleStockSendIn.Algorithm)

	var found = false
	for _, stock := range viableStocks {
		if stock.Symbol == addSingleStockSendIn.SingleStock {
			found = true
			userStocks = append(userStocks, stock)
		}
	}
	if found == false {
		userStocks = append(userStocks, v.Stock{
			Symbol: addSingleStockSendIn.SingleStock,
		})
	}

	db.UpdateUserStocks(addSingleStockSendIn.UserName, userStocks, addSingleStockSendIn.Algorithm)
	log.Println("AddSingleStock succeeded")
}
func AddToUserStocks (w http.ResponseWriter, r *http.Request) {
	var addMultipleStocksSendIn v.MultipleStocksSendIn
	var userStocks []v.Stock
	var decoder *json.Decoder

	log.Println("AddToUserStocks attempted")
	decoder = json.NewDecoder(r.Body)
	err := decoder.Decode(&addMultipleStocksSendIn)

	if err != nil {
		log.Fatal("Error decoding JSON in AddToUserStocks: ", err.Error())
	}

	userStocks = db.GetUserStocks(addMultipleStocksSendIn.UserName, addMultipleStocksSendIn.Algorithm)
	for _, stock := range addMultipleStocksSendIn.UserStocks {
		userStocks = append(userStocks, stock)
	}
	db.UpdateUserStocks(addMultipleStocksSendIn.UserName, userStocks, addMultipleStocksSendIn.Algorithm)

	send, err := json.Marshal(true)
	if err != nil {
		log.Fatal("Error sending json in AddToUserStocks:", err.Error())
	}
	defer r.Body.Close()
	fmt.Fprint(w, string(send))

	log.Println("AddToUserStocks succeeded")

}
func SendUserStocks (w http.ResponseWriter, r *http.Request) {
	var basicUserSendIn v.BasicUserSendIn
	var decoder *json.Decoder

	log.Println("SendUserStocks attempted")

	decoder = json.NewDecoder(r.Body)
	err := decoder.Decode(&basicUserSendIn)

	if err != nil {
		log.Fatal("Error decoding JSON in SendUserStocks: ", err.Error())
	}

	stocks := db.GetUserStocks(basicUserSendIn.UserName, basicUserSendIn.Algorithm)

	stocks = trimBars(stocks)
	send, err := json.Marshal(stocks)
	if err != nil {
		log.Fatal("Error sending json in SendUserStocks:", err.Error())
	}
	defer r.Body.Close()
	fmt.Fprint(w, string(send))

	log.Println("SendUserStocks succeeded")

}
func RemoveFromUserStocks (w http.ResponseWriter, r *http.Request) {
	var multipleStockSendIn v.MultipleStocksSendIn
	var newStockList []v.Stock
	var userStocks []v.Stock
	var decoder *json.Decoder

	log.Println("RemoveFromUserStocks attempted")

	decoder = json.NewDecoder(r.Body)
	err := decoder.Decode(&multipleStockSendIn)

	if err != nil {
		log.Fatal("Error decoding JSON in RemoveFromUserStocks: ", err.Error())
	}

	userStocks = db.GetUserStocks(multipleStockSendIn.UserName, multipleStockSendIn.Algorithm)

	for _, userStock := range userStocks {
		removing := false
		for _, stockToRemove := range multipleStockSendIn.UserStocks {
			if userStock.Symbol == stockToRemove.Symbol {
				removing = true
				break
			}
		}
		if removing == false {
			newStockList = append(newStockList, userStock)
		}
	}

	send, err := json.Marshal(true)
	if err != nil {
		log.Fatal("Error sending json in RemoveFromUserStocks:", err.Error())
	}
	defer r.Body.Close()
	fmt.Fprint(w, string(send))

	db.UpdateUserStocks(multipleStockSendIn.UserName, newStockList, multipleStockSendIn.Algorithm)

	log.Println("RemoveFromUserStocks succeeded")

}

func SendAlgorithms (w http.ResponseWriter, r *http.Request) {
	var userAccount v.UserAccount
	var algorithms []a.UserAlgorithm
	var decoder *json.Decoder

	log.Println("SendAlgorithms attempted")

	decoder = json.NewDecoder(r.Body)
	err := decoder.Decode(&userAccount)

	if err != nil {
		log.Fatal("Error decoding JSON in RemoveFromUserStocks: ", err.Error())
	}

	algorithms = a.GetUserAlgorithms(userAccount.UserName)
	send, err := json.Marshal(algorithms)
	if err != nil {
		log.Fatal("Error sending json in SendAlgorithms:", err.Error())
	}
	defer r.Body.Close()
	fmt.Fprint(w, string(send))
	log.Println("SendAlgorithms succeeded")
}
func SendAlgorithm (w http.ResponseWriter, r *http.Request) {
	var basicUserSendIn v.BasicUserSendIn
	var algorithm string
	var decoder *json.Decoder

	log.Println("SendAlgorithm attempted")

	decoder = json.NewDecoder(r.Body)
	err := decoder.Decode(&basicUserSendIn)

	if err != nil {
		log.Fatal("Error decoding JSON in SendAlgorithm: ", err.Error())
	}

	algorithm = a.GetUserAlgorithm(basicUserSendIn.UserName, basicUserSendIn.Algorithm)
	defer r.Body.Close()
	fmt.Fprint(w, algorithm)
	log.Println("SendAlgorithm succeeded")
}
func CreateAlgorithm (w http.ResponseWriter, r *http.Request) {

	var basicUserSendIn v.BasicUserSendIn
	var decoder *json.Decoder

	log.Println("CreateAlgorithm attempted")

	decoder = json.NewDecoder(r.Body)
	err := decoder.Decode(&basicUserSendIn)

	if err != nil {
		log.Fatal("Error decoding JSON in CreateAlgorithm: ", err.Error())
	}

	blankAlgorithm := a.BlankAlgorithm
	blankAlgorithm.UserName = basicUserSendIn.UserName

	a.InsertAlgorithm(basicUserSendIn.UserName, blankAlgorithm)

	UserAlgorithms := a.GetUserAlgorithms(basicUserSendIn.UserName)

	send, err := json.Marshal(UserAlgorithms)
	if err != nil {
		log.Fatal("Error sending json in CreateAlgorithm:", err.Error())
	}
	defer r.Body.Close()
	fmt.Fprint(w, string(send))
	log.Println("CreateAlgorithm succeeded")

}
func UpdateAlgorithm (w http.ResponseWriter, r *http.Request) {
	var algorithmSendIn a.AlgorithmSendIn
	var decoder *json.Decoder

	log.Println("UpdateAlgorithm attempted")

	decoder = json.NewDecoder(r.Body)
	err := decoder.Decode(&algorithmSendIn)

	if err != nil {
		log.Fatal("Error decoding JSON in UpdateAlgorithm: ", err.Error())
	}

	a.UpdateAlgorithm(algorithmSendIn.Algo)
	a.UpdateAlgorithmMap(algorithmSendIn.Algo)

	UserAlgorithms := a.GetUserAlgorithms(algorithmSendIn.UserName)

	send, err := json.Marshal(UserAlgorithms)
	if err != nil {
		log.Fatal("Error sending json in UpdateAlgorithm:", err.Error())
	}
	defer r.Body.Close()
	fmt.Fprint(w, string(send))
	log.Println("UpdateAlgorithm succeeded")
	return
}
func StartAlgorithm (w http.ResponseWriter, r *http.Request) {

	var algorithm v.BasicUserSendIn
	var decoder *json.Decoder

	log.Println("StartAlgorithm attempted")

	decoder = json.NewDecoder(r.Body)
	err := decoder.Decode(&algorithm)

	if err != nil {
		log.Fatal("Error decoding JSON in StartAlgorithm: ", err.Error())
	}

	if algorithm.UserName != "" {
		a.StartAlgorithm(algorithm.UserName, algorithm.Algorithm)

		send, err := json.Marshal(true)
		if err != nil {
			log.Fatal("Error sending json in StartAlgorithm:", err.Error())
		}
		defer r.Body.Close()
		fmt.Fprint(w, string(send))
		log.Println("StartAlgorithm succeeded")
		return
	}

	send, err := json.Marshal(false)
	if err != nil {
		log.Fatal("Error sending json in StartAlgorithm:", err.Error())
	}
	defer r.Body.Close()
	fmt.Fprint(w, string(send))
	log.Println("StartAlgorithm failed")

}
func StopAlgorithm (w http.ResponseWriter, r *http.Request) {

	var algorithm v.BasicUserSendIn
	var decoder *json.Decoder

	log.Println("StopAlgorithm attempted")

	decoder = json.NewDecoder(r.Body)
	err := decoder.Decode(&algorithm)

	if err != nil {
		log.Fatal("Error decoding JSON in StopAlgorithm: ", err.Error())
	}

	if algorithm.UserName != "" {
		a.ShutDownAlgorithm(algorithm.UserName, algorithm.Algorithm)

		send, err := json.Marshal(true)
		if err != nil {
			log.Fatal("Error sending json in RemoveFromUserStocks:", err.Error())
		}
		defer r.Body.Close()
		fmt.Fprint(w, string(send))
		log.Println("StopAlgorithm succeeded")
		return
	}

	send, err := json.Marshal(false)
	if err != nil {
		log.Fatal("Error sending json in StopAlgorithm:", err.Error())
	}
	defer r.Body.Close()
	fmt.Fprint(w, string(send))
	log.Println("StopAlgorithm failed")

}
func DeleteAlgorithm (w http.ResponseWriter, r *http.Request) {
	var basicUserSendIn v.BasicUserSendIn
	var decoder *json.Decoder

	log.Println("DeleteAlgorithm attempted")

	decoder = json.NewDecoder(r.Body)
	err := decoder.Decode(&basicUserSendIn)

	if err != nil {
		log.Fatal("Error decoding JSON in DeleteAlgorithm: ", err.Error())
	}
	a.DeleteAlgorithm(basicUserSendIn.Algorithm, basicUserSendIn.UserName)
	UserAlgorithms := a.GetUserAlgorithms(basicUserSendIn.UserName)

	send, err := json.Marshal(UserAlgorithms)
	if err != nil {
		log.Fatal("Error sending json in DeleteAlgorithm:", err.Error())
	}
	defer r.Body.Close()
	fmt.Fprint(w, string(send))
	log.Println("DeleteAlgorithm succeeded")
	go a.DeleteAndStopAlgorithm(basicUserSendIn.UserName, basicUserSendIn.Algorithm)
	return

}
func RunViable (w http.ResponseWriter, r *http.Request) {

	var algorithm v.BasicUserSendIn
	var decoder *json.Decoder

	log.Println("RunViable attempted")

	decoder = json.NewDecoder(r.Body)
	err := decoder.Decode(&algorithm)

	if err != nil {
		log.Fatal("Error decoding JSON in RunViable: ", err.Error())
	}

	if algorithm.UserName != "" {
		a.RunViable(algorithm.UserName, algorithm.Algorithm)

		send, err := json.Marshal(true)
		if err != nil {
			log.Fatal("Error sending json in RunViable:", err.Error())
		}
		defer r.Body.Close()
		fmt.Fprint(w, string(send))
		log.Println("RunViable succeeded")
		return
	}

	send, err := json.Marshal(true)
	if err != nil {
		log.Fatal("Error sending json in RunViable:", err.Error())
	}
	defer r.Body.Close()
	fmt.Fprint(w, string(send))
	log.Println("RunViable failed")

}
func RunPurchasable (w http.ResponseWriter, r *http.Request) {

	var algorithm v.BasicUserSendIn
	var decoder *json.Decoder

	log.Println("RunPurchasable attempted")

	decoder = json.NewDecoder(r.Body)
	err := decoder.Decode(&algorithm)

	if err != nil {
		log.Fatal("Error decoding JSON in RunPurchasable: ", err.Error())
	}

	if algorithm.UserName != "" {
		a.RunPurchasable(algorithm.UserName, algorithm.Algorithm)

		send, err := json.Marshal(true)
		if err != nil {
			log.Fatal("Error sending json in RunPurchasable:", err.Error())
		}
		defer r.Body.Close()
		fmt.Fprint(w, string(send))
		log.Println("RunPurchasable succeeded")
		return
	}

	send, err := json.Marshal(true)
	if err != nil {
		log.Fatal("Error sending json in RunPurchasable:", err.Error())
	}
	defer r.Body.Close()
	fmt.Fprint(w, string(send))
	log.Println("RunPurchasable failed")

}
func RunSell (w http.ResponseWriter, r *http.Request) {

	var algorithm v.BasicUserSendIn
	var decoder *json.Decoder

	log.Println("RunSell attempted")

	decoder = json.NewDecoder(r.Body)
	err := decoder.Decode(&algorithm)

	if err != nil {
		log.Fatal("Error decoding JSON in RunSell: ", err.Error())
	}

	if algorithm.UserName != "" {
		a.RunSell(algorithm.UserName, algorithm.Algorithm)

		send, err := json.Marshal(true)
		if err != nil {
			log.Fatal("Error sending json in RunSell:", err.Error())
		}
		defer r.Body.Close()
		fmt.Fprint(w, string(send))
		log.Println("RunSell succeeded")
		return
	}

	send, err := json.Marshal(true)
	if err != nil {
		log.Fatal("Error sending json in RunSell:", err.Error())
	}
	defer r.Body.Close()
	fmt.Fprint(w, string(send))
	log.Println("RunSell failed")
}

func ShutDownProgram (w http.ResponseWriter, r *http.Request) {

	var shutDown bool
	var decoder *json.Decoder

	log.Println("ShutDownProgram attempted")

	decoder = json.NewDecoder(r.Body)
	err := decoder.Decode(&shutDown)

	if err != nil {
		log.Fatal("Error decoding JSON in ShutDownProgram: ", err.Error())
	}

	if shutDown == true {

		send, err := json.Marshal(true)
		if err != nil {
			log.Fatal("Error sending json in ShutDownProgram:", err.Error())
		}
		defer r.Body.Close()
		fmt.Fprint(w, string(send))
		log.Println("ShutDownProgram succeeded")

		ProgramShutDown <- true

		return
	}

	send, err := json.Marshal(true)
	if err != nil {
		log.Fatal("Error sending json in ShutDownProgram:", err.Error())
	}
	defer r.Body.Close()
	fmt.Fprint(w, string(send))
	log.Println("ShutDownProgram failed")

}

func SendNotifications (w http.ResponseWriter, r *http.Request) {
	var basicUserSendIn v.BasicUserSendIn
	var decoder *json.Decoder
	var notifications []v.Notification

	log.Println("SendNotifications attempted")

	decoder = json.NewDecoder(r.Body)
	err := decoder.Decode(&basicUserSendIn)

	if err != nil {
		log.Fatal("Error decoding JSON in SendNotifications: ", err.Error())
	}

	notifications = db.GetNotifications(basicUserSendIn.UserName, basicUserSendIn.Algorithm)

	send, err := json.Marshal(notifications)
	if err != nil {
		log.Fatal("Error sending json in SendNotifications:", err.Error())
	}
	defer r.Body.Close()
	fmt.Fprint(w, string(send))
	log.Println("SendNotifications succeeded")
}
func DeleteNotifications (w http.ResponseWriter, r *http.Request) {
	var deleteNotificationsSendIn v.DeleteNotificationsSendIn
	var decoder *json.Decoder

	log.Println("DeleteNotifications attempted")

	decoder = json.NewDecoder(r.Body)
	err := decoder.Decode(&deleteNotificationsSendIn)

	if err != nil {
		log.Fatal("Error decoding JSON in DeleteNotifications: ", err.Error())
	}

	db.DeleteNotifications(deleteNotificationsSendIn.Notifications, deleteNotificationsSendIn.Algorithm)
	log.Println("DeleteNotifications succeeded")
}

func trimBars (stocks []v.Stock) []v.Stock {
	for i, _ := range stocks {
		stocks[i].WeekBars = nil
		stocks[i].DayBars = nil
		stocks[i].FourHourBars = nil
		stocks[i].HourBars = nil
	}
	return stocks
}