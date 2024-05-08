package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	EXCHANGE_LOGS_TOPIC = "exchange_logs_topic"
	QUEUE_LOGS          = "queue_logs"
)

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		EXCHANGE_LOGS_TOPIC, // name
		"topic",             // type
		true,                // durable?
		false,               // auto-deleted?
		false,               // internal?
		false,               // no-wait?
		nil,                 // arguements?
	)
}
