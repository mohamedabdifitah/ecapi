package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mohamedabdifitah/ecapi/db"
	"github.com/mohamedabdifitah/ecapi/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetMenus(c *gin.Context) {
	// var body *MenuBody
	// err := c.ShouldBindJSON(&body)
	// if err != nil {
	// 	c.String(400, err.Error())
	// 	return
	// }
	menu := &db.Menu{
		// Title:              body.Title,
		// Description:        body.Description,
		// Category:           body.Category,
		// Reciepe:            body.Reciepe,
		// MerchantExternalId: body.MerchantExternalId,
		// Price:              body.Price,
		// Discount:           body.Discount,
		// Status:             body.Status,
	}
	drivers, err := menu.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, drivers)
}
func GetMenu(c *gin.Context) {
	var id string
	id = c.Param("id")
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
	var body *CreateMenuBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	menu := &db.Menu{
		Title:              body.Title,
		Description:        body.Description,
		Status:             body.Status,
		Category:           body.Category,
		Price:              body.Price,
		Attributes:         body.Attributes,
		MerchantExternalId: body.MerchantExternalId,
		Reciepe:            body.Reciepe,
		Barcode:            body.Barcode,
		Discount:           body.Discount,
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
		Id:                 objecId,
		Title:              body.Title,
		Description:        body.Description,
		Category:           body.Category,
		Reciepe:            body.Reciepe,
		MerchantExternalId: body.MerchantExternalId,
		Price:              body.Price,
		Discount:           body.Discount,
		Status:             body.Status,
		Attributes:         body.Attributes,
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
func PutImageMenues(c *gin.Context) {
	// Multipart form
	id := c.Param("id")
	objecid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	file, err := c.FormFile("upload")
	if err != nil {
		c.String(400, err.Error())
		return
	}
	log.Println(file.Filename)

	// Upload the file to specific dst.
	c.SaveUploadedFile(file, fmt.Sprintf("%s/%s", os.Getenv("STATIC_URL"), file.Filename))
	menu := db.Menu{
		Id:     objecid,
		Images: []string{fmt.Sprintf("%s/%s", os.Getenv("STATIC_URL"), file.Filename)},
	}
	res, err := menu.SetImages()
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, res)
}
func AddMenuImages(c *gin.Context) {
	id := c.Param("id")
	objectid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	form, _ := c.MultipartForm()
	files := form.File["photos"]
	res, errres := utils.UploadPhotos(files)
	if errres != nil {
		c.String(errres.StatusCode, errres.Reason.Error())
		return
	}
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(res.Body)
	if err != nil {
		c.String(500, "error uploading file, please try again")
		return
	}
	res.Body.Close()
	var photosUri []string
	err = json.Unmarshal(body.Bytes(), &photosUri)
	if err != nil {
		c.String(500, "error unmarshalling photos , please try again")
		return
	}
	menu := db.Menu{
		Id:     objectid,
		Images: photosUri,
	}
	response, err := menu.SetImages()
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, response)
}
