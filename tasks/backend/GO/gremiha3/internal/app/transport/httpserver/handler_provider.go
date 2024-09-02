package httpserver

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/KozlovNikolai/test-task/internal/app/domain"
	"github.com/gin-gonic/gin"
)

// CreateProvider is ...
// CreateProviderTags		godoc
// @Summary				Добавить поставщика.
// @Description			Save register data of user in Repo.
// @Param				Provider body ProviderRequest true "Create Provider"
// @Produce				application/json
// @Tags				Provider
// @Success				200 {object} ProviderResponse
// @failure				400 {string} err.Error()
// @failure				500 {string} err.Error()
// @Router				/provider [post]
func (h HttpServer) CreateProvider(c *gin.Context) {
	var providerRequest ProviderRequest
	if err := c.ShouldBindJSON(&providerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-json": err.Error()})
		return
	}

	if err := providerRequest.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-request": err.Error()})
		return
	}

	provider, err := toDomainProvider(providerRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error creating domain provider": err.Error()})
		return
	}

	insertedProvider, err := h.providerService.CreateProvider(c, provider)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error DB saving provider": err.Error()})
		return
	}

	response := toResponseProvider(insertedProvider)
	c.JSON(http.StatusCreated, response)
}

// GetProvider is ...
// GetProviderTags 		godoc
// @Summary			Посмотреть постащика по его id.
// @Description		Return Provider with "id" number.
// @Param        id  query   string  false  "id of the provider" example(1)
// @Tags			Provider
// @Success			200 {object} ProviderResponse
// @failure			404 {string} err.Error()
// @Router			/provider [get]
func (h HttpServer) GetProvider(c *gin.Context) {

	providerID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-provider-id": err.Error()})
		return
	}

	provider, err := h.providerService.GetProvider(c, providerID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"provider-not-found": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error-get-provider": err.Error()})
		return
	}

	response := toResponseProvider(provider)

	c.JSON(http.StatusCreated, response)
}

// GetProviders is ...
// GetProvidersTags 		godoc
// @Summary			Получить список всех поставщиков.
// @Description		Return Providers list.
// @Tags			Provider
// @Param        limit  query   string  true  "limit records on page" example(10)
// @Param        offset  query   string  true  "start of record output" example(1)
// @Produce      json
// @Success			200 {object} []ProviderResponse
// @failure			404 {string} err.Error()
// @Router			/providers [get]
func (h HttpServer) GetProviders(c *gin.Context) {
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
	if offset < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"offset-must-be-greater-then-zero": ""})
		return
	}

	providers, err := h.providerService.GetProviders(c, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error get providers": err.Error()})
		return
	}

	response := make([]ProviderResponse, 0, len(providers))
	for _, provider := range providers {
		response = append(response, toResponseProvider(provider))
	}

	c.JSON(http.StatusOK, response)
}
