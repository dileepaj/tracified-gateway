package services

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/dileepaj/tracified-gateway/model"
	"github.com/rabbitmq/amqp091-go"
)

type Person struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Country string `json:"country"`
}

func SubmitData(diliver amqp091.Delivery) {
	log.Printf("Reciver a message: %s", diliver.Body)
	// fmt.Println(diliver.Body)

	// Convert the JSON string to the txnBody struct
	var txnBody model.TransactionCollectionBody
	err := json.Unmarshal(diliver.Body, &txnBody)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Byte slice as string:", txnBody)
	diliver.Ack(false)
	return
}
