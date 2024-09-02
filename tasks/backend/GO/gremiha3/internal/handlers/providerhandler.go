package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/KozlovNikolai/test-task/internal/model"
	"github.com/KozlovNikolai/test-task/internal/pkg/utils"
	"github.com/KozlovNikolai/test-task/internal/store"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ProviderHandler is ...
type ProviderHandler struct {
	logger *zap.Logger
	repoWR store.IProviderRepository
	repoRO store.IProviderRepository
}

// NewProviderHandler is ...
func NewProviderHandler(logger *zap.Logger, repoWR store.IProviderRepository, repoRO store.IProviderRepository) *ProviderHandler {
	return &ProviderHandler{
		logger: logger,
		repoWR: repoWR,
		repoRO: repoRO,
	}
}

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
func (ph *ProviderHandler) CreateProvider(c *gin.Context) {

	authID, authLogin, authRole := utils.GetLevel(c)
	ph.logger.Debug("принятые логин и роль из токена", zap.Int("id", authID), zap.String("login", authLogin), zap.String("role", authRole))
	// если запрос делает не суперпользователь, то выходим с ошибкой
	if authRole != "super" {
		ph.logger.Error("forbidden access level.")
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden access level."})
		return
	}
	var addProvider model.AddProvider
	// Заполняем структуру addProvider данными из запроса
	if err := c.ShouldBindJSON(&addProvider); err != nil {
		ph.logger.Error("Error binding JSON-addProvider", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var Provider model.Provider
	// Заполняем структуру Provider данными из addProvider
	Provider.Name = addProvider.Name
	Provider.Origin = addProvider.Origin
	// Валидация данных поставщика
	if err := Provider.Validate(); err != nil {
		ph.logger.Error("Error creating Provider", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Сохраняем поставщика в БД
	id, err := ph.repoWR.CreateProvider(context.TODO(), Provider)
	if err != nil {
		ph.logger.Error("Error creating Provider", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	Provider.ID = id
	c.JSON(http.StatusCreated, Provider)
}

// GetProvider is ...
// GetProviderTags 		godoc
// @Summary			Посмотреть постащика по его id.
// @Description		Return Provider with "id" number.
// @Param			provider_id path int true "Provider ID"
// @Tags			Provider
// @Success			200 {object} model.Provider
// @failure			404 {string} err.Error()
// @Router			/provider/{provider_id} [get]
func (ph *ProviderHandler) GetProvider(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	provider, err := ph.repoRO.GetProviderByID(context.TODO(), id)
	if err != nil {
		ph.logger.Error("Error getting provider", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
		return
	}
	c.JSON(http.StatusOK, provider)
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
func (ph *ProviderHandler) GetProviders(c *gin.Context) {

	providers, err := ph.repoRO.GetAllProviders(context.TODO())
	if err != nil {
		ph.logger.Error("Error getting providers", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Providers not found"})
		return
	}
	c.JSON(http.StatusOK, providers)
}
