package middleware

import (
	"WB_TEST/config"
	"WB_TEST/internal/session"
	"log/slog"
)

type MiddlewareChief struct {
	sessUC session.SessionUC
	cfg    *config.Config
	logger *slog.Logger
}

func NewMiddlewareChief(sessUC session.SessionUC, cfg *config.Config, logger *slog.Logger) MiddlewareChief {
	return MiddlewareChief{sessUC: sessUC, cfg: cfg, logger: logger}
}
