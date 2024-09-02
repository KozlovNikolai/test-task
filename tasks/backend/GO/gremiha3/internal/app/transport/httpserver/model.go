package httpserver

import (
	"fmt"
	"time"

	"github.com/KozlovNikolai/test-task/internal/app/domain"
)

const passwordLength = 6

// ProviderRequest is ...
type ProviderRequest struct {
	Name   string `json:"name" db:"name" example:"Microsoft"`
	Origin string `json:"origin" db:"origin" example:"Vietnam"`
}

func (p *ProviderRequest) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("%w: name", domain.ErrRequired)
	}
	if p.Origin == "" {
		return fmt.Errorf("%w: origin", domain.ErrRequired)
	}
	return nil
}

type ProviderResponse struct {
	ID     int    `json:"id" db:"id"`
	Name   string `json:"name" db:"name"`
	Origin string `json:"origin" db:"origin"`
}

// #######################################################################################3
// ProductRequest is ...
type ProductRequest struct {
	Name       string  `json:"name" db:"name" example:"синхрофазотрон"`
	ProviderID int     `json:"provider_id" db:"provider_id" example:"1"`
	Price      float64 `json:"price" db:"price" example:"1245.65"`
	Stock      int     `json:"stock" db:"stock" example:"435"`
}

func (p *ProductRequest) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("%w: name", domain.ErrRequired)
	}
	if p.ProviderID <= 0 {
		return fmt.Errorf("%w: provider_id", domain.ErrRequired)
	}
	if p.Price < 0 {
		return fmt.Errorf("%w: price", domain.ErrNegative)
	}
	if p.Stock < 0 {
		return fmt.Errorf("%w: stock", domain.ErrNegative)
	}
	return nil
}

type ProductResponse struct {
	ID         int     `json:"id" db:"id"`
	Name       string  `json:"name" db:"name"`
	ProviderID int     `json:"provider_id" db:"provider_id"`
	Price      float64 `json:"price" db:"price"`
	Stock      int     `json:"stock" db:"stock"`
}

// #######################################################################################
// OrderStateRequest is ...
type OrderStateRequest struct {
	Name string `json:"name" db:"name" example:"в обработке"`
}

func (o *OrderStateRequest) Validate() error {
	if o.Name == "" {
		return fmt.Errorf("%w: name", domain.ErrRequired)
	}
	return nil
}

type OrderStateResponse struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

// #######################################################################################3
// ItemRequest is ...
type ItemRequest struct {
	ProductID int `json:"product_id" db:"product_id" example:"1"`
	Quantity  int `json:"quantity" db:"quantity" example:"3"`
	OrderID   int `json:"order_id" db:"order_id" example:"1"`
}

func (i *ItemRequest) Validate() error {
	if i.ProductID <= 0 {
		return fmt.Errorf("%w: product_id", domain.ErrInvalidValue)
	}
	if i.Quantity <= 0 {
		return fmt.Errorf("%w: quantity", domain.ErrNegative)
	}
	if i.OrderID <= 0 {
		return fmt.Errorf("%w: order_id", domain.ErrInvalidValue)
	}
	return nil
}

type ItemResponse struct {
	ID         int     `json:"id" db:"id"`
	ProductID  int     `json:"product_id" db:"product_id"`
	Quantity   int     `json:"quantity" db:"quantity"`
	TotalPrice float64 `json:"total_price" db:"total_price"`
	OrderID    int     `json:"order_id" db:"order_id"`
}

// #######################################################################################3
// OrderRequest is ...
type OrderRequest struct {
	UserID int `json:"user_id" db:"user_id" example:"1"`
}

func (o *OrderRequest) Validate() error {
	if o.UserID <= 0 {
		return fmt.Errorf("%w: user_id", domain.ErrInvalidValue)
	}
	return nil
}

type OrderResponse struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	StateID     int       `json:"state_id" db:"state_id"`
	TotalAmount float64   `json:"total_amount" db:"total_amount"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// #######################################################################################3
// UserRequest is ...
type UserRequest struct {
	Login    string `json:"login" db:"login" example:"cmd@cmd.ru"`
	Password string `json:"password" db:"password" example:"123456"`
}

func (u *UserRequest) Validate() error {
	if u.Login == "" {
		return fmt.Errorf("%w: login", domain.ErrRequired)
	}
	if !isValidEmail(u.Login) {
		return fmt.Errorf("%w: login", domain.ErrInvalidFormatEmail)
	}
	if u.Password == "" {
		return fmt.Errorf("%w: password", domain.ErrRequired)
	}
	if len(u.Password) < passwordLength {
		return fmt.Errorf("%w: password, must be greater then %d", domain.ErrInvalidLength, passwordLength)
	}
	return nil
}

type UserResponse struct {
	ID       int    `json:"id" db:"id"`
	Login    string `json:"login" db:"login"`
	Password string `json:"password" db:"password"`
	Role     string `json:"role" db:"role"`
	Token    string `json:"token" db:"token"`
}
