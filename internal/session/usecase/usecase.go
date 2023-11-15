package usecase

import (
	"WB_TEST/config"
	"WB_TEST/internal/models"
	"WB_TEST/internal/session"
	"context"
	"log/slog"
	"math/rand"
)

type sessionUC struct {
	cfg            *config.Config
	sessionStorage session.Storage
	logger         *slog.Logger
}

func NewSessionUC(cfg *config.Config, storage session.Storage, logger *slog.Logger) session.SessionUC {
	return &sessionUC{cfg: cfg, sessionStorage: storage, logger: logger}
}

func (uc *sessionUC) CreateSession(ctx context.Context, userID int) (string, error) {
	newRefreshToken, err := uc.genRefreshToken()
	if err != nil {
		return "", err
	}

	createSession, err := uc.sessionStorage.CreateSession(ctx, userID, newRefreshToken, uc.cfg.RefreshTokenExpire)
	if err != nil {
		return "", err
	}

	return createSession, nil
}

func (uc *sessionUC) genRefreshToken() (string, error) {

	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, 64)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b), nil
}

func (uc *sessionUC) RefreshToken(ctx context.Context, token string) (*models.UserIDAndRefreshToken, error) {
	rToken, err := uc.sessionStorage.FindRToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return rToken, err
}
