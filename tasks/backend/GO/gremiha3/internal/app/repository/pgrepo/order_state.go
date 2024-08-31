package pgrepo

import (
	"context"
	"fmt"

	"github.com/KozlovNikolai/test-task/internal/app/domain"
	"github.com/KozlovNikolai/test-task/internal/app/repository/models"
	"github.com/KozlovNikolai/test-task/internal/pkg/pg"
)

type OrderStateRepo struct {
	db *pg.DB
}

func NewOrderStateRepo(db *pg.DB) *OrderStateRepo {
	return &OrderStateRepo{
		db: db,
	}
}

// CreateOrderState implements services.IOrderStateRepository.
func (o *OrderStateRepo) CreateOrderState(ctx context.Context, orderState domain.OrderState) (domain.OrderState, error) {
	dbOrderState := domainToOrderState(orderState)

	var insertedOrderState models.OrderState

	// Начинаем транзакцию
	tx, err := o.db.WR.Begin(ctx)
	if err != nil {
		return domain.OrderState{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Вставка данных о типе статуса и получение ID
	err = tx.QueryRow(ctx, `
			INSERT INTO order_states (name)
			VALUES ($1)
			RETURNING id,name`, dbOrderState.Name).
		Scan(
			&insertedOrderState.ID,
			&insertedOrderState.Name)
	if err != nil {
		return domain.OrderState{}, fmt.Errorf("failed to insert OrderState: %w", err)
	}
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return domain.OrderState{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	domainOrderState, err := orderStateToDomain(insertedOrderState)
	if err != nil {
		return domain.OrderState{}, fmt.Errorf("failed to create domain OrderState: %w", err)
	}

	return domainOrderState, nil
}

// DeleteOrderState implements service.IOrderStateRepository.
func (o *OrderStateRepo) DeleteOrderState(ctx context.Context, id int) error {
	if id == 0 {
		return fmt.Errorf("%w: id", domain.ErrRequired)
	}
	// Начинаем транзакцию
	tx, err := o.db.WR.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	// Проверяем, что тип статуса не связан ни с одним заказом.
	var count int
	err = tx.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM orders
		WHERE state_id = (SELECT id FROM order_states WHERE id = $1)`, id).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to request the OrderStates products: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("error, there are OrderState-related orders.: %w", err)
	}
	// Удаляем поставщика
	_, err = tx.Exec(ctx, `
		DELETE FROM order_states
		WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete OrderState with id %d: %w", id, err)
	}
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// GetOrderStates implements service.IOrderStateRepository.
func (o *OrderStateRepo) GetOrderStates(ctx context.Context, limit, offset int) ([]domain.OrderState, error) {

	query := `
		SELECT id, name
		FROM order_states
		ORDER BY id
		LIMIT $1 OFFSET $2
	`
	// Запрос
	rows, err := o.db.RO.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()
	// Заполняем массив типов статусов
	var orderStates []models.OrderState
	for rows.Next() {
		var orderState models.OrderState
		err := rows.Scan(&orderState.ID, &orderState.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		orderStates = append(orderStates, orderState)
	}

	// Проверка на ошибки, возникшие при итерации по строкам
	if rows.Err() != nil {
		return nil, fmt.Errorf("error occurred during row iteration: %w", rows.Err())
	}
	// мапим модель в домен
	domainOrderStates := make([]domain.OrderState, len(orderStates))
	for i, orderState := range orderStates {
		domainOrderState, err := orderStateToDomain(orderState)
		if err != nil {
			return nil, fmt.Errorf("failed to create domain OrderState: %w", err)
		}

		domainOrderStates[i] = domainOrderState
	}
	return domainOrderStates, nil
}

// GetOrderState implements service.IOrderStateRepository.
func (o *OrderStateRepo) GetOrderState(ctx context.Context, id int) (domain.OrderState, error) {
	if id == 0 {
		return domain.OrderState{}, fmt.Errorf("%w: id", domain.ErrRequired)
	}

	// SQL-запрос на получение данных типа заказа по ID
	query := `
SELECT id, name
FROM order_states
WHERE id = $1
`
	var orderState models.OrderState
	// Выполняем запрос и сканируем результат в структуру orderState
	err := o.db.RO.QueryRow(ctx, query, id).Scan(&orderState.ID, &orderState.Name)
	if err != nil {
		return domain.OrderState{}, fmt.Errorf("failed to get OrderState by id: %w", err)
	}

	domainOrderState, err := orderStateToDomain(orderState)
	if err != nil {
		return domain.OrderState{}, fmt.Errorf("failed to create domain OrderState: %w", err)
	}

	return domainOrderState, nil
}

// UpdateOrderState implements service.IOrderStateRepository.
func (o *OrderStateRepo) UpdateOrderState(ctx context.Context, orderState domain.OrderState) (domain.OrderState, error) {
	dbOrderState := domainToOrderState(orderState)
	//dbOrderState.UpdatedAt = time.Now()
	// Начинаем транзакцию
	tx, err := o.db.WR.Begin(ctx)
	if err != nil {
		return domain.OrderState{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	// SQL-запрос на обновление данных типа заказа
	query := `
		UPDATE order_states
		SET name = $1
		WHERE id = $2
		RETURNING id, name
	`
	var updatedOrderState models.OrderState

	// Выполняем запрос и сканируем обновленный результат в структуру OrderState
	err = tx.QueryRow(ctx, query, dbOrderState.Name, dbOrderState.ID).
		Scan(&updatedOrderState.ID, &updatedOrderState.Name)
	if err != nil {
		return domain.OrderState{}, fmt.Errorf("failed to update OrderState: %w", err)
	}
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return domain.OrderState{}, fmt.Errorf("failed to commit transaction: %w", err)
	}
	domainOrderState, err := orderStateToDomain(updatedOrderState)
	if err != nil {
		return domain.OrderState{}, fmt.Errorf("failed to create domain OrderState: %w", err)
	}

	return domainOrderState, nil
}
