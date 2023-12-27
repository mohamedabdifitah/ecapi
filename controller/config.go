package controller

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func SetupConfig() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	return conf
}
func SiginWithGoogle(c *gin.Context) {
	googleConfig := SetupConfig()
	// Redirect user to Google's consent page to ask for permission
	// for the scopes specified above.
	url := googleConfig.AuthCodeURL("state")
	c.Redirect(http.StatusSeeOther, url)
}
