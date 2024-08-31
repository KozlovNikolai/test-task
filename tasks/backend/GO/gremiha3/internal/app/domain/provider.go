package domain

// Provider is a domain Provider.
type Provider struct {
	id     int
	name   string
	origin string
}

// NewProviderData is a DTO Provider
type NewProviderData struct {
	ID     int
	Name   string
	Origin string
}

// NewProvider creates a new Provider.
func NewProvider(data NewProviderData) (Provider, error) {
	return Provider{
		id:     data.ID,
		name:   data.Name,
		origin: data.Origin,
	}, nil
}

// ID returns the Provider ID.
func (p Provider) ID() int {
	return p.id
}

// Title returns the Provider title.
func (p Provider) Name() string {
	return p.name
}

// Year returns the Provider year.
func (p Provider) Origin() string {
	return p.origin
}
