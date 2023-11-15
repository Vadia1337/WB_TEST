package server

import (
	authHandlers "WB_TEST/internal/auth/http"
	authStorage "WB_TEST/internal/auth/storage"
	authUC "WB_TEST/internal/auth/usecase"
	gameHandlers "WB_TEST/internal/game/http"
	gameStorage "WB_TEST/internal/game/storage"
	gameUC "WB_TEST/internal/game/usecase"
	"WB_TEST/internal/middleware"
	sessionHandlers "WB_TEST/internal/session/http"
	sessionStorage "WB_TEST/internal/session/storage"
	sessionUC "WB_TEST/internal/session/usecase"
	"github.com/gin-gonic/gin"
)

func (s *Server) InitHandlers() *gin.Engine {
	router := gin.New()

	sStorage := sessionStorage.NewSessionStorage(s.db)
	sUC := sessionUC.NewSessionUC(s.cfg, sStorage, s.logger)
	sHandlers := sessionHandlers.NewSessionHandlers(s.cfg, sUC, s.logger)

	aStorage := authStorage.NewAuthStorage(s.db)
	aUC := authUC.NewAuthUseCase(s.cfg, aStorage, s.logger)
	aHandlers := authHandlers.NewAuthHandlers(s.cfg, aUC, sUC, s.logger)

	gStorage := gameStorage.NewGameStorage(s.db)
	gUC := gameUC.NewGameUseCase(s.cfg, gStorage, s.logger)
	gHandlers := gameHandlers.NewGameHandlers(s.cfg, gUC, s.logger)

	mw := middleware.NewMiddlewareChief(sUC, s.cfg, s.logger)

	router.Use(mw.RequestLogMiddleware)

	auth := router.Group("/auth")
	{
		auth.POST("/register", aHandlers.Register)
		auth.POST("/login", aHandlers.Login)
		auth.GET("/refresh", sHandlers.RefreshTokens)
	}

	router.GET("/tasks", gHandlers.CreateTasks)

	game := router.Group("/game", mw.AuthMiddleware)
	{
		game.POST("/tasks", gHandlers.UserTasks)
		game.POST("/me", gHandlers.UserInfo)
		game.POST("/start", gHandlers.StartGame)
	}

	return router
}
