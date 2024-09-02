package pgstore

import (
	"context"

	"github.com/KozlovNikolai/test-task/internal/model"
)

// CreateOrder implements store.IRepository.
func (repo *Repository) CreateOrder(context.Context, model.Order) (int, error) {
	panic("unimplemented")
}

// DeleteOrder implements store.IRepository.
func (repo *Repository) DeleteOrder(context.Context, int) error {
	panic("unimplemented")
}

// GetAllOrders implements store.IRepository.
func (repo *Repository) GetAllOrders(ctx context.Context) ([]model.Order, error) {
	panic("unimplemented")
}

// GetOrderByID implements store.IRepository.
func (repo *Repository) GetOrderByID(context.Context, int) (model.Order, error) {
	panic("unimplemented")
}

// UpdateOrder implements store.IRepository.
func (repo *Repository) UpdateOrder(context.Context, int) error {
	panic("unimplemented")
}
