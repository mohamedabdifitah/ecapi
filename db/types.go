package db

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Status  int
	Message error
	Type    string
}
type TokenResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

var Roles []string = []string{
	"customer",
	"driver",
	"merchant",
}

func (e *ErrorResponse) Error(c *gin.Context) {
	if e.Type == "string" {
		c.String(e.Status, e.Message.Error())
		return
	}
	c.JSON(e.Status, e.Message)
}
