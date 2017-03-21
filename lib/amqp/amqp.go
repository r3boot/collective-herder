package amqp

import (
	"errors"

	"github.com/r3boot/collective-herder/lib/facts"
	"github.com/r3boot/collective-herder/lib/utils"
	"github.com/r3boot/collective-herder/plugins"
)

var (
	Config  AmqpConfig
	Log     utils.Log
	Agents  *plugins.Agents
	Facts   *facts.Facts
	Servers *plugins.Servers
)

func Setup(l utils.Log, c AmqpConfig) error {
	Log = l
	Config = c

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
