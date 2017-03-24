package utils

import (
	"reflect"
	"testing"
)

var VALID_YAML string = `---
TestString: "String value"
TestInt: 666
TestBool: true
TestFloat: 6.66
`

var VALID_JSON string = `{
	"TestString":"String value",
	"TestInt": 666,
	"TestBool": true,
	"TestFloat": 6.66
}`

var INVALID_JSON1 string = "{"
var INVALID_JSON2 string = "{ "
var INVALID_JSON3 string = "{ bladiebla }"

var VALID_KV string = `TestString=String Value
TestInt=666
TestBool=true
TestFloat=6.66`

func validateResults(t *testing.T, content map[string]interface{}) {
	// Parse string value
	if reflect.TypeOf(content["TestString"]).Name() != "string" {
		t.Error("type of yaml[TestString] != string")
	}
	if content["TestString"].(string) != "String value" {
		t.Error("yaml[TestString] != String value")
	}

	// Parse int value
	if reflect.TypeOf(content["TestInt"]).Name() != "int" {
		t.Error("type of yaml[TestInt] != int")
	}
	if content["TestInt"].(int) != 666 {
		t.Error("yaml[TestInt] != 666")
	}

	// Parse bool value
	if reflect.TypeOf(content["TestBool"]).Name() != "bool" {
		t.Error("type of yaml[TestBool] != bool")
	}
	if !content["TestBool"].(bool) {
		t.Error("yaml[TestBool] != true")
	}

	// Parse float value
	if reflect.TypeOf(content["TestFloat"]).Name() != "float64" {
		t.Error("type of yaml[TestFloat] != float")
	}
	if content["TestFloat"] != 6.66 {
		t.Error("yaml[TestFloat] != 6.66")
	}
}

func TestParseAsYaml(t *testing.T) {
	var (
		content map[string]interface{}
		err     error
	)

	if content, err = ParseAsYaml([]byte(VALID_YAML)); err != nil {
		t.Error("Failed to parse VALID_YAML: " + err.Error())
	}

	if content == nil {
		t.Error("Content == nil")
		return
	}

	if len(content) == 0 {
		t.Error("Content size == 0")
		return
	}

	validateResults(t, content)

	if content, err = ParseAsYaml([]byte(VALID_JSON)); err == nil {
		t.Error("Parsed JSON as Yaml")
	}

	if content, err = ParseAsYaml([]byte(VALID_KV)); err == nil {
		t.Error("Parsed KV as Yaml")
	}
}

func TestParseAsJSON(t *testing.T) {
	var (
		content map[string]interface{}
		err     error
	)

	if content, err = ParseAsJSON([]byte(VALID_JSON)); err != nil {
		t.Error("Failed to parse VALID_JSON: " + err.Error())
	}

	if content == nil {
		t.Error("Content == nil")
		return
	}

	if len(content) == 0 {
		t.Error("Content size == 0")
		return
	}

	validateResults(t, content)

	if content, err = ParseAsJSON([]byte(VALID_YAML)); err == nil {
		t.Error("Parsed Yaml as JSON")
	}

	if content, err = ParseAsJSON([]byte(VALID_KV)); err == nil {
		t.Error("Parsed KV as JSON")
	}

	if _, err = ParseAsJSON([]byte(INVALID_JSON1)); err == nil {
		t.Error("Parsed INVALID_JSON1 as JSON")
	}

	if _, err = ParseAsJSON([]byte(INVALID_JSON2)); err == nil {
		t.Error("Parsed INVALID_JSON2 as JSON")
	}

	if _, err = ParseAsJSON([]byte(INVALID_JSON3)); err == nil {
		t.Error("Parsed INVALID_JSON3 as JSON")
	}

	if _, err = ParseAsJSON([]byte("")); err == nil {
		t.Error("Parsed nil as JSON")
	}
}

func TestParseAsKV(t *testing.T) {
	var (
		content map[string]interface{}
		err     error
	)

	if content, err = ParseAsKV([]byte(VALID_KV)); err != nil {
		t.Error("Failed to parse VALID_KV: " + err.Error())
	}

	if content == nil {
		t.Error("Content == nil")
		return
	}

	if len(content) == 0 {
		t.Error("Content size == 0")
		return
	}
}
