package services

import (
	"github.com/NischithB/chirpy/config"
	"github.com/NischithB/chirpy/models"
	"github.com/NischithB/chirpy/utils"
	"golang.org/x/crypto/bcrypt"
)

func Login(email, password string, expiresIn int) (models.UserLoginResponse, error) {
	db := config.Config.DB
	data, err := db.Read()
	if err != nil {
		return models.UserLoginResponse{}, err
	}
	user, exists := getUserByEmail(&data.Users, email)
	if !exists {
		return models.UserLoginResponse{}, utils.ErrUserNotExists
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return models.UserLoginResponse{}, bcrypt.ErrMismatchedHashAndPassword
	}
	token, err := utils.CreateJwt(config.Config.JwtSecret, user.Id, expiresIn)
	if err != nil {
		return models.UserLoginResponse{}, err
	}
	return models.UserLoginResponse{Id: user.Id, Email: user.Email, Token: token}, nil
}
