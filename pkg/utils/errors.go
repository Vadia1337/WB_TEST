package utils

import "errors"

var (
	ThisLoginAlreadyExist    = errors.New("this login already exist in system")
	UnexpectedLogin          = errors.New("unexpected this login in system")
	InvalidPassword          = errors.New("this password is invalid")
	LoginDNOTHavePermissions = errors.New("this login does not have such rights")
	UserIDNOTCustomer        = errors.New("this user id does not have such rights")

	JWTIsDied = errors.New("jwt is died his time is done")
)
