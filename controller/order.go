package controller

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mohamedabdifitah/ecapi/db"
	"github.com/mohamedabdifitah/ecapi/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slices"
)

var CancelReason = []string{"CANCEL_FROM_MERCHANT", "CANCEL_FROM_CUSTOMER", "CANCEL_FROM_ADMIN"}

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
func GetOrderByLocation(c *gin.Context) {
	longtitude, err := strconv.ParseFloat(c.Query("lang"), 64)
	latitude, err := strconv.ParseFloat(c.Query("lat"), 64)
	mindist, err := strconv.ParseInt(c.Query("mindist"), 0, 64)
	maxdist, err := strconv.ParseInt(c.Query("maxdist"), 0, 64)
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
	query := bson.D{
		{
			Key:   "_id",
			Value: objectid,
		},
		{
			Key:   "pickup_external_id",
			Value: merchantid,
		},
	}
	order, err := db.GetOrderBy(query)
	if err != nil {
		c.String(400, "order is not found")
		return
	}
	if order.Stage == "cancel" {
		c.String(400, "order is canceled")
		return
	}
	if order.Stage != "placed" {
		c.String(400, "order is already accepted")
		return
	}
	order.Stage = "accepted"
	filter := bson.M{"_id": order.Id}
	change := bson.D{{Key: "$set", Value: bson.D{
		{Key: "metadata.update_at", Value: time.Now()},
		{Key: "stage", Value: "accepted"},
	}}}
	res, err := db.UpdateOrder(filter, change)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	order, err = db.GetOrderBy(bson.D{{Key: "_id", Value: res.UpsertedID}})
	if err != nil {
		c.String(400, "order is not found")
		return
	}
	err = service.PublishTopic("merchant_accpted_order", order)
	if err != nil {
		c.String(500, "server error , please try again")
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
	var body map[string]string = make(map[string]string)
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	reason, ok := body["reason"]
	if !ok {
		c.String(400, "there is no reason")
		return
	}
	query := bson.D{
		{
			Key:   "_id",
			Value: objectid,
		},
		{
			Key:   "dropoff_external_id",
			Value: customerid,
		},
	}
	order, err := db.GetOrderBy(query)
	if err != nil {
		c.String(400, "order is not found")
		return
	}
	if slices.Contains([]string{"accepted", "placed"}, order.Stage) {

	}
	filter := bson.M{"_id": objectid, "dropoff_external_id": customerid}
	change := bson.D{{Key: "$set", Value: bson.D{
		{Key: "metadata.update_at", Value: time.Now()},
		{Key: "stage", Value: "cancel"},
		{Key: "cancel_reason", Value: CancelReason[1] + " " + reason},
	}}}
	res, err := db.UpdateOrder(filter, change)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	query = bson.D{
		{
			Key:   "_id",
			Value: objectid,
		},
	}
	order, err = db.GetOrderBy(query)
	err = service.PublishTopic("order_canceled", order)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.JSON(200, res)
}
func RejectOrderByMerchant(c *gin.Context) {
	OrderId := c.Param("id")
	var merchantid string = c.GetHeader("ssid")
	var body map[string]string = make(map[string]string)
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	reason, ok := body["reason"]
	if !ok {
		c.String(400, "there is no reason")
		return
	}
	objectid, err := primitive.ObjectIDFromHex(OrderId)
	if err != nil {
		c.String(400, "invalid Order Id")
		return
	}
	query := bson.D{
		{
			Key:   "_id",
			Value: objectid,
		},
		{
			Key:   "pickup_external_id",
			Value: merchantid,
		},
	}
	order, err := db.GetOrderBy(query)
	if err != nil {
		c.String(400, "order is not found")
		return
	}
	if order.Stage == "canceled" {
		c.String(400, "order is already canceled")
		return
	}
	if slices.Contains([]string{"pickuped", "deleivered"}, order.Stage) {
		c.String(400, "can't cancel order because it's already picked or delivered")
		return
	}
	filter := bson.M{"_id": objectid, "pickup_external_id": merchantid}
	change := bson.D{{Key: "$set", Value: bson.D{
		{Key: "metadata.update_at", Value: time.Now()},
		{Key: "stage", Value: "cancel"},
		{Key: "cancel_reason", Value: CancelReason[0] + " " + reason},
	}}}
	res, err := db.UpdateOrder(filter, change)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	query = bson.D{
		{
			Key:   "_id",
			Value: objectid,
		},
	}
	order, err = db.GetOrderBy(query)
	err = service.PublishTopic("order_canceled", order)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, res)
}
func DropOrderByDriver(c *gin.Context) {
	var orderId string = c.Param("id")
	objectid, err := primitive.ObjectIDFromHex(orderId)
	if err != nil {
		c.String(400, "invalid order id")
		return
	}
	var driverId string = c.GetHeader("ssid")
	res, erres := db.DropOrder(objectid, driverId)
	if erres != nil {
		if erres.Type == "string" {
			c.String(erres.Status, erres.Message.Error())
			return
		}
		c.JSON(erres.Status, erres.Message)
		return
	}
	c.JSON(200, res)
}
func ChangeOrderStatus(stage string) gin.HandlerFunc {
	if !slices.Contains([]string{"preparing", "ready", "pickuped"}, stage) {
		log.Fatal("unknow stage")
	}
	return func(c *gin.Context) {
		var id string = c.Param("id")
		var merid string = c.GetHeader("ssid")
		objectid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.String(400, err.Error())
			c.Abort()
			return
		}
		query := bson.M{"_id": objectid, "pickup_external_id": merid}
		change := bson.D{
			{
				Key: "$set", Value: bson.D{
					{
						Key: "stage", Value: stage,
					},
				},
			},
		}
		res, err := db.UpdateOrder(query, change)
		if err != nil {
			c.String(500, err.Error())
			c.Abort()
			return
		}
		order, err := db.GetOrderBy(bson.D{{Key: "_id", Value: objectid}})
		if err != nil {
			c.String(400, err.Error())
			c.Abort()
			return
		}
		err = service.PublishTopic(fmt.Sprintf("order_%s"+stage), order)
		if err != nil {
			c.String(500, "internal error , please try again")
			c.Abort()
			return
		}
		c.JSON(200, res)
		c.Abort()
	}
}
func OrderIsDelivered(c *gin.Context) {
	var id string = c.Param("id")
	var drid string = c.GetHeader("ssid")
	objectid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	query := bson.M{"_id": objectid, "driver_external_id": drid}
	change := bson.D{
		{
			Key: "$set", Value: bson.D{
				{
					Key: "stage", Value: "delivered",
				},
			},
		},
	}
	res, err := db.UpdateOrder(query, change)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	order, err := db.GetOrderBy(bson.D{{Key: "_id", Value: objectid}})
	if err != nil {
		c.String(400, err.Error())
		return
	}
	err = service.PublishTopic("order_delivered", order)
	if err != nil {
		c.String(500, "internal error , please try again")
		return
	}
	c.JSON(200, res)
}
