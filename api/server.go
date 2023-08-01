package api

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	server          *gin.Engine
	CustomerRouter  *gin.RouterGroup
	OrderRouter     *gin.RouterGroup
	ResturantRouter *gin.RouterGroup
	DriverRouter    *gin.RouterGroup
)

func Initserver() {
	server = gin.New()
	port := os.Getenv("PORT")
	server.Use(gin.Recovery(), gin.Logger())
	if os.Getenv("GIN_MODE") == "release" {
		server.Run(fmt.Sprintf(":" + port)) // listen
	}
	InitRoutes(server)
	// this is fix for windows defender popup
	server.Run(fmt.Sprintf("localhost:" + port))
}
func InitRoutes(server *gin.Engine) {
	path := server.RouterGroup
	path.BasePath()
	CustomerRouter = path.Group("/customer")
	DriverRouter = path.Group("/driver")
	ResturantRouter = path.Group("/resturant")
	OrderRouter = path.Group("/order")
	RouterDefinition()
}