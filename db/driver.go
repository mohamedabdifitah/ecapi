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

func (d *Driver) BeforeSave() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(d.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	d.Password = string(hashedPassword)
	d.Metadata.CreatedAt = time.Now().UTC()
	d.Metadata.UpdatedAt = time.Now().UTC()
	return nil
}
func (d *Driver) Save() (*mongo.InsertOneResult, error) {
	err := d.BeforeSave()
	if err != nil {
		return nil, err
	}
	res, err := DriverCollection.InsertOne(Ctx, &d)
	if err != nil {
		return nil, err
	}
	// res.Decode(&d)
	return res, nil
}
func (d *Driver) GetById() error {
	query := bson.M{"_id": d.Id}
	result := DriverCollection.FindOne(
		Ctx, query, options.FindOne().SetProjection(
			ProtectFields(CommonProtoctedFields),
		))
	err := result.Decode(&d)
	return err
}
func (d *Driver) Delete() (*mongo.DeleteResult, error) {
	filter := bson.M{"_id": d.Id}
	result, err := DriverCollection.DeleteOne(Ctx, filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func GetDrivers(query bson.M, projection []string) ([]map[string]interface{}, error) {
	var drivers []map[string]interface{}
	protectedfileds := []string{"password", "metadata.token_version", "metadata.provider"}
	protectedfileds = append(protectedfileds, projection...)
	cursor, err := DriverCollection.Find(Ctx, query, options.Find().SetProjection(ProtectFields(protectedfileds)))
	if err != nil {
		return nil, err
	}
	for cursor.Next(Ctx) {
		var driver map[string]interface{}
		err := cursor.Decode(&driver)
		if err != nil {

			return nil, err

		}
		drivers = append(drivers, driver)
	}
	cursor.Close(Ctx)
	return drivers, nil
}
func (d *Driver) Update() (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": d.Id}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "age", Value: d.Age},
		{Key: "vehicle", Value: d.Vehicle},
		{Key: "given_name", Value: d.GivenName},
		{Key: "address", Value: d.Address},
		{Key: "metadata.updated_at", Value: time.Now().UTC()},
	}}}

	result, err := DriverCollection.UpdateOne(Ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (d *Driver) ChangePassword(OldPassword string, NewPassword string) *ErrorResponse {
	query := bson.M{"_id": d.Id}
	result := DriverCollection.FindOne(Ctx, query)
	err := result.Decode(&d)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return &ErrorResponse{Status: 403, Message: fmt.Errorf("user not found"), Type: "string"}

		}
		return &ErrorResponse{Status: 500, Message: err, Type: "string"}
	}

	err = utils.VerifyPassword(OldPassword, d.Password)
	if err != nil {
		return &ErrorResponse{Status: 401, Message: fmt.Errorf("password is invalid"), Type: "string"}
	}
	d.Password = NewPassword
	err = d.BeforeSave()
	if err != nil {
		return &ErrorResponse{Status: 400, Message: err, Type: "string"}
	}
	filter := bson.M{"_id": d.Id}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "password", Value: d.Password}}}}
	_, err = DriverCollection.UpdateOne(Ctx, filter, update)
	if err != nil {
		return &ErrorResponse{Status: 500, Message: err}
	}
	return nil
}
func (d *Driver) ChangeEmail(OldEmail string, NewEmail string) (*mongo.UpdateResult, *ErrorResponse) {
	filter := bson.M{"_id": d.Id, "email": OldEmail}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "email", Value: NewEmail}}}}
	res, err := DriverCollection.UpdateOne(Ctx, filter, update)
	if err != nil {
		return nil, &ErrorResponse{Status: 500, Message: err}
	}
	return res, nil
}
func (d *Driver) ChangePhone(OldPhone string, NewPhone string) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": d.Id, "phone": OldPhone}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "phone", Value: NewPhone},
	}},
	}
	result, err := DriverCollection.UpdateOne(Ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (d *Driver) ChangeMetadataLogin() error {
	query := bson.M{"_id": d.Id}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "metadata.Last_login", Value: time.Now()},
		{Key: "device", Value: d.Device},
	}}}
	_, err := DriverCollection.UpdateOne(Ctx, query, update)
	return err
}
func DriverLoginCheck(phone string, password string, device Device) (*TokenResponse, *ErrorResponse) {
	query := bson.M{"phone": phone}
	var driver Driver
	result := DriverCollection.FindOne(Ctx, query)
	err := result.Decode(&driver)
	if err != nil {
		res := &ErrorResponse{
			Status:  401,
			Type:    "string",
			Message: fmt.Errorf("You have entered an invalid phone"),
		}
		return nil, res
	}
	err = utils.VerifyPassword(password, driver.Password)
	if err != nil {
		res := &ErrorResponse{
			Status:  401,
			Type:    "string",
			Message: fmt.Errorf("You have entered an invalid password"),
		}
		return nil, res
	}
	idToken, err := utils.GenerateRefreshToken(driver.Phone, driver.Id, driver.Metadata.TokenVersion)
	if err != nil {
		res := &ErrorResponse{
			Status:  500,
			Message: err,
		}
		return nil, res
	}
	token, err := utils.GenerateAccessToken(driver.Phone, driver.Id, Roles[1])
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
	driver.Device = device
	err = driver.ChangeMetadataLogin()
	if err != nil {
		res := &ErrorResponse{
			Status:  500,
			Message: err,
		}
		return nil, res
	}
	return t, nil
}
func UpdateDriver(query bson.M, change bson.D) (*mongo.UpdateResult, error) {
	res, err := DriverCollection.UpdateOne(Ctx, query, change)
	if err != nil {
		return nil, err
	}
	return res, err
}
