package api

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	server         *gin.Engine
	CustomerRouter *gin.RouterGroup
	OrderRouter    *gin.RouterGroup
	MerchantRouter *gin.RouterGroup
	DriverRouter   *gin.RouterGroup
	ReviewRouter   *gin.RouterGroup
	MenuRouter     *gin.RouterGroup
	Router         *gin.RouterGroup
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
		path.Static("/static", "./doc")
	}
	path.Static("/assets", "./assets")
	CustomerRouter = path.Group("/customer")
	DriverRouter = path.Group("/driver")
	MerchantRouter = path.Group("/merchant")
	OrderRouter = path.Group("/order")
	MenuRouter = path.Group("/menu")
	ReviewRouter = path.Group("/review")
	RouterDefinition()
}
