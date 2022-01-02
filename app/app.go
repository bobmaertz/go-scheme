package app

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type SchemaGenerator interface {
	Generate(input []byte) error
}

type Json struct {
	Root Object
}

type Type string

func (j *Json) Generate(input []byte) error {

	b := make(map[string]interface{})
	rawInput := []byte(input)
	err := json.Unmarshal(rawInput, &b)
	//Parse input tree..
	if err != nil {
		return fmt.Errorf("unable to unmarshall raw input into map[string]interface: %v", err)
	}

	root := &Object{
		SchemaName:  "http://json-schema.org/draft-07/schema",
		Id:          "#",
		Type:        "object",
		Description: "The root schema is the schema that comprises the entire JSON document.",
		Title:       "The Root Schema",
		// Properties:  make(map[string]Property),
	}

	out, err := parse(*root, b)
	t, ok := out.(*Object)
	if !ok {
		return fmt.Errorf("unable to set out to object")
	}
	j.Root = *t

	if err != nil {
		return fmt.Errorf("unable to parse input: %v", err)
	}

	//Output result

	//json.Marshal doesnt handle cyclic data structures, so the Object cant be returned here..
	// Need to create a marshaller OR convert Object to a map[string]interface}{}
	a, err := json.MarshalIndent(j.Root, "", "    ")
	if err != nil {
		return fmt.Errorf("unable to marshal output: %v", err)
	}
	fmt.Println(string(a))
	return nil
}

//TODO: Can root be a method reciever?
func parse(root Object, data map[string]interface{}) (interface{}, error) {
	for key, value := range data {

		switch t := value.(type) {
		case map[string]interface{}:
			// parse(value)
			root.Required = append(root.Required, key)
		case string:
			o, _ := createObject(key, value.(string), root.Id, "string")
			if root.Properties == nil {
				root.Properties = make(map[string]temp)
			}
			root.Properties[key] = *o
			root.Required = append(root.Required, key) //TODO: Do a check here
		case int32:
			fmt.Printf("[x] key %s im a int\n", key)
		case float64:
			//Float64 seems to be a default for the marshaller, so check for a decimal to see if
			// the type is really an int..
			var o *temp
			if strings.Contains(fmt.Sprintf("%s", value), ".") {
				o, _ = createObject(key, float64(t), root.Id, "float")
			} else {
				o, _ = createObject(key, int64(t), root.Id, "int")
			}
			if root.Properties == nil {
				root.Properties = make(map[string]temp)
			}
			root.Properties[key] = *o
			root.Required = append(root.Required, key) //TODO: Do a check here
		case bool:
			o, _ := createObject(key, value, root.Id, "bool")
			if root.Properties == nil {
				root.Properties = make(map[string]temp)
			}
			root.Properties[key] = *o
			root.Required = append(root.Required, key) //TODO: Do a check here
		case []interface{}:
			fmt.Printf("[x] key %s im a array\n", key)
		default:
			typ := reflect.TypeOf(value)
			fmt.Printf("[] key %s is a type %v, %v\n", key, t, typ)
		}
	}
	return &root, nil
}

func createObject(key string, value interface{}, path string, typ string) (*temp, error) {
	t := &temp{
		Id:          makeId(path, key),
		Typ:         "bool",
		Title:       fmt.Sprintf("The %s Schema", strings.Title(key)),
		Description: "An explanation about the purpose of this instance.",
		Examples:    Example{value},
	}

	return t, nil
}

type temp struct {
	Id          string  `json:"id,omitempty"`
	Typ         string  `json:"type,omitempty"`
	Title       string  `json:"titke,omitempty"`
	Description string  `json:"d,omitempty"`
	Examples    Example `json:"e,omitempty"`
}

func makeId(path, key string) string {
	return fmt.Sprintf("%s/%s", path, key)
}
