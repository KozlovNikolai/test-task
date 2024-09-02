package services

import (
	"context"
	"time"

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

const createdOrderStateID = 1

func (s OrderService) GetOrder(ctx context.Context, id int) (domain.Order, error) {
	return s.repo.GetOrder(ctx, id)
}

func (s OrderService) GetOrdersByUserID(ctx context.Context, userID, limit, offset int) ([]domain.Order, error) {
	return s.repo.GetOrdersByUserID(ctx, userID, limit, offset)
}

func (s OrderService) CreateOrder(ctx context.Context, order domain.Order) (domain.Order, error) {
	var newOrder = domain.NewOrderData{
		UserID:      order.UserID(),
		StateID:     createdOrderStateID,
		TotalAmount: 0,
		CreatedAt:   time.Now(),
	}
	creatingOrder, err := domain.NewOrder(newOrder)
	if err != nil {
		return domain.Order{}, err
	}
	return s.repo.CreateOrder(ctx, creatingOrder)
}

func (s OrderService) UpdateOrder(ctx context.Context, order domain.Order) (domain.Order, error) {
	return s.repo.UpdateOrder(ctx, order)
}

func (s OrderService) DeleteOrder(ctx context.Context, id int) error {
	return s.repo.DeleteOrder(ctx, id)
}

func (s OrderService) GetOrders(ctx context.Context, limit, offset, userid int) ([]domain.Order, error) {
	return s.repo.GetOrders(ctx, limit, offset, userid)
}
