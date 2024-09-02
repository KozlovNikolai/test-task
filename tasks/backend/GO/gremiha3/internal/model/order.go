package model

import (
	"time"
)

// Order is ...
type Order struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	StateID     int       `json:"state_id" db:"state_id"`
	Items       []int     `json:"items" db:"items"`
	TotalAmount float64   `json:"total_amount" db:"total_amount"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// Order is ...
type AddOrder struct {
	UserID int `json:"user_id" db:"user_id" example:"1"`
}

// Item is ...
type Item struct {
	ID         int     `json:"id" db:"id"`
	ProductID  int     `json:"product_id" db:"product_id"`
	Quantity   int     `json:"quantity" db:"quantity"`
	TotalPrice float64 `json:"total_price" db:"total_price"`
}

// NewOrder is ...
func NewOrder() Order {
	return Order{}
}
