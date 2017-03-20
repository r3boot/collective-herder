package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/r3boot/collective-herder/lib/amqp"
	"github.com/r3boot/collective-herder/lib/config"
	"github.com/r3boot/collective-herder/lib/utils"
	"github.com/r3boot/collective-herder/plugins"
)

const (
	D_DEBUG   bool   = false
	D_CFGFILE string = "/etc/ch.yml"
)

type FactFlags map[string]interface{}

var (
	Amqp      *amqp.AmqpClient
	Agents    *plugins.Agents
	Config    config.Config
	Log       utils.Log
	factFlags FactFlags
	debug     = flag.Bool("d", D_DEBUG, "Enable debug output")
	cfgFile   = flag.String("f", D_CFGFILE, "Path to configuration file")
)

func (f FactFlags) String() string {
	var (
		key      string
		value    interface{}
		response string
	)

	for key, value = range f {
		if response == "" {
			response = key + "=" + value.(string)
		} else {
			response += "," + key + "=" + value.(string)
		}
	}

	return response
}

func (f FactFlags) Set(flag_s string) error {
	var (
		tokens []string
		err    error
	)

	tokens = strings.Split(flag_s, "=")
	if len(tokens) < 2 {
		err = errors.New("Facts needs to be key=value")
		return err
	}

	f[tokens[0]] = strings.Join(tokens[1:], "=")

	return nil
}

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

	if flag.NArg() == 0 {
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

	Amqp.SendToCollective(selectedPlugin, factFlags, nil)

}

func main() {
	var (
		err error
	)

	Agents = plugins.NewAgents()

	flag.Usage = Usage
	factFlags = make(map[string]interface{})
	flag.Var(factFlags, "wf", "Fact to use in key=value form")
	flag.Parse()

	Log = utils.Log{
		UseDebug:     *debug,
		UseVerbose:   *debug,
		UseTimestamp: true,
	}

	if Config, err = config.ReadFile(*cfgFile); err != nil {
		Log.Error(err)
		os.Exit(1)
	}

	amqp.Setup(Log, Config.Amqp)

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
}
