package inmemrepo

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/KozlovNikolai/test-task/internal/app/domain"
	"github.com/KozlovNikolai/test-task/internal/app/repository/models"
)

type ProviderRepo struct {
	providers       map[int]models.Provider
	nextProvidersID int
	mutex           sync.RWMutex
}

func NewProviderRepo() *ProviderRepo {
	return &ProviderRepo{
		providers:       make(map[int]models.Provider),
		nextProvidersID: 1,
	}
}

// CreateProvider implements services.IProviderRepository.
func (repo *ProviderRepo) CreateProvider(_ context.Context, provider domain.Provider) (domain.Provider, error) {

	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	// мапим домен в модель
	dbProvider := domainToProvider(provider)
	dbProvider.ID = repo.nextProvidersID
	// инкрементируем счетчик записей
	repo.nextProvidersID++
	// сохраняем
	repo.providers[dbProvider.ID] = dbProvider
	// мапим модель в домен
	domainProvider, err := providerToDomain(dbProvider)
	if err != nil {
		return domain.Provider{}, fmt.Errorf("failed to create domain provider: %w", err)
	}
	return domainProvider, nil
}

// DeleteProvider implements services.IProviderRepository.
func (repo *ProviderRepo) DeleteProvider(_ context.Context, id int) error {
	if id == 0 {
		return fmt.Errorf("%w: id", domain.ErrRequired)
	}
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	_, exists := repo.providers[id]
	if !exists {
		return fmt.Errorf("provider with id %d - %w", id, domain.ErrNotFound)
	}
	delete(repo.providers, id)
	return nil
}

// GetProvider implements services.IProviderRepository.
func (repo *ProviderRepo) GetProvider(_ context.Context, id int) (domain.Provider, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	provider, exists := repo.providers[id]
	if !exists {
		return domain.Provider{}, fmt.Errorf("provider with id %d - %w", id, domain.ErrNotFound)
	}
	domainProvider, err := providerToDomain(provider)
	if err != nil {
		return domain.Provider{}, fmt.Errorf("failed to create domain provider: %w", err)
	}
	return domainProvider, nil
}

// GetProviders implements services.IProviderRepository.
func (repo *ProviderRepo) GetProviders(_ context.Context, limit int, offset int) ([]domain.Provider, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	// извлекаем все ключи из мапы и сортируем их
	keys := make([]int, 0, len(repo.providers))
	for k := range repo.providers {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	// выбираем записи с нужными ключами
	var providers []models.Provider
	for i := offset; i < offset+limit && i < len(keys); i++ {
		providers = append(providers, repo.providers[i])
	}

	// мапим массив моделей в массив доменов
	domainProviders := make([]domain.Provider, len(providers))
	for i, provider := range providers {
		domainProvider, err := providerToDomain(provider)
		if err != nil {
			return nil, fmt.Errorf("failed to create domain provider: %w", err)
		}
		domainProviders[i] = domainProvider
	}
	return domainProviders, nil
}

// UpdateProvider implements services.IProviderRepository.
func (repo *ProviderRepo) UpdateProvider(_ context.Context, provider domain.Provider) (domain.Provider, error) {
	dbProvider := domainToProvider(provider)
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	// проверяем наличие поставщика
	_, exists := repo.providers[dbProvider.ID]
	if !exists {
		return domain.Provider{}, fmt.Errorf("provider with id %d - %w", dbProvider.ID, domain.ErrNotFound)
	}
	// обновляем запись
	repo.providers[dbProvider.ID] = dbProvider
	domainProvider, err := providerToDomain(dbProvider)
	if err != nil {
		return domain.Provider{}, fmt.Errorf("failed to create domain provider: %w", err)
	}
	return domainProvider, nil
}
