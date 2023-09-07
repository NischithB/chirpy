package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/NischithB/chirpy/services"
	"github.com/NischithB/chirpy/utils"
	"github.com/go-chi/chi"
)

func getChirpRouter() chi.Router {
	chirpRouter := chi.NewRouter()
	chirpRouter.Get("/", handleGetChirps)
	chirpRouter.Post("/", handleCreateChirp)
	chirpRouter.Get("/{chirpID}", handleGetChirpById)
	chirpRouter.Delete("/{chirpID}", handleDeleteChirp)
	return chirpRouter
}

func handleGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := services.GetChirps()
	if err != nil {
		log.Printf("Error fetching chirps: %s", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch chirps")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, chirps)
}

func handleCreateChirp(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Body string `json:"body"`
	}{}

	id, err := services.AuthenticateUser(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Printf("Error decoding request body: %s", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to decode request body")
		return
	}

	if len(body.Body) > 140 {
		utils.RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	cleanBody := utils.CleanChirp(body.Body)

	chirp, err := services.CreateChirp(cleanBody, id)
	if err != nil {
		log.Println("Error creating chirp")
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create chirp")
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, chirp)
}

func handleGetChirpById(w http.ResponseWriter, r *http.Request) {
	chirpID, err := strconv.Atoi(chi.URLParam(r, "chirpID"))
	if err != nil {
		log.Printf("Error extracting chirpID from URL: %s", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to extract chirpID from URL")
		return
	}

	chirp, err := services.GetChirpById(chirpID)
	if errors.Is(err, utils.ErrNotFound) {
		utils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Chirp with id: %d doesn't exist", chirpID))
		return
	}
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch Chirp")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, chirp)
}

func handleDeleteChirp(w http.ResponseWriter, r *http.Request) {
	userId, err := services.AuthenticateUser(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	chirpID, err := strconv.Atoi(chi.URLParam(r, "chirpID"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Failed to get chirpId")
		return
	}

	err = services.DeleteChirp(chirpID, userId)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			utils.RespondWithError(w, http.StatusNotFound, "No chirp with the given id exists")
			return
		}
		if errors.Is(err, utils.ErrForbidden) {
			utils.RespondWithError(w, http.StatusForbidden, "Not authorized to delete this chirp")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete chirp")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "Successfully deleted chirp")
}
