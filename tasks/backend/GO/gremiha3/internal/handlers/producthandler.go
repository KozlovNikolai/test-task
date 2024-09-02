package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/KozlovNikolai/test-task/internal/model"
	"github.com/KozlovNikolai/test-task/internal/pkg/utils"
	"github.com/KozlovNikolai/test-task/internal/store"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ProductHandler is ...
type ProductHandler struct {
	logger *zap.Logger
	repoWR store.IProductRepository
	repoRO store.IProductRepository
}

// NewProductHandler is ...
func NewProductHandler(logger *zap.Logger, repoWR store.IProductRepository, repoRO store.IProductRepository) *ProductHandler {
	return &ProductHandler{
		logger: logger,
		repoWR: repoWR,
		repoRO: repoRO,
	}
}

// CreateProduct is ...
// CreateProductTags		godoc
// @Summary				Добавить товар.
// @Description			Save register data of user in Repo.
// @Param				product body model.AddProduct true "Create product"
// @Produce				application/json
// @Tags				Product
// @Security     	BearerAuth
// @Success				200 {object} model.Product
// @failure				400 {string} err.Error()
// @failure				500 {string} err.Error()
// @Router				/product [post]
func (ph *ProductHandler) CreateProduct(c *gin.Context) {
	authID, authLogin, authRole := utils.GetLevel(c)
	ph.logger.Debug("принятые логин и роль из токена", zap.Int("id", authID), zap.String("login", authLogin), zap.String("role", authRole))
	// если запрос делает суперпользователь, то ему можно всё
	if authRole == "super" {
		var addProduct model.AddProduct
		// Заполняем структуру addProduct данными из запроса
		if err := c.ShouldBindJSON(&addProduct); err != nil {
			ph.logger.Error("Error binding JSON-addProduct", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var product model.Product
		// Заполняем структуру Product данными из addProduct
		product.Name = addProduct.Name
		product.Price = addProduct.Price
		product.ProviderID = addProduct.ProviderID
		product.Stock = addProduct.Stock

		// Валидация данных товара
		if err := product.Validate(); err != nil {
			ph.logger.Error("Error creating product", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Сохраняем товар в БД
		id, err := ph.repoWR.CreateProduct(context.TODO(), product)
		if err != nil {
			ph.logger.Error("Error creating product", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		product.ID = id
		c.JSON(http.StatusCreated, product)
	} else if authRole == "regular" { // если запрос делает обычный пользователь, то не разрешаем:
		ph.logger.Error("forbidden access level.")
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden access level."})
	}
}

// GetProduct is ...
// GetProductTags 		godoc
// @Summary			Посмотреть товар по его id.
// @Description		Return product with "id" number.
// @Param			product_id path int true "Product ID"
// @Tags			Product
// @Success			200 {object} model.Product
// @failure			404 {string} err.Error()
// @Router			/product/{product_id} [get]
func (ph *ProductHandler) GetProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	product, err := ph.repoRO.GetProductByID(context.TODO(), id)
	if err != nil {
		ph.logger.Error("Error getting product", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

// GetProducts is ...
// GetProductsTags 		godoc
// @Summary			Получить список всех товаров.
// @Description		Return products list.
// @Tags			Product
// @Produce      json
// @Success			200 {object} []model.Product
// @failure			404 {string} err.Error()
// @Router			/products [get]
func (ph *ProductHandler) GetProducts(c *gin.Context) {

	users, err := ph.repoRO.GetAllProducts(context.TODO())
	if err != nil {
		ph.logger.Error("Error getting users", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Users not found"})
		return
	}
	c.JSON(http.StatusOK, users)
}
