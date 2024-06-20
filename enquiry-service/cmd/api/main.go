package main

import (
	"database/sql"
	"enquiry-service/config"
	"enquiry-service/handlers"
	"enquiry-service/internal/database"
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

var counts int64

func main() {
	// connect to database
	conn := connectToDB()
	defer conn.Close()

	// setting API configuration
	apiConfig := &config.ApiConfig{
		DB: database.New(conn),
	}

	localApiConfig := &handlers.LocalApiConfig{
		ApiConfig: apiConfig,
	}

	_ = &handlers.LocalApiConfig{
		ApiConfig: apiConfig,
	}

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
