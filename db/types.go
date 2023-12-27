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
type GoogleUserInfoType struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`        // full name
	GivenName     string `json:"given_name"`  // first name
	FamilyName    string `json:"family_name"` // last name
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
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
