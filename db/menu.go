package db

import (
	"fmt"
	"time"

	"github.com/mohamedabdifitah/ecapi/service"
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

func (m *Menu) Create() (*mongo.InsertOneResult, error) {
	m.Metadata.CreatedAt = time.Now().UTC()
	res, err := MenuCollection.InsertOne(Ctx, m)
	if err != nil {
		IsDup(err)
		return nil, err
	}
	m.Id = res.InsertedID.(primitive.ObjectID)
	go service.CreateDocument("menu", m)

	return &mongo.InsertOneResult{}, nil
}
func PrintN(name string) {
	fmt.Println(name)
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
		{Key: "attributes", Value: m.Attributes},
		{Key: "metadata.updated_at", Value: time.Now().UTC()},
	}}}
	result, err := MenuCollection.UpdateOne(Ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func ChangeMenuField(filter bson.D, update bson.D) (*mongo.UpdateResult, error) {
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
func (m Menu) GetFromMerchant() ([]*Menu, error) {
	var menus []*Menu
	query := bson.M{"merchant_external_id": m.MerchantExternalId}
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
func (m Menu) SetImages() (*mongo.UpdateResult, error) {
	query := bson.M{"_id": m.Id}
	change := bson.M{"$push": bson.M{"images": bson.M{"$each": m.Images}}}
	result, err := MenuCollection.UpdateOne(Ctx, query, change)
	if err != nil {
		return nil, err
	}
	return result, nil
}
