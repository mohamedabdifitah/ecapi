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
type SignUpMerchantWithPhoneBody struct {
	BusinessPhone string `json:"business_phone"`
	Password      string `json:"password"`
}
type MerchantBody struct {
	Location           []float64 `json:"location"`
	Address            string    `json:"address"`
	BusinessEmail      string    `json:"business_email"`
	TimeOperationStart int       `json:"time_operation_start"`
	TimeOperationEnd   int       `json:"time_operation_end"`
	BusinessName       string    `json:"business_name"`
}
type ChangePhonedBody struct {
	NewPhone string `json:"new_phone"`
	OldPhone string `json:"old_phone"`
}
type PhoneLoginBody struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}
type DriverBody struct {
	Age         string `json:"age"`
	GivenName   string `json:"given_name"`
	Address     string `json:"address"`
	VehicleType string `json:"vehicle_type"`
}
type SignUpWithDriverBody struct {
	Phone    string `json:"phone"`
	Emai     string `json:"email"`
	Password string `json:"password"`
}
type MenuBody struct {
	Title              string        `json:"title" bson:"title"`             // Burger
	Description        string        `json:"description" bson:"description"` // Chicken fries contains
	Status             string        `json:"status" bson:"status"`           // available , unavailable , banned
	Category           string        `json:"category" bson:"category"`       // fast food , drink ,                     // Images of the product urls
	Price              uint          `json:"price" bson:"price"`             // the price of the product is represented as cents 99 = $0.99
	Attributes         db.Attributes `json:"attributes" bson:"attributes"`
	Discount           uint          `json:"discount" bson:"discount"` // 10%
	MerchantExternalId string        `json:"merchant_external_id" bson:"merchant_external_id"`
	Reciepe            []string      `json:"reciepe" bson:"reciepe"`
}
type CreateMenuBody struct {
	Title              string        `json:"title" bson:"title" binding:"max=30,min=1"` // Burger
	Description        string        `json:"description" bson:"description"`            // Chicken fries contains
	Status             string        `json:"status" bson:"status"`                      // available , unavailable , banned
	Category           string        `json:"category" bson:"category"`                  // fast food , drink ,
	Images             []string      `json:"image" bson:"image"`                        // Images of the product urls
	Price              uint          `json:"price" bson:"price"`                        // the price of the product is represented as cents 99 = $0.99
	Attributes         db.Attributes `json:"attributes" bson:"attributes"`
	Discount           uint          `json:"discount" bson:"discount"` // 10%
	MerchantExternalId string        `json:"merchant_external_id" bson:"merchant_external_id"`
	Reciepe            []string      `json:"reciepe" bson:"reciepe"`
	Barcode            string        `json:"-" bson:"barcode"` // if this needed
}
type ReviewBody struct {
	Rate       uint   `json:"rate"`
	Message    string `json:"message"`
	From       string `json:"from"`
	Type       string `json:"type"`
	ExternalId string `json:"external_id"`
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
}
type Item struct {
	ItemExternalId string `json:"item_external_id"`
	Quantity       uint   `json:"quantity"`
}
