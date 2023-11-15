package session

import (
	"WB_TEST/internal/models"
	"context"
)

type Storage interface {
	CreateSession(ctx context.Context, userID int, refreshToken string, expire int) (string, error)
	FindRToken(ctx context.Context, token string) (*models.UserIDAndRefreshToken, error)
}
