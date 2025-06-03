package storage

import (
	"context"

	"pet-project/internal/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user domain.User) (int, error)
	GetUserByID(ctx context.Context, id int) (domain.User, error)
}

// интерфейс для Product

// интерфейс для Order
