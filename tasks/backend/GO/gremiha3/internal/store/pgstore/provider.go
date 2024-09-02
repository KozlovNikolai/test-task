package pgstore

import (
	"context"

	"github.com/KozlovNikolai/test-task/internal/model"
)

// CreateProvider implements store.IRepository.
func (repo *Repository) CreateProvider(context.Context, model.Provider) (int, error) {
	panic("unimplemented")
}

// DeleteProvider implements store.IRepository.
func (repo *Repository) DeleteProvider(context.Context, int) error {
	panic("unimplemented")
}

// GetAllProviders implements store.IRepository.
func (repo *Repository) GetAllProviders(ctx context.Context) ([]model.Provider, error) {
	panic("unimplemented")
}

// GetProviderByID implements store.IRepository.
func (repo *Repository) GetProviderByID(context.Context, int) (model.Provider, error) {
	panic("unimplemented")
}

// UpdateProvider implements store.IRepository.
func (repo *Repository) UpdateProvider(context.Context, int) error {
	panic("unimplemented")
}
