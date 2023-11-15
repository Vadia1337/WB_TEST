package session

import (
	"WB_TEST/internal/models"
	"context"
)

type SessionUC interface {
	CreateSession(ctx context.Context, userID int) (string, error)
	RefreshToken(ctx context.Context, token string) (*models.UserIDAndRefreshToken, error)
}
