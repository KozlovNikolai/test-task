package httpserver

import (
	"net/http"
	"strings"

	"github.com/KozlovNikolai/test-task/internal/pkg/config"
	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
)

func (h HttpServer) CheckAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(AuthorizationHeader)
		token = strings.TrimPrefix(token, BearerPrefix)
		user, err := h.tokenService.GetUser(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"check-admin-validate-token": err.Error()})
			return
		}
		if user.Login() == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"check-admin-invalid-token": ""})
			return
		}
		if user.Role() != config.AdminRole {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"not-admin": ""})
			return
		}
		c.Set("user", user)
		c.Next()
	}
}

func (h HttpServer) CheckAuthorizedUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(AuthorizationHeader)
		token = strings.TrimPrefix(token, BearerPrefix)
		user, err := h.tokenService.GetUser(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"check-auth-validate-token": err.Error()})
			return
		}
		if user.Login() == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"check-auth-invalid-token": ""})
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
