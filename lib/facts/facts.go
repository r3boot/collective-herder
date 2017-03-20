package facts

import (
	"io/ioutil"
	"os"
	"strconv"

	"github.com/r3boot/collective-herder/lib/utils"
)

var (
	Log utils.Logger
)

func NewFacts(l utils.Logger) *Facts {
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
