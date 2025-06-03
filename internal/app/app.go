package app

import (
	"context"
	"fmt"

	"pet-project/internal/config"
	"pet-project/internal/domain"
	"pet-project/internal/logger"
	"pet-project/internal/service"
	"pet-project/internal/storage"
)

type Application struct {
	Config *config.Config
	Service   service.UserService
	Logger *logger.Logger
}

func New(cfg *config.Config, repo *storage.PostgresStorage, logger *logger.Logger) *Application {
	svc := service.New(repo, logger)
	return &Application{
		Config:  cfg,
		Service: svc,
		Logger:  logger,
	}
}

func (app *Application) Run(ctx context.Context) error {
	
	user := domain.User{
		FirstName: "Иван",
		LastName:  "Иванов",
		Age:       25,
		IsMarried: false,
		Password:  "securepassword123",
	}

	id, err := app.Service.CreateUser(ctx, user)
	if err != nil {
		app.Logger.Fatal(err, "Не удалось создать пользователя", "user", fmt.Sprintf("%+v", user))
	}
	app.Logger.Info("Пользователь создан", "id", id)

	fetchedUser, err := app.Service.GetUserByID(ctx, id)
	if err != nil {
		app.Logger.Fatal(err, "Не удалось получить пользователя", id)
	}
	app.Logger.Info("Пользователь получен", "user", fmt.Sprintf("%+v", fetchedUser))

	_, err = app.Service.CreateUser(ctx, user)
	if err != nil {
		app.Logger.Info("Ожидаемая ошибка при создании дубликата пользователя", "error", err)
	}

	return nil
}
