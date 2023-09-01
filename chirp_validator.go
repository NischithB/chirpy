package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type ReqBody struct {
		Body string `json:"body"`
	}

	type ResBody struct {
		CleanedBody string `json:"cleaned_body"`
	}
	body := ReqBody{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Printf("Error decoding request body: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to decode request body")
		return
	}

	if len(body.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanBody := cleanChirp(body.Body)
	respondWithJSON(w, http.StatusOK, ResBody{CleanedBody: cleanBody})
}

func cleanChirp(chirp string) string {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	substitute := "****"

	tokens := strings.Split(chirp, " ")
	for _, word := range profaneWords {
		for index, token := range tokens {
			if strings.ToLower(token) == word {
				tokens[index] = substitute
			}
		}
	}

	return strings.Join(tokens, " ")
}
