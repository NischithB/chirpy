package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/NischithB/chirpy/services"
	"github.com/NischithB/chirpy/utils"
	"github.com/go-chi/chi"
)

func getPaymentController() chi.Router {
	router := chi.NewRouter()
	router.Post("/webhooks", handlePolkaWebhooks)

	return router
}

func handlePolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}{}

	reqKey, err := utils.ExtractToken(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid key")
		return
	}
	key, exists := os.LookupEnv("POLKA_API_KEY")
	if !exists {
		utils.RespondWithError(w, http.StatusInternalServerError, "")
		return
	}
	if key != reqKey {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid key")
		return
	}

	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Error decoding request body")
		return
	}

	if body.Event != "user.upgraded" {
		utils.RespondWithJSON(w, http.StatusOK, "")
		return
	}

	err = services.UpdateMembership(body.Data.UserID, true)
	if err != nil {
		if errors.Is(err, utils.ErrUserNotExists) {
			utils.RespondWithError(w, http.StatusNotFound, "")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "")
}
