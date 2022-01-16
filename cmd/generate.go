package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
)

//TODO: a placeholder for now, update with proper schema
var (
	Schema_Draft_07 = "http://json-schema.org/draft-07/schema#"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a json schema from the provided input",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(args)
		raw, err := readInput()
		if err != nil {
			return err
		}

		if err := generate(raw); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}

type Object map[string]interface{}

func readInput() ([]byte, error) {

	if inputFile == "" {
		return nil, fmt.Errorf("STDIN not supported")
	}

	raw, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read input file: %v", err)
	}

	return raw, nil
}

//generate - generates graph
func generate(raw []byte) error {

	//Pass / in to denote that this is the root
	output := make(map[string]interface{})
	output["$schema"] = Schema_Draft_07
	output["$id"] = makeId("https://example.com", "root")
	output["title"] = strings.Title("root")
	output["definitions"] = struct{}{}
	output["type"] = "object"
	output["required"] = []string{}

	result, err := traverse("#root", raw)
	if err != nil {
		return fmt.Errorf("unable to traverse graph: %v", err)
	}

	output["properties"] = result
	b, err := json.MarshalIndent(output, "", "\t")
	if err != nil {
		return fmt.Errorf("unable to marshall output: %v", err)
	}

	//Return / Print output
	fmt.Print(string(b))

	return nil
}

func traverse(path string, rawJson []byte) (output map[string]interface{}, err error) {
	var result map[string]interface{}

	json.Unmarshal(rawJson, &result)

	output = make(map[string]interface{})

	for key, value := range result {
		object, err := create(path, key, value)
		if err != nil {
			return nil, err
		}
		output[key] = object

	}

	//TODO: If required keys, add here?
	return output, nil

}
func makeId(path, key string) string {
	return fmt.Sprintf("%s/%s", path, key)
}

func create(path, key string, v interface{}) (object map[string]interface{}, err error) {
	object = make(map[string]interface{})
	object["$id"] = makeId(path, key)
	object["title"] = strings.Title(key)

	switch v.(type) {
	case float64:
		//TODO: Currently assume this is an integer...
		val := fmt.Sprintf("%v", v)
		if isFloat := strings.Contains(val, "."); isFloat {
			object["default"] = 0
			object["type"] = "number"
			return
		}

		object["examples"] = []interface{}{v}
		object["default"] = 0.0
		object["type"] = "integer"
	case bool:
		object["examples"] = []interface{}{v}
		object["default"] = v
		object["type"] = "boolean"
	case string:
		object["examples"] = []interface{}{v}
		object["default"] = ""
		object["pattern"] = "^.*$"
		object["type"] = "string"
	case []interface{}:
		//For each in slice, traverse
		arr := v.([]interface{})
		var out map[string]interface{}
		id := makeId(path, key)
		if len(arr) > 0 {
			out, err = create(id, "items", arr[0])
			if err != nil {
				return nil, err
			}
		}
		object["examples"] = []interface{}{arr[0]}
		object["default"] = []interface{}{}
		object["type"] = "array"
		object["items"] = out

	case map[string]interface{}:
		id := makeId(path, key)

		values, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal value for key %s: %v", key, err)
		}
		o, err := traverse(id, values)
		if err != nil {
			return nil, fmt.Errorf("unable to traverse value in key %s: %v", key, err)
		}
		object["type"] = "object"
		//TODO: Required
		object["properties"] = o
	default:
		fmt.Printf("%s unknown\n", key)
	}
	return object, nil
}
