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
	p.printFunc = make(map[string]func(string, time.Duration, interface{}))
	p.summaryFunc = make(map[string]func())

	p.printFunc[ping.NAME] = ping.RegisterPrint(Log).(func(string, time.Duration, interface{}))
	p.summaryFunc[ping.NAME] = ping.RegisterSummary().(func())
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

func (p *Agents) HasResultSet(uuid string) bool {
	var (
		key string
	)

	for key, _ = range Results {
		if key == uuid {
			return true
		}
	}

	return false
}

func (p *Agents) StartResults(plugin, uuid string, t_start time.Time) {
	Results[uuid] = ResultSet{
		StartTime: t_start,
		Plugin:    plugin,
		Data:      make(map[string]Result),
	}
}

func (p *Agents) CleanResults(uuid string) {
	delete(Results, uuid)
}

func (p *Agents) Print(uuid string, response interface{}) {
	var (
		node           string
		hostUuid       string
		responseResult interface{}
		r              Result
	)

	if !p.HasResultSet(uuid) {
		Log.Warn("Agents.Print: No result set for " + uuid)
		return
	}

	node = response.(Response).Node
	hostUuid = response.(Response).HostUuid
	responseResult = response.(Response).Result

	r = Result{
		Node:     node,
		Duration: time.Since(Results[uuid].StartTime),
		Response: response,
	}

	Results[uuid].Data[hostUuid] = r
	p.printFunc[Results[uuid].Plugin](node, r.Duration, responseResult)
}

func (p *Agents) Summary(plugin string, uuid string) {
	if !p.HasAgent(plugin) {
		Log.Warn("Agents.Summary: No such plugin: " + plugin)
		return
	}
}
