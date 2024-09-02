package inmemory

import (
	"sync"

	"github.com/KozlovNikolai/test-task/internal/model"
	"go.uber.org/zap"
)

// Repository is ...
type Repository struct {
	logger          *zap.Logger
	users           map[int]model.User
	orders          map[int]model.Order
	products        map[int]model.Product
	providers       map[int]model.Provider
	orderState      map[int]model.OrderState
	orderItem       map[int]model.Item
	nextUsersID     int
	nextOrdersID    int
	nextProductsID  int
	nextProvidersID int
	nextOrderState  int
	nextOrderItem   int
	mutex           sync.Mutex
}

// NewRepository is ...
func NewRepository(logger *zap.Logger) *Repository {
	logger = logger.Named("repo")
	return &Repository{
		logger:    logger,
		users:     make(map[int]model.User),
		orders:    make(map[int]model.Order),
		products:  make(map[int]model.Product),
		providers: make(map[int]model.Provider),
		orderState: map[int]model.OrderState{
			0: {ID: 0, Name: "Created"},
			1: {ID: 1, Name: "In progress"},
			2: {ID: 2, Name: "Delivery"},
		},
		orderItem:       make(map[int]model.Item),
		nextUsersID:     1,
		nextOrdersID:    1,
		nextProductsID:  1,
		nextProvidersID: 1,
		nextOrderState:  3,
		nextOrderItem:   1,
	}
}
