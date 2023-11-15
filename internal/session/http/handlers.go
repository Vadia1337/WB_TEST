package http

import (
	"WB_TEST/config"
	"WB_TEST/internal/session"
	"WB_TEST/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"time"
)

type sessionHandlers struct {
	cfg    *config.Config
	sessUC session.SessionUC
	logger *slog.Logger
}

func NewSessionHandlers(cfg *config.Config, sessUC session.SessionUC, logger *slog.Logger) session.Handlers {
	return &sessionHandlers{cfg: cfg, sessUC: sessUC, logger: logger}
}

func (sH *sessionHandlers) RefreshTokens(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second) // Можно заменить любым другим контекстом
	defer cancel()

	cookie, err := c.Cookie("refreshToken")
	if err != nil {
		sH.logger.Error("RefreshTokens Handler :", err.Error())
		utils.SendJsonResponse(c, http.StatusUnauthorized, "where is refresh token")
		return
	}

	userAndToken, err := sH.sessUC.RefreshToken(ctx, cookie)
	if err != nil {
		sH.logger.Error("RefreshTokens Handler :", err.Error())
		utils.SendJsonResponse(c, http.StatusServiceUnavailable, "Ошибка бизнес-логики") // так делать нельзя, но ..
		return
	}

	accessToken, err := utils.JWTGenerate(userAndToken.UserID, sH.cfg)
	if err != nil {
		sH.logger.Error("Login Handler :", err.Error())
		utils.SendJsonResponse(c, http.StatusServiceUnavailable, "Ошибка бизнес-логики") // так делать нельзя, но ..
		return
	}

	refreshToken, err := sH.sessUC.CreateSession(ctx, userAndToken.UserID)
	if err != nil {
		sH.logger.Error("Login Handler :", err.Error())
		utils.SendJsonResponse(c, http.StatusServiceUnavailable, "Ошибка бизнес-логики") // так делать нельзя, но ..
		return
	}

	c.SetCookie(utils.CreateSessionCookie(sH.cfg, refreshToken))

	utils.SendJsonResponse(c, http.StatusOK, accessToken)
}
