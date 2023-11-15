package storage

import (
	"WB_TEST/internal/game"
	"WB_TEST/internal/models"
	"context"
	"github.com/jmoiron/sqlx"
)

type gameStorage struct {
	db *sqlx.DB
}

func NewGameStorage(db *sqlx.DB) game.Storage {
	return &gameStorage{db: db}
}

func (gs *gameStorage) CreateTasks(ctx context.Context, customerID int, weight int, name string) error {
	createTaskQuery := "INSERT INTO jobs (customer_id, name, cargo_weight) VALUES ($1, $2, $3)"

	err := gs.db.QueryRowxContext(
		ctx, createTaskQuery, &customerID, &name, &weight,
	).Err()
	if err != nil {
		return err
	}

	return nil
}

func (gs *gameStorage) GetCustomerByUserID(ctx context.Context, userID int) (*models.Customer, error) {
	getCustomerIDByUserIDQuery := "SELECT id, capital FROM customer WHERE user_id = $1"

	customer := &models.Customer{}
	err := gs.db.QueryRowxContext(
		ctx, getCustomerIDByUserIDQuery, &userID,
	).Scan(&customer.ID, &customer.Capital)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (gs *gameStorage) GetLoaderByUserID(ctx context.Context, userID int) (*models.Loader, error) {
	getLoaderByUserIDQuery := "SELECT id, max_portable_weight, fatigue, salary, drunkenness FROM loader" +
		" WHERE user_id = $1"

	loader := &models.Loader{}
	err := gs.db.QueryRowxContext(
		ctx, getLoaderByUserIDQuery, &userID,
	).Scan(&loader.ID, &loader.MaxPortableWeight, &loader.Fatigue, &loader.Salary, &loader.Drunkenness)
	if err != nil {
		return nil, err
	}

	return loader, nil
}

func (gs *gameStorage) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	getUserQuery := "SELECT id, login, permissions FROM users WHERE id = $1"

	user := &models.User{}
	err := gs.db.QueryRowxContext(
		ctx, getUserQuery, &userID,
	).StructScan(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (gs *gameStorage) GetCustomerTasks(ctx context.Context, customerID int) ([]models.Job, error) {
	getCustomerTasksQuery := "SELECT * FROM jobs WHERE customer_id = $1"

	var Jobs []models.Job

	err := gs.db.SelectContext(ctx, &Jobs, getCustomerTasksQuery, &customerID)
	if err != nil {
		return nil, err
	}

	return Jobs, nil
}

func (gs *gameStorage) GetLoaderDoneTasks(ctx context.Context, loaderID int) ([]models.Job, error) {
	getCustomerTasksQuery := "SELECT id, customer_id, name, cargo_weight " +
		"FROM loader_jobs LEFT JOIN jobs j ON j.id = job_id " +
		"WHERE loader_id = $1"

	var Jobs []models.Job

	err := gs.db.SelectContext(ctx, &Jobs, getCustomerTasksQuery, &loaderID)
	if err != nil {
		return nil, err
	}

	return Jobs, nil
}

func (gs *gameStorage) GetAllLoaders(ctx context.Context) ([]models.Loader, error) {
	getCustomerTasksQuery := "SELECT id, max_portable_weight, fatigue, salary, drunkenness FROM loader"

	var Jobs []models.Loader

	err := gs.db.SelectContext(ctx, &Jobs, getCustomerTasksQuery)
	if err != nil {
		return nil, err
	}

	return Jobs, nil
}

func (gs *gameStorage) GetJobByID(ctx context.Context, jobId int) (*models.Job, error) {
	getJobByIDQuery := "SELECT * FROM jobs WHERE id = $1"

	foundJob := &models.Job{}
	err := gs.db.QueryRowxContext(
		ctx, getJobByIDQuery, &jobId,
	).StructScan(foundJob)
	if err != nil {
		return nil, err
	}

	return foundJob, err
}

func (gs *gameStorage) GetLoaderById(ctx context.Context, loaderID int) (*models.Loader, error) {
	getLoaderByIDQuery := "SELECT id, max_portable_weight, fatigue, salary, drunkenness FROM loader WHERE id = $1"

	foundLoader := &models.Loader{}
	err := gs.db.QueryRowxContext(
		ctx, getLoaderByIDQuery, &loaderID,
	).StructScan(foundLoader)
	if err != nil {
		return nil, err
	}

	return foundLoader, nil
}

func (gs *gameStorage) ChangeFatigueLoader(ctx context.Context, loaderID int, fatigue int) (*models.Loader, error) {
	changeFatigueLoaderQuery := "UPDATE loader SET fatigue = $1 WHERE id = $2" +
		" RETURNING id, max_portable_weight, fatigue, salary, drunkenness"

	updatedLoader := &models.Loader{}
	err := gs.db.QueryRowxContext(
		ctx, changeFatigueLoaderQuery, &fatigue, &loaderID,
	).StructScan(updatedLoader)
	if err != nil {
		return nil, err
	}

	return updatedLoader, nil
}

func (gs *gameStorage) CreateLoaderJob(ctx context.Context, loaderID int, jobID int) error {
	createLoaderJobQuery := "INSERT INTO loader_jobs (loader_id, job_id) VALUES ($1, $2)"

	err := gs.db.QueryRowxContext(
		ctx, createLoaderJobQuery, &loaderID, &jobID,
	).Err()
	if err != nil {
		return err
	}

	return nil
}

func (gs *gameStorage) PayCustomer(ctx context.Context, customerID int, sum int) (int, error) {
	payCustomerQuery := "UPDATE customer SET capital = $1 WHERE id = $2 RETURNING capital"

	var capital int
	err := gs.db.QueryRowxContext(
		ctx, payCustomerQuery, &sum, &customerID,
	).Scan(&capital)
	if err != nil {
		return 0, err
	}

	return capital, nil
}
