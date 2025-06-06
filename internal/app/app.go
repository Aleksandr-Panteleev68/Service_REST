package app

import (
	"context"
	"fmt"

	"pet-project/internal/api"
	"pet-project/internal/config"
	"pet-project/internal/logger"
	"pet-project/internal/service"
	"pet-project/internal/storage"
)

type Application struct {
	Config  *config.Config
	Service service.UserService
	Logger  *logger.Logger
	Handler *api.Handler
}

func New(cfg *config.Config, repo *storage.PostgresStorage, logger *logger.Logger) *Application {
	svc := service.New(repo, logger)
	handler := api.NewHandler(svc, logger, cfg)
	return &Application{
		Config:  cfg,
		Service: svc,
		Logger:  logger,
		Handler: handler,
	}
}

func (app *Application) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := app.Handler.StartServer(ctx); err != nil {
		app.Logger.Error(err, "Failed to run HTTP server")
		return fmt.Errorf("failed to run application: %w", err)
	}

	return nil
}
