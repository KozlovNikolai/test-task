package pgrepo

import (
	"context"
	"fmt"

	"github.com/KozlovNikolai/test-task/internal/app/domain"
	"github.com/KozlovNikolai/test-task/internal/app/repository/models"
	"github.com/KozlovNikolai/test-task/internal/pkg/pg"
)

type ItemRepo struct {
	db *pg.DB
}

func NewItemRepo(db *pg.DB) *ItemRepo {
	return &ItemRepo{
		db: db,
	}
}

// CreateItem implements services.IItemRepository.
func (i *ItemRepo) CreateItem(ctx context.Context, item domain.Item) (domain.Item, error) {
	dbItem := domainToItem(item)

	var insertedItem models.Item

	// Начинаем транзакцию
	tx, err := i.db.WR.Begin(ctx)
	if err != nil {
		return domain.Item{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Вставка данных о добавляемом товаре в заказ
	err = tx.QueryRow(ctx, `
			INSERT INTO items (product_id,quantity,total_price,order_id)
			VALUES ($1, $2, $3, $4)
			RETURNING id,product_id,quantity,total_price,order_id`,
		dbItem.ProductID,
		dbItem.Quantity,
		dbItem.TotalPrice,
		dbItem.OrderID).
		Scan(
			&insertedItem.ID,
			&insertedItem.ProductID,
			&insertedItem.Quantity,
			&insertedItem.TotalPrice,
			&insertedItem.OrderID)
	if err != nil {
		return domain.Item{}, fmt.Errorf("failed to insert item: %w", err)
	}
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return domain.Item{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	domainItem, err := itemToDomain(insertedItem)
	if err != nil {
		return domain.Item{}, fmt.Errorf("failed to create domain item: %w", err)
	}

	return domainItem, nil
}

// DeleteItem implements service.IItemRepository.
func (i *ItemRepo) DeleteItem(ctx context.Context, id int) error {
	if id == 0 {
		return fmt.Errorf("%w: id", domain.ErrRequired)
	}
	// Начинаем транзакцию
	tx, err := i.db.WR.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	// Удаляем товар из заказа
	_, err = tx.Exec(ctx, `
		DELETE FROM items
		WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete item with id %d: %w", id, err)
	}
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// GetItems implements service.IItemRepository.
func (i *ItemRepo) GetItems(ctx context.Context, limit, offset, orderid int) ([]domain.Item, error) {

	query := `
		SELECT id,product_id,quantity,total_price,order_id
		FROM items
		ORDER BY id
		LIMIT $1 OFFSET $2
	`
	// Запрос
	rows, err := i.db.RO.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()
	// Заполняем массив пользователей
	var items []models.Item
	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&item.ID,
			&item.ProductID,
			&item.Quantity,
			&item.TotalPrice,
			&item.OrderID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		items = append(items, item)
	}

	// Проверка на ошибки, возникшие при итерации по строкам
	if rows.Err() != nil {
		return nil, fmt.Errorf("error occurred during row iteration: %w", rows.Err())
	}
	// мапим модель в домен
	domainItems := make([]domain.Item, len(items))
	for i, item := range items {
		domainItem, err := itemToDomain(item)
		if err != nil {
			return nil, fmt.Errorf("failed to create domain item: %w", err)
		}

		domainItems[i] = domainItem
	}
	return domainItems, nil
}

// GetItemByID implements service.IItemRepository.
func (i *ItemRepo) GetItem(ctx context.Context, id int) (domain.Item, error) {
	if id == 0 {
		return domain.Item{}, fmt.Errorf("%w: id", domain.ErrRequired)
	}

	query := `
		SELECT  id,product_id,quantity,total_price,order_id
		FROM items
		WHERE id = $1
	`
	var item models.Item
	// Выполняем запрос и сканируем результат в структуру item
	err := i.db.RO.QueryRow(ctx, query, id).Scan(
		&item.ID,
		&item.ProductID,
		&item.Quantity,
		&item.TotalPrice,
		&item.OrderID)
	if err != nil {
		return domain.Item{}, fmt.Errorf("failed to get item by id: %w", err)
	}

	domainItem, err := itemToDomain(item)
	if err != nil {
		return domain.Item{}, fmt.Errorf("failed to create domain item: %w", err)
	}

	return domainItem, nil
}

// UpdateItem implements service.IItemRepository.
func (i *ItemRepo) UpdateItem(ctx context.Context, item domain.Item) (domain.Item, error) {
	dbItem := domainToItem(item)
	//dbItem.UpdatedAt = time.Now()
	// Начинаем транзакцию
	tx, err := i.db.WR.Begin(ctx)
	if err != nil {
		return domain.Item{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	// SQL-запрос на обновление данных Поставщика
	query := `
		UPDATE items
		SET product_id = $1,quantity = $2,total_price = $3,order_id = $4
		WHERE id = $5
		RETURNING id,product_id,quantity,total_price,order_id
	`
	var updatedItem models.Item

	// Выполняем запрос и сканируем обновленный результат в структуру item
	err = tx.QueryRow(ctx, query,
		dbItem.ProductID,
		dbItem.Quantity,
		dbItem.TotalPrice,
		dbItem.OrderID,
		dbItem.ID).
		Scan(
			&updatedItem.ID,
			&updatedItem.ProductID,
			&updatedItem.Quantity,
			&updatedItem.TotalPrice,
			&updatedItem.OrderID)
	if err != nil {
		return domain.Item{}, fmt.Errorf("failed to update item: %w", err)
	}
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return domain.Item{}, fmt.Errorf("failed to commit transaction: %w", err)
	}
	domainItem, err := itemToDomain(updatedItem)
	if err != nil {
		return domain.Item{}, fmt.Errorf("failed to create domain item: %w", err)
	}

	return domainItem, nil
}
