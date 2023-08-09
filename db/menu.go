package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *Menu) GetById() error {
	query := bson.M{"_id": m.Id}
	result := MenuCollection.FindOne(Ctx, query)
	err := result.Decode(&m)
	return err
}

// func (m *Menu) GetMerchantMenu() ([]*Menu , error) {
// }
func (m *Menu) Create() (*mongo.InsertOneResult, error) {
	m.Metadata.CreatedAt = time.Now().UTC()
	res, err := MenuCollection.InsertOne(Ctx, m)
	if err != nil {
		IsDup(err)
		return nil, err
	}
	// res.Decode(&c)
	return res, nil
}
func (m *Menu) Update() (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": m.Id}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "title", Value: m.Title},
		{Key: "description", Value: m.Description},
		{Key: "category", Value: m.Category},
		{Key: "price", Value: m.Price},
		{Key: "status", Value: m.Status},
		{Key: "reciepe", Value: m.Reciepe},
		{Key: "discount", Value: m.Discount},
		{Key: "attributes", Value: m.Attributes},
		{Key: "merchant_external_id", Value: m.MerchantExternalId},
		{Key: "metadata.updated_at", Value: time.Now().UTC()},
	}}}
	result, err := MenuCollection.UpdateOne(Ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (m *Menu) Delete() (*mongo.DeleteResult, error) {
	filter := bson.M{"_id": m.Id}
	result, err := MenuCollection.DeleteOne(Ctx, filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (m Menu) GetAll() ([]*Menu, error) {
	var menus []*Menu
	_ = bson.D{
		{Key: "description", Value: m.Description},
		{Key: "category", Value: m.Category},
		{Key: "price", Value: m.Price},
		{Key: "status", Value: m.Status},
		{Key: "reciepe", Value: m.Reciepe},
		{Key: "discount", Value: m.Discount},
		// {Key: "attributes", Value: m.Attributes},
	}
	cursor, err := MenuCollection.Find(Ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	for cursor.Next(Ctx) {
		var menu *Menu
		err := cursor.Decode(&menu)
		if err != nil {

			return nil, err

		}
		menus = append(menus, menu)
	}
	cursor.Close(Ctx)
	return menus, nil
}
func GetListMenus(oids []primitive.ObjectID) ([]*Menu, error) {
	var menus []*Menu
	query := bson.M{"_id": bson.M{"$in": oids}}
	cursor, err := MenuCollection.Find(Ctx, query)
	if err != nil {
		return nil, err
	}
	for cursor.Next(Ctx) {
		var menu *Menu
		err := cursor.Decode(&menu)
		if err != nil {

			return nil, err

		}
		menus = append(menus, menu)
	}
	cursor.Close(Ctx)
	return menus, nil
}
