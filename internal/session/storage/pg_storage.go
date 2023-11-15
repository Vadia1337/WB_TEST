package storage

import (
	"WB_TEST/internal/models"
	"WB_TEST/internal/session"
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

type sessionStorage struct {
	db *sqlx.DB
}

func NewSessionStorage(db *sqlx.DB) session.Storage {
	return &sessionStorage{db: db}
}

func (ss *sessionStorage) CreateSession(ctx context.Context, userID int, refreshToken string, expire int) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second) //любой другой контекст
	defer cancel()

	deleteUserSessionsQuery := "DELETE FROM sessions WHERE user_id = $1;"

	err := ss.db.QueryRowxContext(
		ctx, deleteUserSessionsQuery, &userID,
	).Err()
	if err != nil {
		return "", errors.Wrap(err, "sessionStorage.CreateSession")
	}

	tokenExpire := time.Now().Add(time.Second * time.Duration(expire))

	createSessionQuery := "INSERT INTO sessions (user_id, refreshtoken, expires) VALUES ($1, $2, $3) " +
		" RETURNING refreshtoken"

	var token string
	err = ss.db.QueryRowxContext(
		ctx, createSessionQuery, &userID, &refreshToken, &tokenExpire,
	).Scan(&token)
	if err != nil {
		return "", errors.Wrap(err, "sessionStorage.CreateSession")
	}

	return token, nil
}

func (ss *sessionStorage) FindRToken(ctx context.Context, token string) (*models.UserIDAndRefreshToken, error) {
	findTokenQuery := "SELECT user_id, refreshtoken FROM sessions WHERE refreshtoken = $1"

	userIdAndRefreshToken := &models.UserIDAndRefreshToken{}
	err := ss.db.QueryRowxContext(
		ctx, findTokenQuery, &token,
	).StructScan(userIdAndRefreshToken)
	if err != nil {
		return nil, errors.Wrap(err, "sessionStorage.CreateSession")
	}

	return userIdAndRefreshToken, nil
}
