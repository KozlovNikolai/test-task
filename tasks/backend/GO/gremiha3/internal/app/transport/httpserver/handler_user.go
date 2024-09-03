package httpserver

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/KozlovNikolai/test-task/internal/app/domain"
	"github.com/KozlovNikolai/test-task/internal/pkg/config"
	"github.com/gin-gonic/gin"
)

// GetUser is ...
// GetUserTags 		godoc
// @Summary			Посмотреть пользователя по его id или логину.
// @Description		Получить пользователя по его id ли логину.
// @Param        id  query   string  false  "id of the user" example(1) default(1)
// @Param        login  query   string  false  "login of the user" example(cmd@cmd.ru) default(cmd@cmd.ru)
// @Tags			User
// @Security			BearerAuth
// @Success			200 {object} UserResponse
// @failure			404 {string} err.Error()
// @Router			/user [get]
func (h HttpServer) GetUser(c *gin.Context) {
	var userRequest UserRequest
	id_query := c.Query("id")
	login_query := c.Query("login")

	// check auth
	userCtx, err := getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"get orders": domain.ErrNoUserInContext})
		return
	}

	//
	if login_query != "" {
		userRequest.Login = login_query
		userRequest.Password = "fake_password"
		if err := userRequest.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"invalid-user-login": err.Error()})
		}

		domainUser, err := h.userService.GetUserByLogin(c, login_query)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"user-not-found": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"": err.Error()})
			return
		}
		// auth login
		if userCtx.Login() != domainUser.Login() && userCtx.Role() != config.AdminRole {
			c.JSON(http.StatusUnauthorized, gin.H{"invalid user id or role": domain.ErrInvalidUserLogin.Error()})
			return
		}
		response := toResponseUser(domainUser)
		c.JSON(http.StatusOK, response)
		return
	}
	//
	if id_query != "" {
		userID, err := strconv.Atoi(id_query)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"invalid-user-id": err.Error()})
			return
		}
		// auth user id
		if userCtx.ID() != userID && userCtx.Role() != config.AdminRole {
			c.JSON(http.StatusUnauthorized, gin.H{"invalid user id or role": domain.ErrInvalidUserID.Error()})
			return
		}
		user, err := h.userService.GetUserByID(c, userID)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"user-not-found": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"": err.Error()})
			return
		}
		response := toResponseUser(user)
		c.JSON(http.StatusOK, response)
		return
	}
	user, err := h.userService.GetUserByID(c, userCtx.ID())
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"user-not-found": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"": err.Error()})
		return
	}
	response := toResponseUser(user)
	c.JSON(http.StatusOK, response)
}

// GetUsers is ...
// GetUsersTags 		godoc
// @Summary			Получить список всех пользователей.
// @Description		Return users list.
// @Tags			User
// @Security		BearerAuth
// @Param        limit  query   string  true  "limit records on page" example(10) default(10)
// @Param        offset  query   string  true  "start of record output" example(0) default(0)
// @Produce      json
// @Success			200 {object} []UserResponse
// @failure			404 {string} err.Error()
// @Router			/users [get]
func (h HttpServer) GetUsers(c *gin.Context) {
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

	users, err := h.userService.GetUsers(c, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error get users": err.Error()})
		return
	}

	response := make([]UserResponse, 0, len(users))
	for _, user := range users {
		response = append(response, toResponseUser(user))
	}

	c.JSON(http.StatusOK, response)
}
