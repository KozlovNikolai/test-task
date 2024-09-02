package main

import (
	_ "github.com/KozlovNikolai/test-task/docs"
	"github.com/KozlovNikolai/test-task/internal/pkg/config"
	"github.com/KozlovNikolai/test-task/internal/server"
)

// @title 	Shop Service API
// @version	1.0
// @description An Shop service API in Go using Gin framework
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @host 	localhost:8443
// @BasePath /
func main() {
	config.MustLoad()
	server := server.NewServer()

	server.Run()
}
