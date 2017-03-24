package utils

import (
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v2"
	"strconv"
	"strings"
)

func ParseAsYaml(data []byte) (map[string]interface{}, error) {
	var (
		content map[string]interface{}
		err     error
	)

	if data[0] == '{' {
		err = errors.New("data looks like JSON")
		return nil, err
	}

	err = yaml.Unmarshal(data, &content)
	return content, err
}

func ParseAsJSON(data []byte) (map[string]interface{}, error) {
	var (
		rawContent map[string]interface{}
		content    map[string]interface{}
		key        string
		value      interface{}
		value_s    string
		err        error
	)

	rawContent = make(map[string]interface{})
	content = make(map[string]interface{})

	switch len(data) {
	case 0:
		{
			err = errors.New("No data")
			return content, err
		}
	case 1:
		{
			err = errors.New("Too little data")
			return content, err
		}
	}

	if data[0] != '{' {
		err = errors.New("Data does not look like JSON")
		return content, err
	}

	if data[len(data)-1] != '}' {
		err = errors.New("Data does not look like JSON")
		return content, err
	}

	if err = json.Unmarshal(data, &rawContent); err != nil {
		return content, err
	}

	for key, value = range rawContent {

		switch value.(type) {
		case float64:
			{
				value_s = strconv.FormatFloat(value.(float64), 'f', -1, 64)
				if strings.Contains(value_s, ".") {
					content[key] = value.(float64)
				} else {
					content[key] = int(value.(float64))
				}
			}
		default:
			{

				content[key] = value
			}
		}
	}

	return content, err
}

func ParseAsKV(data []byte) (map[string]interface{}, error) {
	var (
		content map[string]interface{}
		line    string
		value   interface{}
		tokens  []string
		err     error
	)

	content = make(map[string]interface{})

	for _, line = range strings.Split(string(data), "\n") {
		tokens = strings.Split(line, "=")
		if len(tokens) >= 2 {
			if value, err = strconv.ParseInt(tokens[1], 10, 64); err == nil {
				content[tokens[0]] = value
				continue
			}

			if value, err = strconv.ParseFloat(tokens[1], 64); err == nil {
				content[tokens[0]] = value
				continue
			}

			if value, err = strconv.ParseBool(tokens[1]); err == nil {
				content[tokens[0]] = value
				continue
			}

			content[tokens[0]] = strings.Join(tokens[1:], "=")
		}
	}

	return content, nil
}
