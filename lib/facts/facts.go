package facts

import (
	"io/ioutil"
	"os"
	"strconv"

	"github.com/r3boot/collective-herder/lib/utils"
)

var (
	Log utils.Log
)

func NewFacts(l utils.Log) *Facts {
	var (
		f *Facts
	)

	Log = l

	f = &Facts{}
	f.facts = make(map[string]interface{})
	f.LoadAllFacts()

	return f
}

func (f *Facts) LoadAllFacts() {
	var (
		dirName    string
		dirEntries []os.FileInfo
		fileInfo   os.FileInfo
		newFacts   map[string]interface{}
		key        string
		value      interface{}
		fileName   string
		err        error
	)

	for _, dirName = range FACT_DIRS {
		if dirEntries, err = ioutil.ReadDir(dirName); err != nil {
			Log.Warn("LoadAllFacts: Failed to read directory " + dirName + ": " + err.Error())
			continue
		}

		for _, fileInfo = range dirEntries {
			fileName = fileInfo.Name()
			newFacts, err = GetFactsInFile(dirName + "/" + fileName)

			for key, value = range newFacts {
				f.facts[key] = value
			}
		}
	}
}

func (f *Facts) NumFactsAsString() string {
	return strconv.Itoa(len(f.facts))
}

func (f *Facts) GetAll() map[string]interface{} {
	return f.facts
}

func (f *Facts) Get(key string) interface{} {
	if _, ok := f.facts[key]; !ok {
		return nil
	}

	return f.facts[key]
}

func (f *Facts) HasFact(reqFacts map[string]interface{}) bool {
	var (
		reqKey    string
		factKey   string
		reqValue  interface{}
		factValue interface{}
		hasFacts  bool
	)

	if len(reqFacts) == 0 {
		return true
	}

	for reqKey, reqValue = range reqFacts {
		hasFacts = false
		for factKey, factValue = range f.facts {
			if reqKey == factKey && reqValue == factValue {
				hasFacts = true
				break
			}
		}
	}

	return hasFacts
}
