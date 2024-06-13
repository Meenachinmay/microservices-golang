package actions

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"log_topics", // name of exchange
		"topic",      // type
		true,
		false,
		false,
		false,
		nil,
	)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"test_queue",
		false,
		false,
		true,
		false,
		nil,
	)
}
