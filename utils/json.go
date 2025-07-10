package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	// Encode the data as JSON and write it to the response
	if jsonData, err := json.Marshal(data); err != nil {
		log.Printf("Error encoding JSON response: %v", data)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.Write(jsonData)
	}
	// Set the content type to application/json
	w.Header().Add("Content-Type", "application/json")
	// Set the HTTP status code
	w.WriteHeader(status)
}

func RespondWithError(w http.ResponseWriter, status int, message string) {
	if status > 499 || status < 600 {
		log.Printf("Server error: %d - %s", status, message)
	} else {
		log.Printf("Client error: %d - %s", status, message)
	}
	type errorResponse struct {
		Message string `json:"error"`
	}
	RespondWithJSON(w, status, errorResponse{Message: message})
}