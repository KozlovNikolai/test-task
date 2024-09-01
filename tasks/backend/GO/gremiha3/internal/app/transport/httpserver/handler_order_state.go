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
}
