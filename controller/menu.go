package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohamedabdifitah/ecapi/db"
	"github.com/mohamedabdifitah/ecapi/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetMenus(c *gin.Context) {

	menu := &db.Menu{}
	menus, err := menu.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, menus)
}
func GetMenu(c *gin.Context) {
	id := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.String(400, "Invalid Id")
		return
	}
	menu := db.Menu{
		Id: objectId,
	}
	err = menu.GetById()
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			c.String(200, "menu not found")
			return
		}
		c.String(200, err.Error())
		return
	}
	c.JSON(200, menu)
}
func CreateMenu(c *gin.Context) {
	var body *MenuBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	var merchantid string = c.GetHeader("ssid")
	menu := &db.Menu{
		Title:              body.Title,
		Description:        body.Description,
		Status:             body.Status,
		Category:           body.Category,
		Price:              body.Price,
		Attributes:         body.Attributes,
		MerchantExternalId: merchantid,
		Reciepe:            body.Reciepe,
		EstimateTime:       body.EstimateTime,
		Images:             body.Images,
	}
	res, err := menu.Create()
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(201, res)
}
func UpdateMenu(c *gin.Context) {
	var body *MenuBody
	var id string = c.Param("id")
	objecId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.String(400, "invalid id")
		return
	}
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	menu := &db.Menu{
		Id:           objecId,
		Title:        body.Title,
		Description:  body.Description,
		Category:     body.Category,
		Reciepe:      body.Reciepe,
		Price:        body.Price,
		Status:       body.Status,
		Attributes:   body.Attributes,
		EstimateTime: body.EstimateTime,
		Images:       body.Images,
	}
	res, err := menu.Update()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if res.MatchedCount > 0 && res.ModifiedCount == 0 {

		c.String(400, "menu is same as original")
		return
	}
	if res.MatchedCount == 0 && res.ModifiedCount == 0 {
		c.String(200, "menu not found")
		return
	}
	c.String(200, "menu updated successfully")
}
func DeleteMenu(c *gin.Context) {
	var id string = c.Param("id")
	objecId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.String(400, "Invalid id")
		return
	}
	menu := db.Menu{
		Id: objecId,
	}
	res, err := menu.Delete()
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, res)
}
func GetMenuFromMerchant(c *gin.Context) {
	id := c.Param("id")
	menu := db.Menu{
		MerchantExternalId: id,
	}
	menues, err := menu.GetFromMerchant()
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, menues)
}
func AddMenuImages(c *gin.Context) {
	id := c.Param("id")
	objectid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.String(400, err.Error())
		return
	}
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
	menu := db.Menu{
		Id:     objectid,
		Images: response,
	}
	confirm, err := menu.SetImages()
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, confirm)
}
