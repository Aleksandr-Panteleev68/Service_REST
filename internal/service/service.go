package service

import (
	"context"
	"errors"
	"fmt"

	"pet-project/internal/domain"
	"pet-project/internal/logger"
	"pet-project/internal/storage"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrValidation = errors.New("validation error")
	ErrConflict   = errors.New("conflict error")
	ErrNotFound   = errors.New("not found error")
)

type UserService interface {
	CreateUser(ctx context.Context, user domain.User) (int64, error)
	GetUserByID(ctx context.Context, id int64) (domain.User, error)
}

type service struct {
	repo   *storage.PostgresStorage
	logger *logger.Logger
}

func New(repo *storage.PostgresStorage, logger *logger.Logger) *service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

func (s *service) CreateUser(ctx context.Context, user domain.User) (int64, error) {
	s.logger.Debug("Validating user", "first_name", user.FirstName, "last_name", user.LastName)

	id, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == storage.ErrCodeUniqueViolation {
			s.logger.Error(nil, "User already exists", "first_name", user.FirstName, "last_name", user.LastName)
			return 0, fmt.Errorf("%w: user with name %s %s already exissts", ErrConflict, user.FirstName, user.LastName)
		}
		s.logger.Error(err, "Failed to create user", "first_name", user.FirstName, "last_name", user.LastName)
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	s.logger.Info("User created successfully", "id", id)
	return id, nil
}

func (s *service) GetUserByID(ctx context.Context, id int64) (domain.User, error) {
	s.logger.Debug("Fetching user", "id", id)
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.Error(nil, "User not found", "id", id)
			return domain.User{}, fmt.Errorf("%w: user with id %d not found", ErrNotFound, id)
		}
		s.logger.Error(err, "Failed to get user", "id", id)
		return domain.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	s.logger.Debug("User fetched succesfully", "id", id)
	return user, nil
}
