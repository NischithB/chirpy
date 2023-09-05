package services

import (
	"github.com/NischithB/chirpy/config"
	"github.com/NischithB/chirpy/models"
)

func CreateUser(email string) (models.User, error) {
	db := config.Config.DB
	data, err := db.Read()
	if err != nil {
		return models.User{}, err
	}

	id := len(data.Users) + 1
	user := models.User{Id: id, Email: email}
	data.Users[id] = user

	if err := db.Write(data); err != nil {
		return models.User{}, err
	}
	return user, nil
}
