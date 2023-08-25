package db

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mohamedabdifitah/ecapi/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo"
)

func (o *Order) GetById() error {
	query := bson.M{"_id": o.Id}
	result := OrderCollection.FindOne(Ctx, query)
	err := result.Decode(&o)
	return err
}
func GetOrdersBy(filter bson.D) ([]*Order, error) {
	var orders []*Order
	cursor, err := OrderCollection.Find(Ctx, filter)
	if err != nil {
		return nil, err
	}
	for cursor.Next(Ctx) {
		var order *Order
		err := cursor.Decode(&order)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	cursor.Close(Ctx)
	return orders, nil
}
func (o *Order) GetByLocation(location []float64, maxdist int64, mindist int64) ([]*Order, error) {
	var orders []*Order
	filter := bson.D{
		{Key: "pickup_location", Value: bson.D{
			{
				Key: "$near", Value: bson.D{
					{
						Key: "$geometry", Value: bson.D{
							{
								Key: "type", Value: "Point",
							},
							{
								Key: "coordinates", Value: location,
							},
						},
					},
					{
						Key: "$maxDistance", Value: maxdist,
					},
					{
						Key: "$minDistance", Value: mindist,
					},
				},
			},
		}},
	}
	cursor, err := OrderCollection.Find(Ctx, filter)
	if err != nil {
		return nil, err
	}
	for cursor.Next(Ctx) {
		var order *Order
		err := cursor.Decode(&order)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	cursor.Close(Ctx)
	return orders, nil
}
func (o *Order) ExtractItems() error {
	oids := make([]primitive.ObjectID, len(o.Items))
	for i := range o.Items {
		ids, err := primitive.ObjectIDFromHex(o.Items[i].ItemExternalId)
		if err != nil {
			return err
		}
		oids[i] = ids
	}
	menus, err := GetListMenus(oids)
	if err != nil {
		return err
	}
	if menus == nil {
		return fmt.Errorf("no menus found")
	}
	if len(menus) == 0 {
		return fmt.Errorf("no menus found")
	}
	if len(menus) != len(o.Items) {
		return fmt.Errorf("some items must be not available")
	}
	for i, menu := range menus {
		o.Items[i].Price = menu.Price * o.Items[i].Quantity
		o.PickupTimeEstimated = o.PickupTimeEstimated + menu.EstimateTime

	}
	return nil
}
func (o *Order) CalculatePrice() {
	for _, order := range o.Items {
		o.OrderValue = o.OrderValue + order.Price
	}
}
func (o *Order) BeforeSave() error {
	o.Stage = "placed"
	o.Metadata.CreatedAt = time.Now().UTC()
	err := o.ExtractItems()
	if err != nil {
		// fmt.Println(err.Error())
		return err
	}
	o.CalculatePrice()
	err = o.PickuPExtract()
	if err != nil {
		return err
	}
	return nil
}
func (o *Order) PickuPExtract() error {
	objectid, err := primitive.ObjectIDFromHex(o.PickUpExternalId)
	merchant := Merchant{
		Id: objectid,
	}
	err = merchant.GetById()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	o.PickUpLocation = merchant.Location
	o.PickupAddress = merchant.Address
	o.PickUpPhone = merchant.BusinessPhone
	return nil
}
func (o *Order) PlaceOrder() (*mongo.InsertOneResult, error) {
	err := o.BeforeSave()
	if err != nil {
		return nil, err
	}
	// res, err := OrderCollection.InsertOne(Ctx, &o)
	// if err != nil {
	// return nil, err
	// }
	// answer from https://www.mongodb.com/community/forums/t/insertone-returns-interface/131263/3
	// o.Id = res.InsertedID.(primitive.ObjectID)
	json, err := json.Marshal(o)

	if err != nil {
		return nil, err
	}
	service.PublishTopic("new-order", json)
	return nil, nil
}
func UpdateOrder(query bson.M, change bson.D) (*mongo.UpdateResult, error) {
	res, err := OrderCollection.UpdateOne(Ctx, query, change)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func AccpetOrderBy(query bson.M, change bson.D, topic string) (*mongo.UpdateResult, error) {
	res, err := UpdateOrder(query, change)
	if err != nil {
		return nil, err
	}
	service.PublishTopic(topic, topic)
	return res, nil
}
func AssignOrderToDriver(orderId primitive.ObjectID, driverId primitive.ObjectID) (*mongo.UpdateResult, *ErrorResponse) {
	order := Order{
		Id: orderId,
	}
	err := order.GetById()
	if err != nil {
		return nil, &ErrorResponse{Message: err, Status: 400}
	}
	if order.Stage == "canceled" {
		return nil, &ErrorResponse{Message: fmt.Errorf("order is canceled"), Status: 403, Type: "string"}
	}
	if order.DriverExternalId != "" && order.DriverPhone != "" {
		return nil, &ErrorResponse{Message: fmt.Errorf("order already assigned to another driver"), Status: 409, Type: "string"}
	}
	driver := Driver{
		Id: driverId,
	}
	err = driver.GetById()
	if err != nil {
		return nil, &ErrorResponse{Message: err, Status: 400}
	}
	query := bson.M{"_id": order.Id}
	change := bson.D{{Key: "$set", Value: bson.D{
		{Key: "metadata.update_at", Value: time.Now()},
		{Key: "driver_external_id", Value: driver.Id.Hex()},
		{Key: "driver_phone", Value: driver.Phone},
	}}}
	res, err := AccpetOrderBy(query, change, "driver-accepted")
	if err != nil {
		return nil, &ErrorResponse{Message: err, Status: 500}
	}
	return res, nil
}
