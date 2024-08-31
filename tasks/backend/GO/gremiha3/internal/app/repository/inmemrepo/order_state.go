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
	orderStates       map[int]models.OrderState
	nextOrderStatesID int
	mutex             sync.RWMutex
}

func NewOrderStateRepo() *OrderStateRepo {
	return &OrderStateRepo{
		orderStates: map[int]models.OrderState{
			0: {ID: 0, Name: "Created"},
			1: {ID: 1, Name: "In progress"},
			2: {ID: 2, Name: "Delivery"},
		},
		nextOrderStatesID: 3,
	}
}

// CreateOrderState implements services.IOrderStateRepository.
func (repo *OrderStateRepo) CreateOrderState(_ context.Context, orderState domain.OrderState) (domain.OrderState, error) {

	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	// мапим домен в модель
	dbOrderState := domainToOrderState(orderState)
	dbOrderState.ID = repo.nextOrderStatesID

	// инкрементируем счетчик записей
	repo.nextOrderStatesID++
	// сохраняем
	repo.orderStates[dbOrderState.ID] = dbOrderState
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
	_, exists := repo.orderStates[id]
	if !exists {
		return fmt.Errorf("OrderState with id %d - %w", id, domain.ErrNotFound)
	}
	delete(repo.orderStates, id)
	return nil
}

// GetOrderState implements services.IOrderStateRepository.
func (repo *OrderStateRepo) GetOrderState(_ context.Context, id int) (domain.OrderState, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	orderState, exists := repo.orderStates[id]
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
	keys := make([]int, 0, len(repo.orderStates))
	for k := range repo.orderStates {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	// выбираем записи с нужными ключами
	var orderStates []models.OrderState
	for i := offset; i < offset+limit && i < len(keys); i++ {
		orderStates = append(orderStates, repo.orderStates[i])
	}

	// мапим массив моделей в массив доменов
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

// UpdateOrderState implements services.IOrderStateRepository.
func (repo *OrderStateRepo) UpdateOrderState(_ context.Context, orderState domain.OrderState) (domain.OrderState, error) {
	dbOrderState := domainToOrderState(orderState)
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	// проверяем наличие записи
	_, exists := repo.orderStates[dbOrderState.ID]
	if !exists {
		return domain.OrderState{}, fmt.Errorf("OrderState with id %d - %w", dbOrderState.ID, domain.ErrNotFound)
	}
	// обновляем запись
	repo.orderStates[dbOrderState.ID] = dbOrderState
	domainOrderState, err := orderStateToDomain(dbOrderState)
	if err != nil {
		return domain.OrderState{}, fmt.Errorf("failed to create domain OrderState: %w", err)
	}
	return domainOrderState, nil
}
