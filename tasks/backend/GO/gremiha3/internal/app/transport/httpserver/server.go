package httpserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/KozlovNikolai/test-task/internal/app/repository/inmemrepo"
	"github.com/KozlovNikolai/test-task/internal/app/repository/pgrepo"
	"github.com/KozlovNikolai/test-task/internal/app/services"
	"github.com/KozlovNikolai/test-task/internal/middlewares"
	"github.com/KozlovNikolai/test-task/internal/pkg/config"
	"github.com/KozlovNikolai/test-task/internal/pkg/pg"
	"github.com/gin-contrib/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Server is ...
type Server struct {
	router *gin.Engine
	logger *zap.Logger
}

// NewServer is ...
func NewServer() *Server {
	// Инициализация логгера Zap
	//	logger, err := zap.NewProduction()
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	var providerRepo services.IProviderRepository
	var productRepo services.IProductRepository
	var orderStateRepo services.IOrderStateRepository
	var userRepo services.IUserRepository
	var itemRepo services.IItemRepository
	var orderRepo services.IOrderRepository
	// Выбор репозитория
	switch config.Cfg.RepoType {
	case "postgres":
		pgDB, err := pg.Dial(config.Cfg.StorageWR, config.Cfg.StorageRO)
		if err != nil {
			logger.Error("pg.Dial failed: %w", zap.Error(err))
		}
		providerRepo = pgrepo.NewProviderRepo(pgDB)
		productRepo = pgrepo.NewProductRepo(pgDB)
		orderStateRepo = pgrepo.NewOrderStateRepo(pgDB)
		userRepo = pgrepo.NewUserRepo(pgDB)
		itemRepo = pgrepo.NewItemRepo(pgDB)
		orderRepo = pgrepo.NewOrderRepo(pgDB)

	case "inmemory":
		providerRepo = inmemrepo.NewProviderRepo()
		productRepo = inmemrepo.NewProductRepo()
		orderStateRepo = inmemrepo.NewOrderStateRepo()
		itemRepo = inmemrepo.NewItemRepo()
		userRepo = inmemrepo.NewUserRepo()
		orderRepo = inmemrepo.NewOrderRepo()

	default:
		logger.Fatal("Invalid repository type")
	}
	// сщздаем сервисы
	providerService := services.NewProviderService(providerRepo)
	productService := services.NewProductService(productRepo)
	userService := services.NewUserService(userRepo)
	orderService := services.NewOrderService(orderRepo)
	itemService := services.NewItemService(itemRepo)
	orderStateService := services.NewOrderStateService(orderStateRepo)
	tokenService := services.NewTokenService(config.Cfg.TokenTimeDuration)
	// создаем сервер
	httpServer := NewHttpServer(
		providerService,
		productService,
		userService,
		orderService,
		itemService,
		orderStateService,
		tokenService,
	)

	// Создание сервера
	server := &Server{
		router: gin.Default(),
		logger: logger,
	}

	// Middleware
	server.router.Use(middlewares.LoggerMiddleware(logger))
	server.router.Use(middlewares.RequestIDMiddleware())

	// CORS
	server.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://localhost:8443", "https://127.0.0.1:8443"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// Инициализация обработчиков
	// userHandler := handlers.NewUserHandler(logger, repoWR, repoRO)
	// productHandler := handlers.NewProductHandler(logger, repoWR, repoRO)
	// providerHandler := handlers.NewProviderHandler(logger, repoWR, repoRO)
	// orderHandler := handlers.NewOrderHandler(logger, repoWR, repoRO)
	// orderStateHandler := handlers.NewOrderStateHandler(logger, repoWR, repoRO)

	// add swagger
	server.router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Открытые маршруты
	openRouter := server.router.Group("/")
	openRouter.POST("/user/register", httpServer.CreateUser)
	openRouter.POST("/user/login", httpServer.LoginUser)
	openRouter.GET("/products", httpServer.GetProducts)
	openRouter.GET("/product/:id", httpServer.GetProduct)
	openRouter.GET("/providers", httpServer.GetProviders)
	openRouter.GET("/provider/:id", httpServer.GetProvider)

	// Закрытые маршруты
	authorized := server.router.Group("/")
	authorized.Use(middlewares.AuthMiddleware())

	// Маршруты для super
	// USER
	authorized.GET("/users", httpServer.GetUsers)
	authorized.GET("/user/:id", httpServer.GetUser)
	authorized.GET("/user/login/:login", httpServer.GetUserByLogin)
	// PRODUCT
	authorized.POST("/product", httpServer.CreateProduct)
	// PROVIDER
	authorized.POST("/provider", httpServer.CreateProvider)
	// ORDERSTATE
	authorized.POST("/orderstate", httpServer.CreateOrderState)
	authorized.GET("/orderstate/:id", httpServer.GetOrderState)
	authorized.GET("/orderstates", httpServer.GetOrderStates)
	// ORDER
	authorized.POST("/order", httpServer.CreateOrder)
	authorized.GET("/order/:id", httpServer.GetOrder)
	authorized.GET("/orders", httpServer.GetOrders)

	return server
}

// Run is ...
func (s *Server) Run() {
	defer func() {
		_ = s.logger.Sync() // flushes buffer, if any
	}()
	// Настройка сервера с таймаутами
	server := &http.Server{
		Addr:         config.Cfg.Address,
		Handler:      s.router,
		ReadTimeout:  config.Cfg.Timeout,
		WriteTimeout: config.Cfg.Timeout,
		IdleTimeout:  config.Cfg.IdleTimout,
	}

	if err := server.ListenAndServeTLS(config.CertFile, config.KeyFile); err != nil && err != http.ErrServerClosed {
		s.logger.Fatal(fmt.Sprintf("Could not listen on %s", config.Cfg.Address), zap.Error(err))
	}
}
