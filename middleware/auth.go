package middleware

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mohamedabdifitah/ecapi/db"
	"github.com/mohamedabdifitah/ecapi/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slices"
)

// Reads The Acces Token and Refresh Token and authorize the roles also authorize the header of ssid of user id
// Reads the refresh token checks the token_v validity generets access token => returns to header => to retry the request
func AuthorizeRolesMiddleware(permissions []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHeader := c.GetHeader("Authorization")
		ReftokenHeader := c.GetHeader("refresh_token") // refresh token header t_v
		id := c.GetHeader("ssid")
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.String(401, "invalid id")
			c.Abort()
			return
		}
		if tokenHeader == "" && len(strings.Split(tokenHeader, " ")) < 2 {
			c.String(401, "authorization key not found")
			c.Abort()
			return
		}
		tokenString := strings.Split(tokenHeader, " ")[1]
		token, err := utils.VerifyAccessToken(tokenString)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				if ReftokenHeader == "" && len(strings.Split(ReftokenHeader, " ")) < 2 {
					c.String(401, "authorization key not found")
					c.Abort()
					return
				}
				reftokenString := strings.Split(ReftokenHeader, " ")[1]
				reftoken, err := utils.VerifyRefereshToken(reftokenString)
				if errors.Is(err, jwt.ErrTokenExpired) {
					c.String(401, "token expired , login again.")
					c.Abort()
					return

				}
				switch token.Role {
				case "customer":
					customer := db.Customer{
						Id: objectId,
					}
					err = customer.GetById()
					if err != nil {
						c.String(401, "user not found")
						c.Abort()
						return
					}
					if reftoken.TokenVersion != customer.Metadata.TokenVersion {
						c.String(403, "Access Denied , please login again")
						c.Abort()
						return
					}
					tokenString, err = utils.GenerateAccessToken(customer.Email, customer.Id, db.Roles[0])
					if err != nil {
						c.String(401, "Error generating access token")
						c.Abort()
						return
					}
					token, err = utils.VerifyAccessToken(tokenString)
					if err != nil {
						c.String(401, err.Error())
						c.Abort()
						return
					}
					c.Header("Authorization", "Bearer "+tokenString)
					return
				case "merchant":
					merchant := db.Merchant{
						Id: objectId,
					}
					err = merchant.GetById()
					if err != nil {
						c.String(401, "user not found")
						c.Abort()
						return
					}
					if reftoken.TokenVersion != merchant.Metadata.TokenVersion {
						c.String(403, "Access Denied , please login again")
						c.Abort()
						return
					}
					tokenString, err = utils.GenerateAccessToken(merchant.BusinessPhone, merchant.Id, db.Roles[2])
					if err != nil {
						c.String(401, "Error generating access token")
						c.Abort()
						return
					}
					token, err = utils.VerifyAccessToken(tokenString)
					if err != nil {
						c.String(401, err.Error())
						c.Abort()
						return
					}
					c.Header("Authorization", "Bearer "+tokenString)
					return
					// break
				}
			} else {
				c.String(401, err.Error())
				c.Abort()
				return

			}
		}
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
func Authenticate(token string) {

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
