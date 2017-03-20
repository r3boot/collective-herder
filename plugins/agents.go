package plugins

import (
	"strconv"
	"time"

	"github.com/r3boot/collective-herder/lib/utils"
	"github.com/r3boot/collective-herder/plugins/ping"
)

type Result struct {
	Node     string
	Response interface{}
	Duration time.Duration
}

type ResultSet struct {
	StartTime time.Time
	Plugin    string
	Data      map[string]Result
}

var (
	Results map[string]ResultSet
)

func NewAgents(l utils.Logger) *Agents {
	var (
		p *Agents
	)
	Log = l

	p = &Agents{}
	p.LoadAllAgents()

	Results = make(map[string]ResultSet)

	return p
}

func (p *Agents) LoadAllAgents() {
	p.printFunc = make(map[string]func(time.Time, interface{}))
	p.summaryFunc = make(map[string]func())

	p.printFunc[ping.NAME] = ping.Print
	p.summaryFunc[ping.NAME] = ping.Summary
}

func (p *Agents) NumAgentsAsString() string {
	return strconv.Itoa(len(p.printFunc))
}

func (p *Agents) HasAgent(name string) bool {
	var (
		key string
	)

	for key, _ = range p.printFunc {
		if key == name {
			return true
		}
	}

	return false
}

func (p *Agents) Print(plugin, uuid string, startTime time.Time, response interface{}) {
	var (
		node           string
		hostUuid       string
		responseResult interface{}
	)

	node = response.(Response).Node
	hostUuid = response.(Response).HostUuid
	responseResult = response.(Response).Result

	switch plugin {
	case ping.NAME:
		{
			result := ping.Result{
				Node:     node,
				Uuid:     hostUuid,
				Duration: time.Since(startTime),
				Response: responseResult.(map[string]interface{})["value"].(string),
			}
			p.printFunc[plugin](startTime, result)
		}
	default:
		{
			Log.Warn("Print: Unknown plugin: " + plugin)
		}
	}
}

func (p *Agents) Summary(plugin string) {
	switch plugin {
	case ping.NAME:
		{
			p.summaryFunc[plugin]()
		}
	default:
		{
			Log.Warn("Summary: Unknown plugin: " + plugin)
		}
	}
}
