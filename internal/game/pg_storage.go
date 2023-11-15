package game

import (
	"WB_TEST/internal/models"
	"context"
)

type Storage interface {
	CreateTasks(ctx context.Context, userID int, weight int, name string) error
	GetCustomerByUserID(ctx context.Context, userID int) (*models.Customer, error)
	GetUserByID(ctx context.Context, userID int) (*models.User, error)
	GetCustomerTasks(ctx context.Context, customerID int) ([]models.Job, error)
	GetLoaderByUserID(ctx context.Context, userID int) (*models.Loader, error)
	GetLoaderDoneTasks(ctx context.Context, loaderID int) ([]models.Job, error)
	GetAllLoaders(ctx context.Context) ([]models.Loader, error)
	GetJobByID(ctx context.Context, jobId int) (*models.Job, error)
	GetLoaderById(ctx context.Context, loaderID int) (*models.Loader, error)
	ChangeFatigueLoader(ctx context.Context, loaderID int, fatigue int) (*models.Loader, error)
	CreateLoaderJob(ctx context.Context, loaderID int, jobID int) error
	PayCustomer(ctx context.Context, customerID int, sum int) (int, error)
}
