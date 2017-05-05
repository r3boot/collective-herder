package amqp

import (
	"os"
	"testing"

	"github.com/r3boot/collective-herder/lib/utils"
	"github.com/r3boot/collective-herder/plugins"
)

const (
	AMQP_TEST_ADDRESS  string = "rabbitmq.service.local:5672"
	AMQP_TEST_USERNAME string = "ch"
	AMQP_TEST_PASSWORD string = "ch"
	AMQP_TEST_SENDEXCH string = "ch-send"
	AMQP_TEST_RECVEXCH string = "ch-recv"
)

func envOrConst(envvar, constValue string) string {
	var (
		value string
	)

	if value = os.Getenv(envvar); value != "" {
		return value
	}

	return constValue
}

func setupAmqpTestClient(t *testing.T, cfg AmqpConfig) *AmqpClient {
	var (
		l      utils.Log
		client *AmqpClient
		err    error
	)

	l = utils.Log{UseDebug: true}

	if err = Setup(l, cfg); err != nil {
		t.Fatalf("Failed to configure amqp library: " + err.Error())
	}

	if client, err = NewAmqpClient(); err != nil {
		t.Fatalf("Failed to initialize AmqpClient: " + err.Error())
	}

	return client
}

func TestValidConnection(t *testing.T) {
	var (
		client *AmqpClient
	)

	client = setupAmqpTestClient(t, AmqpConfig{
		Address:      envOrConst("AMQP_TEST_ADDRESS", AMQP_TEST_ADDRESS),
		Username:     envOrConst("AMQP_TEST_USERNAME", AMQP_TEST_USERNAME),
		Password:     envOrConst("AMQP_TEST_PASSWORD", AMQP_TEST_PASSWORD),
		SendExchange: envOrConst("AMQP_TEST_SENDEXCH", AMQP_TEST_SENDEXCH),
		RecvExchange: envOrConst("AMQP_TEST_SENDEXCH", AMQP_TEST_SENDEXCH),
	})

	if client == nil {
		t.Fatalf("client == nil")
	}
}

func TestConfigureAsAgent(t *testing.T) {
	var (
		client *AmqpClient
		err    error
	)

	// Valid client
	client = setupAmqpTestClient(t, AmqpConfig{
		Address:      envOrConst("AMQP_TEST_ADDRESS", AMQP_TEST_ADDRESS),
		Username:     envOrConst("AMQP_TEST_USERNAME", AMQP_TEST_USERNAME),
		Password:     envOrConst("AMQP_TEST_PASSWORD", AMQP_TEST_PASSWORD),
		SendExchange: envOrConst("AMQP_TEST_SENDEXCH", AMQP_TEST_SENDEXCH),
		RecvExchange: envOrConst("AMQP_TEST_SENDEXCH", AMQP_TEST_SENDEXCH),
	})

	if client == nil {
		t.Fatalf("Failed to initialize AmqpClient")
	}

	if err = client.ConfigureAsAgent(nil); err != nil {
		t.Fatalf("Failed to configure Amqpclient as agent: " + err.Error())
	}
}

func TestConfigureAsServer(t *testing.T) {
	var (
		client *AmqpClient
		err    error
	)

	// Valid client
	client = setupAmqpTestClient(t, AmqpConfig{
		Address:      envOrConst("AMQP_TEST_ADDRESS", AMQP_TEST_ADDRESS),
		Username:     envOrConst("AMQP_TEST_USERNAME", AMQP_TEST_USERNAME),
		Password:     envOrConst("AMQP_TEST_PASSWORD", AMQP_TEST_PASSWORD),
		SendExchange: envOrConst("AMQP_TEST_SENDEXCH", AMQP_TEST_SENDEXCH),
		RecvExchange: envOrConst("AMQP_TEST_SENDEXCH", AMQP_TEST_SENDEXCH),
	})

	if client == nil {
		t.Fatalf("Failed to initialize AmqpClient")
	}

	if err = client.ConfigureAsServer(nil, nil); err != nil {
		t.Fatalf("Failed to configure Amqpclient as agent: " + err.Error())
	}
}

func TestClosedConnection(t *testing.T) {
	var (
		client *AmqpClient
		err    error
	)

	// Invalid send exchange
	client = setupAmqpTestClient(t, AmqpConfig{
		Address:      envOrConst("AMQP_TEST_ADDRESS", AMQP_TEST_ADDRESS),
		Username:     envOrConst("AMQP_TEST_USERNAME", AMQP_TEST_USERNAME),
		Password:     envOrConst("AMQP_TEST_PASSWORD", AMQP_TEST_PASSWORD),
		SendExchange: "nonexisting-exchange",
		RecvExchange: envOrConst("AMQP_TEST_SENDEXCH", AMQP_TEST_SENDEXCH),
	})

	if client == nil {
		t.Fatalf("Failed to initialize AmqpClient")
	}

	client.connection.Close()

	err = client.ConfigureAsAgent(nil)
	if err == nil {
		t.Fatalf("Able to initialize AmqpClient with nonexisting send exchange")
	}
}

func TestSend(t *testing.T) {
	var (
		client *AmqpClient
		err    error
		msg    *plugins.Request
	)

	// Valid client
	client = setupAmqpTestClient(t, AmqpConfig{
		Address:      envOrConst("AMQP_TEST_ADDRESS", AMQP_TEST_ADDRESS),
		Username:     envOrConst("AMQP_TEST_USERNAME", AMQP_TEST_USERNAME),
		Password:     envOrConst("AMQP_TEST_PASSWORD", AMQP_TEST_PASSWORD),
		SendExchange: envOrConst("AMQP_TEST_SENDEXCH", AMQP_TEST_SENDEXCH),
		RecvExchange: envOrConst("AMQP_TEST_SENDEXCH", AMQP_TEST_SENDEXCH),
	})

	if client == nil {
		t.Fatalf("Failed to initialize AmqpClient")
	}

	if err = client.ConfigureAsAgent(nil); err != nil {
		t.Fatalf("Failed to configure AmqpClient as agent: " + err.Error())
	}

	msg = &plugins.Request{
		Uuid:    utils.Uuidgen(),
		MsgType: "ping",
		Facts:   make(map[string]interface{}),
		Opts:    make(map[string]interface{}),
	}

	if err = client.Send(msg); err != nil {
		t.Fatalf("Failed to send valid message: " + err.Error())
	}

}

func TestResponse(t *testing.T) {
	var (
		client *AmqpClient
		err    error
		msg    *plugins.Response
	)

	// Valid client
	client = setupAmqpTestClient(t, AmqpConfig{
		Address:      envOrConst("AMQP_TEST_ADDRESS", AMQP_TEST_ADDRESS),
		Username:     envOrConst("AMQP_TEST_USERNAME", AMQP_TEST_USERNAME),
		Password:     envOrConst("AMQP_TEST_PASSWORD", AMQP_TEST_PASSWORD),
		SendExchange: envOrConst("AMQP_TEST_SENDEXCH", AMQP_TEST_SENDEXCH),
		RecvExchange: envOrConst("AMQP_TEST_SENDEXCH", AMQP_TEST_SENDEXCH),
	})

	if client == nil {
		t.Fatalf("Failed to initialize AmqpClient")
	}

	if err = client.ConfigureAsAgent(nil); err != nil {
		t.Fatalf("Failed to configure AmqpClient as agent: " + err.Error())
	}

	msg = &plugins.Response{
		Uuid:     utils.Uuidgen(),
		Node:     "nonexisting",
		HostUuid: utils.Uuidgen(),
		Result:   make(map[string]interface{}),
	}

	if err = client.Response(msg); err != nil {
		t.Fatalf("Failed to send valid message: " + err.Error())
	}
}
