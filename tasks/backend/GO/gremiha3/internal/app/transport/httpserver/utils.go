package httpserver

import (
	"context"

	"github.com/KozlovNikolai/test-task/internal/app/domain"
)

// #########################################################
func toResponseProvider(provider domain.Provider) ProviderResponse {
	return ProviderResponse{
		ID:     provider.ID(),
		Name:   provider.Name(),
		Origin: provider.Origin(),
	}
}

func toDomainProvider(provider ProviderRequest) (domain.Provider, error) {
	return domain.NewProvider(domain.NewProviderData{
		Name:   provider.Name,
		Origin: provider.Origin,
	})
}

// #########################################################
func toResponseProduct(product domain.Product) ProductResponse {
	return ProductResponse{
		ID:         product.ID(),
		Name:       product.Name(),
		ProviderID: product.ProviderID(),
		Price:      product.Price(),
		Stock:      product.Stock(),
	}
}

func toDomainProduct(product ProductRequest) (domain.Product, error) {
	return domain.NewProduct(domain.NewProductData{
		Name:       product.Name,
		ProviderID: product.ProviderID,
		Price:      product.Price,
		Stock:      product.Stock,
	})
}

// #########################################################
func toResponseOrderState(orderState domain.OrderState) OrderStateResponse {
	return OrderStateResponse{
		ID:   orderState.ID(),
		Name: orderState.Name(),
	}
}

func toDomainOrderState(orderState OrderStateRequest) (domain.OrderState, error) {
	return domain.NewOrderState(domain.NewOrderStateData{
		Name: orderState.Name,
	})
}

// #########################################################
func toResponseOrder(order domain.Order) OrderResponse {
	return OrderResponse{
		ID:          order.ID(),
		UserID:      order.UserID(),
		StateID:     order.StateID(),
		TotalAmount: order.TotalAmount(),
		CreatedAt:   order.CreatedAt(),
	}
}

func toDomainOrder(order OrderRequest) (domain.Order, error) {
	return domain.NewOrder(domain.NewOrderData{
		UserID: order.UserID,
	})
}

// #########################################################
func toResponseItem(item domain.Item) ItemResponse {
	return ItemResponse{
		ID:         item.ID(),
		ProductID:  item.ProductID(),
		Quantity:   item.Quantity(),
		TotalPrice: item.TotalPrice(),
		OrderID:    item.OrderID(),
	}
}

func toDomainItem(item ItemRequest) (domain.Item, error) {
	return domain.NewItem(domain.NewItemData{
		ProductID: item.ProductID,
		Quantity:  item.Quantity,
		OrderID:   item.OrderID,
	})
}

// #########################################################
func toResponseUser(user domain.User) UserResponse {
	return UserResponse{
		ID:       user.ID(),
		Login:    user.Login(),
		Password: user.Password(),
		Role:     user.Role(),
		Token:    user.Token(),
	}
}

func toDomainUser(user UserRequest) (domain.User, error) {
	return domain.NewUser(domain.NewUserData{
		Login:    user.Login,
		Password: user.Password,
	})
}

// #########################################################
func getUserFromContext(ctx context.Context) (domain.User, error) {
	contextUser := ctx.Value(ContextUserKey)
	if contextUser == nil {
		return domain.User{}, domain.ErrNoUserInContext
	}
	user, ok := contextUser.(domain.User)
	if !ok {
		return domain.User{}, domain.ErrNoUserInContext
	}
	return user, nil
}
