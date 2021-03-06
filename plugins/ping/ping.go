package ping

import (
	"fmt"
	"strconv"
	"time"
)

const (
	NAME        string = "ping"
	DESCRIPTION string = "Send a ping to the collective"
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

var (
	resultSet ResultSet
)

func ParseArgs(args []string) map[string]interface{} {
	return make(map[string]interface{})
}

func PreRun(opts map[string]interface{}) {
}

func Run(opts map[string]interface{}) interface{} {
	var (
		response map[string]string
	)

	response = make(map[string]string)
	response["value"] = "PONG"

	return response
}

func Print(startTime time.Time, result interface{}, opts map[string]interface{}) {
	var (
		node     string
		hostUuid string
		duration time.Duration
		response map[string]interface{}
		value    string
	)

	node = result.(Result).Node
	hostUuid = result.(Result).Uuid
	duration = result.(Result).Duration
	response = result.(Result).Response
	value = response["value"].(string)

	if resultSet.Data == nil {
		resultSet = ResultSet{
			StartTime: startTime,
			Data:      make(map[string]Result),
		}
	}

	resultSet.Data[hostUuid] = result.(Result)

	fmt.Println(value + " response from " + node + " in " + duration.String())
}

func Summary(opts map[string]interface{}) {
	var (
		min           time.Duration
		max           time.Duration
		totalDuration time.Duration
		avg           float64
		avg_s         string
		result        Result
	)

	for _, result = range resultSet.Data {
		if (min == 0) || (result.Duration < min) {
			min = result.Duration
		}

		if (max == 0) || (result.Duration > max) {
			max = result.Duration
		}
		totalDuration += result.Duration
	}

	avg = (float64(totalDuration.Nanoseconds()) / float64(len(resultSet.Data)*1000000))
	avg_s = strconv.FormatFloat(avg, 'f', 6, 64)
	fmt.Println("\nSummary: min/avg/max = " + min.String() + "/" + avg_s + "ms/" + max.String())
}
