package domain

import (
	"time"
)

// Order is a domain order
type Order struct {
	id          int
	userID      int
	stateID     int
	totalAmount float64
	createdAt   time.Time
}

// NewOrderData is a domain order
type NewOrderData struct {
	ID          int
	UserID      int
	StateID     int
	TotalAmount float64
	CreatedAt   time.Time
}

func NewOrder(data NewOrderData) (Order, error) {
	return Order{
		id:          data.ID,
		userID:      data.UserID,
		stateID:     data.StateID,
		totalAmount: data.TotalAmount,
		createdAt:   data.CreatedAt,
	}, nil
}

func (o Order) ID() int {
	return o.id
}
func (o Order) UserID() int {
	return o.userID
}
func (o Order) StateID() int {
	return o.stateID
}
func (o Order) TotalAmount() float64 {
	return o.totalAmount
}
func (o Order) CreatedAt() time.Time {
	return o.createdAt
}

// // Order is ...
// type AddOrder struct {
// 	UserID int `json:"user_id" db:"user_id" example:"1"`
// }

// // Item is ...
// type Item struct {
// 	ID         int     `json:"id" db:"id"`
// 	ProductID  int     `json:"product_id" db:"product_id"`
// 	Quantity   int     `json:"quantity" db:"quantity"`
// 	TotalPrice float64 `json:"total_price" db:"total_price"`
// }

// // NewOrder is ...
// func NewOrder() Order {
// 	return Order{}
// }
