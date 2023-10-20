package services

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/configs"
	"github.com/dileepaj/tracified-gateway/services/cache"
	log "github.com/sirupsen/logrus"
)

const queuePrefix = "gateway."
const queueMaxTry = 10
const queueDeadLetter = "dead-letter"
const QueueCacheName = "gateway:current-queues"
const queueCacheTime = 60 * 60 * 3

var queueConnection *amqp.Connection

var channels = make(map[string]*amqp.Channel)
var consumerChannels = make(map[string]*amqp.Channel)

var queueError error

var queueConnectionOnce sync.Once

var queues = make(map[string]amqp.Queue)
var queuesConsumers = make(map[string]bool)

var mu sync.Mutex

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

func GetQueueChannel(queueName string) (*amqp.Channel, error) {
	var err error
	if _, err = GetQueueConnection(); err != nil {
		log.Error(err)
		return nil, err
	}

	ch, ok := channels[queueName]

	if ok {
		return ch, nil
	}

	ch, err = queueConnection.Channel()
	mu.Lock()
	channels[queueName] = ch
	mu.Unlock()

	return ch, err
}

func GetConsumerQueueChannel(queueName string) (*amqp.Channel, error) {
	var err error
	if _, err = GetQueueConnection(); err != nil {
		log.Error(err)
		return nil, err
	}

	mu.Lock()
	ch, ok := consumerChannels[queueName]
	mu.Unlock()

	if ok {
		return ch, nil
	}

	ch, err = queueConnection.Channel()

	mu.Lock()
	consumerChannels[queueName] = ch
	mu.Unlock()

	return ch, err
}

func GetQueue(queueName string) (amqp.Queue, error) {
	queueName = getQueueName(queueName)
	ch, err := GetQueueChannel(queueName)

	if err != nil {
		return amqp.Queue{}, err
	}

	mu.Lock()
	q, ok := queues[queueName]
	mu.Unlock()
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
			log.Error("QueueDeclareError ", err)
			return amqp.Queue{}, err
		}
		err = ch.Qos(1, 0, false)
		if err != nil {
			log.Error("Channel QOS: ", err)
			return amqp.Queue{}, err
		}

		mu.Lock()
		queues[queueName] = q
		mu.Unlock()
	}

	return q, err
}

func PublishToQueue(queueName string, message string, args ...interface{}) error {
	queueName = getQueueName(queueName)
	q, _ := GetQueue(queueName)
	ch, _ := GetQueueChannel(queueName)

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

	cache.InsertSortedSet(QueueCacheName, queueName, float64(time.Now().Unix()+queueCacheTime))

	if registerW != nil {
		err = RegisterWorker(queueName, registerW)
	}

	return err
}

func RegisterWorker(queueName string, cmd func(delivery amqp.Delivery)) error {
	queueName = getQueueName(queueName)
	mu.Lock()
	_, ok := queuesConsumers[queueName]
	mu.Unlock()
	if ok {
		return nil
	}

	ch, _ := GetConsumerQueueChannel(queueName)
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

	mu.Lock()
	queuesConsumers[queueName] = true
	mu.Unlock()
	return err
}

func QueueScheduleWorkers() {
	client := cache.Client()
	client.ZRemRangeByScore(context.Background(), QueueCacheName, "0", strconv.FormatInt(time.Now().Unix(), 10))

	queueNames := client.ZRange(context.Background(), QueueCacheName, 0, -1).Val()

	for _, name := range queueNames {
		queueName := getQueueName(name)
		_, ok := queuesConsumers[queueName]
		if ok {
			continue
		}

		for n, queue := range configs.Queues {
			n = getQueueName(n)
			if queueName == n || strings.HasPrefix(queueName, n) {
				RegisterWorker(queueName, queue().Method)
			}
		}
	}
}
