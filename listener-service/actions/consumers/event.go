package consumers

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"log_topics", // name of exchange
		"direct",     // type
		true,
		false,
		false,
		false,
		nil,
	)
}

func DeclareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
}

func declareMailExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"mail_exchange", // name of exchange
		"direct",        // type
		true,
		false,
		false,
		false,
		nil,
	)
}

func declareMailQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"mail_queue", // Name of the queue
		false,        // Durable
		false,        // Delete when unused
		false,        // Exclusive
		false,        // No-wait
		nil,          // Arguments
	)
}