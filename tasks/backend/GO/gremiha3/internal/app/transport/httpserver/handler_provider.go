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
// @Param				Provider body model.AddProvider true "Create Provider"
// @Produce				application/json
// @Tags				Provider
// @Security     	BearerAuth
// @Success				200 {object} model.Provider
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

// server.RespondOK(response, w, r)
// authID, authLogin, authRole := utils.GetLevel(c)
// ph.logger.Debug("принятые логин и роль из токена", zap.Int("id", authID), zap.String("login", authLogin), zap.String("role", authRole))
// // если запрос делает не суперпользователь, то выходим с ошибкой
//
//	if authRole != "super" {
//		ph.logger.Error("forbidden access level.")
//		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden access level."})
//		return
//	}
//
// var addProvider model.AddProvider
// // Заполняем структуру addProvider данными из запроса
//
//	if err := c.ShouldBindJSON(&addProvider); err != nil {
//		ph.logger.Error("Error binding JSON-addProvider", zap.Error(err))
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
// var Provider model.Provider
// // Заполняем структуру Provider данными из addProvider
// Provider.Name = addProvider.Name
// Provider.Origin = addProvider.Origin
// // Валидация данных поставщика
//
//	if err := Provider.Validate(); err != nil {
//		ph.logger.Error("Error creating Provider", zap.Error(err))
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
// // Сохраняем поставщика в БД
// id, err := ph.repoWR.CreateProvider(context.TODO(), Provider)
//
//	if err != nil {
//		ph.logger.Error("Error creating Provider", zap.Error(err))
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
// Provider.ID = id
// c.JSON(http.StatusCreated, Provider)
// GetProvider is ...
// GetProviderTags 		godoc
// @Summary			Посмотреть постащика по его id.
// @Description		Return Provider with "id" number.
// @Param			provider_id path int true "Provider ID"
// @Tags			Provider
// @Success			200 {object} model.Provider
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

	// id, _ := strconv.Atoi(c.Param("id"))
	// provider, err := ph.repoRO.GetProviderByID(context.TODO(), id)
	// if err != nil {
	// 	ph.logger.Error("Error getting provider", zap.Error(err))
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
	// 	return
	// }
	// c.JSON(http.StatusOK, provider)
}

// GetProviders is ...
// GetProvidersTags 		godoc
// @Summary			Получить список всех поставщиков.
// @Description		Return Providers list.
// @Tags			Provider
// @Produce      json
// @Success			200 {object} []model.Provider
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
	// providers, err := ph.repoRO.GetAllProviders(context.TODO())
	// if err != nil {
	// 	ph.logger.Error("Error getting providers", zap.Error(err))
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "Providers not found"})
	// 	return
	// }
	// c.JSON(http.StatusOK, providers)
}
