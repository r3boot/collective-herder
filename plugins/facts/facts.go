package facts

import (
	"fmt"
	"time"

	sysFacts "github.com/r3boot/collective-herder/lib/facts"
)

const (
	NAME        string = "facts"
	DESCRIPTION string = "Get available facts from the collective"
)

type Result struct {
	Node     string
	Uuid     string
	Response map[string]interface{}
	Duration time.Duration
}

type ResultSet struct {
	StartTime time.Time
	Data      map[string]Result
}

type DisplayFact struct {
	Name  string
	Count map[string]int
}

var (
	resultSet ResultSet
	nodeFacts *sysFacts.Facts
)

func LoadFacts(f *sysFacts.Facts) {
	nodeFacts = f
}

func Run(opts map[string]interface{}) interface{} {
	return nodeFacts.GetAll()
}

func Print(startTime time.Time, result interface{}) {
	var (
		hostUuid string
	)

	hostUuid = result.(Result).Uuid

	if resultSet.Data == nil {
		resultSet = ResultSet{
			StartTime: startTime,
			Data:      make(map[string]Result),
		}
	}
	resultSet.Data[hostUuid] = result.(Result)
}

func Summary() {
	var (
		hostUuid     string
		allFacts     map[string]interface{}
		displayFacts map[string]map[interface{}]int
		key          string
		value        interface{}
		count        int
		values       map[interface{}]int
	)

	displayFacts = make(map[string]map[interface{}]int)
	values = make(map[interface{}]int)

	for hostUuid, _ = range resultSet.Data {
		allFacts = resultSet.Data[hostUuid].Response
		for key, value = range allFacts {
			if _, ok := displayFacts[key]; !ok {
				displayFacts[key] = make(map[interface{}]int)
			}

			displayFacts[key][value] += 1
		}
	}

	for key, values = range displayFacts {
		fmt.Printf("%-22s", key)
		for value, count = range values {
			fmt.Printf("(%2d) %-20v", count, value)
		}
		fmt.Printf("\n")
	}
}
