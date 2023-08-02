package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// all this values are represented as cm or centimeter as unit type
type Attributes struct {
	Volume primitive.Decimal128 `bson:"volume" json:"volume"`
	Height primitive.Decimal128 `bson:"height" json:"height"`
	Width  primitive.Decimal128 `bson:"width" json:"width"`
	Length primitive.Decimal128 `bson:"length" json:"length"`
}

type Menu struct {
	Id                 primitive.ObjectID ` bson:"_id,omitempty" json:"id,omitempty"`       // 63f642ac061b6f5f089b3a65
	Title              string             `json:"title" bson:"title" binding:"max=4,min=1"` // Burger
	Description        string             `json:"description" bson:"description"`           // Chicken fries contains
	Status             string             `json:"status" bson:"status"`                     // available , unavailable , banned
	Category           string             `json:"category" bson:"category"`                 // fast food , drink ,
	Images             []string           `json:"image" bson:"image"`                       // Images of the product urls
	Price              uint               `json:"price" bson:"price"`                       // the price of the product is represented as cents 99 = $0.99
	Attributes         Attributes         `json:"attributes" bson:"attributes"`
	Metadata           Metadata           `json:"metadata" bson:"metadata"`
	Discount           uint               `json:"-" bson:"discount"` // 10%
	MerchantExternalId string             `json:"merchant_external_id" bson:"merchant_external_id"`
	Reciepe            []string           `json:"reciepe" bson:"reciepe"`
	Barcode            string             `json:"-" bson:"barcode"` // if this needed
}
type Metadata struct {
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
type Item struct {
	ItemExternalId string `json:"item_external_id" bson:"item_external_id"`
	Quantity       uint   `json:"quantity" bson:"quantity"`
}
type Order struct {
	Id                    uint      ` bson:"_id,omitempty" json:"id,omitempty"`
	OrderValue            uint      `json:"order_value" bson:"order_value" ` // Order value is represented as cents 199 = $19.9
	Items                 []Item    `json:"items" bson:"items"`
	DropOffExternalId     string    `bson:"dropoff_external_id" json:"dropoff_external_id"`
	DropOffPhone          string    `json:"dropoff_phone" bson:"dropoff_phone"`
	DropOffAddress        string    `json:"dropoff_address" bson:"dropoff_address"`   // address 901 Market Street 6th Floor San Francisco, CA 94103
	DroOffLocation        []float64 `json:"dropoff_location" bson:"dropoff_location"` // location cordinates. float([123.1312343,-37.2144343])
	DropOffInstruction    string    `json:"dropoff_instructions" bson:"dropoff_instructions"`
	Stage                 string    `json:"stage" bson:"stage"`                                      // prepare,pickup,deleivered.
	ActionIfUndeliverable string    `json:"action_if_undeliverable" bson:"action_if_undeliverable"`  // return_to_pickup
	DriverAllowedVehicles []string  `json:"driver_allowed_vehicles" bson:"driver_allowed_vehicles" ` // car , motorcycle , walking
	PickupAddress         string    `json:"pickup_address" bson:"pickup_address"`
	PickupReferenceTag    string    `json:"pickup_reference_tag" bson:"pickup_reference_tag"`
}

type Customer struct {
	Id         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email      string             `bson:"email,omitempty" json:"email"`
	FamilyName string             `bson:"family_name,omitempty" json:"family_name,omitempty"`
	GivenName  string             `bson:"given_name" json:"given_name"`
	Address    string             `bson:"address" json:"address"`
	Metadata   AccountMetadata    `bson:"metadata" json:"metadata"`
	Profile    string             `bson:"profile" json:"profile"`
	Password   string             `bson:"password" json:"-"`
	Device     Device             `bson:"device" json:"-"`
}
type AccountMetadata struct {
	TokenVersion int       `bson:"token_version" json:"-"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
	LasLogin     time.Time `bson:"Last_login" json:"last_login"`
	Provider     string    `bson:"provider" json:"-"` // google , email , facebook
}
type Setting struct {
	ReceiveNotification bool `bson:"receive_notification,omitempty" json:"receive_notification"`
	ReceiveUpdates      bool `bson:"receive_update" json:"receive_update"`
}
type Driver struct {
	Id        string             ` bson:"_id,omitempty" json:"id,omitempty"`
	Role      string             ` bson:"role" json:"role"` // Driver ,...
	Password  string             `json:"-" bson:"password"`
	Email     string             `json:"email" bson:"email"`
	Phone     string             `json:"phone" bson:"phone"`
	GivenName string             `json:"given_name" bson:"given_name"`
	Vehicle   string             `json:"vehicle" bson:"vehicle"`
	Age       primitive.DateTime `json:"age" bson:"age"`
	Address   string             `json:"address" bson:"address"`
}
type Merchant struct {
	Id                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"` // 63f642ac061b6f5f089b3a65
	BusinessName      string             `bson:"business_name" json:"business_name"`
	Location          []float64          `bson:"location" json:"location"`
	Address           string             `bson:"address" json:"address"`
	TimeOperatorStart string             `bson:"operation_time_start" json:"operation_time_start"` // 0710 => 07:19 UTC
	TimeOperatorEnd   string             `bson:"operation_time_end" json:"operation_time_end"`     // 2320 => 23:30 UTC
	PhoneNumber       string             `bson:"phone_number" json:"phone_number"`
	Metadata          Metadata           `bson:"metadata" json:"metadata"`
	Password          string             `json:"-" bson:"password"`
}
type Review struct {
	Id         string   `bson:"_id,omitempty" json:"id,omitempty"`
	Rate       float32  `bson:"rate" json:"rate"`
	Message    string   `bson:"message" json:"message"`
	From       string   `bson:"from" json:"from"`
	Type       string   `bson:"type" json:"type"` // REVIEW_PRODUCT , REVIEW_DRIVER . REVIEW_MERCHANT
	ExternalId string   `bson:"external_id" json:"external_id"`
	Metadata   Metadata `bson:"metadata" json:"metadata"`
}
type Coupon struct {
	Token             string             `json:"token" bson:"token"`
	Rate              uint               `json:"rate" bson:"rate"`                                 // percentage of the discount
	TimeOperatorStart primitive.DateTime `bson:"operation_time_start" json:"operation_time_start"` // Datetime Started
	TimeOperatorEnd   primitive.DateTime `bson:"operation_time_end" json:"operation_time_end"`     // Datetime End coupon expires
}
type Otp struct {
	Code  string `bson:"code" json:"code"`
	Phone string `bson:"phone" json:"phone"`
}
type Device struct {
	DeviceId string `bson:"device_id" json:"device_id"`
	Kind     string `bson:"kind" json:"kind"` // andriod ,ios
}
