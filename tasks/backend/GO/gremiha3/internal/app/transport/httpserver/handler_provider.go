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
// @Security     	BearerAuth
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
// @Param			provider_id path int true "Provider ID"
// @Tags			Provider
// @Success			200 {object} ProviderResponse
// @failure			404 {string} err.Error()
// @Router			/provider/{provider_id} [get]
func (h HttpServer) GetProvider(c *gin.Context) {

	providerID, err := strconv.Atoi(c.Param("id"))
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
		c.JSON(http.StatusInternalServerError, gin.H{"": err.Error()})
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
// @Produce      json
// @Success			200 {object} []ProviderResponse
// @failure			404 {string} err.Error()
// @Router			/providers [get]
func (h HttpServer) GetProviders(c *gin.Context) {
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
	providers, err := h.providerService.GetProviders(c, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"": err.Error()})
		return
	}

	response := make([]ProviderResponse, 0, len(providers))
	for _, provider := range providers {
		response = append(response, toResponseProvider(provider))
	}

	c.JSON(http.StatusCreated, response)
}
