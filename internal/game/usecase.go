package game

import (
	"WB_TEST/internal/models"
	"context"
)

type UseCase interface {
	CreateTasks(ctx context.Context, userID int) error
	UserTasks(ctx context.Context, userID int) (interface{}, error)
	UserInfo(ctx context.Context, userID int) (interface{}, error)
	StartGame(ctx context.Context, userID int, jobsAndLoaders *models.JobAndLoaders) (string, error)
}
