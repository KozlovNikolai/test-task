package pgstore

import (
	"context"

	"github.com/KozlovNikolai/test-task/internal/model"
)

// CreateProduct implements store.IRepository.
func (repo *Repository) CreateProduct(context.Context, model.Product) (int, error) {
	panic("unimplemented")
}

// DeleteProduct implements store.IRepository.
func (repo *Repository) DeleteProduct(context.Context, int) error {
	panic("unimplemented")
}

// GetAllProducts implements store.IRepository.
func (repo *Repository) GetAllProducts(ctx context.Context) ([]model.Product, error) {
	panic("unimplemented")
}

// GetProductByID implements store.IRepository.
func (repo *Repository) GetProductByID(context.Context, int) (model.Product, error) {
	panic("unimplemented")
}

// UpdateProduct implements store.IRepository.
func (repo *Repository) UpdateProduct(context.Context, int) error {
	panic("unimplemented")
}
