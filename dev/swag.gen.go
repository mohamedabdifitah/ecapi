package main

import (
	"fmt"
	"log"
	"os"
	// "github.com/mohamedabdifitah/ecapi/api"
)

func Write(de string) {
	f, err := os.Create("dev/data.json")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(de)

	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println("done")
}

func main() {
	var AllTemplates = make(Templates)
	AllTemplates.LoadTemplates("dev/base.txt", "dev/path.txt", "dev/responses.txt")
	response, err := AllTemplates.TempelateInjector("dev/responses.txt", map[string]string{
		"properties": "\"given_name\": {\n                  \"type\": \"string\"\n                }",
		"type":       "object",
	})
	two, err := AllTemplates.TempelateInjector("dev/path.txt", map[string]string{
		"pathname":  "/drive/test",
		"responses": response,
	})
	msg, err := AllTemplates.TempelateInjector("dev/base.txt", map[string]string{
		"Paths": two,
	})
	if err != nil {
		fmt.Println(err)
	}
	Write(msg)
}
