package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO: support attributes for basic types and custom type for menues.
// all this values are represented as cm or centimeter as unit type
type Attributes struct {
	Volume float64 `bson:"volume" json:"volume"`
	Height float64 `bson:"height" json:"height"`
	Width  float64 `bson:"width" json:"width"`
	Length float64 `bson:"length" json:"length"`
}

type Menu struct {
	Id                  primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title               string             `json:"title" bson:"title" binding:"max=30,min=1"` // Burger
	Description         string             `json:"description" bson:"description"`            // Chicken fries contains
	Status              string             `json:"status" bson:"status"`                      // available , unavailable , banned
	Category            string             `json:"category" bson:"category"`                  // fast food , drink ,
	Images              []string           `json:"images" bson:"images"`                      // Images of the product urls
	Price               uint               `json:"price" bson:"price"`                        // the price of the product is represented as cents 99 = $0.99
	Attributes          Attributes         `json:"attributes" bson:"attributes"`
	Metadata            Metadata           `json:"metadata" bson:"metadata"`
	MerchantExternalId  string             `json:"merchant_external_id" bson:"merchant_external_id"`
	Reciepe             []string           `json:"reciepe" bson:"reciepe"`             // floor , meat , egg etc.
	Barcode             string             `json:"-" bson:"barcode"`                   // if this needed
	EstimateTime        int                `json:"estimate_time" bson:"estimate_time"` // estimate of preparation and cooking time in seconds
	ServiceAvailability []ActiveDays       ` json:"service_availablity" bson:"service_availablity"`
	Likes               uint               `json:"likes" bson:"likes"` //
}
type Category struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Icon        string             `json:"icon" bson:"icon"`
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
	DisplayId             string             `bson:"display_id" json:"display_id"`
	OrderValue            uint               `json:"order_value" bson:"order_value" ` // Order value is represented as cents 199 = $19.9
	Type                  string             `json:"type" bson:"type"`                // pickup or delivery order
	Items                 []Item             `json:"items" bson:"items"`
	DropOffPhone          string             `json:"dropoff_phone" bson:"dropoff_phone"`
	DropOffExteranlId     string             `json:"dropoff_external_id" bson:"dropoff_external_id"`
	DropOffContactName    string             `json:"dropoff_contact_name" bson:"dropoff_contact_name"`
	DropOffTimeEstimated  int                `json:"dropoff_time_estimated" bson:"dropoff_time_estimated"`
	DropOffTime           time.Time          `json:"dropoff_time" bson:"dropoff_time"`
	DropOffAddress        string             `json:"dropoff_address" bson:"dropoff_address"`   // address 901 Market Street 6th Floor San Francisco, CA 94103
	DroOffLocation        Location           `json:"dropoff_location" bson:"dropoff_location"` // location cordinates. float([123.1312343,-37.2144343])
	DropOffInstruction    string             `json:"dropoff_instructions" bson:"dropoff_instructions"`
	Stage                 string             `json:"stage" bson:"stage"`                                     // placed,accepted,preparing,ready,pickuped,deleivered,cancelled.
	ActionIfUndeliverable string             `json:"action_if_undeliverable" bson:"action_if_undeliverable"` // return_to_merchant // destroy
	PickupAddress         string             `json:"pickup_address" bson:"pickup_address"`
	PickUpExternalId      string             `bson:"pickup_external_id" json:"pickup_external_id"`
	PickUpName            string             `bson:"pickup_name" json:"pickup_name"`
	PickUpPhone           string             `bson:"pickup_phone" json:"pickup_phone"`
	PickUpLocation        Location           `bson:"pickup_location" json:"pickup_location"`
	PickupTime            time.Time          `bson:"pickup_time" json:"pickup_time"`
	PickupEstimatedTime   int                `bson:"pickup_estimated_time" json:"pickup_estimated_time"` // seconds
	DriverPhone           string             `bson:"driver_phone" json:"driver_phone"`
	DriverAllowedVehicles []string           `json:"driver_allowed_vehicles" bson:"driver_allowed_vehicles" ` // car , motorcycle , walking
	DriverExternalId      string             `bson:"driver_external_id" json:"driver_external_id"`
	Metadata              Metadata           `bson:"metadata" json:"metadata"`
	CancelReason          string             `bson:"cancel_reason" json:"cancel_reason"`
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
	Preference []string           `bson:"preference" json:"preference"`
}
type AccountMetadata struct {
	TokenVersion    int       `bson:"token_version" json:"-"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time `json:"-" bson:"updated_at"`
	LasLogin        time.Time `bson:"last_login" json:"-"`
	WebhookEndpoint string    `json:"webhook_endpoint" bson:"webhook_endpoint"`
	Provider        string    `bson:"provider" json:"-"` // google , email , facebook
}
type Setting struct {
	ReceiveNotification bool `bson:"receive_notification,omitempty" json:"receive_notification"`
	ReceiveUpdates      bool `bson:"receive_update" json:"receive_update"`
}
type Driver struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Password  string             `json:"-" bson:"password"`
	Email     string             `json:"email" bson:"email"`
	Phone     string             `json:"phone" bson:"phone"`
	GivenName string             `json:"given_name" bson:"given_name"`
	Age       primitive.DateTime `json:"age" bson:"age"`
	Address   string             `json:"address" bson:"address"`
	Metadata  AccountMetadata    `json:"metadata" bson:"metadata"`
	Profile   string             `json:"profile" bson:"profile"`
	Device    Device             `json:"-" bson:"device"`
	Rate      Rate               `json:"rate" bson:"rate"`
	Location  Location           `json:"location" bson:"location"`
	Satus     bool               `json:"status" bson:"status"` // true available, false unavailable
	Vehicle   Vehicle            `json:"vehicle" bson:"vehicle"`
}
type Vehicle struct {
	Model   string  `json:"model" bson:"model"`     // Tesla , bmw ,
	Type    string  `json:"type" bson:"type"`       // car , motorbike , bicycle , truck , etc.
	Payload float64 `json:"payload" bson:"payload"` // how much payload can vehicle carry in kg
}
type Merchant struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"` // 63f642ac061b6f5f089b3a65
	BusinessName  string             `bson:"business_name" json:"business_name"`
	BusinessEmail string             `json:"business_email" bson:"business_email"`
	BusinessPhone string             `bson:"business_phone" json:"business_phone"`
	Location      Location           `bson:"location" json:"location"`
	Address       string             `bson:"address" json:"address"`
	Metadata      AccountMetadata    `bson:"metadata" json:"metadata"`
	Password      string             `json:"-" bson:"password"`
	Profile       string             `json:"profile" bson:"profile"`
	Device        Device             `json:"-" bson:"device"`
	Rate          Rate               `json:"rate" bson:"rate"`
	Badge         Badge              `json:"badge" bson:"badge"`
	Closed        bool               `json:"closed" bson:"closed"`
	ActiveDays    []ActiveDays       `json:"active_days" bson:"active_days"`
	Generes       []string           `bson:"generes" json:"generes"` //Fast food.Fast casual.Casual dining / Slow Casual.Premium casual.Family style.Fine dining.
	Likes         uint               `bson:"likes" json:"likes"`
	Popular       uint               `bson:"popular" json:"popular"`
}
type ActiveDays struct {
	TimeOperatorStart float32 `bson:"time_operation_start" json:"time_operation_start"` // 7.10 => 07:19 UTC
	TimeOperatorEnd   float32 `bson:"time_operation_end" json:"time_operation_end"`
}
type Review struct {
	Id         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"` // 63f642ac061b6f5f089b3a65
	Type       string             `bson:"type" json:"type"`                  // REVIEW_MERCHANT , REVIEW_DRIVER , REVIEW_MENU
	OrderId    string             `bson:"order_id" json:"order_id"`
	From       string             `bson:"from" json:"from"`
	Message    string             `bson:"message" json:"message"`
	Rate       uint               `bson:"rate" json:"rate"`
	ExternalId string             `bson:"external_id" json:"external_id"`
	Options    []string           `bson:"options" json:"options"`
	Metadata   Metadata           `bson:"metadata" json:"metadata"`
}
type Coupon struct {
	Token             string    `json:"token" bson:"token"`
	Rate              uint      `json:"rate" bson:"rate"`                                 // percentage of the discount
	TimeOperatorStart time.Time `bson:"time_operation_start" json:"time_operation_start"` // Datetime Started
	TimeOperatorEnd   time.Time `bson:"time_operation_end" json:"time_operation_end"`     // Datetime End coupon expires
}
type Device struct {
	DeviceId string `bson:"device_id" json:"device_id"`
	Kind     string `bson:"kind" json:"kind"` // andriod ,ios
}
type Location struct {
	Type        string    `bson:"type" json:"type"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates"`
}
type Rate struct {
	Rate         float32 `bson:"rate" json:"rate"`
	ScoreStats   []int   `bson:"stats" json:"-"` // [20,20,0,10,5]
	Participants uint    `bson:"participants" json:"participants"`
}
type Badge struct {
	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
}
