package amqp

import (
	"errors"
	"github.com/streadway/amqp"

	"github.com/r3boot/collective-herder/lib/facts"
	"github.com/r3boot/collective-herder/lib/utils"
	"github.com/r3boot/collective-herder/plugins"
)

func (a *AmqpClient) Connect() error {
	var (
		url   string
		url_d string
		err   error
	)

	url = "amqp://" + Config.Username + ":" + Config.Password + "@" + Config.Address
	url_d = "amqp://" + Config.Username + ":***@" + Config.Address

	// Try to connect to AMQP
	if a.connection, err = amqp.Dial(url); err != nil {
		err = errors.New("AmqpClient.Connect: amqp.Dial failed: " + err.Error())
		a = nil
		return nil
	}
	Log.Debug("AmqpClient.Connect: Connected to " + url_d)

	return nil
}

/*
 * Agent configuration
 */
func (a *AmqpClient) ConfigureAsAgent(agents *plugins.Agents) error {
	var (
		err error
	)

	Agents = agents

	/*
	 * Send channel
	 */
	if a.sendChannel, err = a.connection.Channel(); err != nil {
		err = errors.New("Amqp.ConfigureAsAgent: Failed to setup a.sendChannel: " + err.Error())
		return err
	}

	// Declare the fanout exchange on the newly created channel
	err = a.sendChannel.ExchangeDeclare(
		Config.SendExchange, // Name of the exchange
		"fanout",            // Type of exchange
		true,                // Durable queue
		false,               // Not auto-deleted
		false,               // Not an internal queue
		false,               // No-wait queue
		nil,                 // Arguments
	)
	if err != nil {
		err = errors.New("Amqp.ConfigureAsAgent: Failed to declare an exchange: " + err.Error())
		return err
	}

	// Declare the private queue
	a.sendQueue, err = a.sendChannel.QueueDeclare(
		Config.SendExchange, // Name
		true,                // Durable queue
		false,               // Dont delete when unused
		false,               // Exclusive queue
		false,               // No-wait queue
		nil,                 // No arguments
	)
	if err != nil {
		err = errors.New("Amqp.ConfigureAsAgent: Failed to declare queue: " + err.Error())
		return err
	}

	/*
	 * Receive channel
	 */
	if a.recvChannel, err = a.connection.Channel(); err != nil {
		err = errors.New("Amqp.ConfigureAsAgent: Failed to setup a.recvChannel: " + err.Error())
		return err
	}

	// Declare the fanout exchange on the newly created channel
	err = a.recvChannel.ExchangeDeclare(
		Config.RecvExchange, // Name of the exchange
		"fanout",            // Type of exchange
		true,                // Durable queue
		false,               // Not auto-deleted
		false,               // Not an internal queue
		false,               // No-wait queue
		nil,                 // Arguments
	)
	if err != nil {
		err = errors.New("Amqp.ConfigureAsAgent: Failed to declare an exchange: " + err.Error())
		return err
	}

	// Declare the private queue
	a.recvQueue, err = a.recvChannel.QueueDeclare(
		Config.RecvExchange+"."+utils.Uuidgen(), // Name
		false, // Durable queue
		false, // Dont delete when unused
		true,  // Exclusive queue
		false, // No-wait queue
		nil,   // No arguments
	)
	if err != nil {
		err = errors.New("Amqp.ConfigureAsAgent: Failed to declare queue: " + err.Error())
		return err
	}

	// Bind to the queue
	err = a.recvChannel.QueueBind(
		a.recvQueue.Name,    // Name of queue
		"",                  // Routing key
		Config.RecvExchange, // Exchange
		false,               // No-wait
		nil,                 // Args
	)
	if err != nil {
		err = errors.New("Amqp.ConfigureAsAgent: Failed to bind to queue: " + err.Error())
		return err
	}
	return nil
}

/*
 * Server configuration
 */
func (a *AmqpClient) ConfigureAsServer(p *plugins.Servers, f *facts.Facts) error {
	var (
		err error
	)

	Servers = p
	Facts = f

	if a.sendChannel, err = a.connection.Channel(); err != nil {
		err = errors.New("AmqpClient.Connect: Failed to setup a.sendChannel: " + err.Error())
		return err
	}

	// Declare the fanout exchange on the newly created channel
	err = a.sendChannel.ExchangeDeclare(
		Config.SendExchange, // Name of the exchange
		"fanout",            // Type of exchange
		true,                // Durable queue
		false,               // Not auto-deleted
		false,               // Not an internal queue
		false,               // No-wait queue
		nil,                 // Arguments
	)
	if err != nil {
		err = errors.New("AmqpClient.Connect: Failed to declare an exchange: " + err.Error())
		return err
	}

	// Declare the private queue
	a.sendQueue, err = a.sendChannel.QueueDeclare(
		Config.SendExchange+"."+utils.Uuidgen(), // Name
		false, // Durable queue
		false, // Dont delete when unused
		true,  // Exclusive queue
		false, // No-wait queue
		nil,   // No arguments
	)
	if err != nil {
		err = errors.New("AmqpClient.Connect: Failed to declare queue: " + err.Error())
		return err
	}

	// Bind to the queue
	err = a.sendChannel.QueueBind(
		a.sendQueue.Name,    // Name of queue
		"",                  // Routing key
		Config.SendExchange, // Exchange
		false,               // No-wait
		nil,                 // Args
	)
	if err != nil {
		err = errors.New("AmqpClient.Connect: Failed to bind to queue: " + err.Error())
		return err
	}

	/*
	 * Receive channel
	 */
	if a.recvChannel, err = a.connection.Channel(); err != nil {
		err = errors.New("Amqp.ConfigureAsServer: Failed to setup a.recvChannel: " + err.Error())
		return err
	}

	// Declare the fanout exchange on the newly created channel
	err = a.recvChannel.ExchangeDeclare(
		Config.RecvExchange, // Name of the exchange
		"fanout",            // Type of exchange
		true,                // Durable queue
		false,               // Not auto-deleted
		false,               // Not an internal queue
		false,               // No-wait queue
		nil,                 // Arguments
	)
	if err != nil {
		err = errors.New("Amqp.ConfigureAsServer: Failed to declare an exchange: " + err.Error())
		return err
	}

	// Declare the private queue
	a.recvQueue, err = a.recvChannel.QueueDeclare(
		Config.RecvExchange, // Name
		true,                // Durable queue
		false,               // Dont delete when unused
		false,               // Exclusive queue
		false,               // No-wait queue
		nil,                 // No arguments
	)
	if err != nil {
		err = errors.New("Amqp.ConfigureAsServer: Failed to declare queue: " + err.Error())
		return err
	}
	return nil
}

func (a *AmqpClient) Send(msg *plugins.Request) error {
	var (
		data []byte
		err  error
	)

	if data, err = msg.ToJSON(); err != nil {
		err = errors.New("AmqpClient.Send: Encode to JSON failed: " + err.Error())
		return err
	}

	err = a.sendChannel.Publish(
		Config.SendExchange, // Exchange to use
		"",                  // Key to use for routing
		false,               // Mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)

	return nil
}

func (a *AmqpClient) Response(msg *plugins.Response) error {
	var (
		data []byte
		err  error
	)

	if data, err = msg.ToJSON(); err != nil {
		err = errors.New("AmqpClient.Response: Encode to JSON failed: " + err.Error())
		return err
	}

	err = a.recvChannel.Publish(
		Config.RecvExchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)

	return nil
}
