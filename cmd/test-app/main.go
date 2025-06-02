package main

import (
	"context"
	"log/slog"
	"os"

	"pet-project/internal/app"
	"pet-project/internal/config"
	"pet-project/internal/logger"
	"pet-project/internal/storage"
)

func main() {
	configPath := os.Getenv("PET_CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		slog.Error("Не удалось загрузить конфигурацию", "error", err, "path", configPath)
		os.Exit(1)
	}

	logger := logger.New(cfg.Env)

	ctx := context.Background()
	repo, err := storage.NewDB(ctx, cfg, logger)
	if err != nil {
		logger.Fatal(err, "Не удалось инициализировать БД")
	}
	defer repo.Close()

	application := app.New(cfg, repo, logger)
	if err := application.Run(ctx); err != nil {
		logger.Fatal(err, "Не удалось запустить приложение")
	}
}
