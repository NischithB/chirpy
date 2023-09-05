package utils

import "strings"

func CleanChirp(chirp string) string {
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
