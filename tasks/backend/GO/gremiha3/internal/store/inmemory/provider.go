package inmemory

import (
	"context"
	"errors"

	"github.com/KozlovNikolai/test-task/internal/model"
)

// DeleteProvider implements store.IRepository.
func (repo *Repository) DeleteProvider(context.Context, int) error {
	panic("unimplemented")
}

// GetAllProviders implements store.IRepository.
func (repo *Repository) GetAllProviders(_ context.Context) ([]model.Provider, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	var providers []model.Provider
	for _, provider := range repo.providers {
		providers = append(providers, provider)
	}
	return providers, nil
}

// CreateProvider implements store.IRepository.
func (repo *Repository) CreateProvider(_ context.Context, provider model.Provider) (int, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	provider.ID = repo.nextProvidersID
	repo.nextProvidersID++
	repo.providers[provider.ID] = provider

	return provider.ID, nil
}

// GetProviderByID implements store.IRepository.
func (repo *Repository) GetProviderByID(_ context.Context, providerID int) (model.Provider, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	provider, exists := repo.providers[providerID]
	if !exists {
		return model.Provider{}, errors.New("provider not found")
	}
	return provider, nil
}

// UpdateProvider implements store.IRepository.
func (repo *Repository) UpdateProvider(context.Context, int) error {
	panic("unimplemented")
}
