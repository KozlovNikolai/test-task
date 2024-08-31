package httpserver

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/KozlovNikolai/test-task/internal/app/domain"
	"github.com/gin-gonic/gin"
)

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
func (h HttpServer) CreateOrderState(c *gin.Context) {
	var orderStateRequest OrderStateRequest
	if err := c.ShouldBindJSON(&orderStateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-json": err.Error()})
		return
	}

	if err := orderStateRequest.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-request": err.Error()})
		return
	}

	orderState, err := toDomainOrderState(orderStateRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error creating domain orderState": err.Error()})
		return
	}

	insertedorderState, err := h.orderStateService.CreateOrderState(c, orderState)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error DB saving orderState": err.Error()})
		return
	}

	response := toResponseOrderState(insertedorderState)
	c.JSON(http.StatusCreated, response)

	// authID, authLogin, authRole := utils.GetLevel(c)
	// osh.logger.Debug("принятые логин и роль из токена", zap.Int("id", authID), zap.String("login", authLogin), zap.String("role", authRole))
	// // если запрос делает суперпользователь, то ему можно всё
	// if authRole == "super" {
	// 	var addOrderState model.AddOrderState
	// 	// Заполняем структуру addOrderState данными из запроса
	// 	if err := c.ShouldBindJSON(&addOrderState); err != nil {
	// 		osh.logger.Error("Error binding JSON-addOrderState", zap.Error(err))
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}
	// 	var orderState model.OrderState
	// 	// Заполняем структуру OrderState данными из addOrderState
	// 	orderState.Name = addOrderState.Name
	// 	// Сохраняем имя статуса заказа в БД
	// 	id, err := osh.repoWR.CreateOrderState(context.TODO(), orderState)
	// 	if err != nil {
	// 		osh.logger.Error("Error creating OrderState", zap.Error(err))
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}
	// 	orderState.ID = id
	// 	c.JSON(http.StatusCreated, orderState)
	// } else if authRole == "regular" { // если запрос делает обычный пользователь, то не разрешаем:
	// 	osh.logger.Error("forbidden access level.")
	// 	c.JSON(http.StatusForbidden, gin.H{"error": "forbidden access level."})
	// }
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
func (h HttpServer) GetOrderState(c *gin.Context) {
	orderStateID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-orderState-id": err.Error()})
		return
	}

	orderState, err := h.orderStateService.GetOrderState(c, orderStateID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"orderState-not-found": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"": err.Error()})
		return
	}

	response := toResponseOrderState(orderState)

	c.JSON(http.StatusCreated, response)

	// authID, authLogin, authRole := utils.GetLevel(c)
	// osh.logger.Debug("принятые логин и роль из токена", zap.Int("id", authID), zap.String("login", authLogin), zap.String("role", authRole))
	// // если запрос делает суперпользователь, то ему можно всё
	// if authRole == "super" {
	// 	id, _ := strconv.Atoi(c.Param("id"))
	// 	orderState, err := osh.repoRO.GetOrderStateByID(context.TODO(), id)
	// 	if err != nil {
	// 		osh.logger.Error("Error getting OrderState", zap.Error(err))
	// 		c.JSON(http.StatusNotFound, gin.H{"error": "OrderState not found"})
	// 		return
	// 	}
	// 	c.JSON(http.StatusOK, orderState)
	// } else if authRole == "regular" { // если запрос делает обычный пользователь, то не разрешаем:
	// 	osh.logger.Error("forbidden access level.")
	// 	c.JSON(http.StatusForbidden, gin.H{"error": "forbidden access level."})
	// }
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
func (h HttpServer) GetOrderStates(c *gin.Context) {
	limit, err := strconv.Atoi(c.Param("limit"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-limit": err.Error()})
		return
	}
	offset, err := strconv.Atoi(c.Param("offset"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-offset": err.Error()})
		return
	}
	orderStates, err := h.orderStateService.GetOrderStates(c, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"": err.Error()})
		return
	}

	response := make([]OrderStateResponse, 0, len(orderStates))
	for _, orderState := range orderStates {
		response = append(response, toResponseOrderState(orderState))
	}

	c.JSON(http.StatusCreated, response)
	// authID, authLogin, authRole := utils.GetLevel(c)
	// osh.logger.Debug("принятые логин и роль из токена", zap.Int("id", authID), zap.String("login", authLogin), zap.String("role", authRole))
	// // если запрос делает суперпользователь, то ему можно всё
	// if authRole == "super" {
	// 	orderStates, err := osh.repoRO.GetAllOrderStates(context.TODO())
	// 	if err != nil {
	// 		osh.logger.Error("Error getting order states", zap.Error(err))
	// 		c.JSON(http.StatusNotFound, gin.H{"error": "order states not found"})
	// 		return
	// 	}
	// 	c.JSON(http.StatusOK, orderStates)
	// } else if authRole == "regular" { // если запрос делает обычный пользователь, то не разрешаем:
	// 	osh.logger.Error("forbidden access level.")
	// 	c.JSON(http.StatusForbidden, gin.H{"error": "forbidden access level."})
	// }
}
