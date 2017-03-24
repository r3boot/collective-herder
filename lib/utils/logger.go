package utils

import (
	"fmt"
	"os"
	"time"
)

const (
	MSG_INFO    byte = 0x0
	MSG_WARNING byte = 0x1
	MSG_FATAL   byte = 0x2
	MSG_VERBOSE byte = 0x3
	MSG_DEBUG   byte = 0x4
)

var MSG_STRING = map[byte]string{
	MSG_INFO:    "INFO    ",
	MSG_WARNING: "WARNING ",
	MSG_FATAL:   "FATAL   ",
	MSG_VERBOSE: "VERBOSE ",
	MSG_DEBUG:   "DEBUG   ",
}

type Log struct {
	UseDebug        bool
	UseVerbose      bool
	UseTimestamp    bool
	TimestampFormat string
	TestFd          *os.File
}

func (l Log) Message(log_level byte, message ...interface{}) {
	var (
		msg string
		fd  *os.File
	)

	if l.UseTimestamp {
		if len(l.TimestampFormat) == 0 {
			l.TimestampFormat = time.RFC3339
		}
		timestamp := time.Now().Format(time.RFC3339)
		msg = timestamp + " " + MSG_STRING[log_level] + ":"
	} else {
		msg = MSG_STRING[log_level] + ":"
	}

	all := append([]interface{}{msg}, message...)

	if log_level == MSG_FATAL {
		fd = os.Stderr
	} else {
		fd = os.Stdout
	}

	if l.TestFd != nil {
		fd = l.TestFd
	}

	fmt.Fprintln(fd, all...)
}

func (l Log) Info(message ...interface{}) {
	l.Message(MSG_INFO, message...)
}

func (l Log) Warn(message ...interface{}) {
	l.Message(MSG_WARNING, message...)
}

func (l Log) Error(message ...interface{}) {
	l.Message(MSG_FATAL, message...)
	if l.TestFd == nil {
		os.Exit(1)
	}
}

func (l Log) Verbose(message ...interface{}) {
	if l.UseDebug || l.UseVerbose {
		l.Message(MSG_VERBOSE, message...)
	}
}

func (l Log) Debug(message ...interface{}) {
	if l.UseDebug {
		l.Message(MSG_DEBUG, message...)
	}
}
