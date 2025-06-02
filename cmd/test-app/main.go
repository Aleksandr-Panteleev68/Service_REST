package main

import (
	"context"
	"os"

	"pet-project/internal/app"
	"pet-project/internal/config"
	"pet-project/internal/logger"
	"pet-project/internal/storage"
)

func main() {
	logger := logger.New()

	configPath := os.Getenv("PET_CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
		logger.Info("Переменная PET_CONFIG_PATH не установлена, используется значение по умолчанию", "path", configPath)
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		logger.Fatal(err, "Не удалось загрузить конфигурацию", "path", configPath)
	}

	ctx := context.Background()
	repo, err := storage.NewDB(ctx, cfg, logger)
	if err != nil {
		logger.Fatal(err, "Не удалось инициализировать БД")
	}
	defer repo.Close()
		
	application := app.New(cfg, repo, logger)
	if err := application.Run(); err != nil {
		logger.Fatal(err, "Не удалось запустить приложение")
	}
}
