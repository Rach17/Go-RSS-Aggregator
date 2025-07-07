package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"log"
	"github.com/Rach17/Go-RSS-Aggregator/internal/db"
	"github.com/Rach17/Go-RSS-Aggregator/internal/auth"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}	

func handlerError(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}

func (config *Config) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}


	var params parameters = parameters{}
	var decoder *json.Decoder = json.NewDecoder(r.Body)

	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	hashedPassword, err := HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	_, err = config.DB.CreateUser(r.Context(), db.CreateUserParams{
		Username:     params.Username,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to create user: %v", err)
		respondWithError(w, http.StatusInternalServerError, errorMessage)
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "User created successfully"})
}

func (config *Config) handlerGetUserByAPIKey(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		log.Printf("Error getting API key: %v", err)
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	user, err := config.DB.GetUserByAPIKey(r.Context(), apiKey)
	if err != nil {
		log.Printf("Error retrieving user by API key: %v", err)
		errorMessage := fmt.Sprintf("Failed to get user: %v", err)
		respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}