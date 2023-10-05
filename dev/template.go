package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

type Templates map[string]string

func (templates Templates) LoadTemplates(paths ...string) {
	for _, path := range paths {
		fileContent, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer fileContent.Close()
		byteResult, err := ioutil.ReadAll(fileContent)
		if err != nil {
			fmt.Println(err)
		}
		templates[fileContent.Name()] = fmt.Sprintf("%s", string(byteResult))

	}
}
func (templates Templates) TempelateInjector(key string, definition map[string]string) (string, error) {
	temp := templates.Lookup(key)
	tmpl, err := template.New(key).Parse(temp)
	if err != nil {
		return "", err
	}
	writer := new(bytes.Buffer)
	err = tmpl.Execute(writer, definition)
	if err != nil {
		log.Fatal(err)
	}
	return writer.String(), nil
}
func (templates Templates) Lookup(title string) string {
	return templates[title]
}
