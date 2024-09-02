package httpserver

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SignUp is ...
// SignUpTags		godoc
// @Summary				Загеристрироваться.
// @Description			Sign up a new user in the system.
// @Param				UserRequest body httpserver.UserRequest true "Create user. Логин указывается в формате электронной почты. Пароль не меньше 6 символов. Роль: super или regular"
// @Produce				application/json
// @Tags				Auth
// @Success				201 {object} httpserver.UserResponse
// @failure				400 {string} err.Error()
// @failure				500 {string} string "error-to-create-domain-user"
// @Router				/signup [post]
func (h HttpServer) SignUp(c *gin.Context) {
	var userRequest UserRequest
	var err error
	if err = c.ShouldBindJSON(&userRequest); err != nil {
		log.Println("invalid-json")
		c.JSON(http.StatusBadRequest, gin.H{"invalid-json": err.Error()})
		return
	}

	if err = userRequest.Validate(); err != nil {
		log.Println("invalid-request")
		c.JSON(http.StatusBadRequest, gin.H{"invalid-request": err.Error()})
		return
	}

	userRequest.Password, err = hashPassword(userRequest.Password)
	if err != nil {
		log.Println("error-hashing-password")
		c.JSON(http.StatusUnauthorized, gin.H{"error-hashing-password": err.Error()})
		return
	}

	domainUser, err := toDomainUser(userRequest)
	if err != nil {
		log.Println("error-to-create-domain-user")
		c.JSON(http.StatusInternalServerError, gin.H{"error-to-create-domain-user": err.Error()})
		return
	}

	createdUser, err := h.userService.CreateUser(c, domainUser)
	if err != nil {
		log.Println("error DB saving user")
		c.JSON(http.StatusBadRequest, gin.H{"error DB saving user": err.Error()})
		return
	}

	response := toResponseUser(createdUser)
	c.JSON(http.StatusCreated, response)
}

// SignIn is ...
// SignInTags		godoc
// @Summary				Авторизоваться.
// @Description			Sign in as an existing user.
// @Param				UserRequest body httpserver.UserRequest true "SignIn user. Логин указывается в формате электронной почты. Пароль не меньше 6 символов. Роль: super или regular"
// @Produce				application/json
// @Tags				Auth
// @Success				200 {string} token
// @failure				400 {string} err.Error()
// @failure				500 {string} err.Error()
// @Router				/signin [post]
func (h HttpServer) SignIn(c *gin.Context) {
	var userRequest UserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-json": err.Error()})
		return
	}

	if err := userRequest.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-request": err.Error()})
		return
	}

	domainUser, err := h.userService.GetUserByLogin(c, userRequest.Login)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-request": err.Error()})
		return
	}

	if !checkPasswordHash(userRequest.Password, domainUser.Password()) {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-password": ""})
		return
	}

	token, err := h.tokenService.GenerateToken(domainUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error-generated-token": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
