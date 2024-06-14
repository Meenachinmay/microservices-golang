package main

import (
	"broker/cmd/api/handlers"
	"broker/internal/config"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math"
	"net/http"
	"os"
	"time"
)

const webPort = "80"

func main() {
	// initialize kafka here = producer
	//producer, err := kafka.NewProducer(&kafka.ConfigMap{
	//	"bootstrap.servers": "localhost:9092",
	//})
	//
	//if err != nil {
	//	log.Fatalf("Error creating producer: %s", err)
	//} else {
	//	log.Println("Producer created...:BROKER_SERVICE")
	//}
	//defer producer.Close()

	// initialize the kafka admin client
	//adminClient, err := kafka.NewAdminClientFromProducer(producer)
	//if err != nil {
	//	log.Fatalf("failed to create admin client from producer: %s", err)
	//} else {
	//	log.Println("Kafka Admin client created...:BROKER_SERVICE")
	//}
	//defer adminClient.Close()

	//create the topic if it doesn't exist
	//topic := "new-log"
	//err = utils.CreateKafkaTopic(adminClient, topic)
	//if err != nil {
	//	log.Fatalf("failed to create kafka topic: %s", err)
	//} else {
	//	fmt.Printf("Kafka topic %s created...:BROKER_SERVICE\n", topic)
	//}

	//--------------------------------------------------above this line there is a kafka setup--------------------------------------

	// connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	// setting configuration for global use
	apiConfig := &config.Config{
		Rabbit: rabbitConn,
		//Producer: producer,
	}
	localApiConfig := &handlers.LocalApiConfig{
		Config: apiConfig,
		//Producer: producer,
	}

	// initialize the gin router
	mux := chi.NewRouter()

	// setting cors
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost", "https://*", "http://*"}, // Specify the exact origin of your Next.js app
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true, // Important: Must be true when credentials are included
		MaxAge:           300,
	}))

	// routes
	mux.Use(middleware.Heartbeat("/ping"))
	mux.Post("/", localApiConfig.Broker)
	mux.Post("/handle", localApiConfig.HandleSubmission)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: mux,
	}

	log.Printf("Starting broker service on port %s\n", webPort)
	// start the server
	err = srv.ListenAndServe()
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
