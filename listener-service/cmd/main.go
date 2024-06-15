package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"listener-service/actions/consumers"
	"log"
	"math"
	"math/rand"
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

	// start consumers for services
	go startLogConsumer(rabbitConn)
	go startMailConsumer(rabbitConn)

	select {}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOffTime = 1 * time.Second
	maxBackOffTime := 30 * time.Second
	rand.Seed(time.Now().UnixNano())

	for {
		conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err != nil {
			log.Printf("RabbitMQ not yet ready... (%v)", err)
			counts++
		} else {
			log.Println("Connected to RabbitMQ")
			return conn, nil
		}

		if counts > 5 {
			return nil, fmt.Errorf("failed to connect to RabbitMQ after multiple attempts: %v", err)
		}

		jitter := time.Duration(rand.Int63n(int64(backOffTime)))
		backOffTime = time.Duration(math.Min(float64(maxBackOffTime), float64(backOffTime*2)))
		log.Printf("Backing off for %v seconds (with jitter)...", backOffTime+jitter)
		time.Sleep(backOffTime + jitter)
	}
}

func startLogConsumer(conn *amqp.Connection) {
	consumer, err := consumers.NewLogConsumer(conn)
	if err != nil {
		log.Fatalf("Failed to create log consumer:START_LOG_CONSUMER %v", err)
	}
	err = consumer.ConsumeLogs([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Fatalf("Failed to consume logs: %v", err)
	}
}

func startMailConsumer(conn *amqp.Connection) {
	consumer, err := consumers.NewMailConsumer(conn)
	if err != nil {
		log.Fatalf("Failed to create mail consumer:START_MAIL_CONSUMER %v", err)
	}
	err = consumer.ConsumeMails()
	if err != nil {
		log.Fatalf("Failed to send mail:LISTENER %v", err)
	}
}

//"amqp://guest:guest@rabbitmq:5672/"
