package ping

import (
	"fmt"
	"strconv"
	"time"
)

const (
	NAME string = "ping"
)

type Result struct {
	Node     string
	Uuid     string
	Response string
	Duration time.Duration
}

type ResultSet struct {
	StartTime time.Time
	Data      map[string]Result
}

type Response map[string]string

var (
	resultSet ResultSet
)

func Run(opts map[string]interface{}) interface{} {
	var (
		response Response
	)

	response = make(map[string]string)
	response["value"] = "PONG"

	return response
}

func Print(startTime time.Time, result interface{}) {
	var (
		node     string
		hostUuid string
		duration time.Duration
		value    string
	)

	node = result.(Result).Node
	hostUuid = result.(Result).Uuid
	duration = result.(Result).Duration
	value = result.(Result).Response

	if resultSet.Data == nil {
		resultSet = ResultSet{
			StartTime: startTime,
			Data:      make(map[string]Result),
		}
	}

	resultSet.Data[hostUuid] = result.(Result)

	fmt.Println(value + " response from " + node + " in " + duration.String())
}

func Summary() {
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
