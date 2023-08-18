package db

import (
	"fmt"
	"time"

	"github.com/mohamedabdifitah/ecapi/utils"
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
			ProtectFields("password", "device", "metadata.token_version", "metadata.provider"),
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
	cursor, err := CustomerCollection.Find(Ctx, bson.D{}, options.Find().SetProjection(ProtectFields("password", "device", "metadata.token_version", "metadata.provider")))
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
func (c *Customer) Update() (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": c.Id}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "family_name", Value: c.FamilyName},
		{Key: "given_name", Value: c.GivenName},
		{Key: "address", Value: c.Address},
		{Key: "phone", Value: c.Phone},
		{Key: "metadata.updated_at", Value: time.Now().UTC()},
	}}}

	result, err := CustomerCollection.UpdateOne(Ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (c *Customer) ChangePassword(OldPassword string, NewPassword string) *ErrorResponse {
	query := bson.M{"_id": c.Id}
	result := CustomerCollection.FindOne(Ctx, query)
	err := result.Decode(&c)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return &ErrorResponse{Status: 403, Message: fmt.Errorf("user not found"), Type: "string"}
		}
		return &ErrorResponse{Status: 500, Message: err, Type: "string"}

	}

	err = utils.VerifyPassword(OldPassword, c.Password)
	if err != nil {
		return &ErrorResponse{Status: 401, Message: fmt.Errorf("password is invalid"), Type: "string"}
	}
	c.Password = NewPassword
	err = c.BeforeSave()
	if err != nil {
		return &ErrorResponse{Status: 400, Message: err, Type: "string"}
	}
	filter := bson.M{"_id": c.Id}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "password", Value: c.Password}}}}
	_, err = CustomerCollection.UpdateOne(Ctx, filter, update)
	if err != nil {
		return &ErrorResponse{Status: 500, Message: err}
	}

	return nil
}
func (c *Customer) ChangeEmail(OldEmail string, NewEmail string) (*mongo.UpdateResult, *ErrorResponse) {
	filter := bson.M{"_id": c.Id, "email": OldEmail}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "email", Value: NewEmail}}}}
	res, err := CustomerCollection.UpdateOne(Ctx, filter, update)
	if err != nil {
		return nil, &ErrorResponse{Status: 500, Message: err}
	}
	return res, nil
}
func (c *Customer) ChangeMetadataLogin() error {
	query := bson.M{"_id": c.Id}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "metadata.Last_login", Value: time.Now()},
		{Key: "device", Value: c.Device},
	}}}
	_, err := CustomerCollection.UpdateOne(Ctx, query, update)
	return err
}
func CustomerLoginCheck(email string, password string, device Device) (*TokenResponse, *ErrorResponse) {
	query := bson.M{"email": email}
	var customer Customer
	result := CustomerCollection.FindOne(Ctx, query)
	err := result.Decode(&customer)
	if err != nil {
		res := &ErrorResponse{
			Status:  401,
			Type:    "string",
			Message: fmt.Errorf("You have entered an invalid email"),
		}
		return nil, res
	}
	err = utils.VerifyPassword(password, customer.Password)
	if err != nil {
		res := &ErrorResponse{
			Status:  401,
			Type:    "string",
			Message: fmt.Errorf("You have entered an invalid password"),
		}
		return nil, res
	}
	idToken, err := utils.GenerateRefreshToken(customer.Email, customer.Id, customer.Metadata.TokenVersion)
	if err != nil {
		res := &ErrorResponse{
			Status:  500,
			Message: err,
		}
		return nil, res
	}
	token, err := utils.GenerateAccessToken(customer.Email, customer.Id, Roles[0])
	if err != nil {
		res := &ErrorResponse{
			Status:  500,
			Message: err,
		}
		return nil, res
	}
	t := &TokenResponse{
		RefreshToken: idToken,
		AccessToken:  token,
	}
	customer.Device = device
	err = customer.ChangeMetadataLogin()
	if err != nil {
		res := &ErrorResponse{
			Status:  500,
			Message: err,
		}
		return nil, res
	}
	return t, nil
}
func UpdateCustomer(query bson.M, change bson.D) (*mongo.UpdateResult, error) {
	res, err := CustomerCollection.UpdateOne(Ctx, query, change)
	if err != nil {
		return nil, err
	}
	return res, nil
}
