package service

import (
	"context"
	"fmt"
	"strings"

	"pet-project/internal/domain"
	"pet-project/internal/logger"
	"pet-project/internal/storage"
)

type Service interface {
	CreateUser(ctx context.Context, user domain.User) (int64, error)
	GetUserByID(ctx context.Context, id int64) (domain.User, error)
}

type service struct {
	repo   *storage.PostgresStorage
	logger *logger.Logger
}

func New(repo *storage.PostgresStorage, logger *logger.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

func (s *service) CreateUser(ctx context.Context, user domain.User) (int64, error) {
	s.logger.Info("Validating user", "first_name", user.FirstName, "last_name", user.LastName)

	if strings.TrimSpace(user.FirstName) == "" || strings.TrimSpace(user.LastName) == "" {
		return 0, fmt.Errorf("first_name and last_name cannot be empty")
	}
	if user.Age < 18 {
		return 0, fmt.Errorf("age must be at least 18, got %d", user.Age)
	}
	if len(user.Password) < 8 {
		return 0, fmt.Errorf("password must be at least 8 characters, got %d", len(user.Password))
	}

	id, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		s.logger.Error(err, "Failed to create user", "first_name", user.FirstName, "last_name", user.LastName)
		return 0, err
	}

	s.logger.Info("User created successfully", "id", id)
	return id, nil
}

func (s *service) GetUserByID(ctx context.Context, id int64) (domain.User, error) {
	s.logger.Info("Fetching user", "id", id)
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		s.logger.Error(err, "Failed to get user", "id", id)
		return domain.User{}, err
	}

	s.logger.Info("User fetched successfully", "id", id, "full_name", user.FullName)
	return user, nil
}
