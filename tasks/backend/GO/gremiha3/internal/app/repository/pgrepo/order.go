package pgrepo

import (
	"context"
	"fmt"

	"github.com/KozlovNikolai/test-task/internal/app/domain"
	"github.com/KozlovNikolai/test-task/internal/app/repository/models"
	"github.com/KozlovNikolai/test-task/internal/pkg/pg"
)

type OrderRepo struct {
	db *pg.DB
}

func NewOrderRepo(db *pg.DB) *OrderRepo {
	return &OrderRepo{
		db: db,
	}
}

// CreateOrder implements services.IOrderRepository.
func (u *OrderRepo) CreateOrder(ctx context.Context, order domain.Order) (domain.Order, error) {
	dbOrder := domainToOrder(order)

	var insertedOrder models.Order

	// Начинаем транзакцию
	tx, err := u.db.WR.Begin(ctx)
	if err != nil {
		return domain.Order{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Вставка данных о пользователе и получение ID
	err = tx.QueryRow(ctx, `
			INSERT INTO orders (user_id,state_id,total_amount,created_at)
			VALUES ($1, $2, $3, $4)
			RETURNING id,user_id,state_id,total_amount,created_at`,
		dbOrder.UserID, dbOrder.StateID, dbOrder.TotalAmount, dbOrder.CreatedAt).
		Scan(
			&insertedOrder.ID,
			&insertedOrder.UserID,
			&insertedOrder.StateID,
			&insertedOrder.TotalAmount,
			&insertedOrder.CreatedAt)
	if err != nil {
		return domain.Order{}, fmt.Errorf("failed to insert Order: %w", err)
	}
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return domain.Order{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	domainOrder, err := orderToDomain(insertedOrder)
	if err != nil {
		return domain.Order{}, fmt.Errorf("failed to create domain Order: %w", err)
	}

	return domainOrder, nil
}

// DeleteOrder implements service.IOrderRepository.
func (u *OrderRepo) DeleteOrder(ctx context.Context, id int) error {
	if id == 0 {
		return fmt.Errorf("%w: id", domain.ErrRequired)
	}
	// Начинаем транзакцию
	tx, err := u.db.WR.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	// Проверяем, что пользователь не связан ни с одним заказом.
	var count int
	err = tx.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM orders
		WHERE order_id = (SELECT id FROM orders WHERE id = $1)`, id).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to request the orders Orders: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("error, there are order-related Orders.: %w", err)
	}
	// Удаляем пользователя
	_, err = tx.Exec(ctx, `
		DELETE FROM Orders
		WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete Order with id %d: %w", id, err)
	}
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// GetOrders implements service.IOrderRepository.
func (u *OrderRepo) GetOrders(ctx context.Context, limit, offset int) ([]domain.Order, error) {

	query := `
		SELECT id,user_id,state_id,total_amount,created_at
		FROM orders
		ORDER BY id
		LIMIT $1 OFFSET $2
	`
	// Запрос
	rows, err := u.db.RO.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()
	// Заполняем массив пользователей
	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.StateID,
			&order.TotalAmount,
			&order.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		orders = append(orders, order)
	}

	// Проверка на ошибки, возникшие при итерации по строкам
	if rows.Err() != nil {
		return nil, fmt.Errorf("error occurred during row iteration: %w", rows.Err())
	}
	// мапим модель в домен
	domainOrders := make([]domain.Order, len(orders))
	for i, Order := range orders {
		domainOrder, err := orderToDomain(Order)
		if err != nil {
			return nil, fmt.Errorf("failed to create domain Order: %w", err)
		}

		domainOrders[i] = domainOrder
	}
	return domainOrders, nil
}

// GetOrderByID implements service.IOrderRepository.
func (u *OrderRepo) GetOrderByID(ctx context.Context, id int) (domain.Order, error) {
	if id == 0 {
		return domain.Order{}, fmt.Errorf("%w: id", domain.ErrRequired)
	}

	// SQL-запрос на получение данных заказа по ID
	query := `
		SELECT id,user_id,state_id,total_amount,created_at
		FROM orders
		WHERE id = $1
	`
	var order models.Order
	// Выполняем запрос и сканируем результат в структуру Order
	err := u.db.RO.QueryRow(ctx, query, id).Scan(
		&order.ID,
		&order.UserID,
		&order.StateID,
		&order.TotalAmount,
		&order.CreatedAt)
	if err != nil {
		return domain.Order{}, fmt.Errorf("failed to get Order by id: %w", err)
	}

	domainOrder, err := orderToDomain(order)
	if err != nil {
		return domain.Order{}, fmt.Errorf("failed to create domain Order: %w", err)
	}

	return domainOrder, nil
}

// GetOrderByUserID implements service.IOrderRepository.
func (u *OrderRepo) GetOrdersByUserID(ctx context.Context, userID int, limit, offset int) ([]domain.Order, error) {
	if userID == 0 {
		return nil, fmt.Errorf("%w: user_id", domain.ErrRequired)
	}
	// SQL-запрос на получение заказов Пользователя по логину
	query := `
		SELECT id,user_id,state_id,total_amount,created_at
		FROM orders
		WHERE user_id = $1
		ORDER BY id
		LIMIT $2 OFFSET $3
	`
	rows, err := u.db.RO.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()
	// Заполняем массив пользователей
	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.StateID,
			&order.TotalAmount,
			&order.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		orders = append(orders, order)
	}

	// Проверка на ошибки, возникшие при итерации по строкам
	if rows.Err() != nil {
		return nil, fmt.Errorf("error occurred during row iteration: %w", rows.Err())
	}
	// мапим модель в домен
	domainOrders := make([]domain.Order, len(orders))
	for i, Order := range orders {
		domainOrder, err := orderToDomain(Order)
		if err != nil {
			return nil, fmt.Errorf("failed to create domain Order: %w", err)
		}
		domainOrders[i] = domainOrder
	}
	return domainOrders, nil
}

// UpdateOrder implements service.IOrderRepository.
func (u *OrderRepo) UpdateOrder(ctx context.Context, order domain.Order) (domain.Order, error) {
	dbOrder := domainToOrder(order)
	//dbOrder.UpdatedAt = time.Now()
	// Начинаем транзакцию
	tx, err := u.db.WR.Begin(ctx)
	if err != nil {
		return domain.Order{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	// SQL-запрос на обновление данных Поставщика
	query := `
		UPDATE Orders
		SET user_id = $1, state_id = $2, total_amount = $3, created_at = $4
		WHERE id = $5
		RETURNING id,user_id,state_id,total_amount,created_at
	`
	var updatedOrder models.Order

	// Выполняем запрос и сканируем обновленный результат в структуру Order
	err = tx.QueryRow(ctx, query,
		dbOrder.UserID, dbOrder.StateID, dbOrder.TotalAmount, dbOrder.CreatedAt).
		Scan(
			&updatedOrder.ID,
			&updatedOrder.UserID,
			&updatedOrder.StateID,
			&updatedOrder.TotalAmount,
			&updatedOrder.CreatedAt)
	if err != nil {
		return domain.Order{}, fmt.Errorf("failed to update Order: %w", err)
	}
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return domain.Order{}, fmt.Errorf("failed to commit transaction: %w", err)
	}
	domainOrder, err := orderToDomain(updatedOrder)
	if err != nil {
		return domain.Order{}, fmt.Errorf("failed to create domain Order: %w", err)
	}

	return domainOrder, nil
}
