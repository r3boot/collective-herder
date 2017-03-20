package utils

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"strconv"
	"strings"
)

func ParseAsYaml(data []byte) (map[string]interface{}, error) {
	var (
		content map[string]interface{}
		err     error
	)

	err = yaml.Unmarshal(data, &content)
	return content, err
}

func ParseAsJSON(data []byte) (map[string]interface{}, error) {
	var (
		content map[string]interface{}
		err     error
	)

	err = json.Unmarshal(data, &content)
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
