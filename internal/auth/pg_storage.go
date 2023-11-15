package auth

import (
	"WB_TEST/internal/models"
	"context"
)

type Storage interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	CreateCustomer(ctx context.Context, customer *models.Customer) error
	CreateLoader(ctx context.Context, loader *models.Loader) error
	FindUserByLogin(ctx context.Context, user *models.User) (*models.User, error)
}
