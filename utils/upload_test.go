package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/file", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.String(400, err.Error())
			return
		}
		response, ErrorResponse := UploadFiles(".", file)
		if err != nil {
			c.String(ErrorResponse.StatusCode, ErrorResponse.Reason.Error())
			return
		}
		c.JSON(200, response)
	})
	return r
}
func TestUploadFiles(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	file := strings.NewReader("hello world!")
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "example.txt")
	if err != nil {
		t.Errorf("Error copying file , please try again")
	}
	_, err = io.Copy(part, file)
	if err != nil {
		t.Errorf("this file is broken")
	}
	err = writer.Close()
	if err != nil {
		t.Errorf("there is error , please try again")
	}
	req, _ := http.NewRequest("POST", "/file", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)
	var one []string
	err = json.Unmarshal(w.Body.Bytes(), &one)
	if err != nil {
		t.Error(err)
	}
	savedfile, err := os.Open(one[0])
	if err != nil {
		t.Error("file is not found , after uploading")
	}
	savedfile.Close()
	err = os.Remove(savedfile.Name())
	if err != nil {
		t.Error(err)
	}
}
