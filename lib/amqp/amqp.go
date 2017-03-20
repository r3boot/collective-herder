package amqp

import (
	"errors"
	"os"

	"github.com/r3boot/collective-herder/lib/utils"
	"github.com/r3boot/collective-herder/plugins"
)

var (
	Config   AmqpConfig
	Log      utils.Logger
	Agents   *plugins.Agents
	Servers  *plugins.Servers
	Hostname string
	HostUuid string
)

func Setup(l utils.Logger, c AmqpConfig) error {
	var (
		err error
	)

	Log = l
	Config = c

	if Hostname, err = os.Hostname(); err != nil {
		err = errors.New("amqp.Setup: Failed to get hostname: " + err.Error())
		return err
	}
	Log.Debug("amqp.Setup: My hostname is " + Hostname)

	HostUuid = utils.Uuidgen()
	Log.Debug("amqp.Setup: My uuid is " + HostUuid)

	return nil
}

func NewAmqpClient() (*AmqpClient, error) {
	var (
		amqp *AmqpClient
		err  error
	)

	amqp = &AmqpClient{}
	if err = amqp.Connect(); err != nil {
		err = errors.New("NewAmqpClient: Failed to connect: " + err.Error())
		return nil, err
	}

	return amqp, nil
}
