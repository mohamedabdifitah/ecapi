package db

import (
	"fmt"
	"time"

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
func (o Order) GetAll() ([]*Order, error) {
	var orders []*Order
	cursor, err := OrderCollection.Find(Ctx, bson.D{})
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
func (o *Order) GetByMerchant() ([]*Order, error) {
	var orders []*Order
	filter := bson.D{
		{Key: "pickup _external_id", Value: o.PickUpExternalId},
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
func (o *Order) GetByCustomer() ([]*Order, error) {
	var orders []*Order
	filter := bson.D{
		{Key: "dropoff_external_id", Value: o.DropOffExteranlId},
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
func (o *Order) GetByLocation(location []float64, maxdist int64, mindist int64) ([]*Order, error) {
	var orders []*Order
	filter := bson.D{
		{Key: "location", Value: bson.D{
			{
				Key: "$near", Value: bson.D{
					{
						Key: "$maxDistance", Value: maxdist,
					},
					{
						Key: "$minDistance", Value: mindist,
					},
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
	res, err := OrderCollection.InsertOne(Ctx, &o)
	if err != nil {
		return nil, err
	}
	return res, err
}
