package services

import (
	"github.com/NischithB/chirpy/config"
	"github.com/NischithB/chirpy/models"
	"github.com/NischithB/chirpy/utils"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(email, password string) (models.UserInfo, error) {
	db := config.Config.DB
	data, err := db.Read()
	if err != nil {
		return models.UserInfo{}, err
	}

	if _, exists := getUserByEmail(&data.Users, email); exists {
		return models.UserInfo{}, utils.ErrUserExists
	}

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.UserInfo{}, err
	}
	id := len(data.Users) + 1
	user := models.User{
		UserInfo: models.UserInfo{Id: id, Email: email},
		Password: string(pwdHash),
	}
	data.Users[id] = user

	if err := db.Write(data); err != nil {
		return models.UserInfo{}, err
	}
	return models.UserInfo{Id: user.Id, Email: user.Email}, nil
}

func getUserByEmail(users *map[int]models.User, email string) (models.User, bool) {
	for _, user := range *users {
		if user.Email == email {
			return user, true
		}
	}
	return models.User{}, false
}

func UpdateUser(id int, email, password string) (models.UserInfo, error) {
	db := config.Config.DB
	data, err := db.Read()
	if err != nil {
		return models.UserInfo{}, err
	}

	user, exists := data.Users[id]
	if !exists {
		return models.UserInfo{}, utils.ErrUserNotExists
	}
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.UserInfo{}, err
	}
	user.Email = email
	user.Password = string(pwdHash)
	data.Users[id] = user

	if err := db.Write(data); err != nil {
		return models.UserInfo{}, nil
	}
	return models.UserInfo{Id: user.Id, Email: user.Email}, nil
}

func UpdateMembership(id int, isMember bool) error {
	db := config.Config.DB
	data, err := db.Read()
	if err != nil {
		return err
	}

	user, exists := data.Users[id]
	if !exists {
		return utils.ErrUserNotExists
	}

	user.IsMember = isMember
	data.Users[id] = user
	if err := db.Write(data); err != nil {
		return err
	}
	return nil
}
