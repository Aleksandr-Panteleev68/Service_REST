package storage

type UserRepository interface{
	CreateUser(user User) (int, error)
	GetUserByID(id int) (User, error)
}

//интерфейс для Product

//интерфейс для Order