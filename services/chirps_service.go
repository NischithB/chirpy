package services

import (
	"sort"

	"github.com/NischithB/chirpy/config"
	"github.com/NischithB/chirpy/models"
	"github.com/NischithB/chirpy/utils"
)

func CreateChirp(body string, userId int) (models.Chirp, error) {
	db := config.Config.DB
	data, err := db.Read()
	if err != nil {
		return models.Chirp{}, err
	}

	id := len(data.Chirps) + 1
	chirp := models.Chirp{Id: id, Body: body, AuthorId: userId}
	data.Chirps[id] = chirp

	if err := db.Write(data); err != nil {
		return models.Chirp{}, err
	}
	return chirp, nil
}

func GetChirps(filterByAuthorId int, sortOrder string) ([]models.Chirp, error) {
	db := config.Config.DB
	data, err := db.Read()
	if err != nil {
		return []models.Chirp{}, err
	}

	chirps := []models.Chirp{}
	for _, chirp := range data.Chirps {
		if filterByAuthorId == -1 || filterByAuthorId == chirp.AuthorId {
			chirps = append(chirps, chirp)
		}
	}

	var sortComparator func(int, int) bool
	if sortOrder == "asc" {
		sortComparator = func(i, j int) bool {
			return chirps[i].Id < chirps[j].Id
		}
	} else {
		sortComparator = func(i, j int) bool {
			return chirps[i].Id > chirps[j].Id
		}
	}
	sort.Slice(chirps, sortComparator)

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

func DeleteChirp(chirpID, userId int) error {
	db := config.Config.DB
	data, err := db.Read()
	if err != nil {
		return err
	}

	chirp, exists := data.Chirps[chirpID]
	if !exists {
		return utils.ErrNotFound
	}
	if chirp.AuthorId != userId {
		return utils.ErrForbidden
	}
	delete(data.Chirps, chirpID)
	if err := db.Write(data); err != nil {
		return err
	}
	return nil
}
