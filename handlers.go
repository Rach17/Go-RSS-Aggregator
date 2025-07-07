package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/Rach17/Go-RSS-Aggregator/internal/db"
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