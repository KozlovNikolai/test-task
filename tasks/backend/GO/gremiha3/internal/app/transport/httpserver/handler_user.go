package httpserver

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/KozlovNikolai/test-task/internal/app/domain"
	"github.com/gin-gonic/gin"
)

// GetUser is ...
// GetUserTags 		godoc
// @Summary			Посмотреть пользователя по его id.
// @Description		Return user with "id" number.
// @Param        id  query   string  false  "id of the user"
// @Param        login  query   string  false  "login of the user"
// @Tags			User
// @Success			200 {object} UserResponse
// @failure			404 {string} err.Error()
// @Router			/user [get]
func (h HttpServer) GetUser(c *gin.Context) {
	var userRequest UserRequest
	id_query := c.Query("id")
	login_query := c.Query("login")
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

}

// GetUserByLogin is ...
// GetUserTags 		godoc
// @Summary			Посмотреть пользователя по его логину.
// @Description		Return user with "login" value.
// @Param			login path string true "Login"
// @Tags			User
// @Security     	BearerAuth
// @Success			200 {object} UserResponse
// @failure			404 {string} err.Error()
// @Router			/user/login/{login} [get]
func (h HttpServer) GetUserByLogin(c *gin.Context) {
	// login := c.Param("login")
	// authID, authLogin, authRole := utils.GetLevel(c)
	// u.logger.Debug("принятые логин и роль из токена", zap.Int("id", authID), zap.String("login", authLogin), zap.String("role", authRole))
	// // если запрос делает суперпользователь, то ему можно всё
	// if authRole == "super" {
	// 	user, err := u.repoRO.GetUserByLogin(context.TODO(), login)
	// 	if err != nil {
	// 		u.logger.Error("Error getting user", zap.Error(err))
	// 		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	// 		return
	// 	}
	// 	c.JSON(http.StatusOK, user)
	// } else if authRole == "regular" { // если запрос делает обычный пользователь, то ему можно смотреть только собственные данные
	// 	user, err := u.repoRO.GetUserByLogin(context.TODO(), authLogin)
	// 	if err != nil {
	// 		u.logger.Error("Error getting user", zap.Error(err))
	// 		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	// 		return
	// 	}
	// 	if user.Login != login {
	// 		u.logger.Error("forbidden access level.")
	// 		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden access level."})
	// 		return
	// 	}
	// 	c.JSON(http.StatusOK, user)
	// }
}

// GetUsers is ...
// GetUsersTags 		godoc
// @Summary			Получить список всех пользователей.
// @Description		Return users list.
// @Tags			User
// @Security     	BearerAuth
// @Produce      json
// @Success			200 {object} []UserResponse
// @failure			404 {string} err.Error()
// @Router			/users [get]
func (h HttpServer) GetUsers(c *gin.Context) {
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
	users, err := h.userService.GetUsers(c, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"": err.Error()})
		return
	}

	response := make([]UserResponse, 0, len(users))
	for _, user := range users {
		response = append(response, toResponseUser(user))
	}

	c.JSON(http.StatusCreated, response)
}
