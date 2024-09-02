package inmemrepo

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/KozlovNikolai/test-task/internal/app/domain"
	"github.com/KozlovNikolai/test-task/internal/app/repository/models"
)

type OrderStateRepo struct {
	db    *inMemStore
	mutex sync.RWMutex
}

func NewOrderStateRepo(db *inMemStore) *OrderStateRepo {
	return &OrderStateRepo{
		db: db,
	}
}

// CreateOrderState implements services.IOrderStateRepository.
func (repo *OrderStateRepo) CreateOrderState(_ context.Context, orderState domain.OrderState) (domain.OrderState, error) {

	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	// мапим домен в модель
	dbOrderState := domainToOrderState(orderState)
	dbOrderState.ID = repo.db.nextOrderStatesID

	// инкрементируем счетчик записей
	repo.db.nextOrderStatesID++
	// сохраняем
	repo.db.orderStates[dbOrderState.ID] = dbOrderState
	// мапим модель в домен
	domainOrderState, err := orderStateToDomain(dbOrderState)
	if err != nil {
		return domain.OrderState{}, fmt.Errorf("failed to create domain OrderState: %w", err)
	}
	return domainOrderState, nil
}

// DeleteOrderState implements services.IOrderStateRepository.
func (repo *OrderStateRepo) DeleteOrderState(_ context.Context, id int) error {
	if id == 0 {
		return fmt.Errorf("%w: id", domain.ErrRequired)
	}
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	_, exists := repo.db.orderStates[id]
	if !exists {
		return fmt.Errorf("OrderState with id %d - %w", id, domain.ErrNotFound)
	}
	delete(repo.db.orderStates, id)
	return nil
}

// GetOrderState implements services.IOrderStateRepository.
func (repo *OrderStateRepo) GetOrderState(_ context.Context, id int) (domain.OrderState, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	orderState, exists := repo.db.orderStates[id]
	if !exists {
		return domain.OrderState{}, fmt.Errorf("OrderState with id %d - %w", id, domain.ErrNotFound)
	}
	domainOrderState, err := orderStateToDomain(orderState)
	if err != nil {
		return domain.OrderState{}, fmt.Errorf("failed to create domain OrderState: %w", err)
	}
	return domainOrderState, nil
}

// GetOrderStates implements services.IOrderStateRepository.
func (repo *OrderStateRepo) GetOrderStates(_ context.Context, limit int, offset int) ([]domain.OrderState, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	// извлекаем все ключи из мапы и сортируем их
	keys := make([]int, 0, len(repo.db.orderStates))
	for k := range repo.db.orderStates {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	// выбираем записи с нужными ключами
	var orderStates []models.OrderState
	for i := offset; i < offset+limit && i < len(keys); i++ {
		orderStates = append(orderStates, repo.db.orderStates[keys[i]])
	}

	// мапим массив моделей в массив доменов
	domainorderStates := make([]domain.OrderState, len(orderStates))
	for i, orderState := range orderStates {
		domainorderState, err := orderStateToDomain(orderState)
		if err != nil {
			return nil, fmt.Errorf("failed to create domain User: %w", err)
		}
		domainorderStates[i] = domainorderState
	}
	return domainorderStates, nil
}

// UpdateOrderState implements services.IOrderStateRepository.
func (repo *OrderStateRepo) UpdateOrderState(_ context.Context, orderState domain.OrderState) (domain.OrderState, error) {
	dbOrderState := domainToOrderState(orderState)
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	// проверяем наличие записи
	_, exists := repo.db.orderStates[dbOrderState.ID]
	if !exists {
		return domain.OrderState{}, fmt.Errorf("OrderState with id %d - %w", dbOrderState.ID, domain.ErrNotFound)
	}
	// обновляем запись
	repo.db.orderStates[dbOrderState.ID] = dbOrderState
	domainOrderState, err := orderStateToDomain(dbOrderState)
	if err != nil {
		return domain.OrderState{}, fmt.Errorf("failed to create domain OrderState: %w", err)
	}
	return domainOrderState, nil
}
