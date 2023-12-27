package main

import (
	"rankit/http"
	"rankit/postgres"
	"rankit/service"

	"github.com/alexedwards/scs/v2"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	db, err := postgres.Connect("postgres://rankit:rankit@localhost:5432/rankit")
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	querier := postgres.NewQuerier(db)
	usersvc := service.NewUserService(querier)
	sessionmgr := scs.New()

	logger.Info("Starting server")
	server := http.NewServer(usersvc, nil, sessionmgr, logger)
	if err := server.Listen(":8000"); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
