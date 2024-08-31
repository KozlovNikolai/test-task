package pgrepo

import (
	"context"
	"fmt"

	"github.com/KozlovNikolai/test-task/internal/app/domain"
	"github.com/KozlovNikolai/test-task/internal/app/repository/models"
	"github.com/KozlovNikolai/test-task/internal/pkg/pg"
)

type ProviderRepo struct {
	db *pg.DB
}

func NewProviderRepo(db *pg.DB) *ProviderRepo {
	return &ProviderRepo{
		db: db,
	}
}

// CreateProvider implements services.IProviderRepository.
func (r *ProviderRepo) CreateProvider(ctx context.Context, provider domain.Provider) (domain.Provider, error) {
	dbProvider := domainToProvider(provider)

	var insertedProvider models.Provider

	// Начинаем транзакцию
	tx, err := r.db.WR.Begin(ctx)
	if err != nil {
		return domain.Provider{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Вставка данных о поставщике и получение ID
	err = tx.QueryRow(ctx, `
			INSERT INTO providers (name, origin)
			VALUES ($1, $2)
			RETURNING id,name,origin`, dbProvider.Name, dbProvider.Origin).
		Scan(
			&insertedProvider.ID,
			&insertedProvider.Name,
			&insertedProvider.Origin)
	if err != nil {
		return domain.Provider{}, fmt.Errorf("failed to insert provider: %w", err)
	}
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return domain.Provider{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	domainProvider, err := providerToDomain(insertedProvider)
	if err != nil {
		return domain.Provider{}, fmt.Errorf("failed to create domain provider: %w", err)
	}

	return domainProvider, nil
}

// DeleteProvider implements service.IProviderRepository.
func (r *ProviderRepo) DeleteProvider(ctx context.Context, id int) error {
	if id == 0 {
		return fmt.Errorf("%w: id", domain.ErrRequired)
	}
	// Начинаем транзакцию
	tx, err := r.db.WR.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	// Проверяем, что поставщик не связан ни с одним товаром.
	var count int
	err = tx.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM products
		WHERE provider_id = (SELECT id FROM providers WHERE id = $1)`, id).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to request the providers products: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("error, there are provider-related products.: %w", err)
	}
	// Удаляем поставщика
	_, err = tx.Exec(ctx, `
		DELETE FROM providers
		WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete provider with id %d: %w", id, err)
	}
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// GetProviders implements service.IProviderRepository.
func (r *ProviderRepo) GetProviders(ctx context.Context, limit, offset int) ([]domain.Provider, error) {

	query := `
		SELECT id, name, origin
		FROM providers
		ORDER BY id
		LIMIT $1 OFFSET $2
	`
	// Запрос
	rows, err := r.db.RO.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()
	// Заполняем массив поставщиков
	var providers []models.Provider
	for rows.Next() {
		var provider models.Provider
		err := rows.Scan(&provider.ID, &provider.Name, &provider.Origin)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		providers = append(providers, provider)
	}

	// Проверка на ошибки, возникшие при итерации по строкам
	if rows.Err() != nil {
		return nil, fmt.Errorf("error occurred during row iteration: %w", rows.Err())
	}
	// мапим модель в домен
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

// GetProvider implements service.IProviderRepository.
func (r *ProviderRepo) GetProvider(ctx context.Context, id int) (domain.Provider, error) {
	if id == 0 {
		return domain.Provider{}, fmt.Errorf("%w: id", domain.ErrRequired)
	}

	// SQL-запрос на получение данных Поставщика по ID
	query := `
SELECT id, name, origin
FROM providers
WHERE id = $1
`
	var provider models.Provider
	// Выполняем запрос и сканируем результат в структуру Provider
	err := r.db.RO.QueryRow(ctx, query, id).Scan(&provider.ID, &provider.Name, &provider.Origin)
	if err != nil {
		return domain.Provider{}, fmt.Errorf("failed to get provider by id: %w", err)
	}

	domainProvider, err := providerToDomain(provider)
	if err != nil {
		return domain.Provider{}, fmt.Errorf("failed to create domain provider: %w", err)
	}

	return domainProvider, nil
}

// UpdateProvider implements service.IProviderRepository.
func (r *ProviderRepo) UpdateProvider(ctx context.Context, provider domain.Provider) (domain.Provider, error) {
	dbProvider := domainToProvider(provider)
	//dbProvider.UpdatedAt = time.Now()
	// Начинаем транзакцию
	tx, err := r.db.WR.Begin(ctx)
	if err != nil {
		return domain.Provider{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	// SQL-запрос на обновление данных Поставщика
	query := `
		UPDATE providers
		SET name = $1, origin = $2
		WHERE id = $3
		RETURNING id, name, origin
	`
	var updatedProvider models.Provider

	// Выполняем запрос и сканируем обновленный результат в структуру Provider
	err = tx.QueryRow(ctx, query, dbProvider.Name, dbProvider.Origin, dbProvider.ID).
		Scan(&updatedProvider.ID, &updatedProvider.Name, &updatedProvider.Origin)
	if err != nil {
		return domain.Provider{}, fmt.Errorf("failed to update provider: %w", err)
	}
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return domain.Provider{}, fmt.Errorf("failed to commit transaction: %w", err)
	}
	domainProvider, err := providerToDomain(updatedProvider)
	if err != nil {
		return domain.Provider{}, fmt.Errorf("failed to create domain provider: %w", err)
	}

	return domainProvider, nil
}
