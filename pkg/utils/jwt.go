package utils

import (
	"WB_TEST/config"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

type jwtClaims struct {
	ID string
	jwt.RegisteredClaims
}

func JWTGenerate(userID int, cfg *config.Config) (string, error) {

	claims := jwtClaims{
		strconv.Itoa(userID),
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.AccessTokenExpire * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(cfg.JWTSalt))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func JWTValidate(accessToken string, cfg *config.Config) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(cfg.JWTSalt), nil
	})
	if err != nil {
		return 0, JWTIsDied
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *jwtClaims")
	}

	userID, err := strconv.Atoi(claims.ID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
