package pgstore

import (
	"context"

	"github.com/KozlovNikolai/test-task/internal/model"
)

// CreateOrderState implements store.IRepository.
func (repo *Repository) CreateOrderState(context.Context, model.OrderState) (int, error) {
	panic("unimplemented")
}

// DeleteOrderState implements store.IRepository.
func (repo *Repository) DeleteOrderState(context.Context, int) error {
	panic("unimplemented")
}

// GetAllOrderStates implements store.IRepository.
func (repo *Repository) GetAllOrderStates(ctx context.Context) ([]model.OrderState, error) {
	panic("unimplemented")
}

// GetOrderStateByID implements store.IRepository.
func (repo *Repository) GetOrderStateByID(context.Context, int) (model.OrderState, error) {
	panic("unimplemented")
}

// GetOrderStateByName implements store.IRepository.
func (repo *Repository) GetOrderStateByName(context.Context, string) (int, error) {
	panic("unimplemented")
}

// UpdateOrderState implements store.IRepository.
func (repo *Repository) UpdateOrderState(context.Context, int) error {
	panic("unimplemented")
}

// UpdateUser implements store.IRepository.
func (repo *Repository) UpdateUser(context.Context, int) error {
	panic("unimplemented")
}
