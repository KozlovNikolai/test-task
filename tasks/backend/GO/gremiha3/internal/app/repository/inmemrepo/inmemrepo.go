package inmemrepo

import (
	"github.com/KozlovNikolai/test-task/internal/app/repository/models"
)

type inMemStore struct {
	users             map[int]models.User
	nextUsersID       int
	providers         map[int]models.Provider
	nextProvidersID   int
	products          map[int]models.Product
	nextProductsID    int
	orders            map[int]models.Order
	nextOrdersID      int
	orderStates       map[int]models.OrderState
	nextOrderStatesID int
	items             map[int]models.Item
	nextItemsID       int
}

func NewInMemRepo() *inMemStore {
	return &inMemStore{
		users: map[int]models.User{
			1: {ID: 1, Login: "cmd@cmd.ru", Password: "pass", Role: "regular"},
			2: {ID: 2, Login: "cmd@cmd.org", Password: "pass", Role: "regular"},
			3: {ID: 3, Login: "cmd@cmd.com", Password: "pass", Role: "regular"},
		},
		nextUsersID: 4,
		providers: map[int]models.Provider{
			1: {ID: 1, Name: "Salvatore INC", Origin: "Mexico"},
			2: {ID: 2, Name: "Huaway", Origin: "China"},
			3: {ID: 3, Name: "Roga I Kopyta", Origin: "Russia"},
		},
		nextProvidersID: 4,
		products: map[int]models.Product{
			1:  {ID: 1, Name: "avto", ProviderID: 1, Price: 123.12, Stock: 6745},
			2:  {ID: 2, Name: "submarine", ProviderID: 2, Price: 7356.67, Stock: 12},
			3:  {ID: 3, Name: "phone", ProviderID: 1, Price: 876.23, Stock: 654},
			4:  {ID: 4, Name: "хлеб", ProviderID: 3, Price: 987.23, Stock: 908},
			5:  {ID: 5, Name: "стол", ProviderID: 2, Price: 524.45, Stock: 2765},
			6:  {ID: 6, Name: "стул", ProviderID: 1, Price: 957.23, Stock: 987},
			7:  {ID: 7, Name: "кастрюля", ProviderID: 3, Price: 915.34, Stock: 1987},
			8:  {ID: 8, Name: "боинг", ProviderID: 2, Price: 9756.34, Stock: 9043},
			9:  {ID: 9, Name: "колесо", ProviderID: 1, Price: 809.21, Stock: 126},
			10: {ID: 10, Name: "ручка", ProviderID: 3, Price: 954.09, Stock: 4},
			11: {ID: 11, Name: "ноутбук", ProviderID: 2, Price: 865.03, Stock: 8765},
			12: {ID: 12, Name: "монитор", ProviderID: 1, Price: 136.34, Stock: 908756},
		},
		nextProductsID: 13,
		orders: map[int]models.Order{
			1: {ID: 1, UserID: 1, StateID: 1},
			2: {ID: 2, UserID: 1, StateID: 1},
			3: {ID: 3, UserID: 3, StateID: 1},
		},
		nextOrdersID: 4,
		orderStates: map[int]models.OrderState{
			1: {ID: 1, Name: "Created"},
			2: {ID: 2, Name: "In progress"},
			3: {ID: 3, Name: "Delivery"},
		},
		nextOrderStatesID: 4,
		items: map[int]models.Item{
			1:  {ID: 1, ProductID: 1, Quantity: 23, OrderID: 1},
			2:  {ID: 2, ProductID: 2, Quantity: 5, OrderID: 2},
			3:  {ID: 3, ProductID: 3, Quantity: 7, OrderID: 3},
			4:  {ID: 4, ProductID: 4, Quantity: 9, OrderID: 1},
			5:  {ID: 5, ProductID: 5, Quantity: 32, OrderID: 2},
			6:  {ID: 6, ProductID: 6, Quantity: 65, OrderID: 3},
			7:  {ID: 7, ProductID: 7, Quantity: 2, OrderID: 1},
			8:  {ID: 8, ProductID: 8, Quantity: 1, OrderID: 2},
			9:  {ID: 9, ProductID: 9, Quantity: 76, OrderID: 3},
			10: {ID: 10, ProductID: 10, Quantity: 28, OrderID: 1},
			11: {ID: 11, ProductID: 11, Quantity: 90, OrderID: 2},
			12: {ID: 12, ProductID: 12, Quantity: 2000, OrderID: 3},
			13: {ID: 13, ProductID: 1, Quantity: 23, OrderID: 1},
			14: {ID: 14, ProductID: 2, Quantity: 6, OrderID: 2},
			15: {ID: 15, ProductID: 3, Quantity: 8, OrderID: 3},
			16: {ID: 16, ProductID: 4, Quantity: 234, OrderID: 1},
			17: {ID: 17, ProductID: 5, Quantity: 654, OrderID: 2},
			18: {ID: 18, ProductID: 6, Quantity: 186, OrderID: 3},
			19: {ID: 19, ProductID: 7, Quantity: 908, OrderID: 1},
			20: {ID: 20, ProductID: 8, Quantity: 34, OrderID: 2},
		},
		nextItemsID: 21,
	}
}
