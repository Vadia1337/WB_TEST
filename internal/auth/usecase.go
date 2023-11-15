package auth

import (
	"WB_TEST/internal/models"
	"context"
)

type UseCase interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)
	Login(ctx context.Context, user *models.User) (*models.User, error)
}
