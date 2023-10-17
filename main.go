package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mohamedabdifitah/ecapi/api"
	"github.com/mohamedabdifitah/ecapi/db"
	"github.com/mohamedabdifitah/ecapi/service"
)

func main() {
	if os.Getenv("APP_ENV") == "development" {
		err := godotenv.Load(".env.local")
		if err != nil {
			log.Fatal(err.Error())
		}

	}
	service.InitRedisClient()
	db.ConnectDB()
	service.InitAlgolia()
	api.Initserver()
}
