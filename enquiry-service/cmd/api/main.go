package main

import (
	"database/sql"
	"enquiry-service/config"
	"enquiry-service/handlers"
	"enquiry-service/internal/database"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math"
	"os"
	"time"
)

var counts int64

func main() {
	// connect to database
	conn := connectToDB()
	defer conn.Close()

	// connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	// setting API configuration
	apiConfig := &config.ApiConfig{
		DB:     database.New(conn),
		Rabbit: rabbitConn,
	}
	localApiConfig := &handlers.LocalApiConfig{
		ApiConfig: apiConfig,
	}

	_ = &handlers.LocalApiConfig{
		ApiConfig: apiConfig,
	}

	// Initialize grpc
	go StartGrpcServer()

	// Initialize the router
	router := gin.Default()

	// Configure cors
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "*"}, // Specify the exact origin of your Next.js app
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true, // Important: Must be true when credentials are included
		MaxAge:           12 * time.Hour,
	}))

	// routes
	router.POST("/handle-enquiry", localApiConfig.HandleANewEnquiry)
	router.GET("/fetch-properties", localApiConfig.HandleFetchAllProperties)

	// start the server
	log.Fatal(router.Run(":80"))
}

func openDB() (*sql.DB, error) {
	dbURL := os.Getenv("DATABASE_URL_ENQUIRY_SERVICE")
	if dbURL == "" {
		return nil, errors.New("missing DATABASE_URL")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	for {
		connection, err := openDB()
		if err != nil {
			log.Println("Could not connect to database, Postgres is not ready...")
			counts += 1
		} else {
			log.Println("Connected to database...")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Waiting for database to become ready...")
		time.Sleep(2 * time.Second)
		continue
	}
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
