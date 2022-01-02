package main

import (
	"fmt"
	"io/ioutil"

	"github.com/bobmaertz/go-scheme/app"
)

func main() {

	//TODO: help banner.
	//TODO: read file.
	//TODO: output file
	//TODO: YAML schema from example
	//TODO: JSON schema version
	// var schemaType string
	var generator app.SchemaGenerator

	data, err := ioutil.ReadFile("./testdata/input.json")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	generator = &app.Json{}

	err = generator.Generate(data)

	if err != nil {
		fmt.Println(fmt.Errorf("encountered a problem creating the schema, %v", err))
		return
	}
}
