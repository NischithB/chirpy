package controllers

import (
	"encoding/json"
	"errors"
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
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Printf("Error decoding request body: %s", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Failed to decode request body")
		return
	}

	user, err := services.CreateUser(body.Email, body.Password)
	if errors.Is(err, utils.ErrUserExists) {
		utils.RespondWithError(w, http.StatusBadRequest, "User with given email already exists")
		return
	}
	if err != nil {
		log.Println("Error creating user")
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, user)
}
