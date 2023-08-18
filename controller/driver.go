package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mohamedabdifitah/ecapi/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllDrivers(c *gin.Context) {
	var driver *db.Driver
	drivers, err := driver.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, drivers)
}
func GetDriver(c *gin.Context) {
	var id string
	id = c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.String(400, "Invalid Id")
		return
	}
	driver := db.Driver{
		Id: objectId,
	}
	err = driver.GetById()
	if err != nil {
		c.String(200, err.Error())
		return
	}
	c.JSON(200, driver)
}
func SignUpDriverWithPhone(c *gin.Context) {
	var body *SignUpWithDriverBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(400, err)
		return
	}
	driver := db.Driver{
		Phone:    body.Phone,
		Email:    body.Emai,
		Password: body.Password,
		Metadata: db.AccountMetadata{
			Provider: "phone",
		},
	}
	res, err := driver.Save()
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(201, res)
}
func UpdateDriver(c *gin.Context) {
	var body *DriverBody
	var id string = c.Param("id")
	objecId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.String(400, "Invalid driver id")
		return
	}
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(400, err)
		return
	}
	driver := db.Driver{
		Id:          objecId,
		Address:     body.Address,
		GivenName:   body.GivenName,
		VehicleType: body.VehicleType,
		// Age:         body.Age,
	}
	res, err := driver.Update()
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, res)

}
func DeleteDriver(c *gin.Context) {
	var id string = c.Param("id")
	objecId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.String(400, "Invalid driver id")
		return
	}
	driver := db.Driver{
		Id: objecId,
	}
	res, err := driver.Delete()
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, res)

}
func ChangeDriverPassword(c *gin.Context) {
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
	driver := db.Driver{
		Id: objecId,
	}
	res := driver.ChangePassword(body.OldPassword, body.NewPassword)
	if res != nil {
		res.Error(c)
		return
	}
	c.String(200, "successfully changed password")
}
func ChangeDriverEmail(c *gin.Context) {
	var id string = c.Param("id")
	var body *ChangeEmaildBody
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
	driver := db.Driver{
		Id: objecId,
	}
	res, ErrorResponse := driver.ChangeEmail(body.OldEmail, body.NewEmail)
	if ErrorResponse != nil {
		ErrorResponse.Error(c)
	}
	if res.MatchedCount > 0 && res.ModifiedCount == 0 {

		c.String(400, "email is same as original")
		return
	}
	if res.MatchedCount == 0 && res.ModifiedCount == 0 {
		c.String(200, "email not found")
		return
	}
	c.String(200, "email updated successfully")
}
func DriverPhoneLogin(c *gin.Context) {
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
	tokens, ErrorResponse := db.DriverLoginCheck(body.Phone, body.Password, device)
	if ErrorResponse != nil {
		ErrorResponse.Error(c)
		return
	}
	c.JSON(200, tokens)
}
func ChangeDriverPhone(c *gin.Context) {
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
	driver := db.Driver{
		Id: objecId,
	}
	res, err := driver.ChangePhone(body.OldPhone, body.NewPhone)
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
func ChangeDriverDevice(c *gin.Context) {
	id := c.Param("id")
	var body map[string]interface{}
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
		{Key: "device", Value: body["device"]},
	}}}
	res, err := db.UpdateDriver(query, update)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, res)
}
