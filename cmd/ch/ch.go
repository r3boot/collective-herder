package main

import (
	"flag"
	"os"

	"github.com/r3boot/collective-herder/lib/amqp"
	"github.com/r3boot/collective-herder/lib/utils"
	"github.com/r3boot/collective-herder/plugins"
)

const (
	D_DEBUG bool = false
)

var (
	Amqp  *amqp.AmqpClient
	a     *plugins.Agents
	Log   utils.Logger
	debug = flag.Bool("d", D_DEBUG, "Enable debug output")
)

func main() {
	var (
		err error
	)

	flag.Parse()

	Log = utils.NewLogger(*debug)

	a = plugins.NewAgents(Log)

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

	if err = Amqp.ConfigureAsAgent(a); err != nil {
		Log.Error(err)
		os.Exit(1)
	}

	Log.Debug("Amqp client initialized")

	Amqp.SendToCollective("ping", nil, nil)
}
