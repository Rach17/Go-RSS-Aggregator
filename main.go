package main

import (
	"fmt"      // Importing fmt for formatted I/O
	"log"      // Importing log for logging errors and messages
	"net/http" // Importing net/http for HTTP server functionality
	"os"       // Importing os for environment variable access
	"strconv"  // Importing strconv for string conversion

	"github.com/joho/godotenv" // Importing godotenv to load environment variables from .env file
	"github.com/rs/cors"       // Importing rs/cors for handling CORS (Cross-Origin Resource Sharing) in HTTP requests
)

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

func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
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
		log.Fatal("Invalid PORT value: %v", err)
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
	routerWithMiddleware := Chain(router, corsMiddleware(corsOptions)) // Apply CORS middleware to the router

	server := createServer(routerWithMiddleware, port) // Create a new server with the specified port

	// Log the server address and port
	fmt.Printf("Starting server on port %d\n", port)

	// Start the HTTP server
	if err := startServer(server); err != nil {
		log.Fatalf("Failed to start server: %v", err) // Log an error if the server fails to start
	}
}
