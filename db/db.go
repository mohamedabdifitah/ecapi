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

var (
	Mongoclient        *mongo.Client
	Ctx                context.Context
	CustomerCollection *mongo.Collection
	DriverCollection   *mongo.Collection
	OrderCollection    *mongo.Collection
	MerchantCollection *mongo.Collection
)

func ConnectDB() {
	Ctx = context.TODO()
	Mongoclient, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("DB_URI")))
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
	ResturantDB := Mongoclient.Database("resturant")

	CustomerCollection = ResturantDB.Collection("customer")
	DriverCollection = ResturantDB.Collection("driver")
	MerchantCollection = ResturantDB.Collection("merchant")
	CreateIndex("username", CustomerCollection)
	CreateIndex("email", CustomerCollection)
	CreateIndex("email", DriverCollection)
	CreateIndex("phone", DriverCollection)
	CreateIndex("given_name", DriverCollection)
	CreateIndex("business_name", MerchantCollection)
	CreateIndex("phone_number", MerchantCollection)

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
