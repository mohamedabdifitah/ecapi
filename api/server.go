package api

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	server *gin.Engine
)

func Initserver() {
	server = gin.New()
	port := os.Getenv("PORT")
	server.Use(gin.Recovery(), gin.Logger())
	InitRoutes(server)
	if os.Getenv("GIN_MODE") == "release" {
		server.Run(fmt.Sprintf(":" + port)) // listen
	}
	// this is fix for windows defender popup
	server.Run(fmt.Sprintf("localhost:" + port))
}
func InitRoutes(server *gin.Engine) {
	path := server.RouterGroup
	path.BasePath()
	if os.Getenv("APP_ENV") == "development" {
		path.Static("/doc", "./doc")
	}
	path.Static("/assets", "./assets")
	for _, r := range Routes {
		r.register()
	}
}
