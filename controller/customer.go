package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllCustomers(c *gin.Context) {
	c.JSON(http.StatusOK, "hello world")
}
