package storage

import(
	"pet-project/internal/domain"
)

type UserRepository interface{
	CreateUser(user domain.User) (int, error)
	GetUserByID(id int) (domain.User, error)
	Close() error
}

//интерфейс для Product

//интерфейс для Order