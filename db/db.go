package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// TODO: support db validation and error handling which almost return description what error happened
var (
	Mongoclient *mongo.Client
	Ctx         context.Context
	// cancel             context.CancelFunc
	err                error
	CustomerCollection *mongo.Collection
	DriverCollection   *mongo.Collection
	OrderCollection    *mongo.Collection
	MerchantCollection *mongo.Collection
	MenuCollection     *mongo.Collection
	ReviewCollection   *mongo.Collection
	CategoryCollection *mongo.Collection
)

func ConnectDB() {
	Ctx = context.TODO()
	Mongoclient, err = mongo.NewClient(options.Client().ApplyURI(os.Getenv("DB_URI")))
	if err != nil {
		log.Fatal(err)
	}
	Ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err = Mongoclient.Connect(Ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = Mongoclient.Ping(Ctx, readpref.Primary())
	if err != nil {
		log.Fatal("error while trying to ping mongo", err)
	}
	fmt.Println("mongo connection established")
	ResturantDB := Mongoclient.Database("resys")

	CustomerCollection = ResturantDB.Collection("customer")
	DriverCollection = ResturantDB.Collection("driver")
	MerchantCollection = ResturantDB.Collection("merchant")
	MenuCollection = ResturantDB.Collection("menu")
	ReviewCollection = ResturantDB.Collection("review")
	OrderCollection = ResturantDB.Collection("order")
	CategoryCollection = ResturantDB.Collection("category")
	CreateIndex("email", CustomerCollection)
	CreateIndex("phone", DriverCollection)
	CreateIndex("business_name", MerchantCollection)
	CreateIndex("business_email", MerchantCollection)
	CreateGeoIndex("location", MerchantCollection)
	CreateGeoIndex("dropoff_location", OrderCollection)
	CreateGeoIndex("pickup_location", OrderCollection)
	CreateGeoIndex("location", DriverCollection)
	CreateIndex("order_id", ReviewCollection)
	CreateIndex("title", CategoryCollection)
}
func CloseDB() error {
	// Ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	err := Mongoclient.Disconnect(Ctx)

	return err
}
func CreateIndex(name string, Collection *mongo.Collection) {
	_, err := Collection.Indexes().CreateOne(
		Ctx,
		mongo.IndexModel{
			Keys:    bson.D{{Key: name, Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		panic(err)
	}
}
func CreateGeoIndex(name string, Collection *mongo.Collection) {
	_, err := Collection.Indexes().CreateOne(
		Ctx,
		mongo.IndexModel{
			Keys: bson.D{{Key: name, Value: "2dsphere"}},
		},
	)
	if err != nil {
		panic(err)
	}
}
