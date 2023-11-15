package usecase

import (
	"WB_TEST/config"
	"WB_TEST/internal/game"
	"WB_TEST/internal/models"
	"WB_TEST/pkg/utils"
	"context"
	"log/slog"
	"math/rand"
)

type gameUC struct {
	cfg         *config.Config
	gameStorage game.Storage
	logger      *slog.Logger
}

func NewGameUseCase(cfg *config.Config, storage game.Storage, logger *slog.Logger) game.UseCase {
	return &gameUC{cfg: cfg, gameStorage: storage, logger: logger}
}

func (uc *gameUC) CreateTasks(ctx context.Context, userID int) error {

	// проверить user id на принадлежность к группе customer и вернуть customer
	customer, err := uc.gameStorage.GetCustomerByUserID(ctx, userID)
	if err != nil {
		return utils.UserIDNOTCustomer
	}

	weightMin := uc.cfg.CargoWeightMin
	weightMax := uc.cfg.CargoWeightMax

	tasksName := "Новое задание!"

	for i := uc.cfg.CountJobsToGenerate; i > 0; i-- {
		cargoWeight := rand.Intn(weightMax-weightMin) + weightMin

		err = uc.gameStorage.CreateTasks(ctx, customer.ID, cargoWeight, tasksName)
		if err != nil {
			return err
		}
	}

	return nil
}

func (uc *gameUC) UserTasks(ctx context.Context, userID int) (interface{}, error) {
	user, err := uc.gameStorage.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	switch user.Permissions {
	case "customer":
		customer, err := uc.gameStorage.GetCustomerByUserID(ctx, user.ID)
		if err != nil {
			return nil, utils.UserIDNOTCustomer
		}

		tasks, err := uc.gameStorage.GetCustomerTasks(ctx, customer.ID)
		if err != nil {
			return nil, err
		}

		return tasks, nil

	case "loader":
		loader, err := uc.gameStorage.GetLoaderByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}

		tasks, err := uc.gameStorage.GetLoaderDoneTasks(ctx, loader.ID)
		if err != nil {
			return nil, err
		}

		return tasks, nil

	default:
		return nil, nil
	}
}

func (uc *gameUC) UserInfo(ctx context.Context, userID int) (interface{}, error) {
	user, err := uc.gameStorage.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	switch user.Permissions {
	case "customer":
		customer, err := uc.gameStorage.GetCustomerByUserID(ctx, user.ID)
		if err != nil {
			return nil, utils.UserIDNOTCustomer
		}

		loaders, err := uc.gameStorage.GetAllLoaders(ctx)
		if err != nil {
			return nil, err
		}

		CustomerAndAllLoaders := models.CustomerAndAllLoaders{
			Customer: customer,
			Loaders:  loaders,
		}

		return CustomerAndAllLoaders, nil

	case "loader":
		loader, err := uc.gameStorage.GetLoaderByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}

		return loader, nil

	default:
		return nil, nil
	}
}

func (uc *gameUC) StartGame(ctx context.Context, userID int, jobsAndLoaders *models.JobAndLoaders) (string, error) {
	user, err := uc.gameStorage.GetCustomerByUserID(ctx, userID)
	if err != nil {
		return "", err
	}

	if user.Capital <= 0 {
		return "У вас закончился капитал, вы проиграли :-(", nil
	}

	job, err := uc.gameStorage.GetJobByID(ctx, jobsAndLoaders.JobId)
	if err != nil {
		return "", err
	}

	// проверить принадлежность кастомера к этому заданию

	sum := 0
	weight := 0.0

	userCapital := user.Capital
	for _, v := range jobsAndLoaders.Loaders {
		loader, err := uc.gameStorage.GetLoaderById(ctx, v)
		if err != nil {
			return "", err
		}

		sum += loader.Salary
		weight += loaderWeightFormula(float64(loader.MaxPortableWeight), loader.Drunkenness, float64(loader.Fatigue))

		if weight > 0 { // меняем усталось у тех грузчиков, у кого это возможно, чтобы не уперется в 100 усталости
			fatigue := loader.Fatigue + uc.cfg.LoaderFatigue
			if fatigue > 100 {
				fatigue = 100
			}

			_, err := uc.gameStorage.ChangeFatigueLoader(ctx, loader.ID, fatigue)
			if err != nil {
				return "", err
			}
		}

		//так как грузчика закрепили за работой (хоть и не смог т.к пьяный или уставший) укажем, что работу выполнял
		err = uc.gameStorage.CreateLoaderJob(ctx, loader.ID, job.ID)
		if err != nil {
			return "", err
		}

		userCapital = userCapital - loader.Salary

		userCapital, err = uc.gameStorage.PayCustomer(ctx, user.ID, userCapital)
		if err != nil {
			return "", err
		}

	}

	uc.logger.Info("", weight, sum, userCapital)

	if float64(job.CargoWeight) > weight {
		return "Вы проиграли! Вес задания был больше, чем грузоподъемность грузчиков", nil
	}

	if user.Capital < sum {
		return "Вы проиграли! Вы не смогли выплатить зарплату грузчикам", nil
	}

	return "Отлично, вы выполнили задание!", nil
}

func loaderWeightFormula(weight float64, drunkenness int, fatigue float64) float64 {

	if drunkenness == 1 {
		fatigue += 50.0
	}

	portableWeight := weight * ((100 - fatigue) / 100)

	if portableWeight < 0.0 {
		return 0
	}

	return portableWeight
}
