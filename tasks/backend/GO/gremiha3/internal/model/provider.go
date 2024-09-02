package model

import (
	"errors"
)

// Provider is ...
type Provider struct {
	ID     int    `json:"id" db:"id"`
	Name   string `json:"name" db:"name"`
	Origin string `json:"origin" db:"origin"`
}

// Provider is ...
type AddProvider struct {
	Name   string `json:"name" db:"name" example:"Microsoft"`
	Origin string `json:"origin" db:"origin" example:"Kazakhstan"`
}

// NewProvider is ...
func NewProvider() Provider {
	return Provider{}
}

// Validate проверяет правильность данных поставщика
func (p *Provider) Validate() error {
	if p.Name == "" {
		return errors.New("name is required")
	}
	if p.Origin == "" {
		return errors.New("origin is required")
	}
	return nil
}
