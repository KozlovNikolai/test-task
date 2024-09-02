package inmemory

import (
	"context"
	"errors"

	"github.com/KozlovNikolai/test-task/internal/model"
)

// CreateOrder implements store.IRepository.
func (repo *Repository) CreateOrder(_ context.Context, order model.Order) (int, error) {

	// существует ли пользователь
	_, exists := repo.users[order.UserID]
	if !exists {
		return 0, errors.New("такого пользователя не существует")
	}
	// существует ли статус
	_, exists = repo.orderState[order.StateID]
	if !exists {
		return 0, errors.New("такого статуса не существует")
	}
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	order.ID = repo.nextOrdersID
	repo.nextOrdersID++
	repo.orders[order.ID] = order
	return order.ID, nil

}

// DeleteOrder implements store.IRepository.
func (repo *Repository) DeleteOrder(context.Context, int) error {
	panic("unimplemented")
}

// GetAllOrders implements store.IRepository.
func (repo *Repository) GetAllOrders(ctx context.Context) ([]model.Order, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	var orders []model.Order
	for _, order := range repo.orders {
		orders = append(orders, order)
	}
	return orders, nil
}

// GetOrderByID implements store.IRepository.
func (repo *Repository) GetOrderByID(_ context.Context, orderID int) (model.Order, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	order, exists := repo.orders[orderID]
	if !exists {
		return model.Order{}, errors.New("order not found")
	}
	return order, nil
}

// UpdateOrder implements store.IRepository.
func (repo *Repository) UpdateOrder(context.Context, int) error {
	panic("unimplemented")
}
