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
	}
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(400, err)
		return
	}
	customer := db.Customer{
		Id:         objecId,
		Profile:    body.Profile,
		Address:    body.Address,
		GivenName:  body.GivenName,
		FamilyName: body.FamilyName,
	}
	res, err := customer.Update()
	if err != nil {
		c.JSON(500, err)
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
	}
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(400, err)
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