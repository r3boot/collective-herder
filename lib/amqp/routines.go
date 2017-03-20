package amqp

import (
	"encoding/json"
	"errors"
	"github.com/streadway/amqp"
	"os"
	"time"

	"github.com/r3boot/collective-herder/plugins"
)

func (a *AmqpClient) RequestHandler() error {
	var (
		data       []byte
		msg        plugins.Request
		msgChannel <-chan amqp.Delivery
		result     interface{}
		response   *plugins.Response
		err        error
	)

	msgChannel, err = a.sendChannel.Consume(
		a.sendQueue.Name, // Queue to use
		"",               // Name of this consumer
		true,             // Auto-acknowledge event
		false,            // Non-exclusive consumer
		false,            // No-local consumer
		false,            // No-wait consumer
		nil,              // No arguments
	)
	if err != nil {
		err = errors.New("AmqpClient.RequestHandler: Failed to start consumer: " + err.Error())
		return err
	}

	for {
		select {
		case msgReceived := <-msgChannel:
			{
				data = msgReceived.Body
				if err = json.Unmarshal(data, &msg); err != nil {
					err = errors.New("AmqpClient.RequestHandler: Failed to decode message as JSON: " + err.Error())
					Log.Warn(err)
					continue
				}
				if !Servers.HasServer(msg.MsgType) {
					err = errors.New("AmqpClient.RequestHandler: Unknown message type: " + msg.MsgType)
					Log.Warn(err)
					continue
				}
				result = Servers.RunServer(msg.MsgType, msg.Opts)
				response = plugins.NewResponse(msg.Uuid, result)
				a.Response(response)
			}
		}
	}

	return nil
}

func (a *AmqpClient) ResponseHandler(uuid string, t_start time.Time) {
	var (
		data       []byte
		response   plugins.Response
		msgChannel <-chan amqp.Delivery
		err        error
	)

	msgChannel, err = a.recvChannel.Consume(
		a.recvQueue.Name, // Queue to use
		"",               // Name of this consumer
		true,             // Auto-acknowledge event
		false,            // Non-exclusive consumer
		false,            // No-local consumer
		false,            // No-wait consumer
		nil,              // No arguments
	)
	if err != nil {
		err = errors.New("AmqpClient.ResponseHandler: Failed to start consumer: " + err.Error())
		Log.Error(err)
		os.Exit(1)
	}

	for {
		select {
		case msgReceived := <-msgChannel:
			{
				data = msgReceived.Body
				if err = json.Unmarshal(data, &response); err != nil {
					err = errors.New("AmqpClient.RequestHandler: Failed to decode message as JSON: " + err.Error())
					Log.Warn(err)
					continue
				}
				Agents.Print(uuid, response)
			}
		}
	}
}
