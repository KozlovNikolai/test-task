package models

import (
	"time"
)

// OrderState is a domain provider.
type OrderState struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
