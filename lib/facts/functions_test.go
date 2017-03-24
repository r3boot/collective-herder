package facts

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/r3boot/collective-herder/lib/utils"
)

var YAML_TEST_DATA = `---
TestKey: "testvalue"`

var JSON_TEST_DATA = `{"TestKey":"testvalue"}`

var KV_TEST_DATA = `TestKey=testvalue`

var INVALID_TEST_DATA = `nothing really useful

and more unuseful stuff`

func createFactsFile(t *testing.T, content string) (string, string) {
	var (
		dirName string
		fd      *os.File
		err     error
	)

	if dirName, err = ioutil.TempDir("/tmp", "ch-functions_test"); err != nil {
		t.Fatalf("Failed to create tempdir: " + err.Error())
	}

	if fd, err = ioutil.TempFile(dirName, "ch-functions_test"); err != nil {
		t.Fatalf("Failed to create tempfile: " + err.Error())
	}
	fd.Close()

	if err = ioutil.WriteFile(fd.Name(), []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create tempfile: " + err.Error())
	}

	return dirName, fd.Name()
}

func createLogFile(t *testing.T) utils.Log {
	var (
		fd  *os.File
		err error
	)

	if fd, err = ioutil.TempFile("/tmp", "ch-functions_test"); err != nil {
		t.Fatalf("Failed to create tempfile: " + err.Error())
	}

	return utils.Log{
		TestFd: fd,
	}
}

func clearFactsFile(dirName, fileName string) {
	os.Remove(fileName)
	os.RemoveAll(dirName)
}

func clearLogFile(fd *os.File) {
	fd.Close()
	os.Remove(fd.Name())
}

func runTestsOnData(t *testing.T, content string) {
	var (
		dirName    string
		fileName   string
		value      interface{}
		otherValue interface{}
		reqFacts   map[string]interface{}
		l          utils.Log
		f          *Facts
	)

	dirName, fileName = createFactsFile(t, content)
	l = createLogFile(t)
	defer clearFactsFile(dirName, fileName)
	defer clearLogFile(l.TestFd)

	FACT_DIRS = []string{dirName}

	f = NewFacts(l)

	if f.facts == nil {
		t.Fatalf("f.facts == nil")
	}

	if f.NumFactsAsString() != "1" {
		t.Fatalf("f.NumFactsAsString != 1")
	}

	if len(f.GetAll()) != 1 {
		t.Fatalf("len(f.GetAll()) != 1")
	}

	value = f.Get("TestKey")
	if value == nil {
		t.Fatalf("value == nil")
	}

	otherValue = f.Get("somerandomkey")
	if otherValue != nil {
		t.Fatalf("f.Get: somerandomkey found")
	}

	if reflect.TypeOf(value).String() != "string" {
		t.Fatalf("type of value != string")
	}

	if value.(string) != "testvalue" {
		t.Fatalf("value != testvalue")
	}

	reqFacts = make(map[string]interface{})
	reqFacts["TestKey"] = "testvalue"
	if !f.HasFact(reqFacts) {
		t.Fatalf("f.HasFact(): TestKey not found")
	}

	reqFacts = make(map[string]interface{})
	if !f.HasFact(reqFacts) {
		t.Fatalf("f.HasFact(): no fallback to true")
	}
}

func runCornerCases(t *testing.T, content string) {
	var (
		dirName  string
		fileName string
		err      error
		l        utils.Log
		f        *Facts
	)

	if _, err = GetFactsInFile("./nonexistent"); err == nil {
		t.Fatalf("Able to read facts from nonexistent file")
	}

	dirName, fileName = createFactsFile(t, content)
	l = createLogFile(t)
	defer clearFactsFile(dirName, fileName)
	defer clearLogFile(l.TestFd)

	FACT_DIRS = []string{"/some/random/nonexisting/directory"}
	f = NewFacts(l)
	if len(f.facts) != 0 {
		t.Fatalf("len(f.facts) != 0")
	}
}

func TestGetFactsInFile(t *testing.T) {
	runTestsOnData(t, YAML_TEST_DATA)
	runTestsOnData(t, JSON_TEST_DATA)
	runTestsOnData(t, KV_TEST_DATA)
	runCornerCases(t, INVALID_TEST_DATA)
}
