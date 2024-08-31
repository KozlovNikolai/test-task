package httpserver

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/KozlovNikolai/test-task/internal/app/domain"
	"github.com/gin-gonic/gin"
)

// CreateUser is ...
// CreateUserTags		godoc
// @Summary				Добавить пользователя.
// @Description			Save register data of user in Repo.
// @Param				user body model.AddUser true "Create user. Логин указывается в формате электронной почты. Пароль не меньше 6 символов. Роль: super или regular"
// @Produce				application/json
// @Tags				Auth
// @Success				200 {object} model.User
// @Success				201 {string} string "Пользователь успешно зарегистрирован."
// @failure				400 {string} err.Error()
// @failure				500 {string} err.Error()
// @Router				/user/register [post]
func (h HttpServer) CreateUser(c *gin.Context) {
	var userRequest UserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-json": err.Error()})
		return
	}

	if err := userRequest.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid-request": err.Error()})
		return
	}

	user, err := toDomainUser(userRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error creating domain user": err.Error()})
		return
	}

	inserteduser, err := h.userService.CreateUser(c, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error DB saving user": err.Error()})
		return
	}

	response := toResponseUser(inserteduser)
	c.JSON(http.StatusCreated, response)
	// var addUser model.AddUser
	// // Заполняем структуру addUser данными из запроса
	// if err := c.ShouldBindJSON(&addUser); err != nil {
	// 	u.logger.Error("Error binding JSON-addUser", zap.Error(err))
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	// var user model.User
	// // Заполняем структуру User данными из addUser
	// user.Login = addUser.Login
	// user.Password = addUser.Password
	// user.Role = addUser.Role
	// // Валидация данных пользователя
	// if err := user.Validate(); err != nil {
	// 	u.logger.Error("Error creating user", zap.Error(err))
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	// // Хеширование пароля
	// if err := user.HashPassword(); err != nil {
	// 	u.logger.Error("Error creating user", zap.Error(err))
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
	// 	return
	// }
	// // Сохраняем пользователя в базе с проверкой на уникальность email
	// id, err := u.repoWR.CreateUser(context.TODO(), user)
	// if err != nil {
	// 	u.logger.Error("Error creating user", zap.Error(err))
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	// user.ID = id
	// c.JSON(http.StatusCreated, user)
}

// GetUser is ...
// GetUserTags 		godoc
// @Summary			Посмотреть пользователя по его id.
// @Description		Return user with "id" number.
// @Param			user_id path int true "User ID"
// @Tags			User
// @Security     	BearerAuth
// @Success			200 {object} model.User
// @failure			404 {string} err.Error()
// @Router			/user/{user_id} [get]
func (h HttpServer) GetUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
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

	c.JSON(http.StatusCreated, response)
	// authID, authLogin, authRole := utils.GetLevel(c)
	// u.logger.Debug("принятые логин и роль из токена", zap.Int("id", authID), zap.String("login", authLogin), zap.String("role", authRole))
	// // если запрос делает суперпользователь, то ему можно всё
	// if authRole == "super" {
	// 	id, _ := strconv.Atoi(c.Param("id"))
	// 	user, err := u.repoRO.GetUserByID(context.TODO(), id)
	// 	if err != nil {
	// 		u.logger.Error("Error getting user", zap.Error(err))
	// 		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	// 		return
	// 	}
	// 	c.JSON(http.StatusOK, user)
	// } else if authRole == "regular" { // если запрос делает обычный пользователь, то ему можно смотреть только собственные данные
	// 	id, _ := strconv.Atoi(c.Param("id"))
	// 	user, err := u.repoRO.GetUserByID(context.TODO(), authID)
	// 	if err != nil {
	// 		u.logger.Error("Error getting user", zap.Error(err))
	// 		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	// 		return
	// 	}
	// 	if user.ID != id {
	// 		u.logger.Error("forbidden access level.")
	// 		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden access level."})
	// 		return
	// 	}
	// 	c.JSON(http.StatusOK, user)
	// }
}

// GetUserByLogin is ...
// GetUserTags 		godoc
// @Summary			Посмотреть пользователя по его логину.
// @Description		Return user with "login" value.
// @Param			login path string true "Login"
// @Tags			User
// @Security     	BearerAuth
// @Success			200 {object} model.User
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

// LoginUser is ...
// LoginUserTags 		godoc
// @Summary				Авторизоваться по логину и паролю.
// @Description			Returns the authorization status
// @Param				user body model.LoginUser true "Login user. Логин указывается в формате электронной почты. Пароль не меньше 6 символов. Роль: super или regular"
// @Accept       		json
// @Produce				json
// @Tags				Auth
// @Success      		200 {string} string "Logged in"
// @failure				400 {string} err.Error()
// @failure				500 {string} err.Error()
// @Router				/user/login [post]
func (h HttpServer) LoginUser(c *gin.Context) {
	// var loginUser model.LoginUser
	// if err := c.BindJSON(&loginUser); err != nil {
	// 	u.logger.Error("Error binding JSON-loginUser", zap.Error(err))
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	// // Find login in DB
	// user, err := u.repoRO.GetUserByLogin(context.TODO(), loginUser.Login)
	// if err != nil {
	// 	u.logger.Error("Error. User not find", zap.Error(err))
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	// // Check username and password
	// err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password))
	// if err != nil {
	// 	u.logger.Error("Error. Password is wrong", zap.Error(err))
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Error. Password is wrong"})
	// 	return
	// }
	// // Create JWT token
	// // expirationTime := time.Now().Add(config.Cfg.TokenTimeDuration)
	// tokenString, err := utils.GenerateJWT(user.ID, user.Login, user.Role)
	// if err != nil {
	// 	u.logger.Error("Error. Could not create token", zap.Error(err))
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
	// }
	// msg := fmt.Sprintf("Logged in with id: %d, login: %s, role: %s", user.ID, user.Login, user.Role)
	// c.JSON(http.StatusOK, gin.H{"message": msg, "token": tokenString})
}

// GetUsers is ...
// GetUsersTags 		godoc
// @Summary			Получить список всех пользователей.
// @Description		Return users list.
// @Tags			User
// @Security     	BearerAuth
// @Produce      json
// @Success			200 {object} []model.User
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
	// authID, authLogin, authRole := utils.GetLevel(c)
	// u.logger.Debug("принятые логин и роль из токена", zap.Int("id", authID), zap.String("login", authLogin), zap.String("role", authRole))
	// // если запрос делает суперпользователь, то ему можно всё
	// if authRole == "super" {
	// 	users, err := u.repoRO.GetAllUsers(context.TODO())
	// 	if err != nil {
	// 		u.logger.Error("Error getting users", zap.Error(err))
	// 		c.JSON(http.StatusNotFound, gin.H{"error": "Users not found"})
	// 		return
	// 	}
	// 	c.JSON(http.StatusOK, users)
	// 	return
	// }
	// u.logger.Error("forbidden access level.")
	// c.JSON(http.StatusForbidden, gin.H{"error": "forbidden access level."})

}
