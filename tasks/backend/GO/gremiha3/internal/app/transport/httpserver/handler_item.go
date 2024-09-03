package httpserver

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/KozlovNikolai/test-task/internal/app/domain"
	"github.com/KozlovNikolai/test-task/internal/pkg/config"
	"github.com/gin-gonic/gin"
)

// CreateItem is ...
// CreateItemTags		godoc
// @Summary				Добавить товар в заказ.
// @Description			Add product to order.
// @Param				item body ItemRequest true "Create item"
// @Produce				application/json
// @Tags				Item
// @Security			BearerAuth
// @Success				200 {object} ItemResponse
// @failure				400 {string} err.Error()
// @failure				500 {string} err.Error()
// @Router				/item [post]
func (h HttpServer) CreateItem(c *gin.Context) {
	// Auth
	// check auth
	userCtx, err := getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"get orders": domain.ErrNoUserInContext})
		return
	}

	//#####################
	var itemRequest ItemRequest
	if err := c.ShouldBindJSON(&itemRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-json": err.Error()})
		return
	}

	if err := itemRequest.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-request": err.Error()})
		return
	}

	item, err := toDomainItem(itemRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error creating domain item": err.Error()})
		return
	}
	// получаем заказ
	order, err := h.orderService.GetOrder(c, item.OrderID())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error get order": err.Error()})
		return
	}
	// получаем пользователя
	user, err := h.userService.GetUserByID(c, order.UserID())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error get user": err.Error()})
		return
	}
	// если пользователь не владелец и не админ, то выходим
	if userCtx.ID() != user.ID() && userCtx.Role() != config.AdminRole {
		c.JSON(http.StatusUnauthorized, gin.H{"invalid user id or role": domain.ErrInvalidUserLogin.Error()})
		return
	}

	insertedItem, err := h.itemService.CreateItem(c, item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error DB saving item": err.Error()})
		return
	}

	response := toResponseItem(insertedItem)
	c.JSON(http.StatusCreated, response)
}

// GetItem is ...
// GetItemTags 		godoc
// @Summary			Посмотреть запись в заказе по ее id.
// @Description		Return item with "id" number.
// @Param        id  query   string  false  "id of the item" example(1)  default(1)
// @Tags			Item
// @Security			BearerAuth
// @Success			200 {object} ItemResponse
// @failure			404 {string} err.Error()
// @Router			/item [get]
func (h HttpServer) GetItem(c *gin.Context) {

	itemID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-item-id": err.Error()})
		return
	}
	// Auth
	// check auth
	userCtx, err := getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"get orders": domain.ErrNoUserInContext})
		return
	}
	// получаем запись заказа
	item, err := h.itemService.GetItem(c, itemID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"item-not-found": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error-get-item": err.Error()})
		return
	}
	// получаем заказ
	order, err := h.orderService.GetOrder(c, item.OrderID())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error get order": err.Error()})
		return
	}
	// получаем пользователя
	user, err := h.userService.GetUserByID(c, order.UserID())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error get user": err.Error()})
		return
	}
	// если пользователь не владелец и не админ, то выходим
	if userCtx.ID() != user.ID() && userCtx.Role() != config.AdminRole {
		c.JSON(http.StatusUnauthorized, gin.H{"invalid user id or role": domain.ErrInvalidUserLogin.Error()})
		return
	}
	response := toResponseItem(item)

	c.JSON(http.StatusCreated, response)
}

// GetItems is ...
// GetItemsTags 		godoc
// @Summary			Получить список товаров в заказе.
// @Description		Return items list by order id.
// @Tags			Item
// @Param        limit  query   string  true  "limit records on page" example(10)  default(10)
// @Param        offset  query   string  true  "start of record output" example(0)  default(0)
// @Param        orderid  query   string  true  "filter by order id" example(1)  default(1)
// @Produce      json
// @Security			BearerAuth
// @Success			200 {object} []ItemResponse
// @failure			404 {string} err.Error()
// @Router			/items [get]
func (h HttpServer) GetItems(c *gin.Context) {
	limit_query := c.Query("limit")
	offset_query := c.Query("offset")
	orderid_query := c.Query("orderid")

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
	orderid, err := strconv.Atoi(orderid_query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-orderid": err.Error()})
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
	if orderid < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"userid-must-be-greater-then-zero": ""})
		return
	}
	// Auth
	// check auth
	userCtx, err := getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"get orders": domain.ErrNoUserInContext})
		return
	}
	// получаем заказ
	order, err := h.orderService.GetOrder(c, orderid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error get order": err.Error()})
		return
	}
	// получаем пользователя
	user, err := h.userService.GetUserByID(c, order.ID())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error get user": err.Error()})
		return
	}
	// если пользователь не владелец и не админ, то выходим
	if userCtx.ID() != user.ID() && userCtx.Role() != config.AdminRole {
		c.JSON(http.StatusUnauthorized, gin.H{"invalid user id or role": domain.ErrInvalidUserLogin.Error()})
		return
	}
	items, err := h.itemService.GetItems(c, limit, offset, orderid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error get items": err.Error()})
		return
	}

	response := make([]ItemResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toResponseItem(item))
	}

	c.JSON(http.StatusOK, response)
}
