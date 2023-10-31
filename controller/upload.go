package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mohamedabdifitah/ecapi/utils"
)

func UploadFile(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.String(400, err.Error())
		return
	}
	files := form.File["photos"]
	response, ErrorResponse := utils.UploadFiles("", files...)
	if err != nil {
		c.String(ErrorResponse.StatusCode, ErrorResponse.Reason.Error())
		return
	}
	c.JSON(200, response)
}
