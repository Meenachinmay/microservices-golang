package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"listener-service/actions"
	"log"
	"math"
	"os"
	"time"
)

func main() {
	// initialize kafka consumer
	//kafkaConsumer, err := kafka.NewConsumer(&kafka.ConfigMap{
	//	"bootstrap.servers": "localhost:9092",
	//	"group.id":          "log-service.consumer",
	//	"auto.offset.reset": "earliest",
	//})
	//
	//if err != nil {
	//	fmt.Printf("Failed to create consumer: %s\n", err)
	//	os.Exit(1)
	//} else {
	//	fmt.Printf("consumer created...:LISTENER \n")
	//}
	//defer kafkaConsumer.Close()

	// initialize kafka consumer admin client
	//adminClientConsumer, err := kafka.NewAdminClientFromConsumer(kafkaConsumer)
	//if err != nil {
	//	fmt.Printf("Failed to create consumer admin client: %s\n", err)
	//	os.Exit(1)
	//}
	//defer adminClientConsumer.Close()
	//
	//// Subscribe to the topic = new-user-signup
	//err = kafkaConsumer.SubscribeTopics([]string{"new-log"}, nil)
	//if err != nil {
	//	fmt.Printf("Failed to subscribe to topics: %s\n", err)
	//	os.Exit(1)
	//} else {
	//	fmt.Printf("subscribed to producer topics:LISTENER \n")
	//}

	// get the list of all topics
	//topicMetadata, err := adminClientConsumer.GetMetadata(nil, true, 10000)
	//if err != nil {
	//	fmt.Printf("failed to get topic metadata: %s\n", err)
	//	os.Exit(1)
	//}
	//fmt.Printf("all topics in the cluster: \n")
	//for _, topic := range topicMetadata.Topics {
	//	fmt.Println(topic)
	//}
	//
	//go handlers.ConsumeMessages(kafkaConsumer)
	//
	//signChan := make(chan os.Signal, 1)
	//signal.Notify(signChan, syscall.SIGINT, syscall.SIGTERM)
	//<-signChan
	//
	//fmt.Println("shutdown down gracefully")

	//// connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	log.Println("listening for and consuming messages")
	consumer, err := actions.NewConsumer(rabbitConn)
	if err != nil {
		log.Panic(err)
	}

	err = consumer.ConsumeLogs([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Panic(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOffTime = 1 * time.Second

	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOffTime = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOffTime)
		continue
	}

	log.Println("Connected to RabbitMQ...")
	return connection, nil
}
