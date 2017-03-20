package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func ReadFile(fileName string) (Config, error) {
	var (
		config Config
		data   []byte
		err    error
	)

	if data, err = ioutil.ReadFile(fileName); err != nil {
		err = errors.New("Failed to read config: " + err.Error())
		return Config{}, err
	}

	if err = yaml.Unmarshal(data, &config); err != nil {
		err = errors.New("Failed to parse config: " + err.Error())
		return Config{}, err
	}

	return config, nil
}
