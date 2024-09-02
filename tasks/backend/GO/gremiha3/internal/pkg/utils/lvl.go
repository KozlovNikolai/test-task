package utils

import "github.com/gin-gonic/gin"

func GetLevel(c *gin.Context) (int, string, string) {
	id, _ := c.Get("auth_id")
	login, _ := c.Get("auth_login")
	role, _ := c.Get("auth_role")
	return id.(int), login.(string), role.(string)
}
