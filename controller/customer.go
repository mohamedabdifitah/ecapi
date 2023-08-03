package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohamedabdifitah/ecapi/db"
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
func SingUpCustomerWithEmail(c *gin.Context) {
	var body *SingUpCustomerWithEmailBody
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
		c.JSON(500, err.Error())
		return
	}
	c.JSON(201, res)
}
func UpdateCustomer(c *gin.Context) {
	var body *CustomerBody
	var id string = c.Param("id")
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
	customer := db.Customer{
		Id: objecId,
	}
	res, ErrorResponse := customer.ChangeEmail(body.OldEmail, body.NewEmail)
	if ErrorResponse != nil {
		ErrorResponse.Error(c)
	}
	if res.ModifiedCount == 0 {
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
