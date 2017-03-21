package facts

import (
	"fmt"
	"os"
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

func ParseArgs(args []string) map[string]interface{} {
	var (
		query []string
		opts  map[string]interface{}
		value string
	)

	opts = make(map[string]interface{})
	opts["help"] = false

	for _, value = range args {
		switch value {
		case "--help", "-h":
			{
				opts["help"] = true
			}
		default:
			{
				query = append(query, value)
			}
		}
	}

	opts["query"] = query

	return opts
}

func PreRun(opts map[string]interface{}) {
	if opts == nil {
		return
	}

	if opts["help"].(bool) {
		fmt.Fprintf(os.Stderr, "Usage of facts plugin: ch facts [-h, --help] [<fact> ...]\n\n")
		fmt.Fprintf(os.Stderr, "Query the collective for its facts. By specifying no\n")
		fmt.Fprintf(os.Stderr, "options, all facts and their values will be retrieved and\n")
		fmt.Fprintf(os.Stderr, "displayed. You can query for one or more facts by adding\n")
		fmt.Fprintf(os.Stderr, "them as space-separated parameters.\n")
		os.Exit(2)
	}
}

func Run(opts map[string]interface{}) interface{} {
	var (
		fact   interface{}
		value  interface{}
		query  []interface{}
		result map[string]interface{}
	)

	if opts == nil {
		return nodeFacts.GetAll()
	}

	if _, ok := opts["query"]; !ok {
		return nil
	}

	query = opts["query"].([]interface{})
	result = make(map[string]interface{})

	for _, fact = range query {
		if value = nodeFacts.Get(fact.(string)); value == nil {
			continue
		}
		result[fact.(string)] = value
	}

	return result
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
