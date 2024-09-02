package services

import (
	"context"

	"github.com/KozlovNikolai/test-task/internal/app/domain"
)

// IUserRepository is ...
type IUserRepository interface {
	CreateUser(context.Context, domain.User) (domain.User, error)
	GetUsers(context.Context, int, int) ([]domain.User, error)
	GetUserByID(context.Context, int) (domain.User, error)
	GetUserByLogin(context.Context, string) (domain.User, error)
	UpdateUser(context.Context, domain.User) (domain.User, error)
	DeleteUser(context.Context, int) error
}

// IProviderRepository is ...
type IProviderRepository interface {
	CreateProvider(context.Context, domain.Provider) (domain.Provider, error)
	GetProviders(context.Context, int, int) ([]domain.Provider, error)
	GetProvider(context.Context, int) (domain.Provider, error)
	UpdateProvider(context.Context, domain.Provider) (domain.Provider, error)
	DeleteProvider(context.Context, int) error
}

// IProductRepository is ...
type IProductRepository interface {
	CreateProduct(context.Context, domain.Product) (domain.Product, error)
	GetProducts(context.Context, int, int) ([]domain.Product, error)
	GetProduct(context.Context, int) (domain.Product, error)
	UpdateProduct(context.Context, domain.Product) (domain.Product, error)
	DeleteProduct(context.Context, int) error
}

// IOrderRepository is ...
type IOrderRepository interface {
	CreateOrder(context.Context, domain.Order) (domain.Order, error)
	GetOrders(context.Context, int, int, int) ([]domain.Order, error)
	GetOrder(context.Context, int) (domain.Order, error)
	GetOrdersByUserID(context.Context, int, int, int) ([]domain.Order, error)
	UpdateOrder(context.Context, domain.Order) (domain.Order, error)
	DeleteOrder(context.Context, int) error
}

// IOrderStateRepository is ...
type IOrderStateRepository interface {
	CreateOrderState(context.Context, domain.OrderState) (domain.OrderState, error)
	GetOrderStates(context.Context, int, int) ([]domain.OrderState, error)
	GetOrderState(context.Context, int) (domain.OrderState, error)
	UpdateOrderState(context.Context, domain.OrderState) (domain.OrderState, error)
	DeleteOrderState(context.Context, int) error
}

// IOrderRepository is ...
type IItemRepository interface {
	CreateItem(context.Context, domain.Item) (domain.Item, error)
	GetItems(context.Context, int, int, int) ([]domain.Item, error)
	GetItem(context.Context, int) (domain.Item, error)
	UpdateItem(context.Context, domain.Item) (domain.Item, error)
	DeleteItem(context.Context, int) error
}
