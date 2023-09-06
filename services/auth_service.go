package services

import (
	"github.com/NischithB/chirpy/config"
	"github.com/NischithB/chirpy/models"
	"github.com/NischithB/chirpy/utils"
	"golang.org/x/crypto/bcrypt"
)

func Login(email, password string) (models.UserInfo, error) {
	db := config.Config.DB
	data, err := db.Read()
	if err != nil {
		return models.UserInfo{}, err
	}
	user, exists := getUserByEmail(&data.Users, email)
	if !exists {
		return models.UserInfo{}, utils.ErrUserNotExists
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return models.UserInfo{}, bcrypt.ErrMismatchedHashAndPassword
	}
	return models.UserInfo{Id: user.Id, Email: user.Email}, nil
}
