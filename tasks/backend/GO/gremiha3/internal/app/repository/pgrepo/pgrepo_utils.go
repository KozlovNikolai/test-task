package pgrepo

import (
	"github.com/KozlovNikolai/test-task/internal/app/domain"
	"github.com/KozlovNikolai/test-task/internal/app/repository/models"
)

func domainToProvider(provider domain.Provider) models.Provider {
	return models.Provider{
		ID:     provider.ID(),
		Name:   provider.Name(),
		Origin: provider.Origin(),
	}
}

func providerToDomain(provider models.Provider) (domain.Provider, error) {
	return domain.NewProvider(domain.NewProviderData{
		ID:     provider.ID,
		Name:   provider.Name,
		Origin: provider.Origin,
	})
}

func domainToProduct(product domain.Product) models.Product {
	return models.Product{
		ID:         product.ID(),
		Name:       product.Name(),
		ProviderID: product.ProviderID(),
		Price:      product.Price(),
		Stock:      product.Stock(),
	}
}

func productToDomain(product models.Product) (domain.Product, error) {
	return domain.NewProduct(domain.NewProductData{
		ID:         product.ID,
		Name:       product.Name,
		ProviderID: product.ProviderID,
		Price:      product.Price,
		Stock:      product.Stock,
	})
}
func domainToOrderState(orderState domain.OrderState) models.OrderState {
	return models.OrderState{
		ID:   orderState.ID(),
		Name: orderState.Name(),
	}
}

func orderStateToDomain(orderState models.OrderState) (domain.OrderState, error) {
	return domain.NewOrderState(domain.NewOrderStateData{
		ID:   orderState.ID,
		Name: orderState.Name,
	})
}

func domainToUser(user domain.User) models.User {
	return models.User{
		ID:       user.ID(),
		Login:    user.Login(),
		Password: user.Password(),
		Role:     user.Role(),
		Token:    user.Token(),
	}
}

func userToDomain(user models.User) (domain.User, error) {
	return domain.NewUser(domain.NewUserData{
		ID:       user.ID,
		Login:    user.Login,
		Password: user.Password,
		Role:     user.Role,
		Token:    user.Token,
	})
}

func domainToOrder(order domain.Order) models.Order {
	return models.Order{
		ID:          order.ID(),
		UserID:      order.UserID(),
		StateID:     order.StateID(),
		TotalAmount: order.TotalAmount(),
		CreatedAt:   order.CreatedAt(),
	}
}

func orderToDomain(order models.Order) (domain.Order, error) {
	return domain.NewOrder(domain.NewOrderData{
		ID:          order.ID,
		UserID:      order.UserID,
		StateID:     order.StateID,
		TotalAmount: order.TotalAmount,
		CreatedAt:   order.CreatedAt,
	})
}

func domainToItem(item domain.Item) models.Item {
	return models.Item{
		ID:         item.ID(),
		ProductID:  item.ProductID(),
		Quantity:   item.Quantity(),
		TotalPrice: item.TotalPrice(),
		OrderID:    item.OrderID(),
	}
}

func itemToDomain(item models.Item) (domain.Item, error) {
	return domain.NewItem(domain.NewItemData{
		ID:         item.ID,
		ProductID:  item.ProductID,
		Quantity:   item.Quantity,
		TotalPrice: item.TotalPrice,
		OrderID:    item.OrderID,
	})
}
