package httpserver

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/KozlovNikolai/test-task/internal/app/domain"
	"github.com/gin-gonic/gin"
)

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
func (h HttpServer) CreateOrder(c *gin.Context) {
	var orderRequest OrderRequest
	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-json": err.Error()})
		return
	}

	if err := orderRequest.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-request": err.Error()})
		return
	}

	order, err := toDomainOrder(orderRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error creating domain order": err.Error()})
		return
	}

	insertedorder, err := h.orderService.CreateOrder(c, order)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error DB saving order": err.Error()})
		return
	}

	response := toResponseOrder(insertedorder)
	c.JSON(http.StatusCreated, response)
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
func (h HttpServer) GetOrder(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-order-id": err.Error()})
		return
	}

	order, err := h.orderService.GetOrderByID(c, orderID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"order-not-found": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"": err.Error()})
		return
	}

	response := toResponseOrder(order)

	c.JSON(http.StatusCreated, response)
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
func (h HttpServer) GetOrders(c *gin.Context) {
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
	orders, err := h.orderService.GetOrders(c, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"": err.Error()})
		return
	}

	response := make([]OrderResponse, 0, len(orders))
	for _, order := range orders {
		response = append(response, toResponseOrder(order))
	}

	c.JSON(http.StatusCreated, response)
}
