package facts

import (
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
	"strings"
)

func GetFactsInFile(fileName string) (map[string]interface{}, error) {
	var (
		newFacts map[string]interface{}
		data     []byte
		err      error
	)

	Log.Debug("Reading data from " + fileName)

	if data, err = ioutil.ReadFile(fileName); err != nil {
		err = errors.New("GetFactsInFile: failed to read file: " + err.Error())
		return nil, err
	}

	if newFacts, err = ParseAsJSON(data); err == nil {
		Log.Debug("GetFactsInFile: parsed " + fileName + " as JSON")
		return newFacts, nil
	}

	if newFacts, err = ParseAsYaml(data); err == nil {
		Log.Debug("GetFactsInFile: parsed " + fileName + " as YAML")
		return newFacts, nil
	}

	if newFacts, err = ParseAsKV(data); err == nil {
		Log.Debug("GetFactsInFile: parsed " + fileName + " as KV")
		return newFacts, nil
	}

	err = errors.New("GetFactsInFile: Failed to parse " + fileName)
	return nil, err
}

func ParseAsYaml(data []byte) (map[string]interface{}, error) {
	var (
		newFacts map[string]interface{}
		err      error
	)

	err = yaml.Unmarshal(data, &newFacts)
	return newFacts, err
}

func ParseAsJSON(data []byte) (map[string]interface{}, error) {
	var (
		newFacts map[string]interface{}
		err      error
	)

	err = json.Unmarshal(data, &newFacts)
	return newFacts, err
}

func ParseAsKV(data []byte) (map[string]interface{}, error) {
	var (
		newFacts map[string]interface{}
		line     string
		value    interface{}
		tokens   []string
		err      error
	)

	newFacts = make(map[string]interface{})
	for _, line = range strings.Split(string(data), "\n") {
		tokens = strings.Split(line, "=")
		if len(tokens) >= 2 {
			if value, err = strconv.ParseInt(tokens[1], 10, 64); err == nil {
				newFacts[tokens[0]] = value
				continue
			}

			if value, err = strconv.ParseFloat(tokens[1], 64); err == nil {
				newFacts[tokens[0]] = value
				continue
			}

			if value, err = strconv.ParseBool(tokens[1]); err == nil {
				newFacts[tokens[0]] = value
				continue
			}

			newFacts[tokens[0]] = strings.Join(tokens[1:], "=")
		}
	}

	return newFacts, nil
}
