package services

import (
	"errors"

	"github.com/Rzero6/self-checkout-api/internal/repositories"
	"github.com/Rzero6/self-checkout-api/internal/utils"
)

func Login(username, password string) (string, error) {
	user, err := repositories.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	if !utils.CheckPassword(user.Password, password) {
		return "", errors.New("invalid password")
	}

	token, err := utils.GenerateJWT(int64(user.ID), user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}
