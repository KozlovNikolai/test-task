package store

import (
	"context"

	"github.com/KozlovNikolai/test-task/internal/model"
)

type IRepository interface {
	IUserRepository
	IOrderRepository
	IProductRepository
	IProviderRepository
	IOrderStateRepository
}

// IUserRepository is ...
type IUserRepository interface {
	CreateUser(context.Context, model.User) (int, error)
	GetAllUsers(ctx context.Context) ([]model.User, error)
	GetUserByID(context.Context, int) (model.User, error)
	GetUserByLogin(context.Context, string) (model.User, error)
	UpdateUser(context.Context, int) error
	DeleteUser(context.Context, int) error
}

// IProviderRepository is ...
type IProviderRepository interface {
	CreateProvider(context.Context, model.Provider) (int, error)
	GetAllProviders(ctx context.Context) ([]model.Provider, error)
	GetProviderByID(context.Context, int) (model.Provider, error)
	UpdateProvider(context.Context, int) error
	DeleteProvider(context.Context, int) error
}

// IProductRepository is ...
type IProductRepository interface {
	CreateProduct(context.Context, model.Product) (int, error)
	GetAllProducts(ctx context.Context) ([]model.Product, error)
	GetProductByID(context.Context, int) (model.Product, error)
	UpdateProduct(context.Context, int) error
	DeleteProduct(context.Context, int) error
}

// IOrderRepository is ...
type IOrderRepository interface {
	CreateOrder(context.Context, model.Order) (int, error)
	GetAllOrders(ctx context.Context) ([]model.Order, error)
	GetOrderByID(context.Context, int) (model.Order, error)
	UpdateOrder(context.Context, int) error
	DeleteOrder(context.Context, int) error
}

// IOrderStateRepository is ...
type IOrderStateRepository interface {
	CreateOrderState(context.Context, model.OrderState) (int, error)
	GetAllOrderStates(ctx context.Context) ([]model.OrderState, error)
	GetOrderStateByID(context.Context, int) (model.OrderState, error)
	UpdateOrderState(context.Context, int) error
	DeleteOrderState(context.Context, int) error
	GetOrderStateByName(context.Context, string) (int, error)
}
