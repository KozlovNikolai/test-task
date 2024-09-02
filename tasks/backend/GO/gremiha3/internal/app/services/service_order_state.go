package services

import (
	"context"

	"github.com/KozlovNikolai/test-task/internal/app/domain"
)

// OrderStateService is a OrderState service
type OrderStateService struct {
	repo IOrderStateRepository
}

// NewOrderStateService creates a new OrderState service
func NewOrderStateService(repo IOrderStateRepository) OrderStateService {
	return OrderStateService{
		repo: repo,
	}
}

func (s OrderStateService) GetOrderState(ctx context.Context, id int) (domain.OrderState, error) {
	return s.repo.GetOrderState(ctx, id)
}

func (s OrderStateService) CreateOrderState(ctx context.Context, OrderState domain.OrderState) (domain.OrderState, error) {
	return s.repo.CreateOrderState(ctx, OrderState)
}

func (s OrderStateService) UpdateOrderState(ctx context.Context, OrderState domain.OrderState) (domain.OrderState, error) {
	return s.repo.UpdateOrderState(ctx, OrderState)
}

func (s OrderStateService) DeleteOrderState(ctx context.Context, id int) error {
	return s.repo.DeleteOrderState(ctx, id)
}

func (s OrderStateService) GetOrderStates(ctx context.Context, limit, offset int) ([]domain.OrderState, error) {
	return s.repo.GetOrderStates(ctx, limit, offset)
}
