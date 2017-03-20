package ping

import (
	"fmt"
	"time"

	"github.com/r3boot/collective-herder/lib/utils"
)

const (
	NAME string = "ping"
)

type Response map[string]string

var (
	Log utils.Logger
)

func RegisterPrint(l utils.Logger) interface{} {
	Log = l
	return Print
}

func RegisterSummary() interface{} {
	return Summary
}

func RegisterServer(l utils.Logger) interface{} {
	Log = l
	return Run
}

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
