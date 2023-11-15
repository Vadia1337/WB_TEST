package main

import (
	"WB_TEST/config"
	"WB_TEST/internal/server"
	"WB_TEST/pkg/logger"
	"WB_TEST/pkg/storage/postgres"
	"os"
)

func main() {
	// init config
	cfg := config.GetConfig()

	//init default logger
	gameLogger := logger.NewLogger()
	pgDB, err := postgres.NewConnToPostgres(cfg)
	if err != nil {
		gameLogger.Error("PostgresSQL init:", err.Error())
		os.Exit(1)
	}
	defer pgDB.Close()

	gameLogger.Info("PostgresSQL connected!")

	s := server.NewServer(cfg, pgDB, gameLogger)
	err = s.Run()
	if err != nil {
		gameLogger.Error("Server Run:", err.Error())
		os.Exit(1)
	}
}
