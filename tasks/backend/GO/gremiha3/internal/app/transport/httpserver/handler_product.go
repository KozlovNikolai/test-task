package httpserver

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/KozlovNikolai/test-task/internal/app/domain"
	"github.com/gin-gonic/gin"
)

// CreateProduct is ...
// CreateProductTags		godoc
// @Summary				Добавить товар.
// @Description			Save register data of user in Repo.
// @Param				product body ProductRequest true "Create product"
// @Produce				application/json
// @Tags				Product
// @Security			BearerAuth
// @Success				200 {object} ProductResponse
// @failure				400 {string} err.Error()
// @failure				500 {string} err.Error()
// @Router				/product [post]
func (h HttpServer) CreateProduct(c *gin.Context) {
	var productRequest ProductRequest
	if err := c.ShouldBindJSON(&productRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-json": err.Error()})
		return
	}

	if err := productRequest.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-request": err.Error()})
		return
	}

	product, err := toDomainProduct(productRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error creating domain product": err.Error()})
		return
	}

	insertedproduct, err := h.productService.CreateProduct(c, product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error DB saving product": err.Error()})
		return
	}

	response := toResponseProduct(insertedproduct)
	c.JSON(http.StatusCreated, response)
}

// GetProduct is ...
// GetProductTags 		godoc
// @Summary			Посмотреть товар по его id.
// @Description		Return product with "id" number.
// @Param        id  query   string  false  "id of the product" example(1) default(1)
// @Tags			Product
// @Success			200 {object} ProductResponse
// @failure			404 {string} err.Error()
// @Router			/product [get]
func (h HttpServer) GetProduct(c *gin.Context) {
	productID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-product-id": err.Error()})
		return
	}

	product, err := h.productService.GetProduct(c, productID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"product-not-found": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error-get-product": err.Error()})
		return
	}

	response := toResponseProduct(product)

	c.JSON(http.StatusCreated, response)
}

// GetProducts is ...
// GetProductsTags 		godoc
// @Summary			Получить список всех товаров.
// @Description		Return products list.
// @Tags			Product
// @Param        limit  query   string  true  "limit records on page" example(10) default(10)
// @Param        offset  query   string  true  "start of record output" example(0) default(0)
// @Produce      json
// @Success			200 {object} []ProductResponse
// @failure			404 {string} err.Error()
// @Router			/products [get]
func (h HttpServer) GetProducts(c *gin.Context) {
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

	products, err := h.productService.GetProducts(c, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error get products": err.Error()})
		return
	}

	response := make([]ProductResponse, 0, len(products))
	for _, product := range products {
		response = append(response, toResponseProduct(product))
	}

	c.JSON(http.StatusOK, response)
}
