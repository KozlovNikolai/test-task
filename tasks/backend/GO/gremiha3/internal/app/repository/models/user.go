package models

import (
	"time"
)

// User is a domain user.
type User struct {
	ID        int
	Login     string
	Password  string
	Role      string
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
