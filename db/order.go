package db

import (
	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
)

func CreateOrder() {

}
func CancelOrder() {
}
func (o *Order) GetById() error {
	query := bson.M{"_id": o.Id}
	result := OrderCollection.FindOne(Ctx, query)
	err := result.Decode(&o)
	return err
}
func (o *Order) GetByMerchant() ([]*Order, error) {
	var orders []*Order
	filter := bson.D{
		{Key: "pickup _external_id", Value: o.PickUpExternalId},
	}
	cursor, err := MenuCollection.Find(Ctx, filter)
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
func (o *Order) GetByCustomer() {
}
func (o *Order) GetByLocation() {
}
