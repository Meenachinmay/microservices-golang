package main

import (
	"broker/cmd/api/handlers"
	"broker/internal/config"
	"broker/middlewares"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
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

	// Initialize the REDIS here
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	// connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	// setting configuration for global use
	apiConfig := &config.Config{
		Rabbit:      rabbitConn,
		RedisClient: redisClient,
		//Producer: producer,
	}
	localApiConfig := &handlers.LocalApiConfig{
		Config: apiConfig,
		//Producer: producer,
	}

	// initialize the gin router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost", "https://*", "http://*"}, // Specify the exact origin of your Next.js app
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true, // Important: Must be true when credentials are included
		MaxAge:           12 * time.Hour,
	}))

	// apply rate limiting middleware here
	router.Use(middlewares.RateLimitMiddleware(redisClient, 2, time.Minute))

	// routes
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	router.POST("/", localApiConfig.Broker)
	router.POST("/handle", localApiConfig.HandleSubmission)

	// routes for calling grpc services
	router.POST("/log-grpc", localApiConfig.LogViaGRPC)
	router.POST("/payment", localApiConfig.PaymentViaGRPC)

	// start the server
	log.Printf("Starting broker service on port %s\n", webPort)
	log.Fatal(router.Run(":" + webPort))
}

// connect to rabbitmq api-gateway
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
