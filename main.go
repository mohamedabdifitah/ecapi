package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mohamedabdifitah/ecapi/api"
	"github.com/mohamedabdifitah/ecapi/db"
	"github.com/mohamedabdifitah/ecapi/service"
	"github.com/mohamedabdifitah/ecapi/template"
)

func hell() {
	if os.Getenv("APP_ENV") == "development" {
		err := godotenv.Load(".env.local")
		if err != nil {
			log.Fatal(err.Error())
		}

	}
	service.InitRedisClient()
	db.ConnectDB()
	service.InitMelliClient()
	api.Initserver()
}
func main() {

	template.AllTemplates.LoadTemplates("./template/template.json")
	// fmt.Println(template.AllTemplates)
	temp, err := template.AllTemplates.TempelateInjector("SwaggGenPath", map[string]string{
		"Path":        "/custom/swagger",
		"method":      "get",
		"tags":        "drivers",
		"summary":     "h",
		"description": "h",
		"operationId": "",
		"responses":   "",
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(temp)
}
