package controller

import (
	"fmt"
	"strconv"
	"time"

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
			return
		}
		items = append(items, db.Item{
			Quantity:       item.Quantity,
			ItemExternalId: item.ItemExternalId,
		})
	}
	_, err = primitive.ObjectIDFromHex(body.PickUpExternalId)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	_, err = primitive.ObjectIDFromHex(body.DropOffExteranlId)
	if err != nil {
		c.String(400, err.Error())
		return
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
	res, errres := order.PlaceOrder()
	if errres != nil {
		if errres.Type == "string" {
			c.String(errres.Status, errres.Message.Error())
			return
		}
		c.JSON(errres.Status, errres.Message)
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
func GetOrderByLocation(c *gin.Context) {
	longtitude, err := strconv.ParseFloat(c.Query("lang"), 64)
	latitude, err := strconv.ParseFloat(c.Query("lat"), 64)
	mindist, err := strconv.ParseInt(c.Query("mindist"), 0, 64)
	maxdist, err := strconv.ParseInt(c.Query("maxdist"), 0, 64)
	fmt.Println(maxdist, mindist)
	if err != nil {
		c.String(400, err.Error())
	}
	var order db.Order
	location := []float64{
		longtitude,
		latitude,
	}
	orders, err := order.GetByLocation(location, maxdist, mindist)
	if err != nil {
		c.String(500, err.Error())
	}
	c.JSON(200, orders)
}
func MerchantOrderAccept(c *gin.Context) {
	OrderId := c.Param("id")
	var merchantid string = c.GetHeader("ssid")
	objectid, err := primitive.ObjectIDFromHex(OrderId)
	if err != nil {
		c.String(400, "invalid Order Id")
		return
	}
	query := bson.M{"_id": objectid, "pickup_external_id": merchantid}
	change := bson.D{{Key: "$set", Value: bson.D{
		{Key: "metadata.update_at", Value: time.Now()},
		{Key: "stage", Value: "accepted"},
	}}}
	res, err := db.AccpetOrderBy(query, change, "merchant-accpted")
	if err != nil {
		c.String(400, err.Error())
		return
	}
	c.JSON(200, res)
}
func AssignOrderToDriver(c *gin.Context) {
	OrderId := c.Param("oid")
	DriverId := c.Param("did")
	DriverObjId, err := primitive.ObjectIDFromHex(DriverId)
	if err != nil {
		c.String(400, "invalid driver Id")
		return
	}

	objectid, err := primitive.ObjectIDFromHex(OrderId)
	if err != nil {
		c.String(400, "invalid Order Id")
		return
	}
	res, erres := db.AssignOrderToDriver(objectid, DriverObjId)
	if erres != nil {
		if erres.Type == "string" {
			c.String(erres.Status, erres.Message.Error())
			return
		}
		c.JSON(erres.Status, erres.Message.Error())
		return
	}
	c.JSON(200, res)
}
func AccpetOrderByDriver(c *gin.Context) {
	DriverId := c.GetHeader("ssid")
	OrderId := c.Param("id")
	DriverObjId, err := primitive.ObjectIDFromHex(DriverId)
	if err != nil {
		c.String(400, "invalid driver Id")
		return
	}

	objectid, err := primitive.ObjectIDFromHex(OrderId)
	if err != nil {
		c.String(400, "invalid Order Id")
		return
	}
	res, erres := db.AssignOrderToDriver(objectid, DriverObjId)
	if erres != nil {
		if erres.Type == "string" {
			c.String(erres.Status, erres.Message.Error())
			return
		}
		c.JSON(erres.Status, erres.Message.Error())
		return
	}
	c.JSON(200, res)
}
func CancelOrder(c *gin.Context) {
	OrderId := c.Param("id")
	var customerid string = c.GetHeader("ssid")
	objectid, err := primitive.ObjectIDFromHex(OrderId)
	if err != nil {
		c.String(400, "invalid Order Id")
		return
	}
	query := bson.M{"_id": objectid, "dropoff_external_id": customerid}
	change := bson.D{{Key: "$set", Value: bson.D{
		{Key: "metadata.update_at", Value: time.Now()},
		{Key: "stage", Value: "cancel"},
		{Key: "cancel_reason", Value: ""},
	}}}
	res, err := db.AccpetOrderBy(query, change, "order_canceled")
	if err != nil {
		c.String(400, err.Error())
		return
	}
	c.JSON(200, res)
}
func RejectOrderByMerchant(c *gin.Context) {
	OrderId := c.Param("id")
	var merchantid string = c.GetHeader("ssid")
	var reason string = "CANCEL_FROM_MERCHANT"
	objectid, err := primitive.ObjectIDFromHex(OrderId)
	if err != nil {
		c.String(400, "invalid Order Id")
		return
	}
	query := bson.M{"_id": objectid, "pickup_external_id": merchantid}
	change := bson.D{{Key: "$set", Value: bson.D{
		{Key: "metadata.update_at", Value: time.Now()},
		{Key: "stage", Value: "cancel"},
		{Key: "cancel_reason", Value: reason},
	}}}
	var info map[string]interface{}
	info["id"] = objectid
	info["cancel_reason"] = reason
	info["merchant_id"] = merchantid
	res, err := db.CancelOrder(query, change, "order_canceled", info)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	c.JSON(200, res)
}
