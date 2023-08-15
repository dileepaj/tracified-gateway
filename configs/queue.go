package configs

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type QueueConfig struct {
	Name   string
	Method func(delivery amqp.Delivery)
	Prefix string
}

var QueueBackLinks = QueueConfig{
	Name: "backlinks",
}

var QueueTransaction = QueueConfig{
	Prefix: "transaction.",
}

var Queues = map[string]QueueConfig{
	"backlinks":    QueueBackLinks,
	"transaction.": QueueTransaction,
}
