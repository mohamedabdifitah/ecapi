package db

import (
	"fmt"
	"time"

	"github.com/mohamedabdifitah/ecapi/service"
	"github.com/mohamedabdifitah/ecapi/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (o *Order) GetById() error {
	query := bson.M{"_id": o.Id}
	result := OrderCollection.FindOne(Ctx, query)
	err := result.Decode(&o)
	return err
}
func GetOrderBy(query bson.D) (*Order, error) {
	var order *Order
	result := OrderCollection.FindOne(Ctx, query)
	err := result.Decode(&order)
	if err != nil {
		return nil, err
	}
	return order, nil
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
func (o *Order) ExtractItems() *ErrorResponse {
	oids := make([]primitive.ObjectID, len(o.Items))
	for i, item := range o.Items {
		ids, err := primitive.ObjectIDFromHex(item.ItemExternalId)
		if err != nil {
			return &ErrorResponse{Status: 400, Message: fmt.Errorf("invalid id of item %s", item.ItemExternalId), Type: "string"}
		}
		oids[i] = ids
	}
	// get all items from the collection of menu
	menus, err := GetListMenus(oids)
	if err != nil {
		return DBErrorHandler(err)
	}
	if menus == nil && len(menus) == 0 {
		return &ErrorResponse{Status: 400, Message: fmt.Errorf("no menus found"), Type: "string"}
	}

	if len(menus) != len(o.Items) {
		return &ErrorResponse{Status: 400, Message: fmt.Errorf("some items must be not available"), Type: "string"}
	}
	for i, menu := range menus {
		o.Items[i].Price = menu.Price * o.Items[i].Quantity
		o.PickupEstimatedTime = o.PickupEstimatedTime + menu.EstimateTime
	}
	return nil
}
func (o *Order) CalculatePrice() {
	for _, order := range o.Items {
		o.OrderValue = o.OrderValue + order.Price
	}
}
func (o *Order) BeforeSave() *ErrorResponse {
	o.Stage = "placed"
	o.Metadata.CreatedAt = time.Now().UTC()
	err := o.ExtractItems()
	if err != nil {
		return err
	}
	o.CalculatePrice()
	erres := o.PickuPExtract()
	if erres != nil {
		return erres
	}
	o.ActionIfUndeliverable = "RETURN_TO_MERCHANT"
	o.DisplayId = utils.GenerateIDs(8)
	return nil
}
func (o *Order) PickuPExtract() *ErrorResponse {
	objectid, err := primitive.ObjectIDFromHex(o.PickUpExternalId)
	if err != nil {
		return &ErrorResponse{Status: 400, Message: fmt.Errorf("invalid id"), Type: "string"}
	}
	merchant := Merchant{
		Id: objectid,
	}
	err = merchant.GetById()
	if err != nil {
		return &ErrorResponse{Status: 400, Message: fmt.Errorf("merchant not found"), Type: "string"}
	}
	if merchant.Closed {
		return &ErrorResponse{Status: 400, Message: fmt.Errorf("merchant is currently closed"), Type: "string"}
	}
	weekday := time.Now().Weekday()
	today := merchant.ActiveDays[int(weekday)]
	ok := utils.IsTimeBetween(time.Now(), today.TimeOperatorStart, today.TimeOperatorEnd)
	if !ok {
		return &ErrorResponse{Status: 400, Message: fmt.Errorf("merchant is closed"), Type: "string"}
	}
	o.PickUpLocation = merchant.Location
	o.PickUpName = merchant.BusinessName
	o.PickupAddress = merchant.Address
	o.PickUpPhone = merchant.BusinessPhone
	return nil
}

// places order
func (o *Order) PlaceOrder() (*mongo.InsertOneResult, *ErrorResponse) {
	errres := o.BeforeSave()
	if errres != nil {
		return nil, errres
	}
	res, err := OrderCollection.InsertOne(Ctx, &o)
	if err != nil {
		return nil, DBErrorHandler(err)
	}
	// answer from https://www.mongodb.com/community/forums/t/insertone-returns-interface/131263/3
	o.Id = res.InsertedID.(primitive.ObjectID)
	order := o.GetById()

	if err != nil {
		return nil, &ErrorResponse{Status: 400, Message: err, Type: "string"}
	}
	go service.ProduceMessage("", "new_order", "", order)
	return res, nil
}
func UpdateOrder(query bson.M, change bson.D) (*mongo.UpdateResult, error) {
	res, err := OrderCollection.UpdateOne(Ctx, query, change, options.Update().SetUpsert(true))
	if err != nil {
		return nil, err
	}
	return res, nil
}
func AssignOrderToDriver(orderId primitive.ObjectID, driverId primitive.ObjectID) (*mongo.UpdateResult, *ErrorResponse) {
	order := &Order{
		Id: orderId,
	}
	err := order.GetById()
	if err != nil {
		return nil, &ErrorResponse{Message: fmt.Errorf("order is not found"), Status: 400, Type: "string"}
	}
	if order.Stage == "canceled" {
		return nil, &ErrorResponse{Message: fmt.Errorf("cannot accpet , order is already canceled"), Status: 403, Type: "string"}
	}
	if order.DriverExternalId == driverId.Hex() {
		return nil, &ErrorResponse{Message: fmt.Errorf("you already accepted the order"), Status: 409, Type: "string"}
	}
	if order.DriverExternalId != "" && order.DriverPhone != "" {
		return nil, &ErrorResponse{Message: fmt.Errorf("order already assigned to another driver"), Status: 409, Type: "string"}
	}
	driver := Driver{
		Id: driverId,
	}
	err = driver.GetById()
	if err != nil {
		return nil, &ErrorResponse{Message: fmt.Errorf("driver is not found"), Status: 400, Type: "string"}
	}
	query := bson.M{"_id": order.Id}
	change := bson.D{{Key: "$set", Value: bson.D{
		{Key: "metadata.update_at", Value: time.Now()},
		{Key: "driver_external_id", Value: driver.Id.Hex()},
		{Key: "driver_phone", Value: driver.Phone},
	}}}
	res, err := UpdateOrder(query, change)
	if err != nil {
		return nil, &ErrorResponse{Message: err, Status: 500}
	}
	order, err = GetOrderBy(bson.D{{Key: "_id", Value: res.UpsertedID}})
	if err != nil {
		return nil, &ErrorResponse{Message: fmt.Errorf("order is not found"), Status: 400, Type: "string"}
	}
	go service.ProduceMessage("", "driver_accepted_order", "", order)
	return res, nil
}

// drop the order like driver decline order after he accept
func DropOrder(orderId primitive.ObjectID, driverId string) (*mongo.UpdateResult, *ErrorResponse) {
	order := &Order{
		Id: orderId,
	}
	err := order.GetById()
	if err != nil {
		return nil, &ErrorResponse{Message: fmt.Errorf("order is not found"), Status: 400, Type: "string"}
	}
	if order.DriverExternalId != driverId {
		return nil, &ErrorResponse{Message: fmt.Errorf("you cannot drop order , that you're driver"), Status: 403, Type: "string"}
	}
	if order.Stage == "deleivered" {
		return nil, &ErrorResponse{Message: fmt.Errorf("order is already delivered"), Status: 403, Type: "string"}
	}
	if order.Stage == "cancelled" {
		return nil, &ErrorResponse{Message: fmt.Errorf("order is already cancelled"), Status: 403, Type: "string"}
	}
	if order.Stage == "pickuped" {
		return nil, &ErrorResponse{Message: fmt.Errorf("you cannot drop order if you already have a pickuped"), Status: 403, Type: "string"}
	}
	query := bson.M{"_id": order.Id}
	change := bson.D{{Key: "$set", Value: bson.D{
		{Key: "metadata.update_at", Value: time.Now()},
		{Key: "driver_external_id", Value: ""},
		{Key: "driver_phone", Value: ""},
	}}}
	res, err := UpdateOrder(query, change)
	if err != nil {
		return nil, &ErrorResponse{Message: err, Status: 500}
	}
	order, err = GetOrderBy(bson.D{{Key: "_id", Value: res.UpsertedID}})
	if err != nil {
		return nil, &ErrorResponse{Message: fmt.Errorf("order is not found"), Status: 400, Type: "string"}
	}
	go service.ProduceMessage("", "order_dropped", "", order)

	return res, nil
}
