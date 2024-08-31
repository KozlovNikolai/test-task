package services

import (
	"context"

	"github.com/KozlovNikolai/test-task/internal/app/domain"
)

// OrderService is a Order service
type OrderService struct {
	repo IOrderRepository
}

// NewOrderService creates a new Order service
func NewOrderService(repo IOrderRepository) OrderService {
	return OrderService{
		repo: repo,
	}
}

func (s OrderService) GetOrderByID(ctx context.Context, id int) (domain.Order, error) {
	return s.repo.GetOrderByID(ctx, id)
}

func (s OrderService) GetOrdersByUserID(ctx context.Context, userID, limit, offset int) ([]domain.Order, error) {
	return s.repo.GetOrdersByUserID(ctx, userID, limit, offset)
}

func (s OrderService) CreateOrder(ctx context.Context, Order domain.Order) (domain.Order, error) {
	return s.repo.CreateOrder(ctx, Order)
}

func (s OrderService) UpdateOrder(ctx context.Context, Order domain.Order) (domain.Order, error) {
	return s.repo.UpdateOrder(ctx, Order)
}

func (s OrderService) DeleteOrder(ctx context.Context, id int) error {
	return s.repo.DeleteOrder(ctx, id)
}

func (s OrderService) GetOrders(ctx context.Context, limit, offset int) ([]domain.Order, error) {
	return s.repo.GetOrders(ctx, limit, offset)
}
