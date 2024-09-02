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
	db    *inMemStore
	mutex sync.RWMutex
}

func NewItemRepo(db *inMemStore) *ItemRepo {
	return &ItemRepo{
		db: db,
	}
}

// CreateItem implements services.IItemRepository.
func (repo *ItemRepo) CreateItem(_ context.Context, item domain.Item) (domain.Item, error) {

	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	// проверяем, существует ли заказ
	if _, exists := repo.db.orders[item.OrderID()]; !exists {
		return domain.Item{}, fmt.Errorf("order with id %d does not exist", item.OrderID())
	}
	// проверяем, существует ли товар
	if _, exists := repo.db.products[item.ProductID()]; !exists {
		return domain.Item{}, fmt.Errorf("product with id %d does not exist", item.ProductID())
	}

	// мапим домен в модель
	dbItem := domainToItem(item)
	dbItem.ID = repo.db.nextItemsID

	// инкрементируем счетчик записей
	repo.db.nextItemsID++
	// сохраняем
	repo.db.items[dbItem.ID] = dbItem
	// мапим модель в домен
	domainItem, err := itemToDomain(dbItem)
	if err != nil {
		return domain.Item{}, fmt.Errorf("failed to create domain Item: %w", err)
	}
	return domainItem, nil
}

// GetItems implements services.IItemRepository.
func (repo *ItemRepo) GetItems(_ context.Context, limit, offset, orderid int) ([]domain.Item, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	// извлекаем все ключи из мапы и сортируем их
	keys := make([]int, 0, len(repo.db.items))
	for k := range repo.db.items {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	// выбираем все записи выбранного заказа
	var itemsByOrderID []models.Item
	for i := 0; i < len(keys); i++ {
		item := repo.db.items[keys[i]]
		if item.OrderID == orderid {
			itemsByOrderID = append(itemsByOrderID, item)
		}
	}

	// выбираем записи с нужными ключами
	var items []models.Item
	for i := offset; i < offset+limit && i < len(itemsByOrderID); i++ {
		items = append(items, itemsByOrderID[i])
	}

	// мапим массив моделей в массив доменов
	domainItems := make([]domain.Item, len(items))
	for i, item := range items {
		domainItem, err := itemToDomain(item)
		if err != nil {
			return nil, fmt.Errorf("failed to create domain Order: %w", err)
		}
		domainItems[i] = domainItem
	}
	return domainItems, nil
}

// GetItemByID implements services.IItemRepository.
func (repo *ItemRepo) GetItem(_ context.Context, id int) (domain.Item, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	item, exists := repo.db.items[id]
	if !exists {
		return domain.Item{}, fmt.Errorf("item with id %d - %w", id, domain.ErrNotFound)
	}
	domainItem, err := itemToDomain(item)
	if err != nil {
		return domain.Item{}, fmt.Errorf("failed to create domain Item: %w", err)
	}
	return domainItem, nil
}

// UpdateItem implements services.IItemRepository.
func (repo *ItemRepo) UpdateItem(_ context.Context, item domain.Item) (domain.Item, error) {
	dbItem := domainToItem(item)
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	// проверяем наличие записи
	_, exists := repo.db.items[dbItem.ID]
	if !exists {
		return domain.Item{}, fmt.Errorf("item with id %d - %w", dbItem.ID, domain.ErrNotFound)
	}
	// обновляем запись
	repo.db.items[dbItem.ID] = dbItem
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
	_, exists := repo.db.items[id]
	if !exists {
		return fmt.Errorf("item with id %d - %w", id, domain.ErrNotFound)
	}
	delete(repo.db.items, id)
	return nil
}
