package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func (c *Customer) BeforeSave() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	c.Password = string(hashedPassword)
	c.Metadata.CreatedAt = time.Now().UTC()
	c.Metadata.UpdatedAt = time.Now().UTC()
	return nil
}
func (c *Customer) Save() (*mongo.InsertOneResult, error) {
	err := c.BeforeSave()
	if err != nil {
		return nil, err
	}
	res, err := CustomerCollection.InsertOne(Ctx, &c)
	if err != nil {
		return nil, err
	}
	// res.Decode(&c)
	return res, nil
}
func (c *Customer) GetById() error {
	query := bson.M{"_id": c.Id}
	result := CustomerCollection.FindOne(
		Ctx, query, options.FindOne().SetProjection(
			ProtectFields("password", "devices"),
		))
	err := result.Decode(&c)
	return err
}
func (c *Customer) Delete() (*mongo.DeleteResult, error) {
	filter := bson.M{"_id": c.Id}
	result, err := CustomerCollection.DeleteOne(Ctx, filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (c *Customer) GetAll() ([]*Customer, error) {
	var customers []*Customer
	cursor, err := CustomerCollection.Find(Ctx, bson.D{}, options.Find().SetProjection(ProtectFields("password", "devices")))
	if err != nil {
		return nil, err
	}
	for cursor.Next(Ctx) {
		var customer *Customer
		err := cursor.Decode(&customer)
		if err != nil {

			return nil, err

		}
		customers = append(customers, customer)
	}
	cursor.Close(Ctx)
	return customers, nil
}
