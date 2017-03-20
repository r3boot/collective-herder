package amqp

import (
	"encoding/json"
	"errors"
	"github.com/streadway/amqp"
	"os"
	"time"

	"github.com/r3boot/collective-herder/plugins"
)

func (a *AmqpClient) SendToCollective(plugin string, facts map[string]interface{}, opts map[string]interface{}) {
	var (
		msg     *plugins.Request
		err     error
		t_start time.Time
	)

	if !Agents.HasAgent(plugin) {
		err = errors.New("AmqpClient.SendToCollective: No such plugin: " + plugin)
		Log.Warn(err)
		return
	}

	t_start = time.Now()
	msg = plugins.NewRequest(plugin, facts, opts)
	Log.Debug("SendToCollective: Sending " + plugin + " message to collective")
	a.Send(msg)
	a.ResponseHandler(plugin, msg.Uuid, t_start)
}

func (a *AmqpClient) RequestHandler() error {
	var (
		data       []byte
		msgChannel <-chan amqp.Delivery
		result     interface{}
		request    plugins.Request
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
				Log.Debug("Received message")

				if err = json.Unmarshal(data, &request); err != nil {
					err = errors.New("AmqpClient.RequestHandler: Failed to decode message as JSON: " + err.Error())
					Log.Warn(err)
					continue
				}

				if !Servers.HasServer(request.MsgType) {
					err = errors.New("AmqpClient.RequestHandler: Unknown message type: " + request.MsgType)
					Log.Warn(err)
					continue
				}

				if !Facts.HasFact(request.Facts) {
					err = errors.New("AmqpClient.RequestHandler: No fact matches a locally available fact")
					Log.Warn(err)
					continue
				}

				result = Servers.RunServer(request.MsgType, request.Opts)
				response = plugins.NewResponse(request.Uuid, result)
				a.Response(response)
			}
		}
	}

	return nil
}

func (a *AmqpClient) ResponseHandler(plugin string, uuid string, startTime time.Time) {
	var (
		data       []byte
		response   plugins.Response
		msgChannel <-chan amqp.Delivery
		timeout    time.Duration
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

	if timeout, err = time.ParseDuration(Config.RequestTimeout); err != nil {
		Log.Error(err)
		os.Exit(1)
	}

	timeoutTimer := make(chan bool, 1)
	go func() {
		time.Sleep(timeout)
		timeoutTimer <- true
	}()

	stop_loop := false

	for {
		if stop_loop {
			break
		}

		select {
		case msgReceived := <-msgChannel:
			{
				data = msgReceived.Body
				if err = json.Unmarshal(data, &response); err != nil {
					err = errors.New("AmqpClient.RequestHandler: Failed to decode message as JSON: " + err.Error())
					Log.Warn(err)
					continue
				}
				Agents.Print(plugin, uuid, startTime, response)
			}
		case <-timeoutTimer:
			{
				stop_loop = true
				break
			}
		}
	}

	Agents.Summary(plugin)
}

func (a *AmqpClient) ProbeResponseHandler(uuid string) []string {
	var (
		data           []byte
		listeningUuids []string
		response       plugins.Response
		msgChannel     <-chan amqp.Delivery
		timeout        time.Duration
		err            error
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

	if timeout, err = time.ParseDuration(Config.RequestTimeout); err != nil {
		Log.Error(err)
		os.Exit(1)
	}

	timeoutTimer := make(chan bool, 1)
	go func() {
		time.Sleep(timeout)
		timeoutTimer <- true
	}()

	stop_loop := false

	for {
		if stop_loop {
			break
		}

		select {
		case msgReceived := <-msgChannel:
			{
				data = msgReceived.Body
				if err = json.Unmarshal(data, &response); err != nil {
					err = errors.New("AmqpClient.RequestHandler: Failed to decode message as JSON: " + err.Error())
					Log.Warn(err)
					continue
				}
				listeningUuids = append(listeningUuids, response.HostUuid)
			}
		case <-timeoutTimer:
			{
				stop_loop = true
				break
			}
		}
	}

	return listeningUuids
}
