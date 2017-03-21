package plugins

import (
	"time"
)

type Agents struct {
	argsFunc    map[string]func([]string) map[string]interface{}
	preRunFunc  map[string]func(map[string]interface{})
	printFunc   map[string]func(time.Time, interface{})
	summaryFunc map[string]func()
	Meta        map[string]string
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
