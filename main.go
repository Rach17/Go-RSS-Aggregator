package main

import (
	"fmt"      // Importing fmt for formatted I/O
	"log"      // Importing log for logging errors and messages
	"net/http" // Importing net/http for HTTP server functionality
	"os"       // Importing os for environment variable access
	"strconv"  // Importing strconv for string conversion
	"database/sql" // Importing database/sql for SQL database operations

	"github.com/Rach17/Go-RSS-Aggregator/internal/db" // Importing the db package for database queries

	"github.com/joho/godotenv" // Importing godotenv to load environment variables from .env file
	"github.com/rs/cors"       // Importing rs/cors for handling CORS (Cross-Origin Resource Sharing) in HTTP requests
	_ "github.com/lib/pq" // Importing the PostgreSQL driver for database connection (underscore means we don't use it directly)
)

// config 
type Config struct {
	DB *db.Queries
}

// Middleware type
type Middleware func(http.Handler) http.Handler

func createRouter() *http.ServeMux {
	// Create a new HTTP request multiplexer (router)
	router := http.NewServeMux()
	return router
}

func createServer(router http.Handler, port int) *http.Server {
	// Create a new HTTP server with the specified port
	return &http.Server{
		Addr:    ":" + strconv.Itoa(port), // Set the server address to the specified port
		Handler: router,                   // Use the provided router for handling requests
	}
}

func startServer(server *http.Server) error {
	// Start the HTTP server and listen for incoming requests
	return server.ListenAndServe()
}

func corsMiddleware(corsOptions cors.Options) Middleware {
	// Create a new CORS handler with the specified options
	corsHandler := cors.New(corsOptions)
	return func(next http.Handler) http.Handler {
		// Wrap the next handler with the CORS handler
		return corsHandler.Handler(next)
	}
}

func ChainMiddleware(handler http.Handler, middlewares ...Middleware) http.Handler {
    for i := len(middlewares) - 1; i >= 0; i-- {
        handler = middlewares[i](handler)
    }
    return handler
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found, using default environment variables")
	}
	// Get the PORT environment variable
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT environment variable is not set")
	}



	// Convert the PORT string to an integer
	port, err := strconv.Atoi(portString)
	if err != nil {
		log.Fatalf("Invalid PORT value: %v", err)
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}
	connection, err := sql.Open("postgres", dbURL) // Open a connection to the PostgreSQL database using the provided DB_URL
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err) // Log an
	}
	defer connection.Close()

	queries := db.New(connection) // Initialize db.Queries with the connection

	config := &Config{
		DB: queries, // Set the DB field in the config to the queries instance
	}

	corsOptions := cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins (⚠️ for dev only)
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"}, // Allow all headers (e.g. Authorization)
		AllowCredentials: true,          // Allow cookies or Authorization headers
		MaxAge:           600,           // Cache preflight response for 10 minutes
		Debug:            true,          // Log CORS decisions to console (very useful)
	}

	router := createRouter() // Create a new router for handling HTTP requests
	router.HandleFunc("GET /readiness", handlerReadiness) // Register the readiness handler
	router.HandleFunc("GET /error", handlerError) // Register the error handler
	router.HandleFunc("POST /users", config.handlerCreateUser) // Register the user creation handler

	// Swagger UI endpoint
	routerWithMiddleware := ChainMiddleware(router, corsMiddleware(corsOptions)) // Apply CORS middleware to the router

	server := createServer(routerWithMiddleware, port) // Create a new server with the specified port

	// Log the server address and port
	fmt.Printf("Starting server on port %d\n", port)

	// Start the HTTP server
	if err := startServer(server); err != nil {
		log.Fatalf("Failed to start server: %v", err) // Log an error if the server fails to start
	}
}
