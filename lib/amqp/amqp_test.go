package amqp

import (
	"testing"

	"github.com/r3boot/collective-herder/lib/utils"
)

func TestSetup(t *testing.T) {
	var (
		l   utils.Log
		c   AmqpConfig
		err error
	)

	l = utils.Log{}
	c = AmqpConfig{}

	if err = Setup(l, c); err != nil {
		t.Fatalf("Setup returned error: " + err.Error())
	}
}

func runNewAmqpClientWith(t *testing.T, address string) {
	var (
		l      utils.Log
		cfg    AmqpConfig
		client *AmqpClient
		err    error
	)

	l = utils.Log{
		UseDebug: true,
	}

	cfg = AmqpConfig{
		Address:      address,
		Username:     "test",
		Password:     "test",
		SendExchange: "test-send",
		RecvExchange: "test-recv",
	}

	Setup(l, cfg)

	client, err = NewAmqpClient()

	switch address {
	case AMQP_TEST_ADDRESS1:
		{
			if err != nil {
				t.Fatalf("AmqpClient(AMQP_TEST_ADDRESS1): err != nil")
			}

			if client.connection != nil {
				t.Fatalf("AmqpClient: got initialized??")
			}
		}
	case AMQP_TEST_ADDRESS2:
		{
			if err == nil {
				t.Fatalf("AmqpClient(AMQP_TEST_ADDRESS2): err == nil")
			}
		}
	}

}

func TestNewAmqpClient(t *testing.T) {
	runNewAmqpClientWith(t, AMQP_TEST_ADDRESS1)
	runNewAmqpClientWith(t, AMQP_TEST_ADDRESS2)
}
