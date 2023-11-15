package server

import (
	"WB_TEST/config"
	"context"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	httpServer *http.Server
	cfg        *config.Config
	db         *sqlx.DB
	logger     *slog.Logger
}

func NewServer(cfg *config.Config, db *sqlx.DB, logger *slog.Logger) *Server {
	return &Server{cfg: cfg, db: db, logger: logger}
}

func (s *Server) Run() error {

	handlers := s.InitHandlers()

	s.httpServer = &http.Server{
		Addr:    s.cfg.HttpServerPort,
		Handler: handlers,
	}

	go func() {
		err := s.httpServer.ListenAndServe()
		if err == http.ErrServerClosed {
			s.logger.Info(err.Error())
			os.Exit(0)
		}
		if err != nil {
			s.logger.Error(err.Error())
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	return s.Shutdown(context.Background())
}

func (s *Server) Shutdown(ctx context.Context) error {
	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
