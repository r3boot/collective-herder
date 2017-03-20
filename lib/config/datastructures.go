package config

import (
	"github.com/r3boot/collective-herder/lib/amqp"
)

type Config struct {
	Amqp amqp.AmqpConfig `yaml:"amqp"`
}
