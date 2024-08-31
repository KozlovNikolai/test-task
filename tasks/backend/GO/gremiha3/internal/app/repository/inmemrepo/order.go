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
	orders       map[int]models.Order
	nextOrdersID int
	mutex        sync.RWMutex
}

func NewOrderRepo() *OrderRepo {
	return &OrderRepo{
		orders:       make(map[int]models.Order),
		nextOrdersID: 1,
	}
}

// CreateOrder implements services.IOrderRepository.
func (repo *OrderRepo) CreateOrder(_ context.Context, order domain.Order) (domain.Order, error) {

	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	// мапим домен в модель
	dbOrder := domainToOrder(order)
	dbOrder.ID = repo.nextOrdersID

	// инкрементируем счетчик записей
	repo.nextOrdersID++
	// сохраняем
	repo.orders[dbOrder.ID] = dbOrder
	// мапим модель в домен
	domainOrder, err := orderToDomain(dbOrder)
	if err != nil {
		return domain.Order{}, fmt.Errorf("failed to create domain Order: %w", err)
	}
	return domainOrder, nil
}

// GetOrders implements services.IOrderRepository.
func (repo *OrderRepo) GetOrders(_ context.Context, limit int, offset int) ([]domain.Order, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	// извлекаем все ключи из мапы и сортируем их
	keys := make([]int, 0, len(repo.orders))
	for k := range repo.orders {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	// выбираем записи с нужными ключами
	var orders []models.Order
	for i := offset; i < offset+limit && i < len(keys); i++ {
		orders = append(orders, repo.orders[i])
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
func (repo *OrderRepo) GetOrderByID(_ context.Context, id int) (domain.Order, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	order, exists := repo.orders[id]
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
	for _, order := range repo.orders {
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
	_, exists := repo.orders[dbOrder.ID]
	if !exists {
		return domain.Order{}, fmt.Errorf("order with id %d - %w", dbOrder.ID, domain.ErrNotFound)
	}
	// обновляем запись
	repo.orders[dbOrder.ID] = dbOrder
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
	_, exists := repo.orders[id]
	if !exists {
		return fmt.Errorf("order with id %d - %w", id, domain.ErrNotFound)
	}
	delete(repo.orders, id)
	return nil
}
