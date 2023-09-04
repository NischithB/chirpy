package chirps

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/NischithB/chirpy/utils"
	"github.com/go-chi/chi"
)

var chirpRouter chi.Router
var repository *DataRepository

func GetChirpRouter() *chi.Router {
	chirpRouter = chi.NewRouter()
	chirpRouter.Get("/", HandleGetChirps)
	chirpRouter.Post("/", HandleCreateChirp)

	return &chirpRouter
}

func InitRepository() {
	var err error
	repository, err = NewChirpRepository()
	if err != nil {
		log.Printf("Error initializing repository: %s", err)
		os.Exit(0)
	}
}

func HandleCreateChirp(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Body string `json:"body"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Printf("Error decoding request body: %s", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to decode request body")
		return
	}

	if len(body.Body) > 140 {
		utils.RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	cleanBody := cleanChirp(body.Body)

	chirp, err := repository.CreateChirp(cleanBody)
	if err != nil {
		log.Println("Error creating chirp")
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create chirp")
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, chirp)
}

func HandleGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := repository.GetChirps()
	if err != nil {
		log.Printf("Error fetching chirps: %s", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch chirps")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, chirps)
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
