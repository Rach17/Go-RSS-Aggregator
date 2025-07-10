package main

import (
	"log"      // Importing log for logging errors and messages
	"os"       // Importing os for environment variable access
	"strconv"  // Importing strconv for string conversion
	"database/sql" // Importing database/sql for SQL database operations
	"github.com/Rach17/Go-RSS-Aggregator/api" // Importing the api package for API server functionality
	"github.com/Rach17/Go-RSS-Aggregator/repository" // Importing the repository package for database interactions
	"github.com/Rach17/Go-RSS-Aggregator/service" // Importing the service package for business logic

	"github.com/joho/godotenv" // Importing godotenv to load environment variables from .env file
	_ "github.com/lib/pq" // Importing the PostgreSQL driver for database connection (underscore means we don't use it directly)
)





func main() {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found, using default environment variables")
	}
	// Get the PORT environment variable
	portString := os.Getenv("PORT")
	if portString == "" || portString == "0"{
		log.Fatal("PORT environment variable is not set")
	}
	port, err := strconv.Atoi(portString) // Convert the port string to an integer
	if err != nil {
		log.Fatalf("Invalid PORT environment variable: %v", err) // Log an error if conversion fails
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}
	connection, err := sql.Open("postgres", dbURL) // Open a connection to the PostgreSQL database using the provided DB_URL
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err) // Log an error if the connection fails
	}
	defer connection.Close() // Ensure the database connection is closed when the function exits
	userRepo := repository.NewDBUserRepository(connection) // Create a new user repository using the database connection
	userService := service.NewUserService(userRepo) // Create a new user service using the user
	authService := service.NewAuthService(userRepo) // Create a new authentication service using the user repository
	feedRepo := repository.NewDBFeedRepository(connection) // Create a new feed repository using the database connection
	feedService := service.NewFeedService(feedRepo) // Create a new feed service using the feed repository

	server := api.NewServer(port, userService, authService, feedService) // Create a new API server with the specified port and services
	server.Start()
}
