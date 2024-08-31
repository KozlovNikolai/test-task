package inmemrepo

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/KozlovNikolai/test-task/internal/app/domain"
	"github.com/KozlovNikolai/test-task/internal/app/repository/models"
)

type ItemRepo struct {
	items       map[int]models.Item
	nextItemsID int
	mutex       sync.RWMutex
}

func NewItemRepo() *ItemRepo {
	return &ItemRepo{
		items:       make(map[int]models.Item),
		nextItemsID: 1,
	}
}

// CreateItem implements services.IItemRepository.
func (repo *ItemRepo) CreateItem(_ context.Context, item domain.Item) (domain.Item, error) {

	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	// мапим домен в модель
	dbItem := domainToItem(item)
	dbItem.ID = repo.nextItemsID

	// инкрементируем счетчик записей
	repo.nextItemsID++
	// сохраняем
	repo.items[dbItem.ID] = dbItem
	// мапим модель в домен
	domainItem, err := itemToDomain(dbItem)
	if err != nil {
		return domain.Item{}, fmt.Errorf("failed to create domain Item: %w", err)
	}
	return domainItem, nil
}

// GetItems implements services.IItemRepository.
func (repo *ItemRepo) GetItems(_ context.Context, limit int, offset int) ([]domain.Item, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	// извлекаем все ключи из мапы и сортируем их
	keys := make([]int, 0, len(repo.items))
	for k := range repo.items {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	// выбираем записи с нужными ключами
	var items []models.Item
	for i := offset; i < offset+limit && i < len(keys); i++ {
		items = append(items, repo.items[i])
	}

	// мапим массив моделей в массив доменов
	domainItems := make([]domain.Item, len(items))
	for i, Item := range items {
		domainItem, err := itemToDomain(Item)
		if err != nil {
			return nil, fmt.Errorf("failed to create domain Item: %w", err)
		}
		domainItems[i] = domainItem
	}
	return domainItems, nil
}

// GetItemByID implements services.IItemRepository.
func (repo *ItemRepo) GetItemByID(_ context.Context, id int) (domain.Item, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	item, exists := repo.items[id]
	if !exists {
		return domain.Item{}, fmt.Errorf("item with id %d - %w", id, domain.ErrNotFound)
	}
	domainItem, err := itemToDomain(item)
	if err != nil {
		return domain.Item{}, fmt.Errorf("failed to create domain Item: %w", err)
	}
	return domainItem, nil
}

// GetItemsByOrderID implements services.IItemRepository.
func (repo *ItemRepo) GetItemsByOrderID(_ context.Context, orderID int) ([]domain.Item, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	var items []models.Item
	for _, item := range repo.items {
		if item.OrderID == orderID {
			items = append(items, item)
		}
	}
	// сортируем по полю id
	sort.Slice(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})
	// мапим массив моделей в массив доменов
	domainItems := make([]domain.Item, len(items))
	for i, item := range items {
		domainItem, err := itemToDomain(item)
		if err != nil {
			return nil, fmt.Errorf("failed to create domain Item: %w", err)
		}
		domainItems[i] = domainItem
	}

	return domainItems, nil
}

// UpdateItem implements services.IItemRepository.
func (repo *ItemRepo) UpdateItem(_ context.Context, item domain.Item) (domain.Item, error) {
	dbItem := domainToItem(item)
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	// проверяем наличие записи
	_, exists := repo.items[dbItem.ID]
	if !exists {
		return domain.Item{}, fmt.Errorf("item with id %d - %w", dbItem.ID, domain.ErrNotFound)
	}
	// обновляем запись
	repo.items[dbItem.ID] = dbItem
	domainItem, err := itemToDomain(dbItem)
	if err != nil {
		return domain.Item{}, fmt.Errorf("failed to create domain Item: %w", err)
	}
	return domainItem, nil
}

// DeleteItem implements services.IItemRepository.
func (repo *ItemRepo) DeleteItem(_ context.Context, id int) error {
	if id == 0 {
		return fmt.Errorf("%w: id", domain.ErrRequired)
	}
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	_, exists := repo.items[id]
	if !exists {
		return fmt.Errorf("item with id %d - %w", id, domain.ErrNotFound)
	}
	delete(repo.items, id)
	return nil
}
