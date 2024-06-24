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
	}
	localApiConfig := &handlers.LocalApiConfig{
		Config: apiConfig,
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
	router.Use(middlewares.RateLimitMiddleware(redisClient, 50, time.Minute))

	// HTTP routes
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	router.POST("/", localApiConfig.Broker)
	router.POST("/handle", localApiConfig.HandleSubmission)

	// GRPC routes for calling grpc services
	router.POST("/log-grpc", localApiConfig.WriteLog)
	router.POST("/payment-grpc", localApiConfig.PaymentViaGRPC)
	router.POST("/enquiry-grpc", localApiConfig.EnquiryViaGRPC)

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
