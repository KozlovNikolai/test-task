package domain

// OrderState is a domain OrderState
type OrderState struct {
	id   int
	name string
}

// OrderState is a domain OrderState
type NewOrderStateData struct {
	ID   int
	Name string
}

func NewOrderState(data NewOrderStateData) (OrderState, error) {
	return OrderState{
		id:   data.ID,
		name: data.Name,
	}, nil
}

func (o OrderState) ID() int {
	return o.id
}
func (o OrderState) Name() string {
	return o.name
}
