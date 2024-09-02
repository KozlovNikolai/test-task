package model

// OrderState is ...
type OrderState struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

// AddOrderState is ...
type AddOrderState struct {
	Name string `json:"name" db:"name" example:"Закрыт"`
}
