package services

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/pools"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

func ReciverRmq() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			var queue model.SendToQueue
			var response string
			if err := json.Unmarshal(d.Body, &queue); err != nil {
				logrus.Error(err)
			}
			logrus.Info("Recivered ", queue)
			//--------------------------------------------------------------------------
			if queue.Type == "Pool" {
				response, err = pools.PoolCreateHandle(queue.EqationJson, queue.CoinMap, queue.PoolCreationArray)
				if err != nil {
					logrus.Error(err)
				}
				logrus.Info("Pools Created")
			} else if queue.Type == "CionConvert" {
				response, err = pools.PathPaymentHandle(queue.CoinConvert)
				if err != nil {
					logrus.Error(err)
				}
				logrus.Info("Coin Converted")
			} else {
				logrus.Error("Queue Request Type Error")
			}
			logrus.Info(response)
			logrus.Info("QUEUE TASK DONE")
		}
	}()

	logrus.Info(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func SendToQueue(queue model.SendToQueue) string {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(queue)
	body := reqBodyBytes.Bytes()
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent ")
	return "sent"
}

func failOnError(err error, msg string) {
	if err != nil {
		logrus.Error("%s: %s", msg, err)
	}
}
