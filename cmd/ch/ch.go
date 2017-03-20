package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/r3boot/collective-herder/lib/amqp"
	"github.com/r3boot/collective-herder/lib/utils"
	"github.com/r3boot/collective-herder/plugins"
)

const (
	D_DEBUG bool = false
)

var (
	Amqp   *amqp.AmqpClient
	Agents *plugins.Agents
	Log    utils.Log
	debug  = flag.Bool("d", D_DEBUG, "Enable debug output")
)

func Usage() {
	var (
		myName      string
		plugin      string
		description string
	)

	myName = path.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, "Usage: %s [flags] [plugin] [opts]\n", myName)
	fmt.Fprintf(os.Stderr, "\nAvailable plugins:\n")
	for plugin, description = range Agents.Meta {
		fmt.Fprintf(os.Stderr, "%-20s%s\n", plugin, description)
	}
	fmt.Fprintf(os.Stderr, "\nAvailable flags:\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func AgentSelector() {
	var (
		wantedPlugin   string
		selectedPlugin string
		plugin         string
	)

	if len(flag.Args()) == 0 {
		Usage()
	}
	wantedPlugin = flag.Args()[0]
	for plugin, _ = range Agents.Meta {
		if plugin == wantedPlugin {
			selectedPlugin = wantedPlugin
		}
	}

	if selectedPlugin == "" {
		fmt.Fprintf(os.Stderr, "ERROR: Unknown plugin: "+wantedPlugin+"\n\n")
		Usage()
	}

	Amqp.SendToCollective(selectedPlugin, nil, nil)

}

func main() {
	var (
		err error
	)

	Agents = plugins.NewAgents()

	flag.Usage = Usage
	flag.Parse()

	Log = utils.Log{
		UseDebug:     *debug,
		UseVerbose:   *debug,
		UseTimestamp: true,
	}

	amqp.Setup(Log, amqp.AmqpConfig{
		Address:        "rabbitmq.service.local:5672",
		Username:       "ch",
		Password:       "ch",
		SendExchange:   "ch-send",
		RecvExchange:   "ch-recv",
		RequestTimeout: "1s",
	})

	if Amqp, err = amqp.NewAmqpClient(); err != nil {
		Log.Error(err)
		os.Exit(1)
	}

	if err = Amqp.ConfigureAsAgent(Agents); err != nil {
		Log.Error(err)
		os.Exit(1)
	}

	Log.Debug("Amqp client initialized")

	AgentSelector()
	// Amqp.SendToCollective("ping", nil, nil)

}
