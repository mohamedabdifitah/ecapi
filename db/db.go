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
	cancel             context.CancelFunc
	err                error
	CustomerCollection *mongo.Collection
	DriverCollection   *mongo.Collection
	OrderCollection    *mongo.Collection
	MerchantCollection *mongo.Collection
	MenuCollection     *mongo.Collection
	ReviewCollection   *mongo.Collection
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
	// dbList, err := Mongoclient.ListDatabaseNames(Ctx, bson.D{})
	// if err != nil {
	// 	fmt.Println("error while trying to list database names", err)
	// }
	// fmt.Println(dbList)
	fmt.Println("mongo connection established")
	ResturantDB := Mongoclient.Database("resys")

	CustomerCollection = ResturantDB.Collection("customer")
	DriverCollection = ResturantDB.Collection("driver")
	MerchantCollection = ResturantDB.Collection("merchant")
	MenuCollection = ResturantDB.Collection("menu")
	ReviewCollection = ResturantDB.Collection("review")
	CreateIndex("email", CustomerCollection)
	CreateIndex("phone", DriverCollection)
	CreateIndex("business_name", MerchantCollection)
	CreateIndex("business_phone", MerchantCollection)
	CreateIndex("barcode", MenuCollection)
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
