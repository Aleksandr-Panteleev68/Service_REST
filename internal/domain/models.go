package domain

import "time"

type User struct {
	ID        int
	FirstName string
	LastName  string
	Age       int
	IsMarried bool
	Password  string
}

type Product struct {
	ID          int
	Description string
	Tags        []string
	Quantity    int
}

type Order struct {
	ID         int
	UserID     int
	CreatedAt  time.Time
	TotalPrice float64
}

type OrderProduct struct {
	OrderID   int
	ProductID int
	Quantity  int
	Price     float64
}
