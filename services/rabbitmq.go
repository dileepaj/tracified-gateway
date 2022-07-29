package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/pools"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

var (
	USER     = commons.GoDotEnvVariable("RABBITUSER")
	PASSWORD = commons.GoDotEnvVariable("RABBITPASSWORD")
	HOSTNAME = commons.GoDotEnvVariable("RABBITHOSTNAME")
	PORT     = commons.GoDotEnvVariable("RABBITPORT")
)

func ReciverRmq() error{
	rabbitCoonection := `amqp://` + USER + `:` + PASSWORD + `@` + HOSTNAME + `:` + PORT + `/`
	conn, err := amqp.Dial(rabbitCoonection)
	if err != nil {
		logrus.Error("rabbitmq connection issue ", err)
		return err
	}
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"socialimpact", // name
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
			if queue.Type == "POOL" {
				response, err = pools.PoolCreateHandle(queue.EqationJson, queue.CoinMap, queue.PoolCreationArray)
				if err != nil {
					logrus.Error(err)
				}
				logrus.Info("Pools Created")
			} else if queue.Type == "COINCONVERT" {
				response, err = pools.PathPaymentHandle(queue.CoinConvert)
				if err != nil {
					logrus.Error(err)
				}
				logrus.Info("Coin Converted")
				logrus.Info("QUEUE TASK DONE")
			} else {
				logrus.Error("Queue Request Type Error")
			}
			out, err := json.Marshal(response)
			if err != nil {
				logrus.Error(err)
			}
	
			fmt.Println("--------------- Pool response ---------  ", string(out))
		}
	}()

	logrus.Info(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	return nil
}

func SendToQueue(queue model.SendToQueue) (string,error) {
	rabbitConnection := `amqp://` + USER + `:` + PASSWORD + `@` + HOSTNAME + `:` + PORT + `/`
	conn, err := amqp.Dial(rabbitConnection)
	if err != nil {
		logrus.Error("rabbitmq connection issue ", err)
		return "",err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"socialimpact", // name
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
	return "sent",nil
}

func failOnError(err error, msg string) {
	if err != nil {
		logrus.Error("%s: %s", msg, err)
	}
}
