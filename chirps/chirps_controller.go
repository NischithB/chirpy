package chirps

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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
	chirpRouter.Get("/{chirpID}", HandleGetChirpById)

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

func HandleGetChirpById(w http.ResponseWriter, r *http.Request) {
	chirpID, err := strconv.Atoi(chi.URLParam(r, "chirpID"))
	if err != nil {
		log.Printf("Error extracting chirpID from URL: %s", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to extract chirpID from URL")
		return
	}

	chirp, err := repository.GetChirpById(chirpID)
	if errors.Is(err, ErrNotFound) {
		utils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Chirp with id: %d doesn't exist", chirpID))
		return
	}
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch Chirp")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, chirp)
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
