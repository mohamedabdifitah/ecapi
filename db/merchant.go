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

func (m *Merchant) BeforeSave() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	m.Password = string(hashedPassword)
	m.Metadata.CreatedAt = time.Now().UTC()
	m.Metadata.UpdatedAt = time.Now().UTC()
	return nil
}
func (m *Merchant) Save() (*mongo.InsertOneResult, *ErrorResponse) {
	err := m.BeforeSave()
	if err != nil {
		return nil, &ErrorResponse{Status: 400, Message: err, Type: "string"}
	}
	res, err := MerchantCollection.InsertOne(Ctx, &m)
	if err != nil {
		return nil, DBErrorHandler(err)
	}
	return res, nil
}
func (m *Merchant) GetById() error {
	query := bson.M{"_id": m.Id}
	result := MerchantCollection.FindOne(
		Ctx, query, options.FindOne().SetProjection(
			ProtectFields(CommonProtoctedFields),
		))
	err := result.Decode(&m)
	return err
}
func (m *Merchant) Delete() (*mongo.DeleteResult, error) {
	filter := bson.M{"_id": m.Id}
	result, err := MerchantCollection.DeleteOne(Ctx, filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (m *Merchant) GetAll() ([]*Merchant, error) {
	var merchants []*Merchant
	cursor, err := MerchantCollection.Find(Ctx, bson.D{}, options.Find().SetProjection(ProtectFields(CommonProtoctedFields)))
	if err != nil {
		return nil, err
	}
	for cursor.Next(Ctx) {
		var merchant *Merchant
		err := cursor.Decode(&merchant)
		if err != nil {

			return nil, err

		}
		merchants = append(merchants, merchant)
	}
	cursor.Close(Ctx)
	return merchants, nil
}
func (m *Merchant) Update() (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": m.Id}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "location", Value: m.Location},
		{Key: "business_name", Value: m.BusinessName},
		{Key: "time_operation_start", Value: m.TimeOperatorStart},
		{Key: "time_operation_end", Value: m.TimeOperatorEnd},
		{Key: "business_email", Value: m.BusinessEmail},
		{Key: "address", Value: m.Address},
		{Key: "category", Value: m.Category},
		{Key: "metadata.updated_at", Value: time.Now().UTC()},
	}}}

	result, err := MerchantCollection.UpdateOne(Ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (m *Merchant) ChangeBusinessPhone(OldPhone string, NewPhone string) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": m.Id, "business_phone": OldPhone}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "business_phone", Value: NewPhone},
	}},
	}
	result, err := MerchantCollection.UpdateOne(Ctx, filter, update)
	if err != nil {

		return nil, err
	}
	return result, nil
}
func (m *Merchant) ChangePassword(OldPassword string, NewPassword string) *ErrorResponse {
	query := bson.M{"_id": m.Id}
	result := MerchantCollection.FindOne(Ctx, query)
	err := result.Decode(&m)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return &ErrorResponse{Status: 403, Message: fmt.Errorf("user not found"), Type: "string"}
		}
		return &ErrorResponse{Status: 500, Message: err, Type: "string"}
	}
	err = utils.VerifyPassword(OldPassword, m.Password)
	if err != nil {
		return &ErrorResponse{Status: 401, Message: fmt.Errorf("password is invalid"), Type: "string"}
	}
	m.Password = NewPassword
	err = m.BeforeSave()
	if err != nil {
		return &ErrorResponse{Status: 400, Message: err, Type: "string"}
	}
	filter := bson.M{"_id": m.Id}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "password", Value: m.Password}}}}
	_, err = MerchantCollection.UpdateOne(Ctx, filter, update)
	if err != nil {
		return &ErrorResponse{Status: 500, Message: err}
	}
	return nil
}
func (m *Merchant) ChangeMetadataLogin() error {
	query := bson.M{"_id": m.Id}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "metadata.Last_login", Value: time.Now()},
		{Key: "device", Value: m.Device},
	}}}
	_, err := MerchantCollection.UpdateOne(Ctx, query, update)
	return err
}
func MerchantLoginCheck(phone string, password string, device Device) (*TokenResponse, *ErrorResponse) {
	query := bson.M{"business_phone": phone}
	var merchant Merchant
	result := MerchantCollection.FindOne(Ctx, query)
	err := result.Decode(&merchant)
	if err != nil {
		res := &ErrorResponse{
			Status:  401,
			Type:    "string",
			Message: fmt.Errorf("You have entered an invalid business phone number"),
		}
		return nil, res
	}
	err = utils.VerifyPassword(password, merchant.Password)
	if err != nil {
		res := &ErrorResponse{
			Status:  401,
			Type:    "string",
			Message: fmt.Errorf("You have entered an invalid password"),
		}
		return nil, res
	}
	idToken, err := utils.GenerateRefreshToken(merchant.BusinessPhone, merchant.Id, merchant.Metadata.TokenVersion)
	if err != nil {
		res := &ErrorResponse{
			Status:  500,
			Message: err,
		}
		return nil, res
	}
	token, err := utils.GenerateAccessToken(merchant.BusinessPhone, merchant.Id, Roles[2])
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
	merchant.Device = device
	err = merchant.ChangeMetadataLogin()
	if err != nil {
		res := &ErrorResponse{
			Status:  500,
			Message: err,
		}
		return nil, res
	}
	return t, nil
}
func (m *Merchant) GetMerchantByLocation(location []float64, maxdist int64, mindist int64) ([]*Merchant, error) {
	var merchants []*Merchant
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
	cursor, err := MerchantCollection.Find(Ctx, filter, options.Find().SetProjection(ProtectFields(CommonProtoctedFields)))
	if err != nil {
		return nil, err
	}
	for cursor.Next(Ctx) {
		var merchant *Merchant
		err := cursor.Decode(&merchant)
		if err != nil {

			return nil, err

		}
		merchants = append(merchants, merchant)
	}
	cursor.Close(Ctx)
	return merchants, nil
}
func UpdateMerchant(query bson.M, change bson.D) (*mongo.UpdateResult, error) {
	res, err := MerchantCollection.UpdateOne(Ctx, query, change)
	if err != nil {
		return nil, err
	}
	return res, err
}
