package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mohamedabdifitah/ecapi/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetOrderByid(c *gin.Context) {
	var id string = c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.String(400, "Invalid Id")
		return
	}
	order := db.Order{
		Id: objectId,
	}
	err = order.GetById()
	if err != nil {
		c.String(200, err.Error())
		return
	}
	c.JSON(200, order)

}
func GetAllOrders(c *gin.Context) {
	filter := bson.D{}
	orders, err := db.GetOrdersBy(filter)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, orders)

}
func PlaceOrder(c *gin.Context) {
	var body *PlaceOrderBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	var items []db.Item
	for _, item := range body.Items {
		_, err = primitive.ObjectIDFromHex(item.ItemExternalId)
		if err != nil {
			c.String(400, err.Error())
		}
		items = append(items, db.Item{
			Quantity:       item.Quantity,
			ItemExternalId: item.ItemExternalId,
		})
	}
	_, err = primitive.ObjectIDFromHex(body.PickUpExternalId)
	if err != nil {
		c.String(400, err.Error())
	}
	_, err = primitive.ObjectIDFromHex(body.DropOffExteranlId)
	if err != nil {
		c.String(400, err.Error())
	}
	order := db.Order{
		Items:              items,
		DropOffPhone:       body.DropOffPhone,
		DropOffContactName: body.DropOffContactName,
		DropOffExteranlId:  body.DropOffExteranlId,
		DropOffAddress:     body.DropOffAddress,
		DroOffLocation: db.Location{
			Type:        "Point",
			Coordinates: body.DroOffLocation,
		},
		DropOffInstruction: body.DropOffInstruction,
		PickUpExternalId:   body.PickUpExternalId,
		Type:               body.Type,
	}
	res, err := order.PlaceOrder()
	if err != nil {
		c.String(400, err.Error())
		return
	}
	c.JSON(201, res)

}
func GetOrderByCustomer(c *gin.Context) {
	id := c.Param("id")
	filter := bson.D{
		{Key: "dropoff_external_id", Value: id},
	}
	orders, err := db.GetOrdersBy(filter)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, orders)
}
func GetOrderByMerchant(c *gin.Context) {
	id := c.Param("id")
	filter := bson.D{
		{Key: "pickup_external_id", Value: id},
	}
	orders, err := db.GetOrdersBy(filter)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, orders)
}
func GetOrderByDriver(c *gin.Context) {
	id := c.Param("id")
	filter := bson.D{
		{Key: "driver_external_id", Value: id},
	}
	orders, err := db.GetOrdersBy(filter)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, orders)
}
func DriverAcceptOrder(c *gin.Context) {
	// driver information
	// order information => weighting order information => vehicle type
	// notifying redis subscribers
}
