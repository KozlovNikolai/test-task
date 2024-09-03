package httpserver

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/KozlovNikolai/test-task/internal/app/domain"
	"github.com/KozlovNikolai/test-task/internal/pkg/config"
	"github.com/gin-gonic/gin"
)

// CreateOrder is ...
// CreateOrderTags		godoc
// @Summary				Добавить заказ.
// @Description			Создание заказа для дальнейшего его заполнения.
// @Param				order body OrderRequest true "Create Order"
// @Produce				application/json
// @Tags				Order
// @Security			BearerAuth
// @Success				201 {object} OrderResponse
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
// @Summary			Посмотреть заказ по его id или по логину пользователя.
// @Description		Return Order with "id" number.
// @Param        id  query   string  false  "id of the order" example(1)   default(1)
// @Tags			Order
// @Security			BearerAuth
// @Success			200 {object} OrderResponse
// @failure			404 {string} err.Error()
// @Router			/order [get]
func (h HttpServer) GetOrder(c *gin.Context) {
	// check auth
	userCtx, err := getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"get orders": domain.ErrNoUserInContext})
		return
	}

	orderID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-order-id": err.Error()})
		return
	}

	order, err := h.orderService.GetOrder(c, orderID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"order-not-found": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"": err.Error()})
		return
	}

	if userCtx.ID() != order.UserID() && userCtx.Role() != config.AdminRole {
		c.JSON(http.StatusUnauthorized, gin.H{"invalid user id or role": domain.ErrInvalidUserLogin.Error()})
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
// @Param        limit  query   string  true  "limit records on page" example(10)  default(10)
// @Param        offset  query   string  true  "start of record output" example(0) default(0)
// @Param        userid  query   string  true  "filter by user id" example(1) default(1)
// @Produce      json
// @Security			BearerAuth
// @Success			200 {object} []OrderResponse
// @failure			404 {string} err.Error()
// @Router			/orders [get]
func (h HttpServer) GetOrders(c *gin.Context) {
	limit_query := c.Query("limit")
	offset_query := c.Query("offset")
	userid_query := c.Query("userid")

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
	userid, err := strconv.Atoi(userid_query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-userid": err.Error()})
		return
	}
	// check auth
	userCtx, err := getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"get orders": domain.ErrNoUserInContext})
		return
	}

	if limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"limit-must-be-greater-then-zero": ""})
		return
	}
	if offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"offset-must-be-equal-or-greater-then-zero": ""})
		return
	}
	if userid < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"userid-must-be-greater-then-zero": ""})
		return
	}
	if userCtx.ID() != userid && userCtx.Role() != config.AdminRole {
		c.JSON(http.StatusUnauthorized, gin.H{"invalid user id or role": domain.ErrInvalidUserLogin.Error()})
		return
	}
	orders, err := h.orderService.GetOrders(c, limit, offset, userid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error get orders": err.Error()})
		return
	}

	response := make([]OrderResponse, 0, len(orders))
	for _, order := range orders {
		response = append(response, toResponseOrder(order))
	}

	c.JSON(http.StatusOK, response)
}
