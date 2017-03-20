package ping

import (
	"fmt"
	"time"
)

const (
	NAME string = "ping"
)

type Response map[string]string

func Run(opts map[string]interface{}) interface{} {
	var (
		response Response
	)

	response = make(map[string]string)
	response["value"] = "PONG"

	return response
}

func Print(node string, duration time.Duration, response interface{}) {
	var (
		value string
	)

	value = response.(map[string]interface{})["value"].(string)
	fmt.Println(value + " response from " + node + " in " + duration.String())
}

func Summary() {
}
