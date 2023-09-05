package services

import (
	"github.com/NischithB/chirpy/config"
	"github.com/NischithB/chirpy/models"
	"github.com/NischithB/chirpy/utils"
)

func CreateChirp(body string) (models.Chirp, error) {
	db := config.Config.DB
	data, err := db.Read()
	if err != nil {
		return models.Chirp{}, err
	}

	id := len(data.Chirps) + 1
	chirp := models.Chirp{Id: id, Body: body}
	data.Chirps[id] = chirp

	if err := db.Write(data); err != nil {
		return models.Chirp{}, err
	}
	return chirp, nil
}

func GetChirps() ([]models.Chirp, error) {
	db := config.Config.DB
	data, err := db.Read()
	if err != nil {
		return []models.Chirp{}, err
	}

	chirps := []models.Chirp{}
	for _, chirp := range data.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func GetChirpById(chirpID int) (models.Chirp, error) {
	db := config.Config.DB
	data, err := db.Read()
	if err != nil {
		return models.Chirp{}, err
	}

	chirp, exists := data.Chirps[chirpID]
	if !exists {
		return models.Chirp{}, utils.ErrNotFound
	}
	return chirp, nil
}
