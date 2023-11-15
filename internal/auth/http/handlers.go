package http

import (
	"WB_TEST/config"
	"WB_TEST/internal/auth"
	"WB_TEST/internal/models"
	"WB_TEST/internal/session"
	"WB_TEST/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"time"
)

type authHandlers struct {
	cfg       *config.Config
	authUC    auth.UseCase
	sessionUC session.SessionUC
	logger    *slog.Logger
}

func NewAuthHandlers(cfg *config.Config, authUC auth.UseCase, sUC session.SessionUC, logger *slog.Logger) auth.Handlers {
	return &authHandlers{cfg: cfg, authUC: authUC, sessionUC: sUC, logger: logger}
}

func (h *authHandlers) Register(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second) // Можно заменить любым другим контекстом
	defer cancel()

	user := &models.User{}

	err := utils.ReadBody(c, user)
	if err != nil {
		h.logger.Error("Register Handler :", err.Error())
		utils.SendJsonResponse(c, http.StatusBadRequest, "Your request contains incorrect data")
		return
	}

	err = utils.UserRegisterValidate(user)
	if err != nil {
		h.logger.Error("Register Handler :", err.Error())
		utils.SendJsonResponse(c, http.StatusBadRequest, "Your request contains incorrect data")
		return
	}

	registeredUser, err := h.authUC.Register(ctx, user) // нужно систематизировать работу над ошибками
	if err == utils.ThisLoginAlreadyExist {
		h.logger.Error("Register Handler :", err.Error())
		utils.SendJsonResponse(c, http.StatusBadRequest, "Your request contains incorrect data")
		return
	}
	if err != nil {
		h.logger.Error("Register Handler :", err.Error())
		utils.SendJsonResponse(c, http.StatusServiceUnavailable, "Ошибка бизнес-логики") // так делать нельзя, но ..
		return
	}

	accessToken, err := utils.JWTGenerate(registeredUser.ID, h.cfg)
	if err != nil {
		h.logger.Error("Register Handler :", err.Error())
		utils.SendJsonResponse(c, http.StatusServiceUnavailable, "Ошибка бизнес-логики") // так делать нельзя, но ..
		return
	}

	refreshToken, err := h.sessionUC.CreateSession(ctx, registeredUser.ID)
	if err != nil {
		h.logger.Error("Register Handler :", err.Error())
		utils.SendJsonResponse(c, http.StatusServiceUnavailable, "Ошибка бизнес-логики") // так делать нельзя, но ..
		return
	}

	c.SetCookie(utils.CreateSessionCookie(h.cfg, refreshToken))

	userAndAccessToken := models.UserAndAccessToken{
		User:  registeredUser,
		Token: accessToken,
	}

	utils.SendJsonResponse(c, http.StatusOK, userAndAccessToken)
}

func (h *authHandlers) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second) // Можно заменить любым другим контекстом
	defer cancel()

	user := &models.User{}

	err := utils.ReadBody(c, user)
	if err != nil {
		h.logger.Error("Login Handler :", err.Error())
		utils.SendJsonResponse(c, http.StatusBadRequest, "Your request contains incorrect data")
		return
	}

	err = utils.UserLoginValidate(user)
	if err != nil {
		h.logger.Error("Login Handler :", err.Error())
		utils.SendJsonResponse(c, http.StatusBadRequest, "Your request contains incorrect data")
		return
	}

	userAfterLogin, err := h.authUC.Login(ctx, user)
	if err == utils.UnexpectedLogin || err == utils.InvalidPassword || err == utils.LoginDNOTHavePermissions {
		h.logger.Error("Login Handler :", err.Error())
		utils.SendJsonResponse(c, http.StatusBadRequest, "Your request contains incorrect data")
		return
	}
	if err != nil {
		h.logger.Error("Login Handler :", err.Error())
		utils.SendJsonResponse(c, http.StatusServiceUnavailable, "Ошибка бизнес-логики") // так делать нельзя, но ..
		return
	}

	accessToken, err := utils.JWTGenerate(userAfterLogin.ID, h.cfg)
	if err != nil {
		h.logger.Error("Login Handler :", err.Error())
		utils.SendJsonResponse(c, http.StatusServiceUnavailable, "Ошибка бизнес-логики") // так делать нельзя, но ..
		return
	}

	refreshToken, err := h.sessionUC.CreateSession(ctx, userAfterLogin.ID)
	if err != nil {
		h.logger.Error("Login Handler :", err.Error())
		utils.SendJsonResponse(c, http.StatusServiceUnavailable, "Ошибка бизнес-логики") // так делать нельзя, но ..
		return
	}

	c.SetCookie(utils.CreateSessionCookie(h.cfg, refreshToken))

	userAndAccessToken := models.UserAndAccessToken{
		User:  userAfterLogin,
		Token: accessToken,
	}

	utils.SendJsonResponse(c, http.StatusOK, userAndAccessToken)
}
