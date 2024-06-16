package config

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Rabbit      *amqp.Connection
	RedisClient *redis.Client
	//Producer *kafka.Producer
}
