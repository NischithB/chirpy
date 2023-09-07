package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/NischithB/chirpy/services"
	"github.com/NischithB/chirpy/utils"
	"github.com/go-chi/chi"
	"golang.org/x/crypto/bcrypt"
)

func getAuthController() chi.Router {
	authRouter := chi.NewRouter()
	authRouter.Post("/login", handleLogin)

	return authRouter
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Expiry   int    `json:"expires_in_seconds"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Printf("Error decoding request body: %s", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Failed to decode request body")
		return
	}

	if body.Expiry == 0 {
		body.Expiry = 60 * 60 * 24
	}

	user, err := services.Login(body.Email, body.Password, body.Expiry)
	if errors.Is(err, utils.ErrUserNotExists) {
		utils.RespondWithError(w, http.StatusBadRequest, "No user exists with the given email")
		return
	}
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		utils.RespondWithError(w, http.StatusUnauthorized, "Failed to login, Invalid credentials")
		return
	}
	if err != nil {
		log.Printf("Error logging in: %s", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to login")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, user)
}
