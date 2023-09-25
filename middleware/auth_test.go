package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mohamedabdifitah/ecapi/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/hello", AuthorizeRolesMiddleware([]string{}), func(c *gin.Context) {
		c.String(200, "hello world")
	})
	return r
}
func TestAuthorizeRolesMiddleware(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	w2 := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/hello", nil)
	req2, _ := http.NewRequest("GET", "/hello", nil)
	var id primitive.ObjectID = primitive.NewObjectID()
	token, err := utils.GenerateAccessToken("test@gmail.com", id, "")
	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req2.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("ssid", id.Hex())
	router.ServeHTTP(w, req)
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Body.String(), "hello world")
	router.ServeHTTP(w2, req2)
	assert.Equal(t, w2.Body.String(), "not authorized , ssid not found")
	assert.Equal(t, w2.Code, 403)

}
