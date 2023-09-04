package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, statusCode int, payload any) {
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

func RespondWithError(w http.ResponseWriter, statusCode int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	RespondWithJSON(w, statusCode, errorResponse{Error: msg})
}
