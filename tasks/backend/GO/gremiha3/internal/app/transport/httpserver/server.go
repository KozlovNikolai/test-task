package httpserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
		imDB := inmemrepo.NewInMemRepo()
		providerRepo = inmemrepo.NewProviderRepo(imDB)
		productRepo = inmemrepo.NewProductRepo(imDB)
		orderStateRepo = inmemrepo.NewOrderStateRepo(imDB)
		itemRepo = inmemrepo.NewItemRepo(imDB)
		userRepo = inmemrepo.NewUserRepo(imDB)
		orderRepo = inmemrepo.NewOrderRepo(imDB)

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

	// add swagger
	server.router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	root := server.router.Group("/")
	root.POST("signup", httpServer.SignUp)
	root.POST("signin", httpServer.SignIn)

	root.GET("user", httpServer.GetUser)
	root.GET("users", httpServer.GetUsers)

	root.POST("provider", httpServer.CreateProvider)
	root.GET("provider", httpServer.GetProvider)
	root.GET("providers", httpServer.GetProviders)

	root.POST("product", httpServer.CreateProduct)
	root.GET("product", httpServer.GetProduct)
	root.GET("products", httpServer.GetProducts)

	root.POST("orderstate", httpServer.CreateOrderState)
	root.GET("orderstate", httpServer.GetOrderState)
	root.GET("orderstates", httpServer.GetOrderStates)

	root.POST("order", httpServer.CreateOrder)
	root.GET("order", httpServer.GetOrder)
	root.GET("orders", httpServer.GetOrders)

	root.POST("item", httpServer.CreateItem)
	root.GET("item", httpServer.GetItem)
	root.GET("items", httpServer.GetItems)
	//#################################################################################

	// Закрытые маршруты
	authorized := server.router.Group("/")
	authorized.Use(middlewares.AuthMiddleware())

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
	// listen to OS signals and gracefully shutdown HTTP server
	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
		close(stopped)
	}()
	if err := server.ListenAndServeTLS(config.CertFile, config.KeyFile); err != nil && err != http.ErrServerClosed {
		s.logger.Fatal(fmt.Sprintf("Could not listen on %s", config.Cfg.Address), zap.Error(err))
	}
	<-stopped

	log.Printf("Bye!")
}
