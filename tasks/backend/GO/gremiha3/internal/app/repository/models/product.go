package models

import (
	"time"
)

// Product is a domain product.
type Product struct {
	ID         int
	Name       string
	ProviderID int
	Price      float64
	Stock      int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
