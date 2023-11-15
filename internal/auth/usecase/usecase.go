package usecase

import (
	"WB_TEST/config"
	"WB_TEST/internal/auth"
	"WB_TEST/internal/models"
	"WB_TEST/pkg/utils"
	"context"
	"log/slog"
	"math/rand"
)

type authUC struct {
	cfg         *config.Config
	authStorage auth.Storage
	logger      *slog.Logger
}

func NewAuthUseCase(cfg *config.Config, storage auth.Storage, logger *slog.Logger) auth.UseCase {
	return &authUC{cfg: cfg, authStorage: storage, logger: logger}
}

func (uc *authUC) Register(ctx context.Context, user *models.User) (*models.User, error) {
	//контекст

	err := user.PrepareCreate()
	if err != nil {
		return nil, err
	}

	login, err := uc.authStorage.FindUserByLogin(ctx, user)
	if login != nil || err == nil {
		return nil, utils.ThisLoginAlreadyExist
	}

	userAfterCreate, err := uc.authStorage.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	switch userAfterCreate.Permissions {
	case "customer":
		minimum := uc.cfg.GameRules.CustomerCapitalMin
		maximum := uc.cfg.GameRules.CustomerCapitalMax

		capital := rand.Intn(maximum-minimum) + minimum

		customer := &models.Customer{
			User:    userAfterCreate,
			Capital: capital,
		}

		err = uc.authStorage.CreateCustomer(ctx, customer)
		if err != nil {
			return nil, err
		}

	case "loader":
		MinPortableWeight := uc.cfg.GameRules.LoaderMinPortableWeight
		MaxPortableWeight := uc.cfg.GameRules.LoaderMaxPortableWeight
		PortableWeight := rand.Intn(MaxPortableWeight-MinPortableWeight) + MinPortableWeight

		SalaryMin := uc.cfg.GameRules.LoaderSalaryMin
		SalaryMax := uc.cfg.GameRules.LoaderSalaryMax
		Salary := rand.Intn(SalaryMax-SalaryMin) + SalaryMin

		loader := &models.Loader{
			User:              userAfterCreate,
			MaxPortableWeight: PortableWeight,
			Fatigue:           rand.Intn(100),
			Salary:            Salary,
			Drunkenness:       rand.Intn(2),
		}

		err = uc.authStorage.CreateLoader(ctx, loader)
		if err != nil {
			return nil, err
		}
	}

	return userAfterCreate, nil
}

func (uc *authUC) Login(ctx context.Context, user *models.User) (*models.User, error) {
	// контекст

	foundUser, err := uc.authStorage.FindUserByLogin(ctx, user)
	if user == nil || err != nil {
		return nil, utils.UnexpectedLogin
	}

	if foundUser.Permissions != user.Permissions {
		return nil, utils.LoginDNOTHavePermissions
	}

	err = foundUser.CheckPassword(user.Password)
	if err != nil {
		return nil, utils.InvalidPassword
	}

	return foundUser, nil
}
