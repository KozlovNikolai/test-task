package inmemory

import (
	"context"
	"errors"

	"github.com/KozlovNikolai/test-task/internal/model"
)

// UpdateProduct implements store.IRepository.
func (repo *Repository) UpdateProduct(context.Context, int) error {
	panic("unimplemented")
}

// GetProductByID implements store.IRepository.
func (repo *Repository) GetProductByID(_ context.Context, productID int) (model.Product, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	product, exists := repo.products[productID]
	if !exists {
		return model.Product{}, errors.New("product not found")
	}
	return product, nil
}

// GetAllProducts implements store.IRepository.
func (repo *Repository) GetAllProducts(ctx context.Context) ([]model.Product, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	var products []model.Product
	for _, product := range repo.products {
		products = append(products, product)
	}
	return products, nil
}

// DeleteProduct implements store.IRepository.
func (repo *Repository) DeleteProduct(context.Context, int) error {
	panic("unimplemented")
}

// CreateProduct implements store.IRepository.
func (repo *Repository) CreateProduct(_ context.Context, product model.Product) (int, error) {

	_, exists := repo.providers[product.ProviderID]
	if !exists {
		return 0, errors.New("такого поставщика не существует")
	}
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	product.ID = repo.nextProductsID
	repo.nextProductsID++
	repo.products[product.ID] = product
	return product.ID, nil

}
