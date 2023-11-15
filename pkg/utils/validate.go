package utils

import (
	"WB_TEST/internal/models"
	"errors"
	"strconv"
)

func UserRegisterValidate(user *models.User) error {
	if user.Login == "" {
		return errors.New("invalid login")
	}

	if user.Password == "" {
		return errors.New("invalid password")
	}

	if len(user.Password) < 8 {
		return errors.New("password len < 8")
	}

	if user.Permissions != "loader" && user.Permissions != "customer" {
		return errors.New("invalid permissions")
	}

	return nil
}

func UserLoginValidate(user *models.User) error {
	if user.Login == "" {
		return errors.New("invalid login")
	}

	if user.Password == "" {
		return errors.New("invalid password")
	}

	if len(user.Password) < 8 {
		return errors.New("password len < 8")
	}

	if user.Permissions != "loader" && user.Permissions != "customer" {
		return errors.New("invalid permissions")
	}

	return nil
}

func CreateTasksValidate(id string) (int, error) {
	idToInt, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}

	if idToInt <= 0 {
		return 0, errors.New("invalid user id")
	}

	return idToInt, nil
}

func StartGameValidate(jobLoaders *models.JobAndLoaders) error {
	if jobLoaders.JobId <= 0 {
		return errors.New("invalid job id")
	}

	if len(jobLoaders.Loaders) <= 0 {
		return errors.New("where is loaders")
	}

	return nil
}
