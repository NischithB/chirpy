package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")

	responseBody, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
	w.Write(responseBody)
}

func respondWithError(w http.ResponseWriter, statusCode int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, statusCode, errorResponse{Error: msg})
}
