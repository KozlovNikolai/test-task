package inmemrepo

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/KozlovNikolai/test-task/internal/app/domain"
	"github.com/KozlovNikolai/test-task/internal/app/repository/models"
)

type OrderRepo struct {
	db    *inMemStore
	mutex sync.RWMutex
}

func NewOrderRepo(db *inMemStore) *OrderRepo {
	return &OrderRepo{
		db: db,
	}
}

// CreateOrder implements services.IOrderRepository.
func (repo *OrderRepo) CreateOrder(_ context.Context, order domain.Order) (domain.Order, error) {

	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	// проверяем, существует ли пользователь
	if _, exists := repo.db.users[order.UserID()]; !exists {
		return domain.Order{}, fmt.Errorf("user with id %d does not exist", order.UserID())
	}
	// проверяем, существует ли статус
	if _, exists := repo.db.orderStates[order.StateID()]; !exists {
		return domain.Order{}, fmt.Errorf("order state with id %d does not exist", order.StateID())
	}
	// мапим домен в модель
	dbOrder := domainToOrder(order)
	dbOrder.ID = repo.db.nextOrdersID

	// инкрементируем счетчик записей
	repo.db.nextOrdersID++
	// сохраняем
	repo.db.orders[dbOrder.ID] = dbOrder
	// мапим модель в домен
	domainOrder, err := orderToDomain(dbOrder)
	if err != nil {
		return domain.Order{}, fmt.Errorf("failed to create domain Order: %w", err)
	}
	return domainOrder, nil
}

// GetOrders implements services.IOrderRepository.
func (repo *OrderRepo) GetOrders(_ context.Context, limit, offset, userid int) ([]domain.Order, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	// извлекаем все ключи из мапы и сортируем их
	keys := make([]int, 0, len(repo.db.orders))
	for k := range repo.db.orders {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	// выбираем все заказы выбранного пользователя
	var ordersByUserID []models.Order
	for i := 0; i < len(keys); i++ {
		order := repo.db.orders[keys[i]]
		if order.UserID == userid {
			ordersByUserID = append(ordersByUserID, order)
		}
	}

	// выбираем записи с нужными ключами
	var orders []models.Order
	for i := offset; i < offset+limit && i < len(ordersByUserID); i++ {
		orders = append(orders, ordersByUserID[i])
	}

	// мапим массив моделей в массив доменов
	domainOrders := make([]domain.Order, len(orders))
	for i, order := range orders {
		domainOrder, err := orderToDomain(order)
		if err != nil {
			return nil, fmt.Errorf("failed to create domain Order: %w", err)
		}
		domainOrders[i] = domainOrder
	}
	return domainOrders, nil
}

// GetOrderByID implements services.IOrderRepository.
func (repo *OrderRepo) GetOrder(_ context.Context, id int) (domain.Order, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	order, exists := repo.db.orders[id]
	if !exists {
		return domain.Order{}, fmt.Errorf("order with id %d - %w", id, domain.ErrNotFound)
	}
	domainOrder, err := orderToDomain(order)
	if err != nil {
		return domain.Order{}, fmt.Errorf("failed to create domain Order: %w", err)
	}
	return domainOrder, nil
}

// GetOrdersByUserID implements services.IOrderRepository.
func (repo *OrderRepo) GetOrdersByUserID(_ context.Context, userID int, limit int, offset int) ([]domain.Order, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	var orders []models.Order
	for _, order := range repo.db.orders {
		if order.UserID == userID {
			orders = append(orders, order)
		}
	}
	// сортируем по полю id
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].ID < orders[j].ID
	})

	// выбираем записи с нужными ключами
	var limitOrders []models.Order
	for i := offset; i < offset+limit && i < len(orders); i++ {
		limitOrders = append(limitOrders, orders[i])
	}

	// мапим массив моделей в массив доменов
	domainOrders := make([]domain.Order, len(orders))
	for i, order := range limitOrders {
		domainOrder, err := orderToDomain(order)
		if err != nil {
			return nil, fmt.Errorf("failed to create domain Item: %w", err)
		}
		domainOrders[i] = domainOrder
	}

	return domainOrders, nil
}

// UpdateOrder implements services.IOrderRepository.
func (repo *OrderRepo) UpdateOrder(_ context.Context, order domain.Order) (domain.Order, error) {
	dbOrder := domainToOrder(order)
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	// проверяем наличие записи
	_, exists := repo.db.orders[dbOrder.ID]
	if !exists {
		return domain.Order{}, fmt.Errorf("order with id %d - %w", dbOrder.ID, domain.ErrNotFound)
	}
	// обновляем запись
	repo.db.orders[dbOrder.ID] = dbOrder
	domainOrder, err := orderToDomain(dbOrder)
	if err != nil {
		return domain.Order{}, fmt.Errorf("failed to create domain Order: %w", err)
	}
	return domainOrder, nil
}

// DeleteOrder implements services.IOrderRepository.
func (repo *OrderRepo) DeleteOrder(_ context.Context, id int) error {
	if id == 0 {
		return fmt.Errorf("%w: id", domain.ErrRequired)
	}
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	_, exists := repo.db.orders[id]
	if !exists {
		return fmt.Errorf("order with id %d - %w", id, domain.ErrNotFound)
	}
	delete(repo.db.orders, id)
	return nil
}
