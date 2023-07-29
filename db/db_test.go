package db

import (
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func TestConnectDB(t *testing.T) {
	err := Mongoclient.Ping(Ctx, readpref.Primary())
	if err != nil {
		log.Fatal("error while trying to ping mongo", err)
	}
}
