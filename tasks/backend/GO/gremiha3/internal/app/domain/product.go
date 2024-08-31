package domain

// Product is a domain Product
type Product struct {
	id         int
	name       string
	providerID int
	price      float64
	stock      int
}

// NewProductData is a DTO Product
type NewProductData struct {
	ID         int
	Name       string
	ProviderID int
	Price      float64
	Stock      int
}

// NewProduct is ...
func NewProduct(data NewProductData) (Product, error) {
	return Product{
		id:         data.ID,
		name:       data.Name,
		providerID: data.ProviderID,
		price:      data.Price,
		stock:      data.Stock,
	}, nil
}

func (p Product) ID() int {
	return p.id
}
func (p Product) Name() string {
	return p.name
}
func (p Product) ProviderID() int {
	return p.providerID
}
func (p Product) Price() float64 {
	return p.price
}
func (p Product) Stock() int {
	return p.stock
}
