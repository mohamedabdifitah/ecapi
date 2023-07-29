package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mohamedabdifitah/ecapi/api"
	"github.com/mohamedabdifitah/ecapi/db"
)

func main() {
	if os.Getenv("APP_ENV") == "development" {
		err := godotenv.Load(".env.local")
		if err != nil {
			log.Fatal(err.Error())
		}

	}
	db.ConnectDB()
	api.Initserver()

}
