package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Rach17/Go-RSS-Aggregator/service"
	"github.com/Rach17/Go-RSS-Aggregator/utils"
	"github.com/rs/cors"
)

// Middleware type
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Chain multiple middleware together
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	corsOptions := cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins (⚠️ for dev only)
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"}, // Allow all headers (e.g. Authorization)
		AllowCredentials: true,          // Allow cookies or Authorization headers
		MaxAge:           600,           // Cache preflight response for 10 minutes
		Debug:            true,          // Log CORS decisions to console (very useful)
	}
	// Create a new CORS handler with the specified options
	corsHandler := cors.New(corsOptions)
	return func(w http.ResponseWriter, r *http.Request) {
		// Wrap the next handler with the CORS handler
		corsHandler.ServeHTTP(w, r, next)
	}

}

type AuthMiddleware struct {
	authService *service.AuthService
}

func NewAuthMiddleware(authService *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

func (am *AuthMiddleware) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	// Middleware to verify API key in the request header
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := utils.GetAPIKey(r.Header) // Get the API key from the request header
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, err.Error()) // Respond with an error if API key is missing or invalid
			return
		}

		user, err := am.authService.IsAuth(r.Context(), apiKey) // Get the user associated with the API key
		if err != nil {
			errorMessage := fmt.Sprintf("Failed to get user: %v", err)     // Log the error if user retrieval fails
			utils.RespondWithError(w, http.StatusBadRequest, errorMessage) // Respond with a bad request error
			return
		}
		const userContextKey contextKey = "user"
		// Store the user in the request context for further processing
		ctx := context.WithValue(r.Context(), userContextKey, user) // Store the user in the request context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
