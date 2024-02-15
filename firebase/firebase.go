package firebase

import (
	"context"
	"log"
	"time"

	"../db"
	"../variables"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

func InitFirebase() {
	var err error
	opt := option.WithCredentialsFile("")

	variables.FireBase, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatal("Error initializing firebase:", err.Error())
	}
}
func PingStockToSell(username string, stock string, algorithm int) {
	userAccount := db.GetUser(username)
	messagingClient, err := variables.FireBase.Messaging(context.Background())
	if err != nil {
		log.Fatal("Error initializing firebase messaging client in PingStockToSell:", err.Error())
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: "Time to sell",
			Body:  stock,
		},
		Token: userAccount.FirebaseId,
	}
	response, err := messagingClient.Send(context.Background(), message)
	if err != nil {
		log.Fatal("Error sending firebase message in PingStockToSell:", err.Error())
	}
	notification := variables.Notification{
		UserName: username,
		Date:     time.Now().Format("2006-01-02 15:04:05"),
		Title:    "Time to sell",
		Message:  stock,
	}
	db.StoreNotification(notification, algorithm)
	log.Println("Successfully sent message: ", response)
}
func PingPossibleTop(username string, stock string, algorithm int) {
	userAccount := db.GetUser(username)
	messagingClient, err := variables.FireBase.Messaging(context.Background())
	if err != nil {
		log.Fatal("Error initializing firebase messaging client in PingStockToSell:", err.Error())
	}
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: "Possible top",
			Body:  stock,
		},
		Token: userAccount.FirebaseId,
	}
	response, err := messagingClient.Send(context.Background(), message)
	if err != nil {
		log.Fatal("Error sending firebase message in PingStockToSell:", err.Error())
	}
	notification := variables.Notification{
		UserName: username,
		Date:     time.Now().Format("2006-01-02 15:04:05"),
		Title:    "Possible top",
		Message:  stock,
	}
	db.StoreNotification(notification, algorithm)
	log.Println("Successfully sent message: ", response)

}
func PingUpdate(username string, title string, body string, algorithm int) {
	userAccount := db.GetUser(username)
	messagingClient, err := variables.FireBase.Messaging(context.Background())
	if err != nil {
		log.Fatal("Error initializing firebase messaging client in PingUpdate:", err.Error())
	}
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Token: userAccount.FirebaseId,
	}
	response, err := messagingClient.Send(context.Background(), message)
	if err != nil {
		log.Fatal("Error sending firebase message in PingUpdate:", err.Error())
	}

	notification := variables.Notification{
		UserName: username,
		Date:     time.Now().Format("2006-01-02 15:04:05"),
		Title:    title,
		Message:  body,
	}
	db.StoreNotification(notification, algorithm)

	log.Println("Successfully sent message:", response)
}
