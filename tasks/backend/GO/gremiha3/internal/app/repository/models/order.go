package models

import "time"

// Order is a domain order
type Order struct {
	ID          int
	UserID      int
	StateID     int
	TotalAmount float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
