package http

import (
	"WB_TEST/config"
	"WB_TEST/internal/game"
	"WB_TEST/internal/models"
	"WB_TEST/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type gameHandlers struct {
	cfg    *config.Config
	gameUC game.UseCase
	logger *slog.Logger
}

func NewGameHandlers(cfg *config.Config, uc game.UseCase, logger *slog.Logger) game.Handlers {
	return &gameHandlers{cfg: cfg, gameUC: uc, logger: logger}
}

func (gh *gameHandlers) UserTasks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second) // Можно заменить любым другим контекстом
	defer cancel()

	userID := c.Value("userID").(int)

	tasks, err := gh.gameUC.UserTasks(ctx, userID)
	if err != nil {
		gh.logger.Error("UserTasks Handler: ", err.Error())
		utils.SendJsonResponse(c, http.StatusServiceUnavailable, "Ошибка бизнес-логики")
		return
	}

	utils.SendJsonResponse(c, http.StatusOK, tasks)
}

func (gh *gameHandlers) CreateTasks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second) // Можно заменить любым другим контекстом
	defer cancel()

	id := c.Query("user_id")

	userId, err := utils.CreateTasksValidate(id)
	if err != nil {
		gh.logger.Error("CreateTasks Handler: ", err.Error())
		utils.SendJsonResponse(c, http.StatusBadRequest, "Your request contains incorrect data")
		return
	}

	err = gh.gameUC.CreateTasks(ctx, userId)
	if err == utils.UserIDNOTCustomer {
		gh.logger.Error("CreateTasks Handler: ", err.Error())
		utils.SendJsonResponse(c, http.StatusBadRequest, "Your request contains incorrect data")
		return
	}
	if err != nil {
		gh.logger.Error("CreateTasks Handler: ", err.Error())
		utils.SendJsonResponse(c, http.StatusServiceUnavailable, "Ошибка бизнес-логики")
		return
	}

	utils.SendJsonResponse(c, http.StatusOK, "Create tasks in count: "+strconv.Itoa(gh.cfg.CountJobsToGenerate))
}

func (gh *gameHandlers) UserInfo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second) // Можно заменить любым другим контекстом
	defer cancel()

	userID := c.Value("userID").(int)

	info, err := gh.gameUC.UserInfo(ctx, userID)
	if err != nil {
		gh.logger.Error("UserInfo Handler: ", err.Error())
		utils.SendJsonResponse(c, http.StatusServiceUnavailable, "Ошибка бизнес-логики")
		return
	}

	utils.SendJsonResponse(c, http.StatusOK, info)
}

func (gh *gameHandlers) StartGame(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second) // Можно заменить любым другим контекстом
	defer cancel()

	userID := c.Value("userID").(int)

	jobAndLoaders := &models.JobAndLoaders{}
	err := utils.ReadBody(c, jobAndLoaders)
	if err != nil {
		gh.logger.Error("StartGame Handler :", err.Error())
		utils.SendJsonResponse(c, http.StatusBadRequest, "Your request contains incorrect data")
		return
	}

	err = utils.StartGameValidate(jobAndLoaders)
	if err != nil {
		gh.logger.Error("StartGame Handler :", err.Error())
		utils.SendJsonResponse(c, http.StatusBadRequest, "Your request contains incorrect data")
		return
	}

	gameResult, err := gh.gameUC.StartGame(ctx, userID, jobAndLoaders)
	if err != nil {
		gh.logger.Error("StartGame Handler :", err.Error())
		utils.SendJsonResponse(c, http.StatusServiceUnavailable, "ошибка бизнес-логики")
		return
	}

	utils.SendJsonResponse(c, http.StatusOK, gameResult)

}
