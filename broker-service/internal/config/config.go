package config

import amqp "github.com/rabbitmq/amqp091-go"

type Config struct {
	Rabbit *amqp.Connection
}
