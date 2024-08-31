package domain

type Item struct {
	id         int
	productID  int
	quantity   int
	totalPrice float64
	orderID    int
}

type NewItemData struct {
	ID         int
	ProductID  int
	Quantity   int
	TotalPrice float64
	OrderID    int
}

func NewItem(data NewItemData) (Item, error) {
	return Item{
		id:         data.ID,
		productID:  data.ProductID,
		quantity:   data.Quantity,
		totalPrice: data.TotalPrice,
		orderID:    data.OrderID,
	}, nil
}

func (i Item) ID() int {
	return i.id
}
func (i Item) ProductID() int {
	return i.productID
}
func (i Item) Quantity() int {
	return i.quantity
}
func (i Item) TotalPrice() float64 {
	return i.totalPrice
}
func (i Item) OrderID() int {
	return i.orderID
}
