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
		c.JSON(500, err)
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
	}
	res, err := customer.Save()
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(201, res)
}
