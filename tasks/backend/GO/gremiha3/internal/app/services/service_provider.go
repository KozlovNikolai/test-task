package services

import (
	"context"

	"github.com/KozlovNikolai/test-task/internal/app/domain"
)

// ProviderService is a Provider service
type ProviderService struct {
	repo IProviderRepository
}

// NewProviderService creates a new Provider service
func NewProviderService(repo IProviderRepository) ProviderService {
	return ProviderService{
		repo: repo,
	}
}

func (s ProviderService) GetProvider(ctx context.Context, id int) (domain.Provider, error) {
	return s.repo.GetProvider(ctx, id)
}

func (s ProviderService) CreateProvider(ctx context.Context, Provider domain.Provider) (domain.Provider, error) {
	return s.repo.CreateProvider(ctx, Provider)
}

func (s ProviderService) UpdateProvider(ctx context.Context, Provider domain.Provider) (domain.Provider, error) {
	return s.repo.UpdateProvider(ctx, Provider)
}

func (s ProviderService) DeleteProvider(ctx context.Context, id int) error {
	return s.repo.DeleteProvider(ctx, id)
}

// GetProviders is getting list of providers. args: limit, offset
func (s ProviderService) GetProviders(ctx context.Context, limit, offset int) ([]domain.Provider, error) {
	return s.repo.GetProviders(ctx, limit, offset)
}
