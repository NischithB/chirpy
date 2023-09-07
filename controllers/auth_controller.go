package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/NischithB/chirpy/config"
	"github.com/NischithB/chirpy/services"
	"github.com/NischithB/chirpy/utils"
	"github.com/go-chi/chi"
	"golang.org/x/crypto/bcrypt"
)

func getAuthController() chi.Router {
	authRouter := chi.NewRouter()
	authRouter.Post("/login", handleLogin)
	authRouter.Post("/refresh", handleRefresh)
	authRouter.Post("/revoke", handleRevoke)

	return authRouter
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
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

	user, err := services.Login(body.Email, body.Password)
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

func handleRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	refToken, err := utils.ExtractToken(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Refresh token is missing")
		return
	}
	id, err := utils.ValidateJwt(refToken, config.Config.JwtSecret, false)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}
	revoked, err := services.IsTokenRevoked(refToken)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to refresh token")
		return
	}
	if revoked {
		utils.RespondWithError(w, http.StatusUnauthorized, "Inalid token")
		return
	}
	newtoken, err := utils.CreateJwt(config.Config.JwtSecret, id, true)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to refresh token")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, response{Token: newtoken})
}

func handleRevoke(w http.ResponseWriter, r *http.Request) {
	refToken, err := utils.ExtractToken(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Refresh token is missing")
		return
	}
	err = services.RevokeToken(refToken)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Failed to revoke token")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "Successfully revoked token")
}
