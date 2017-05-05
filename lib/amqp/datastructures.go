package amqp

import (
	"github.com/streadway/amqp"
)

const (
	AMQP_BUFSIZE       int    = 16
	CONTROL_BUFSIZE    int    = 1
	DONE_BUFSIZE       int    = 1
	AMQP_TEST_ADDRESS1 string = "ch-amqp_test"
	AMQP_TEST_ADDRESS2 string = "ch-amqp_test-nonexisting"
)

type AmqpClient struct {
	connection  *amqp.Connection
	sendChannel *amqp.Channel
	sendQueue   amqp.Queue
	recvChannel *amqp.Channel
	recvQueue   amqp.Queue
}

type AmqpConfig struct {
	Address        string `yaml:"address"`
	Username       string `yaml:"username"`
	Password       string `yaml:"password"`
	SendExchange   string `yaml:"send_exchange"`
	RecvExchange   string `yaml:"recv_exchange"`
	RequestTimeout string `yaml:"request_timeout"`
}
