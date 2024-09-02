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

// OrderHandler is ...
type OrderStateHandler struct {
	logger *zap.Logger
	repoWR store.IOrderStateRepository
	repoRO store.IOrderStateRepository
}

// NewOrderHandler is ...
func NewOrderStateHandler(logger *zap.Logger, repoWR store.IOrderStateRepository, repoRO store.IOrderStateRepository) *OrderStateHandler {
	return &OrderStateHandler{
		logger: logger,
		repoWR: repoWR,
		repoRO: repoRO,
	}
}

// CreateOrderState is ...
// CreateOrderStateTags		godoc
// @Summary				Добавить тип статуса заказа.
// @Description			Создание типа статуса заказа.
// @Param				orderState body model.AddOrderState true "Create Order State type"
// @Produce				application/json
// @Tags				OrderState
// @Security     	BearerAuth
// @Success				200 {object} model.OrderState
// @failure				400 {string} err.Error()
// @failure				500 {string} err.Error()
// @Router				/orderstate [post]
func (osh *OrderStateHandler) CreateOrderState(c *gin.Context) {
	authID, authLogin, authRole := utils.GetLevel(c)
	osh.logger.Debug("принятые логин и роль из токена", zap.Int("id", authID), zap.String("login", authLogin), zap.String("role", authRole))

	// если запрос делает суперпользователь, то ему можно всё
	if authRole == "super" {
		var addOrderState model.AddOrderState
		// Заполняем структуру addOrderState данными из запроса
		if err := c.ShouldBindJSON(&addOrderState); err != nil {
			osh.logger.Error("Error binding JSON-addOrderState", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var orderState model.OrderState
		// Заполняем структуру OrderState данными из addOrderState
		orderState.Name = addOrderState.Name

		// Сохраняем имя статуса заказа в БД
		id, err := osh.repoWR.CreateOrderState(context.TODO(), orderState)
		if err != nil {
			osh.logger.Error("Error creating OrderState", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		orderState.ID = id
		c.JSON(http.StatusCreated, orderState)
	} else if authRole == "regular" { // если запрос делает обычный пользователь, то не разрешаем:
		osh.logger.Error("forbidden access level.")
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden access level."})
	}
}

// GetOrderState is ...
// GetOrderStateTags 		godoc
// @Summary			Посмотреть тип статуса по его id.
// @Description		Return OrderState with "id" number.
// @Param			id path int true "OrderState ID"
// @Tags			OrderState
// @Security     	BearerAuth
// @Success			200 {object} model.OrderState
// @failure			404 {string} err.Error()
// @Router			/orderstate/{id} [get]
func (osh *OrderStateHandler) GetOrderState(c *gin.Context) {
	authID, authLogin, authRole := utils.GetLevel(c)
	osh.logger.Debug("принятые логин и роль из токена", zap.Int("id", authID), zap.String("login", authLogin), zap.String("role", authRole))

	// если запрос делает суперпользователь, то ему можно всё
	if authRole == "super" {
		id, _ := strconv.Atoi(c.Param("id"))

		orderState, err := osh.repoRO.GetOrderStateByID(context.TODO(), id)
		if err != nil {
			osh.logger.Error("Error getting OrderState", zap.Error(err))
			c.JSON(http.StatusNotFound, gin.H{"error": "OrderState not found"})
			return
		}
		c.JSON(http.StatusOK, orderState)
	} else if authRole == "regular" { // если запрос делает обычный пользователь, то не разрешаем:
		osh.logger.Error("forbidden access level.")
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden access level."})
	}
}

// GetOrderStates is ...
// GetOrderStatesTags 		godoc
// @Summary			Получить список всех статусов.
// @Description		Return OrderStates list.
// @Tags			OrderState
// @Security     	BearerAuth
// @Produce      json
// @Success			200 {object} []model.OrderState
// @failure			404 {string} err.Error()
// @Router			/orderstates [get]
func (osh *OrderStateHandler) GetOrderStates(c *gin.Context) {
	authID, authLogin, authRole := utils.GetLevel(c)
	osh.logger.Debug("принятые логин и роль из токена", zap.Int("id", authID), zap.String("login", authLogin), zap.String("role", authRole))

	// если запрос делает суперпользователь, то ему можно всё
	if authRole == "super" {
		orderStates, err := osh.repoRO.GetAllOrderStates(context.TODO())
		if err != nil {
			osh.logger.Error("Error getting order states", zap.Error(err))
			c.JSON(http.StatusNotFound, gin.H{"error": "order states not found"})
			return
		}
		c.JSON(http.StatusOK, orderStates)
	} else if authRole == "regular" { // если запрос делает обычный пользователь, то не разрешаем:
		osh.logger.Error("forbidden access level.")
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden access level."})
	}
}
