package db

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Status  int
	Message error
	Type    string
}

func (e *ErrorResponse) Error(c *gin.Context) {
	if e.Type == "string" {
		c.String(e.Status, e.Message.Error())
		return
	}
	c.JSON(e.Status, e.Message)
}
