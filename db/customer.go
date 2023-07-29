package db

import (
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
