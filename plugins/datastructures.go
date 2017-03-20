package plugins

import (
	"time"
)

type Agents struct {
	printFunc   map[string]func(string, time.Duration, interface{})
	summaryFunc map[string]func()
}

type Servers struct {
	runFunc map[string]func(map[string]interface{}) interface{}
}

type Request struct {
	Uuid    string
	MsgType string
	Facts   map[string]interface{}
	Opts    map[string]interface{}
}

type Response struct {
	Uuid     string
	Node     string
	HostUuid string
	Result   interface{}
}
