package controller

import (
	"github.com/mohamedabdifitah/ecapi/db"
)

type SignUpWithEmailBody struct {
	Email    string `json:"email" `
	Password string `json:"password"`
}
type CustomerBody struct {
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Address    string `json:"address"`
	Phone      string `json:"phone"`
}
type ChangePasswordBody struct {
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
}
type ChangeEmaildBody struct {
	OldEmail string `json:"old_email"`
	NewEmail string `json:"new_email"`
}

type EmailLoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	DeviceId string `json:"device_id"`
	Kind     string `json:"kind"`
}
type SignUpMerchantBody struct {
	BusinessEmail string    `json:"business_email" binding:"required"`
	BusinessName  string    `json:"business_name" binding:"required"`
	Password      string    `json:"password" binding:"required"`
	Location      []float64 `json:"location" binding:"required"`
}
type MerchantBody struct {
	Location           []float64       `json:"location"`
	Address            string          `json:"address"`
	BusinessEmail      string          `json:"business_email"`
	TimeOperationStart float32         `json:"time_operation_start"`
	TimeOperationEnd   float32         `json:"time_operation_end"`
	BusinessName       string          `json:"business_name"`
	ActiveDays         []db.ActiveDays `json:"active_days"`
}
type ChangePhonedBody struct {
	NewPhone string `json:"new_phone"`
	OldPhone string `json:"old_phone"`
}
type PhoneLoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type DriverBody struct {
	Age       string `json:"age"`
	GivenName string `json:"given_name"`
	Address   string `json:"address"`
}
type SignUpWithDriverBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type MenuBody struct {
	Title        string        `json:"title" bson:"title"`             // Burger
	Description  string        `json:"description" bson:"description"` // Chicken fries contains
	Status       string        `json:"status" bson:"status"`           // available , unavailable , banned
	Category     string        `json:"category" bson:"category"`       // fast food , drink ,                     // Images of the product urls
	Price        uint          `json:"price" bson:"price"`             // the price of the product is represented as cents 99 = $0.99
	Attributes   db.Attributes `json:"attributes" bson:"attributes"`
	Reciepe      []string      `json:"reciepe" bson:"reciepe"`
	EstimateTime int           `json:"estimate_time"`
	Images       []string      `json:"images"`
}
type ReviewBody struct {
	OrderId    string   `json:"order_id"`
	Rate       uint     `json:"rate"`
	Message    string   `json:"message"`
	ExternalId string   `json:"external_id"`
	From       string   `json:"from"`
	Options    []string `json:"options"`
}
type ReviewExternalBody struct {
}
type PlaceOrderBody struct {
	Items              []Item    `json:"items"`
	DropOffPhone       string    `json:"dropoff_phone" `
	DropOffExteranlId  string    `json:"dropoff_external_id"`
	DropOffContactName string    `json:"dropoff_contact_name"`
	DropOffAddress     string    `json:"dropoff_address"`  // address 901 Market Street 6th Floor San Francisco, CA 94103
	DroOffLocation     []float64 `json:"dropoff_location"` // location cordinates. float([123.1312343,-37.2144343])
	DropOffInstruction string    `json:"dropoff_instructions"`
	PickUpExternalId   string    `json:"pickup_external_id"`
	Type               string    `json:"type"`
}
type Item struct {
	ItemExternalId string `json:"item_external_id"`
	Quantity       uint   `json:"quantity"`
}
type AccOrderMerchantBody struct {
	OrderId    string `json:"order_id"`
	MerchantId string `json:"merchant_id"`
}
type FilterMerchantsBody struct {
	Location [2]float64 `json:"location"`
	Near     int        `json:"near"`
	Popular  bool       `json:"popular"`
	Rate     float32    `json:"rate"`
}
type Filter struct {
	On bool `json:"on"`
}
