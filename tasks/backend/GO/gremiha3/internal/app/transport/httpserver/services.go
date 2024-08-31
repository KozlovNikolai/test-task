package httpserver

// HttpServer is a HTTP server for ports
type HttpServer struct {
	providerService   IProviderService
	productService    IProductService
	userService       IUserService
	orderService      IOrderService
	itemService       IItemService
	orderStateService IOrderStateService
	tokenService      ITokenService
}

// NewHttpServer creates a new HTTP server for ports
func NewHttpServer(
	providerService IProviderService,
	productService IProductService,
	userService IUserService,
	orderService IOrderService,
	itemService IItemService,
	orderStateService IOrderStateService,
	tokenService ITokenService,
) HttpServer {
	return HttpServer{
		providerService:   providerService,
		productService:    productService,
		userService:       userService,
		orderService:      orderService,
		itemService:       itemService,
		orderStateService: orderStateService,
		tokenService:      tokenService,
	}
}
