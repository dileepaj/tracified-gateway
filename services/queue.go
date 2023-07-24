package services

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"

	"github.com/dileepaj/tracified-gateway/commons"
)

const queuePrefix = "gateway."

var queueConnection *amqp.Connection

var queueChannel *amqp.Channel

var queueError error

var queueConnectionOnce sync.Once
var queueChannelOnce sync.Once

var queues = make(map[string]amqp.Queue)
var queuesConsumers = make(map[string]bool)

func getQueueName(queueName string) string {
	if strings.HasPrefix(queueName, queuePrefix) {
		return queueName
	}

	return queuePrefix + queueName
}

func GetQueueConnection() (*amqp.Connection, error) {
	queueConnectionOnce.Do(func() {
		conn, err := amqp.Dial(commons.GoDotEnvVariable("RABBITMQ_SERVER_URI"))

		queueConnection = conn
		queueError = err

	})

	if queueError != nil {
		log.Error(queueError)
	}
	return queueConnection, queueError
}

func GetQueueChannel() (*amqp.Channel, error) {
	if _, err := GetQueueConnection(); err != nil {
		log.Error(err)
		return nil, err
	}

	queueChannelOnce.Do(func() {
		ch, err := queueConnection.Channel()

		if err != nil {
			log.Error(err)
		}

		queueChannel = ch
		queueError = err

	})

	return queueChannel, queueError
}

func GetQueue(queueName string) (amqp.Queue, error) {
	queueName = getQueueName(queueName)
	ch, err := GetQueueChannel()

	if err != nil {
		return amqp.Queue{}, err
	}

	q, ok := queues[queueName]
	if !ok {
		q, err = ch.QueueDeclare(
			queueName,
			true,
			false,
			false,
			false,
			amqp.Table{
				"x-single-active-consumer": true,
			},
		)
		if err != nil {
			log.Error(err)
			return amqp.Queue{}, err
		}
		err = ch.Qos(1, 0, false)
		if err != nil {
			log.Error(err)
			return amqp.Queue{}, err
		}

		queues[queueName] = q
	}

	return q, err
}

func PublishToQueue(queueName string, message string) error {
	queueName = getQueueName(queueName)
	q, _ := GetQueue(queueName)
	ch, _ := GetQueueChannel()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(message),
		})

	return err
}

func RegisterWorker(queueName string, cmd func(delivery amqp.Delivery)) error {
	queueName = getQueueName(queueName)
	_, ok := queuesConsumers[queueName]
	if ok {
		return nil
	}

	ch, _ := GetQueueChannel()
	q, _ := GetQueue(queueName)

	messages, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Error(err)
		return err
	}
	var forever chan struct{}

	go func() {
		for d := range messages {
			cmd(d)
		}
	}()
	<-forever

	queuesConsumers[queueName] = true
	return err
}
