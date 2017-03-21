package plugins

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/r3boot/collective-herder/plugins/facts"
	"github.com/r3boot/collective-herder/plugins/ping"
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
	p.LoadAllAgents()

	return p
}

func (p *Agents) LoadAllAgents() {
	p.argsFunc = make(map[string]func([]string) map[string]interface{})
	p.preRunFunc = make(map[string]func(map[string]interface{}))
	p.printFunc = make(map[string]func(time.Time, interface{}))
	p.summaryFunc = make(map[string]func())

	// Ping agent
	p.argsFunc[ping.NAME] = ping.ParseArgs
	p.preRunFunc[ping.NAME] = ping.PreRun
	p.printFunc[ping.NAME] = ping.Print
	p.summaryFunc[ping.NAME] = ping.Summary
	p.Meta[ping.NAME] = ping.DESCRIPTION

	// Facts agent
	p.argsFunc[facts.NAME] = facts.ParseArgs
	p.preRunFunc[facts.NAME] = facts.PreRun
	p.printFunc[facts.NAME] = facts.Print
	p.summaryFunc[facts.NAME] = facts.Summary
	p.Meta[facts.NAME] = facts.DESCRIPTION
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
	case facts.NAME:
		{
			result := facts.Result{
				Node:     node,
				Uuid:     hostUuid,
				Duration: time.Since(startTime),
				Response: responseResult.(map[string]interface{}),
			}
			p.printFunc[plugin](startTime, result)
		}
	default:
		{
			fmt.Fprintf(os.Stderr, "Print: Unknown plugin: "+plugin)
		}
	}
}

func (p *Agents) Summary(plugin string) {
	switch plugin {
	case ping.NAME:
		{
			p.summaryFunc[plugin]()
		}
	case facts.NAME:
		{
			p.summaryFunc[plugin]()
		}
	default:
		{
			fmt.Fprintf(os.Stderr, "Summary: Unknown plugin: "+plugin)
		}
	}
}
