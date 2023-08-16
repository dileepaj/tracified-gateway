package services

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/configs"
	"github.com/dileepaj/tracified-gateway/services/cache"
)

const queuePrefix = "gateway."
const queueMaxTry = 10
const queueDeadLetter = "dead-letter"
const QueueCacheName = "gateway:current-queues"
const queueCacheTime = 60 * 60 * 3

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
				"x-consumer-timeout":       int(2 * time.Minute),
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

func PublishToQueue(queueName string, message string, args ...interface{}) error {
	queueName = getQueueName(queueName)
	q, _ := GetQueue(queueName)
	ch, _ := GetQueueChannel()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	payload := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         []byte(message),
	}

	var registerW func(delivery amqp.Delivery)
	// TODO add support for others to be passed
	for _, arg := range args {
		switch t := arg.(type) {
		case amqp.Table:
			payload.Headers = t
		case func(delivery amqp.Delivery):
			registerW = t
		}
	}

	err := ch.PublishWithContext(ctx, "", q.Name, false, false, payload)

	if err != nil {
		log.Error(err)
		return err
	}

	log.Info("Message published " + queueName + "cache time" + strconv.FormatInt(time.Now().Unix()+queueCacheTime, 10))
	cache.InsertSortedSet(QueueCacheName, queueName, float64(time.Now().Unix()+queueCacheTime))

	if registerW != nil {
		err = RegisterWorker(queueName, registerW)
	}

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

	go func() {
		for d := range messages {
			if d.Redelivered {
				d.Ack(false)
				var deliveryCount int32 = 1
				if d.Headers["x-delivery-count"] != nil {
					deliveryCount += d.Headers["x-delivery-count"].(int32)
				}
				if deliveryCount >= queueMaxTry {
					PublishToQueue(queueDeadLetter, string(d.Body), amqp.Table{
						"x-delivery-count":      deliveryCount,
						"x-delivery-queue-name": queueName,
					})
				} else {
					PublishToQueue(queueName, string(d.Body), amqp.Table{"x-delivery-count": deliveryCount})
				}
			} else {
				cmd(d)
			}
		}
	}()

	queuesConsumers[queueName] = true
	return err
}

func QueueScheduleWorkers() {
	client := cache.Client()
	client.ZRemRangeByScore(context.Background(), QueueCacheName, "0", strconv.FormatInt(time.Now().Unix(), 10))

	queueNames := client.ZRange(context.Background(), QueueCacheName, 0, -1).Val()

	log.Info("Current Queues: " + strings.Join(queueNames, ", "))
	for _, name := range queueNames {
		queueName := getQueueName(name)
		_, ok := queuesConsumers[queueName]
		if ok {
			continue
		}

		for n, queue := range configs.Queues {
			n = getQueueName(n)
			if queueName == n || strings.HasPrefix(queueName, n) {
				RegisterWorker(queueName, queue.Method)
			}
		}
	}
}
