package plugins

import (
	"strconv"
	"time"
)

var (
	AgentMeta map[string]string
)

func NewAgents() *Agents {
	var (
		p *Agents
	)

	p = &Agents{
		Meta: make(map[string]string),
	}
	p.argsFunc = make(map[string]func([]string) map[string]interface{})
	p.preRunFunc = make(map[string]func(map[string]interface{}))
	p.printFunc = make(map[string]func(time.Time, interface{}, map[string]interface{}))
	p.summaryFunc = make(map[string]func(map[string]interface{}))

	p.LoadAllAgents()

	return p
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

func (p *Agents) ParseArgs(plugin string, args []string) map[string]interface{} {
	return p.argsFunc[plugin](args)
}

func (p *Agents) PreRun(plugin string, opts map[string]interface{}) {
	p.preRunFunc[plugin](opts)
}
