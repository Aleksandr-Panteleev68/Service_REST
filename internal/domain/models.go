package domain

import "time"

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	FullName  string `json:"full_name"`
	Age       int    `json:"age"`
	IsMarried bool   `json:"is_married"`
	Password  string `json:"password"`
}

type Product struct {
	ID          int64    `json:"id"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Quantity    int      `json:"quantity"`
	Price       float64  `json:"price"`
}

type Order struct {
	ID           int            `json:"id"`
	UserID       int            `json:"user_id"`
	CreatedAt    time.Time      `json:"created_at"`
	TotalPrice   float64        `json:"total_price"`
	OrderProduct []OrderProduct `json:"order_products"`
}

type OrderProduct struct {
	OrderID   int     `json:"order_id"`
	ProductID int     `json:"prooduct_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
