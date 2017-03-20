package facts

import (
	"errors"
	"io/ioutil"

	"github.com/r3boot/collective-herder/lib/utils"
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

	if newFacts, err = utils.ParseAsJSON(data); err == nil {
		return newFacts, nil
	}

	if newFacts, err = utils.ParseAsYaml(data); err == nil {
		return newFacts, nil
	}

	if newFacts, err = utils.ParseAsKV(data); err == nil {
		return newFacts, nil
	}

	err = errors.New("GetFactsInFile: Failed to parse " + fileName)
	return nil, err
}
