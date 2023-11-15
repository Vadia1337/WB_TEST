package storage

import (
	"WB_TEST/internal/auth"
	"WB_TEST/internal/models"
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

type authStorage struct {
	db *sqlx.DB
}

func NewAuthStorage(db *sqlx.DB) auth.Storage {
	return &authStorage{db: db}
}

func (as *authStorage) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second) //любой другой контекст
	defer cancel()

	createUserQuery := "INSERT INTO users (login, password, permissions) VALUES ($1, $2, $3) RETURNING *"

	u := &models.User{}
	err := as.db.QueryRowxContext(
		ctx, createUserQuery, &user.Login, &user.Password, &user.Permissions,
	).StructScan(u)
	if err != nil {
		return nil, errors.Wrap(err, "authStorage.CreateUser")
	}

	return u, nil
}

func (as *authStorage) CreateCustomer(ctx context.Context, customer *models.Customer) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second) //любой другой контекст
	defer cancel()

	createCustomerQuery := "INSERT INTO customer (user_id, capital) VALUES ($1, $2)"

	err := as.db.QueryRowxContext(
		ctx, createCustomerQuery, &customer.User.ID, &customer.Capital,
	).Err()
	if err != nil {
		return err
	}

	return nil
}

func (as *authStorage) CreateLoader(ctx context.Context, loader *models.Loader) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second) //любой другой контекст
	defer cancel()

	createLoaderQuery := "INSERT INTO loader (user_id, max_portable_weight, fatigue, salary, drunkenness) " +
		"VALUES ($1, $2, $3, $4, $5)"

	err := as.db.QueryRowxContext(
		ctx, createLoaderQuery, &loader.User.ID, &loader.MaxPortableWeight,
		&loader.Fatigue, &loader.Salary, &loader.Drunkenness,
	).Err()
	if err != nil {
		return err
	}

	return nil
}

func (as *authStorage) FindUserByLogin(ctx context.Context, user *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second) //любой другой контекст
	defer cancel()

	findUserByLoginQuery := "SELECT * FROM users WHERE login = $1"

	u := &models.User{}
	err := as.db.QueryRowxContext(
		ctx, findUserByLoginQuery, &user.Login,
	).StructScan(u)
	if err != nil {
		return nil, err
	}

	return u, nil
}
