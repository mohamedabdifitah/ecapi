package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// all this values are represented as cm or centimeter as unit type
type Attributes struct {
	Volume float64 `bson:"volume" json:"volume"`
	Height float64 `bson:"height" json:"height"`
	Width  float64 `bson:"width" json:"width"`
	Length float64 `bson:"length" json:"length"`
}

type Menu struct {
	Id                 primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title              string             `json:"title" bson:"title" binding:"max=30,min=1"` // Burger
	Description        string             `json:"description" bson:"description"`            // Chicken fries contains
	Status             string             `json:"status" bson:"status"`                      // available , unavailable , banned
	Category           string             `json:"category" bson:"category"`                  // fast food , drink ,
	Images             []string           `json:"images" bson:"images"`                      // Images of the product urls
	Price              uint               `json:"price" bson:"price"`                        // the price of the product is represented as cents 99 = $0.99
	Attributes         Attributes         `json:"attributes" bson:"attributes"`
	Metadata           Metadata           `json:"metadata" bson:"metadata"`
	Discount           uint               `json:"discount" bson:"discount"` // 10%
	MerchantExternalId string             `json:"merchant_external_id" bson:"merchant_external_id"`
	Reciepe            []string           `json:"reciepe" bson:"reciepe"`
	Barcode            string             `json:"-" bson:"barcode"`                   // if this needed
	EstimateTime       int                `json:"estimate_time" bson:"estimate_time"` // in seconds
}
type Metadata struct {
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
type Item struct {
	ItemExternalId string `json:"item_external_id" bson:"item_external_id"`
	Quantity       uint   `json:"quantity" bson:"quantity"`
	Price          uint   `json:"price" bson:"price"`
}
type Order struct {
	Id                    primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	OrderValue            uint               `json:"order_value" bson:"order_value" ` // Order value is represented as cents 199 = $19.9
	Items                 []Item             `json:"items" bson:"items"`
	DropOffPhone          string             `json:"dropoff_phone" bson:"dropoff_phone"`
	DropOffExteranlId     string             `json:"dropoff_external_id" bson:"dropoff_external_id"`
	DropOffContactName    string             `json:"dropoff_contact_name" bson:"dropoff_contact_name"`
	DropOffTimeEstimated  time.Time          `json:"dropoff_time_estimated" bson:"dropoff_time_estimated"`
	DropOffAddress        string             `json:"dropoff_address" bson:"dropoff_address"`   // address 901 Market Street 6th Floor San Francisco, CA 94103
	DroOffLocation        Location           `json:"dropoff_location" bson:"dropoff_location"` // location cordinates. float([123.1312343,-37.2144343])
	DropOffInstruction    string             `json:"dropoff_instructions" bson:"dropoff_instructions"`
	Stage                 string             `json:"stage" bson:"stage"`                                     // pendding,accepted,preparing,ready,pickuped,deleivered.
	ActionIfUndeliverable string             `json:"action_if_undeliverable" bson:"action_if_undeliverable"` // return_to_pickup
	PickupAddress         string             `json:"pickup_address" bson:"pickup_address"`
	PickUpExternalId      string             `bson:"pickup_external_id" json:"pickup_external_id"`
	PickUpPhone           string             `bson:"pickup_phone" json:"pickup_phone"`
	PickUpLocation        Location           `bson:"pickup_location" json:"pickup_location"`
	PickupTime            time.Time          `bson:"pickup_time" json:"pickup_time"`
	PickupTimeEstimated   int                `bson:"pickup_time_estimated" json:"pickup_time_estimated"` // seconds
	PickupReferenceTag    string             `json:"pickup_reference_tag" bson:"pickup_reference_tag"`
	DriverPhone           string             `bson:"driver_phone" json:"driver_phone"`
	DriverAllowedVehicles []string           `json:"driver_allowed_vehicles" bson:"driver_allowed_vehicles" ` // car , motorcycle , walking
	DriverExternalId      string             `bson:"driver_external_id" json:"driver_external_id" bson:"driver_external_id"`
	Metadata              Metadata           `bson:"metadata" json:"metadata" bson:"metadata"`
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
	Phone      string             `bson:"phone" json:"phone"`
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
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Password    string             `json:"-" bson:"password"`
	Email       string             `json:"email" bson:"email"`
	Phone       string             `json:"phone" bson:"phone"`
	GivenName   string             `json:"given_name" bson:"given_name"`
	VehicleType string             `json:"vehicle_type" bson:"vehicle_type"`
	Age         primitive.DateTime `json:"age" bson:"age"`
	Address     string             `json:"address" bson:"address"`
	Metadata    AccountMetadata    `json:"metadata" bson:"metadata"`
	Profile     string             `json:"profile" bson:"profile"`
	Device      Device             `json:"-" bson:"device"`
}
type Merchant struct {
	Id                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"` // 63f642ac061b6f5f089b3a65
	BusinessName      string             `bson:"business_name" json:"business_name"`
	BusinessEmail     string             `json:"business_email" bson:"business_email"`
	BusinessPhone     string             `bson:"business_phone" json:"business_phone"`
	Location          Location           `bson:"location" json:"location"`
	Address           string             `bson:"address" json:"address"`
	TimeOperatorStart int                `bson:"time_operation_start" json:"time_operation_start"` // 0710 => 07:19 UTC
	TimeOperatorEnd   int                `bson:"time_operation_end" json:"time_operation_end"`     // 2320 => 23:30 UTC
	Metadata          AccountMetadata    `bson:"metadata" json:"metadata"`
	Password          string             `json:"-" bson:"password"`
	Profile           string             `json:"profile" bson:"profile"`
	Device            Device             `json:"-" bson:"device"`
}
type Review struct {
	Id         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"` // 63f642ac061b6f5f089b3a65
	Rate       uint               `bson:"rate" json:"rate"`
	Message    string             `bson:"message" json:"message"`
	From       string             `bson:"from" json:"from"`
	Type       string             `bson:"type" json:"type"` // REVIEW_MENU , REVIEW_DRIVER . REVIEW_MERCHANT
	ExternalId string             `bson:"external_id" json:"external_id"`
	Metadata   Metadata           `bson:"metadata" json:"metadata"`
}
type Coupon struct {
	Token             string             `json:"token" bson:"token"`
	Rate              uint               `json:"rate" bson:"rate"`                                 // percentage of the discount
	TimeOperatorStart primitive.DateTime `bson:"time_operation_start" json:"time_operation_start"` // Datetime Started
	TimeOperatorEnd   primitive.DateTime `bson:"time_operation_end" json:"time_operation_end"`     // Datetime End coupon expires
}
type Otp struct {
	Code  string `bson:"code" json:"code"`
	Phone string `bson:"phone" json:"phone"`
}
type Device struct {
	DeviceId string `bson:"device_id" json:"device_id"`
	Kind     string `bson:"kind" json:"kind"` // andriod ,ios
}
type Location struct {
	Type        string    `bson:"type" json:"type"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates"`
}
