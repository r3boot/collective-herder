package plugins
/*
 * WARNING: Modify at your own risk, this file is auto-generated ...
 */

import (
	"fmt"
	"os"
	"time"


	"github.com/r3boot/collective-herder/plugins/facts"

	"github.com/r3boot/collective-herder/plugins/ping"

	"github.com/r3boot/collective-herder/plugins/run"

)

func (p *Agents) LoadAllAgents() {

	// Glue code for facts plugin
	p.argsFunc[facts.NAME] = facts.ParseArgs
	p.preRunFunc[facts.NAME] = facts.PreRun
	p.printFunc[facts.NAME] = facts.Print
	p.summaryFunc[facts.NAME] = facts.Summary
	p.Meta[facts.NAME] = facts.DESCRIPTION


	// Glue code for ping plugin
	p.argsFunc[ping.NAME] = ping.ParseArgs
	p.preRunFunc[ping.NAME] = ping.PreRun
	p.printFunc[ping.NAME] = ping.Print
	p.summaryFunc[ping.NAME] = ping.Summary
	p.Meta[ping.NAME] = ping.DESCRIPTION


	// Glue code for run plugin
	p.argsFunc[run.NAME] = run.ParseArgs
	p.preRunFunc[run.NAME] = run.PreRun
	p.printFunc[run.NAME] = run.Print
	p.summaryFunc[run.NAME] = run.Summary
	p.Meta[run.NAME] = run.DESCRIPTION


}

func (p *Agents) Print(plugin, uuid string, startTime time.Time, response interface{}, opts map[string]interface{}) {
	var (
		node           string
		hostUuid       string
		responseResult interface{}
	)

	node = response.(Response).Node
	hostUuid = response.(Response).HostUuid
	responseResult = response.(Response).Result

	switch plugin {

	case facts.NAME:
		{
			result := facts.Result{
				Node: node,
				Uuid: hostUuid,
				Duration: time.Since(startTime),
				Response: responseResult.(map[string]interface{}),
			}
			p.printFunc[plugin](startTime, result, opts)
		}

	case ping.NAME:
		{
			result := ping.Result{
				Node: node,
				Uuid: hostUuid,
				Duration: time.Since(startTime),
				Response: responseResult.(map[string]interface{}),
			}
			p.printFunc[plugin](startTime, result, opts)
		}

	case run.NAME:
		{
			result := run.Result{
				Node: node,
				Uuid: hostUuid,
				Duration: time.Since(startTime),
				Response: responseResult.(map[string]interface{}),
			}
			p.printFunc[plugin](startTime, result, opts)
		}

	default:
		{
			fmt.Fprintf(os.Stderr, "Error: Print: unknown plugin: %s", plugin)
			os.Exit(1)
		}
	}
}

func (p *Agents) Summary(plugin string, opts map[string]interface{}) {
	switch plugin {

	case facts.NAME:
		{
			p.summaryFunc[plugin](opts)
		}

	case ping.NAME:
		{
			p.summaryFunc[plugin](opts)
		}

	case run.NAME:
		{
			p.summaryFunc[plugin](opts)
		}

	default:
		{
			fmt.Fprintf(os.Stderr, "Error: Summary: unknown plugin: %s", plugin)
			os.Exit(1)
		}
	}
}
