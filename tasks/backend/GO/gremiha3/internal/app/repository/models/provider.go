package models

import (
	"time"
)

// Provider is a domain provider.
type Provider struct {
	ID        int
	Name      string
	Origin    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
