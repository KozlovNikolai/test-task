package inmemory

import (
	"context"
	"errors"

	"github.com/KozlovNikolai/test-task/internal/model"
)

// DeleteOrderState implements store.IRepository.
func (repo *Repository) DeleteOrderState(context.Context, int) error {
	panic("unimplemented")
}

// GetAllOrderStates implements store.IRepository.
func (repo *Repository) GetAllOrderStates(_ context.Context) ([]model.OrderState, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	var OrderStates []model.OrderState
	for _, OrderState := range repo.orderState {
		OrderStates = append(OrderStates, OrderState)
	}
	return OrderStates, nil
}

// CreateOrderState implements store.IRepository.
func (repo *Repository) CreateOrderState(_ context.Context, OrderState model.OrderState) (int, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	OrderState.ID = repo.nextOrderState
	repo.nextOrderState++
	repo.orderState[OrderState.ID] = OrderState

	return OrderState.ID, nil
}

// GetOrderStateByID implements store.IRepository.
func (repo *Repository) GetOrderStateByID(_ context.Context, OrderStateID int) (model.OrderState, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	OrderState, exists := repo.orderState[OrderStateID]
	if !exists {
		return model.OrderState{}, errors.New("OrderState not found")
	}
	return OrderState, nil
}

// GetOrderStateByID implements store.IRepository.
func (repo *Repository) GetOrderStateByName(_ context.Context, orderStateName string) (int, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	for _, orderState := range repo.orderState {
		if orderState.Name == orderStateName {
			return orderState.ID, nil
		}
	}
	return 0, errors.New("Order state not found")
}

// UpdateOrderState implements store.IRepository.
func (repo *Repository) UpdateOrderState(context.Context, int) error {
	panic("unimplemented")
}
