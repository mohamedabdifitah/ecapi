package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mohamedabdifitah/ecapi/db"
	"github.com/mohamedabdifitah/ecapi/service"
	"github.com/mohamedabdifitah/ecapi/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllCustomers(c *gin.Context) {
	var customer *db.Customer
	customers, err := customer.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, customers)
}
func GetCustomer(c *gin.Context) {
	var id string
	id = c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.String(400, "Invalid Id")
		return
	}
	customer := db.Customer{
		Id: objectId,
	}
	err = customer.GetById()
	if err != nil {
		c.String(200, err.Error())
		return
	}
	c.JSON(200, customer)
}
func SignUpCustomerWithEmail(c *gin.Context) {
	var body *SignUpWithEmailBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(400, err)
		return
	}
	customer := db.Customer{
		Email:    body.Email,
		Password: body.Password,
		Metadata: db.AccountMetadata{
			Provider: "email",
		},
	}
	res, err := customer.Save()
	if err != nil {
		eres := utils.HandlerError(err)
		c.JSON(eres.Status, eres.Message)
		return
	}
	var info map[string]string = make(map[string]string)
	info["rec"] = customer.Email
	info["type"] = "email"
	go service.ProduceMessage("", "verification", "", info)
	c.JSON(201, res.InsertedID)
}
func CompleteSignUp(c *gin.Context) {
	var body *CustomerBody
	var id string = c.Param("ssid")
	objecId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.String(400, "Invalid customer id")
		return
	}
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(400, err)
		return
	}
	customer := db.Customer{
		Id:         objecId,
		Address:    body.Address,
		GivenName:  body.GivenName,
		FamilyName: body.FamilyName,
		Phone:      body.Phone,
	}
	res, err := customer.Update()
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, res)
}
func UpdateCustomer(c *gin.Context) {
	var body *CustomerBody
	var id string = c.GetHeader("ssid")
	objecId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.String(400, "Invalid customer id")
		return
	}
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(400, err)
		return
	}
	customer := db.Customer{
		Id:         objecId,
		Address:    body.Address,
		GivenName:  body.GivenName,
		FamilyName: body.FamilyName,
		Phone:      body.Phone,
	}
	res, err := customer.Update()
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, res)

}
func DeleteCustomer(c *gin.Context) {
	var id string = c.Param("id")
	objecId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.String(400, "Invalid customer id")
		return
	}
	customer := db.Customer{
		Id: objecId,
	}
	res, err := customer.Delete()
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, res)

}
func ChangeCustomerPassword(c *gin.Context) {
	var id string = c.GetHeader("ssid")
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
	customer := db.Customer{
		Id: objecId,
	}
	res := customer.ChangePassword(body.OldPassword, body.NewPassword)
	if res != nil {
		res.Error(c)
		return
	}
	c.String(200, "successfully changed password")
}
func ChangeCustomerEmail(c *gin.Context) {
	var id string = c.GetHeader("ssid")
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
	customer := db.Customer{
		Id: objecId,
	}
	res, ErrorResponse := customer.ChangeEmail(body.OldEmail, body.NewEmail)
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
func CustomerEmailLogin(c *gin.Context) {
	var body *EmailLoginBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(400, err)
		return
	}
	device := db.Device{
		DeviceId: c.GetHeader("device_id"),
		Kind:     c.GetHeader("device_kind"),
	}
	tokens, ErrorResponse := db.CustomerLoginCheck(body.Email, body.Password, device)
	if ErrorResponse != nil {
		ErrorResponse.Error(c)
		return
	}
	c.JSON(200, tokens)
}
func ChangeCustomerDevice(c *gin.Context) {
	var id string = c.GetHeader("ssid")
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
	res, err := db.UpdateCustomer(query, update)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, res)
}
func ChangeCustomerProfile(c *gin.Context) {
	var id string = c.GetHeader("ssid")
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
	confirm, err := db.UpdateCustomer(query, update)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, confirm)
}
func ChangeCustomerWebhooks(c *gin.Context) {
	var id string = c.GetHeader("ssid")
	objectid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.String(400, "invalid id")
		return
	}
	var body map[string]string
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	webhook, ok := body["webhook"]
	if !ok {
		c.String(400, "Required webhook")
		return
	}
	query := bson.M{"_id": objectid}
	// no set
	change := bson.D{
		{
			Key: "$set", Value: bson.D{
				{
					Key: "metadata.webhook_endpoint", Value: webhook,
				},
			},
		},
	}
	res, err := db.UpdateCustomer(query, change)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, res)
}
