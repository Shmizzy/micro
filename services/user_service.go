package services

import (
	"errors"
	"service/models"
)

func RegisterUser(user models.User) error {
	return nil
}

func LoginUser(user models.Credentials) (string, error) {
	return "", errors.New("not implemented")
}
