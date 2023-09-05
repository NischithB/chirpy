package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/NischithB/chirpy/services"
	"github.com/NischithB/chirpy/utils"
	"github.com/go-chi/chi"
)

func getUsersRouter() chi.Router {
	usersRouter := chi.NewRouter()
	usersRouter.Post("/", handleCreateUser)

	return usersRouter
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Email string `json:"email"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Printf("Error decoding request body: %s", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to decode request body")
		return
	}

	user, err := services.CreateUser(body.Email)
	if err != nil {
		log.Println("Error creating user")
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create chirp")
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, user)
}
