package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

// DB is a shortcut structure to a Postgres DB
type DB struct {
	WR *pgxpool.Pool
	RO *pgxpool.Pool
}

// Dial creates new database connection to postgres
func Dial(dsnWR string, dsnRepl ...string) (*DB, error) {
	// создаем подключение к основной базе данных
	if dsnWR == "" {
		return nil, errors.New("no postgres DSN provided")
	}
	WR, err := pgxpool.Connect(context.Background(), dsnWR)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to master DB: %w", err)
	}
	if err := WR.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("db.Ping to master DB failed: %w", err)
	}
	// создаем подключение к реплике:
	var RO *pgxpool.Pool
	switch len(dsnRepl) {
	case 0:
		RO = WR
	default:
		RO, err = pgxpool.Connect(context.Background(), dsnRepl[0])
		if err != nil {
			return nil, fmt.Errorf("unable to connect to replica DB: %w", err)
		}
		if err := WR.Ping(context.Background()); err != nil {
			return nil, fmt.Errorf("db.Ping to replica DB failed: %w", err)
		}
	}

	return &DB{
		WR: WR,
		RO: RO,
	}, nil
}

// func saveProviderAndProduct(ctx context.Context, pool *pgxpool.Pool, provider *Provider, product *Product) error {
// 	// Начинаем транзакцию
// 	tx, err := pool.Begin(ctx)
// 	if err != nil {
// 		return fmt.Errorf("failed to begin transaction: %w", err)
// 	}
// 	defer tx.Rollback(ctx)

// 	// Вставка данных о провайдере и получение ID
// 	err = tx.QueryRow(ctx, `
// 		INSERT INTO providers (name)
// 		VALUES ($1)
// 		RETURNING id`, provider.Name).Scan(&provider.ID)
// 	if err != nil {
// 		return fmt.Errorf("failed to insert provider: %w", err)
// 	}

// 	// Вставка данных о продукте с использованием ProviderID
// 	err = tx.QueryRow(ctx, `
// 		INSERT INTO products (name, provider_id)
// 		VALUES ($1, $2)
// 		RETURNING id`, product.Name, provider.ID).Scan(&product.ID)
// 	if err != nil {
// 		return fmt.Errorf("failed to insert product: %w", err)
// 	}

// 	// Фиксация транзакции
// 	if err := tx.Commit(ctx); err != nil {
// 		return fmt.Errorf("failed to commit transaction: %w", err)
// 	}

// 	return nil
// }
