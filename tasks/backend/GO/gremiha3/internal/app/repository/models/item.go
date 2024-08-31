package models

import "time"

// Item is a domain item
type Item struct {
	ID         int
	ProductID  int
	Quantity   int
	TotalPrice float64
	OrderID    int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
