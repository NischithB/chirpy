package services

import (
	"net/http"
	"time"

	"github.com/NischithB/chirpy/config"
	"github.com/NischithB/chirpy/models"
	"github.com/NischithB/chirpy/utils"
	"golang.org/x/crypto/bcrypt"
)

func Login(email, password string) (models.UserLoginResponse, error) {
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
	access, err := utils.CreateJwt(config.Config.JwtSecret, user.Id, true)
	if err != nil {
		return models.UserLoginResponse{}, err
	}
	refresh, err := utils.CreateJwt(config.Config.JwtSecret, user.Id, false)
	if err != nil {
		return models.UserLoginResponse{}, err
	}
	return models.UserLoginResponse{
		UserInfo: models.UserInfo{
			Id:       user.Id,
			Email:    user.Email,
			IsMember: user.IsMember,
		},
		Token:   access,
		Refresh: refresh,
	}, nil
}

func RevokeToken(token string) error {
	db := config.Config.DB
	data, err := db.Read()
	if err != nil {
		return err
	}
	_, err = utils.ValidateJwt(token, config.Config.JwtSecret, false)
	if err != nil {
		return err
	}
	data.RevokedTokens[token] = time.Now().UTC().String()
	if err := db.Write(data); err != nil {
		return err
	}
	return nil
}

func IsTokenRevoked(token string) (bool, error) {
	db := config.Config.DB
	data, err := db.Read()
	if err != nil {
		return false, err
	}

	if _, revoked := data.RevokedTokens[token]; revoked {
		return true, nil
	}
	return false, nil
}

func AuthenticateUser(r *http.Request) (int, error) {
	token, err := utils.ExtractToken(r)
	if err != nil {
		return -1, err
	}

	id, err := utils.ValidateJwt(token, config.Config.JwtSecret, true)
	if err != nil {
		return -1, err
	}
	return id, nil
}
