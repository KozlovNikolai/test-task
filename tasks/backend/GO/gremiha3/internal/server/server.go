package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/KozlovNikolai/test-task/internal/handlers"
	"github.com/KozlovNikolai/test-task/internal/middlewares"
	"github.com/KozlovNikolai/test-task/internal/pkg/config"
	"github.com/KozlovNikolai/test-task/internal/store"
	"github.com/KozlovNikolai/test-task/internal/store/inmemory"
	"github.com/KozlovNikolai/test-task/internal/store/pgstore"
	"github.com/gin-contrib/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
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

	var repoWR store.IRepository
	var repoRO store.IRepository

	// Выбор репозитория
	switch config.Cfg.RepoType {
	case "postgres":
		poolWR, err := pgxpool.Connect(context.Background(), config.Cfg.StorageWR)
		if err != nil {
			logger.Fatal("Unable to connect to database master", zap.Error(err))
		}
		repoWR = pgstore.NewRepository(poolWR, logger)

		poolRO, err := pgxpool.Connect(context.Background(), config.Cfg.StorageRO)
		if err != nil {
			logger.Fatal("Unable to connect to database replica", zap.Error(err))
		}
		repoRO = pgstore.NewRepository(poolRO, logger)
	case "inmemory":
		repoWR = inmemory.NewRepository(logger)
		repoRO = repoWR
	default:
		logger.Fatal("Invalid repository type")
	}

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
	userHandler := handlers.NewUserHandler(logger, repoWR, repoRO)
	productHandler := handlers.NewProductHandler(logger, repoWR, repoRO)
	providerHandler := handlers.NewProviderHandler(logger, repoWR, repoRO)
	orderHandler := handlers.NewOrderHandler(logger, repoWR, repoRO)
	orderStateHandler := handlers.NewOrderStateHandler(logger, repoWR, repoRO)

	// add swagger
	server.router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Открытые маршруты
	openRouter := server.router.Group("/")

	openRouter.POST("/user/register", userHandler.CreateUser)
	openRouter.POST("/user/login", userHandler.LoginUser)
	openRouter.GET("/products", productHandler.GetProducts)
	openRouter.GET("/product/:id", productHandler.GetProduct)
	openRouter.GET("/providers", providerHandler.GetProviders)
	openRouter.GET("/provider/:id", providerHandler.GetProvider)

	// Закрытые маршруты
	authorized := server.router.Group("/")
	authorized.Use(middlewares.AuthMiddleware())

	// Маршруты для super
	// USER
	authorized.GET("/users", userHandler.GetUsers)
	authorized.GET("/user/:id", userHandler.GetUser)
	authorized.GET("/user/login/:login", userHandler.GetUserByLogin)
	// PRODUCT
	authorized.POST("/product", productHandler.CreateProduct)
	// PROVIDER
	authorized.POST("/provider", providerHandler.CreateProvider)
	// ORDERSTATE
	authorized.POST("/orderstate", orderStateHandler.CreateOrderState)
	authorized.GET("/orderstate/:id", orderStateHandler.GetOrderState)
	authorized.GET("/orderstates", orderStateHandler.GetOrderStates)
	// ORDER
	authorized.POST("/order", orderHandler.CreateOrder)
	authorized.GET("/order/:id", orderHandler.GetOrder)
	authorized.GET("/orders", orderHandler.GetOrders)

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
