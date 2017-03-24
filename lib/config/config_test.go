package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

var VALID_CONFIG string = `---
amqp:
  address: "localhost:5672"
  username: "ch"
  password: "ch"
  send_exchange: "ch-send"
  recv_exchange: "ch-recv"

commands_dir: "/etc/ch/commands.d"`

var INVALID_CONFIG string = `---
invalid configuration`

func createTestConfig(t *testing.T, content string) (string, error) {
	var (
		fd         *os.File
		numWritten int
		err        error
	)

	if fd, err = ioutil.TempFile("/tmp", "ch-config_test"); err != nil {
		t.Error("Failed to create tmpfile: " + err.Error())
		return "", err
	}
	defer fd.Close()

	if numWritten, err = fd.Write([]byte(content)); err != nil {
		t.Error("Failed to write to tmpfile: " + err.Error())
		return "", err
	}

	if numWritten != len(content) {
		t.Error("Failed to write to tmpfile: short write")
		return "", err
	}

	return fd.Name(), nil
}

func cleanTestConfig(fileName string) {
	os.Remove(fileName)
}

func TestReadFile(t *testing.T) {
	var (
		fileName string
		config   Config
		err      error
	)

	if fileName, err = createTestConfig(t, VALID_CONFIG); err != nil {
		return
	}

	if config, err = ReadFile(fileName); err != nil {
		t.Error("Failed to read config: " + err.Error())
	}

	// Parse amqp.address
	if reflect.TypeOf(config.Amqp.Address).String() != "string" {
		t.Error("type of config.Amqp.Address != string")
	}
	if config.Amqp.Address != "localhost:5672" {
		t.Error("config.Amqp.Address != localhost:5672")
	}

	// Parse amqp.username
	if reflect.TypeOf(config.Amqp.Username).String() != "string" {
		t.Error("type of config.Amqp.Username != string")
	}
	if config.Amqp.Username != "ch" {
		t.Error("config.Amqp.Username != ch")
	}

	// Parse amqp.password
	if reflect.TypeOf(config.Amqp.Password).String() != "string" {
		t.Error("type of config.Amqp.Password != string")
	}
	if config.Amqp.Password != "ch" {
		t.Error("config.Amqp.Password != ch")
	}

	// Parse amqp.send_exchange
	if reflect.TypeOf(config.Amqp.SendExchange).String() != "string" {
		t.Error("type of config.Amqp.SendExchange != string")
	}
	if config.Amqp.SendExchange != "ch-send" {
		t.Error("config.Amqp.SendExchange != ch-send")
	}

	// Parse amqp.recv_exchange
	if reflect.TypeOf(config.Amqp.RecvExchange).String() != "string" {
		t.Error("type of config.Amqp.RecvExchange != string")
	}
	if config.Amqp.RecvExchange != "ch-recv" {
		t.Error("config.Amqp.RecvExchange != ch-recv")
	}

	// Parse commands_dir
	if reflect.TypeOf(config.CommandsDir).String() != "string" {
		t.Error("type of config.CommandsDir != string")
	}
	if config.CommandsDir != "/etc/ch/commands.d" {
		t.Error("config.commandsDir != /etc/ch/commands.d")
	}

	cleanTestConfig(fileName)

	if _, err = ReadFile(""); err == nil {
		t.Error("Cannot parse non-existing file")
	}

	if fileName, err = createTestConfig(t, INVALID_CONFIG); err != nil {
		fmt.Println("here")
		return
	}

	if _, err = ReadFile(fileName); err == nil {
		t.Error("Parsed invalid yaml")
	}

	cleanTestConfig(fileName)
}
