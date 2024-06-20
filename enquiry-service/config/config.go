package config

import (
	"enquiry-service/internal/database"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ApiConfig struct {
	DB     *database.Queries
	Rabbit *amqp.Connection
}
