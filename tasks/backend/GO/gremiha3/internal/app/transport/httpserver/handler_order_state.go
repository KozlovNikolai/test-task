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
// @Param				orderState body OrderStateRequest true "Create Order State type"
// @Produce				application/json
// @Tags				OrderState
// @Success				200 {object} OrderStateResponse
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

	insertedOrderState, err := h.orderStateService.CreateOrderState(c, orderState)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error DB saving orderState": err.Error()})
		return
	}

	response := toResponseOrderState(insertedOrderState)
	c.JSON(http.StatusCreated, response)
}

// GetOrderState is ...
// GetOrderStateTags 		godoc
// @Summary			Посмотреть тип статуса по его id.
// @Description		Return OrderState with "id" number.
// @Param        id  query   string  false  "id of the order state" example(1) default(1)
// @Tags			OrderState
// @Success			200 {object} OrderStateResponse
// @failure			404 {string} err.Error()
// @Router			/orderstate [get]
func (h HttpServer) GetOrderState(c *gin.Context) {
	orderStateID, err := strconv.Atoi(c.Query("id"))
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
		c.JSON(http.StatusInternalServerError, gin.H{"error-get-orderState": err.Error()})
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
// @Param        limit  query   string  true  "limit records on page" example(10) default(10)
// @Param        offset  query   string  true  "start of record output" example(0) default(0)
// @Produce      json
// @Success			200 {object} []OrderStateResponse
// @failure			404 {string} err.Error()
// @Router			/orderstates [get]
func (h HttpServer) GetOrderStates(c *gin.Context) {
	limit_query := c.Query("limit")
	offset_query := c.Query("offset")

	limit, err := strconv.Atoi(limit_query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-limit": err.Error()})
		return
	}

	offset, err := strconv.Atoi(offset_query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-offset": err.Error()})
		return
	}
	if limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"limit-must-be-greater-then-zero": ""})
		return
	}
	if offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"offset-must-be-greater-then-zero": ""})
		return
	}

	orderStates, err := h.orderStateService.GetOrderStates(c, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error get orderStates": err.Error()})
		return
	}

	response := make([]OrderStateResponse, 0, len(orderStates))
	for _, orderState := range orderStates {
		response = append(response, toResponseOrderState(orderState))
	}

	c.JSON(http.StatusOK, response)
}
