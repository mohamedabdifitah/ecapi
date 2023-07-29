package main

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/mohamedabdifitah/ecapi/db"
)

func TestMain(m *testing.M) {
	if os.Getenv("APP_ENV") == "development" {
		err := godotenv.Load(".env.local")
		if err != nil {
			log.Fatal(err.Error())
		}

	}
	db.ConnectDB()
	// api.Initserver()
	m.Run()
	db.CloseDB()
}
