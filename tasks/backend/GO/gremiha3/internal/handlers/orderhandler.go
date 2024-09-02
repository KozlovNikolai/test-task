package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/KozlovNikolai/test-task/internal/model"
	"github.com/KozlovNikolai/test-task/internal/pkg/utils"
	"github.com/KozlovNikolai/test-task/internal/store"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// OrderHandler is ...
type OrderHandler struct {
	logger *zap.Logger
	repoWR store.IRepository
	repoRO store.IRepository
}

// NewOrderHandler is ...
func NewOrderHandler(logger *zap.Logger, repoWR store.IRepository, repoRO store.IRepository) *OrderHandler {
	return &OrderHandler{
		logger: logger,
		repoWR: repoWR,
		repoRO: repoRO,
	}
}

// CreateOrder is ...
// CreateOrderTags		godoc
// @Summary				Добавить заказ.
// @Description			Создание заказа для дальнейшего его заполнения.
// @Param				order body model.AddOrder true "Create Order"
// @Produce				application/json
// @Tags				Order
// @Security     	BearerAuth
// @Success				200 {object} model.Order
// @failure				400 {string} err.Error()
// @failure				500 {string} err.Error()
// @Router				/order [post]
func (oh *OrderHandler) CreateOrder(c *gin.Context) {
	authID, authLogin, authRole := utils.GetLevel(c)
	oh.logger.Debug("принятые логин и роль из токена", zap.Int("id", authID), zap.String("login", authLogin), zap.String("role", authRole))

	// если запрос делает суперпользователь, то ему можно всё
	if authRole == "super" {
		_, err := oh.repoRO.GetUserByLogin(context.TODO(), authLogin)
		if err != nil {
			oh.logger.Error("Error getting user", zap.Error(err))
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		var addOrder model.AddOrder
		// Заполняем структуру addOrder данными из запроса
		if err := c.ShouldBindJSON(&addOrder); err != nil {
			oh.logger.Error("Error binding JSON-addOrder", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var order model.Order
		// Заполняем структуру Order данными из addOrder
		order.CreatedAt = time.Now()
		order.UserID = addOrder.UserID
		order.StateID, _ = oh.repoRO.GetOrderStateByName(context.TODO(), "Created")
		order.TotalAmount = 0

		// Сохраняем товар в БД
		id, err := oh.repoWR.CreateOrder(context.TODO(), order)
		if err != nil {
			oh.logger.Error("Error creating Order", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		order.ID = id
		c.JSON(http.StatusCreated, order)
	} else if authRole == "regular" { // если запрос делает обычный пользователь, то ему можно создавать только собственные данные
		user, err := oh.repoRO.GetUserByLogin(context.TODO(), authLogin)
		if err != nil {
			oh.logger.Error("Error getting user", zap.Error(err))
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		var addOrder model.AddOrder
		// Заполняем структуру addOrder данными из запроса
		if err := c.ShouldBindJSON(&addOrder); err != nil {
			oh.logger.Error("Error binding JSON-addOrder", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var order model.Order
		// Заполняем структуру Order данными из addOrder
		order.CreatedAt = time.Now()
		order.UserID = user.ID
		order.StateID, _ = oh.repoRO.GetOrderStateByName(context.TODO(), "Created")
		order.TotalAmount = 0

		// Сохраняем товар в БД
		id, err := oh.repoWR.CreateOrder(context.TODO(), order)
		if err != nil {
			oh.logger.Error("Error creating Order", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		order.ID = id
		c.JSON(http.StatusCreated, order)
	}
}

// GetOrder is ...
// GetOrderTags 		godoc
// @Summary			Посмотреть товар по его id.
// @Description		Return Order with "id" number.
// @Param			id path int true "Order ID"
// @Tags			Order
// @Security     	BearerAuth
// @Success			200 {object} model.Order
// @failure			404 {string} err.Error()
// @Router			/order/{id} [get]
func (oh *OrderHandler) GetOrder(c *gin.Context) {
	authID, authLogin, authRole := utils.GetLevel(c)
	oh.logger.Debug("принятые логин и роль из токена", zap.Int("id", authID), zap.String("login", authLogin), zap.String("role", authRole))

	// если запрос делает суперпользователь, то ему можно всё
	if authRole == "super" {
		id, _ := strconv.Atoi(c.Param("id"))

		Order, err := oh.repoRO.GetOrderByID(context.TODO(), id)
		if err != nil {
			oh.logger.Error("Error getting Order", zap.Error(err))
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		c.JSON(http.StatusOK, Order)
	} else if authRole == "regular" { // если запрос делает обычный пользователь, и его ID совпадает с user.ID заказа, то ему можно с ним работать
		id, _ := strconv.Atoi(c.Param("id"))

		order, err := oh.repoRO.GetOrderByID(context.TODO(), id)
		if err != nil {
			oh.logger.Error("Error getting Order", zap.Error(err))
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		if order.UserID != authID {
			oh.logger.Error("Это заказ другого пользователя")
			c.JSON(http.StatusNotFound, gin.H{"error": "Это заказ другого пользователя"})
			return
		}
		c.JSON(http.StatusOK, order)

	}

}

// GetOrders is ...
// GetOrdersTags 		godoc
// @Summary			Получить список всех заказов.
// @Description		Return Orders list.
// @Tags			Order
// @Security     	BearerAuth
// @Produce      json
// @Success			200 {object} []model.Order
// @failure			404 {string} err.Error()
// @Router			/orders [get]
func (oh *OrderHandler) GetOrders(c *gin.Context) {
	authID, authLogin, authRole := utils.GetLevel(c)
	oh.logger.Debug("принятые логин и роль из токена", zap.Int("id", authID), zap.String("login", authLogin), zap.String("role", authRole))

	// если запрос делает суперпользователь, то ему можно всё
	if authRole == "super" {
		users, err := oh.repoRO.GetAllOrders(context.TODO())
		if err != nil {
			oh.logger.Error("Error getting users", zap.Error(err))
			c.JSON(http.StatusNotFound, gin.H{"error": "Users not found"})
			return
		}
		c.JSON(http.StatusOK, users)
		return
	}
	oh.logger.Error("forbidden access level.")
	c.JSON(http.StatusForbidden, gin.H{"error": "forbidden access level."})
}
