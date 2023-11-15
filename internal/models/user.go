package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          int    `json:"id,omitempty"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	Permissions string `json:"permissions"`
}

type UserAndAccessToken struct {
	*User
	Token string
}

type UserIDAndRefreshToken struct {
	UserID       int    `db:"user_id"`
	RefreshToken string `db:"refreshtoken"`
}

func (u *User) PrepareCreate() error {
	err := u.HashPassword()
	if err != nil {
		return err
	}
	// можно выполнить любые преобразования перед созданием юзера
	return nil
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	return nil
}

func (u *User) CheckPassword(pass string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pass)); err != nil {
		return err
	}
	return nil
}
