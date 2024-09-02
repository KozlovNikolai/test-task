package services

import (
	"context"

	"github.com/KozlovNikolai/test-task/internal/app/domain"
)

// ItemService is a Item service
type ItemService struct {
	repo IItemRepository
}

// GetItem implements httpserver.IItemService.
func (s ItemService) GetItem(ctx context.Context, id int) (domain.Item, error) {
	return s.repo.GetItem(ctx, id)
}

// NewItemService creates a new Item service
func NewItemService(repo IItemRepository) ItemService {
	return ItemService{
		repo: repo,
	}
}

func (s ItemService) CreateItem(ctx context.Context, Item domain.Item) (domain.Item, error) {
	return s.repo.CreateItem(ctx, Item)
}

func (s ItemService) UpdateItem(ctx context.Context, Item domain.Item) (domain.Item, error) {
	return s.repo.UpdateItem(ctx, Item)
}

func (s ItemService) DeleteItem(ctx context.Context, id int) error {
	return s.repo.DeleteItem(ctx, id)
}

func (s ItemService) GetItems(ctx context.Context, limit, offset, orderid int) ([]domain.Item, error) {
	return s.repo.GetItems(ctx, limit, offset, orderid)
}
