package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/r3boot/collective-herder/lib/amqp"
	"github.com/r3boot/collective-herder/lib/config"
	"github.com/r3boot/collective-herder/lib/facts"
	"github.com/r3boot/collective-herder/lib/utils"
	"github.com/r3boot/collective-herder/plugins"
)

const (
	D_DEBUG   bool   = false
	D_CFGFILE string = "/etc/chd.yml"
)

var (
	Amqp    *amqp.AmqpClient
	Log     utils.Log
	Config  config.Config
	f       *facts.Facts
	p       *plugins.Servers
	debug   = flag.Bool("d", D_DEBUG, "Enable debug output")
	cfgFile = flag.String("f", D_CFGFILE, "Path to configuration file")
)

func Usage() {
	var (
		myName string
	)

	myName = path.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, "Usage: %s [flags] [plugin] [opts]\n", myName)
	fmt.Fprintf(os.Stderr, "\nAvailable flags:\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	var (
		err error
	)

	flag.Usage = Usage
	flag.Parse()

	Log = utils.Log{
		UseDebug:     *debug,
		UseVerbose:   *debug,
		UseTimestamp: false,
	}

	if Config, err = config.ReadFile(*cfgFile); err != nil {
		Log.Error(err)
		os.Exit(1)
	}

	f = facts.NewFacts(Log)
	Log.Debug("Loaded " + f.NumFactsAsString() + " facts")

	p = plugins.NewServers(Log, f)
	Log.Debug("Loaded " + p.NumServersAsString() + " servers")

	amqp.Setup(Log, Config.Amqp)

	if Amqp, err = amqp.NewAmqpClient(); err != nil {
		Log.Error(err)
		os.Exit(1)
	}

	if err = Amqp.ConfigureAsServer(p, f); err != nil {
		Log.Error(err)
		os.Exit(1)
	}

	Log.Debug("Amqp client initialized")

	if err = Amqp.RequestHandler(); err != nil {
		Log.Error(err)
		os.Exit(1)
	}

}
