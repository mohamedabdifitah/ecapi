package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mohamedabdifitah/ecapi/db"
	"github.com/mohamedabdifitah/ecapi/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllMerchants(c *gin.Context) {
	var merchant *db.Merchant
	merchants, err := merchant.GetAll()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, merchants)
}
func GetMerchant(c *gin.Context) {
	var id string
	id = c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.String(400, "Invalid Id")
		return
	}
	merchant := db.Merchant{
		Id: objectId,
	}
	err = merchant.GetById()
	if err != nil {
		c.String(200, err.Error())
		return
	}
	c.JSON(200, merchant)
}
func SignUpMerchantWithPhone(c *gin.Context) {
	var body *SignUpMerchantWithPhoneBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	merchant := db.Merchant{
		BusinessPhone: body.BusinessPhone,
		Password:      body.Password,
		BusinessName:  body.BusinessName,
		Location: db.Location{
			Type:        "Point",
			Coordinates: body.Location,
		},
		Metadata: db.AccountMetadata{
			Provider: "phone",
		},
	}
	res, erres := merchant.Save()
	if erres != nil {
		if erres.Type == "string" {
			c.String(erres.Status, erres.Message.Error())
			return
		}
		c.JSON(erres.Status, erres.Message)
		// c.JSON(500, err.Error())
		return
	}
	c.JSON(201, res)
}
func UpdateMerchant(c *gin.Context) {
	var body *MerchantBody
	var id string = c.Param("id")
	objecId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.String(400, "Invalid merchant id")
		return
	}
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(400, err)
		return
	}
	merchant := db.Merchant{
		Id: objecId,
		Location: db.Location{
			Type:        "Point",
			Coordinates: body.Location,
		},
		Address:           body.Address,
		BusinessName:      body.BusinessName,
		BusinessEmail:     body.BusinessEmail,
		TimeOperatorStart: body.TimeOperationStart,
		TimeOperatorEnd:   body.TimeOperationEnd,
		Category:          body.Category,
	}
	res, err := merchant.Update()
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, res)

}
func DeleteMerchant(c *gin.Context) {
	var id string = c.Param("id")
	objecId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.String(400, "Invalid merchant id")
		return
	}
	merchant := db.Merchant{
		Id: objecId,
	}
	res, err := merchant.Delete()
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, res)

}
func ChangeMerchantPassword(c *gin.Context) {
	var id string = c.Param("id")
	var body *ChangePasswordBody
	objecId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.String(400, "invalid id")
		return
	}
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(400, err)
		return
	}
	if body.OldPassword == body.NewPassword {
		c.String(400, "password is same as original")
		return
	}
	merchant := db.Merchant{
		Id: objecId,
	}
	res := merchant.ChangePassword(body.OldPassword, body.NewPassword)
	if res != nil {
		res.Error(c)
		return
	}
	c.String(200, "successfully changed password")
}
func ChangeMerchantPhone(c *gin.Context) {
	var id string = c.Param("id")
	var body *ChangePhonedBody
	objecId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.String(400, "invalid id")
		return
	}
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(400, err)
		return
	}
	merchant := db.Merchant{
		Id: objecId,
	}
	res, err := merchant.ChangeBusinessPhone(body.OldPhone, body.NewPhone)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	if res.MatchedCount > 0 && res.ModifiedCount == 0 {

		c.String(400, "phone is same as original")
		return
	}
	if res.MatchedCount == 0 && res.ModifiedCount == 0 {
		c.String(200, "phone not found")
		return
	}
	c.String(200, "phone updated successfully")
}
func MerchantPhoneLogin(c *gin.Context) {
	var body *PhoneLoginBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(400, err)
		return
	}
	device := db.Device{
		DeviceId: c.GetHeader("device_id"),
		Kind:     c.GetHeader("device_kind"),
	}
	tokens, ErrorResponse := db.MerchantLoginCheck(body.Phone, body.Password, device)
	if ErrorResponse != nil {
		ErrorResponse.Error(c)
		return
	}
	c.JSON(200, tokens)
}
func GetMerchantByLocation(c *gin.Context) {
	longtitude, err := strconv.ParseFloat(c.Query("lang"), 64)
	latitude, err := strconv.ParseFloat(c.Query("lat"), 64)
	mindist, err := strconv.ParseInt(c.Query("mindist"), 0, 64)
	maxdist, err := strconv.ParseInt(c.Query("maxdist"), 0, 64)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	var merchant db.Merchant
	location := []float64{
		longtitude,
		latitude,
	}
	merchants, err := merchant.GetMerchantByLocation(location, maxdist, mindist)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, merchants)
}
func ChangeMerchantDevice(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Device db.Device `json:"device" binding:"required"`
	}
	objectid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.String(400, "Invalid id")
		return
	}
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	query := bson.M{
		"_id": objectid,
	}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "metadata.update_at", Value: time.Now()},
		{Key: "device", Value: body.Device},
	}}}
	res, err := db.UpdateMerchant(query, update)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, res)
}
func ChangeMerchantProfile(c *gin.Context) {
	id := c.Param("id")
	objectid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.String(400, "Invalid id ")
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.String(400, err.Error())
		return
	}
	response, ErrorResponse := utils.UploadFiles("", file)
	if err != nil {
		c.String(ErrorResponse.StatusCode, ErrorResponse.Reason.Error())
		return
	}
	query := bson.M{"_id": objectid}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "metadata.update_at", Value: time.Now()},
		{Key: "profile", Value: response[0]},
	}}}
	confirm, err := db.UpdateMerchant(query, update)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, confirm)
}
