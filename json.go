package main

import (
	"encoding/json" // Importing encoding/json for JSON encoding and decoding
	"log"      // Importing log for logging errors and messages
	"net/http" // Importing net/http for HTTP server functionality
)

func respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	// Encode the data as JSON and write it to the response
	if jsonData, err := json.Marshal(data); err != nil {
		log.Println("Error encoding JSON response: %v", data)
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

func respondWithError(w http.ResponseWriter, status int, message string) {
	if status > 499 || status < 600 {
		log.Printf("Server error: %d - %s", status, message)
	} else {
		log.Printf("Client error: %d - %s", status, message)
	}

	type errorResponse struct {
		Message string `json:"error"`
	}

	respondWithJSON(w, status, errorResponse{Message: message})
}