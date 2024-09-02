package services

import (
	"context"

	"github.com/KozlovNikolai/test-task/internal/app/domain"
)

// ProductService is a Product service
type ProductService struct {
	repo IProductRepository
}

// NewProductService creates a new Product service
func NewProductService(repo IProductRepository) ProductService {
	return ProductService{
		repo: repo,
	}
}

func (s ProductService) GetProduct(ctx context.Context, id int) (domain.Product, error) {
	return s.repo.GetProduct(ctx, id)
}

func (s ProductService) CreateProduct(ctx context.Context, Product domain.Product) (domain.Product, error) {
	return s.repo.CreateProduct(ctx, Product)
}

func (s ProductService) UpdateProduct(ctx context.Context, Product domain.Product) (domain.Product, error) {
	return s.repo.UpdateProduct(ctx, Product)
}

func (s ProductService) DeleteProduct(ctx context.Context, id int) error {
	return s.repo.DeleteProduct(ctx, id)
}

func (s ProductService) GetProducts(ctx context.Context, limit, offset int) ([]domain.Product, error) {
	return s.repo.GetProducts(ctx, limit, offset)
}
