package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mohamedabdifitah/ecapi/utils"
	"golang.org/x/exp/slices"
)

// Reads The Acces Token => expires  then Searches Role
// Reads the refresh token checks the token_v validity generets access token => returns to header => to retry the request
func AuthorizeRolesMiddleware(permissions []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHeader := c.GetHeader("Authorization")
		id := c.GetHeader("ssid")
		if tokenHeader == "" && len(strings.Split(tokenHeader, " ")) < 2 {
			c.String(401, "authorization key not found")
			c.Abort()
			return
		}
		tokenString := strings.Split(tokenHeader, " ")[1]
		token, err := utils.VerifyAccessToken(tokenString)
		if err != nil {
			c.String(401, err.Error())
			c.Abort()
			return
		}
		fmt.Println(token.Id)
		if id != token.Id {
			c.String(403, "Authentication Error")
			c.Abort()
			return
		}
		if len(permissions) == 0 {
			c.Next()
			return

		}
		if slices.Contains(permissions, token.Role) {
			c.Next()
			return
		} else {
			c.String(403, "Access Denied")
			c.Abort()
			return
		}
	}
}
func TokenErrorHandler(c *gin.Context, err error) {

	switch err {
	// case jwt.ErrTokenExpired:
	case jwt.ErrInvalidKey:
		c.String(401, "Invalid token")
		c.Abort()
		break
	// case jwt.ErrTokenExpired:
	default:
		c.JSON(401, err)
		// break
		c.Abort()
		break
	}
}
