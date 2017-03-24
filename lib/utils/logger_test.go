package utils

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"
)

var (
	TEST_MSGS_NOTS map[byte]string = map[byte]string{
		MSG_INFO:    "INFO    : test info",
		MSG_WARNING: "WARNING : test warning",
		MSG_FATAL:   "FATAL   : test fatal",
		MSG_VERBOSE: "VERBOSE : test verbose",
		MSG_DEBUG:   "DEBUG   : test debug",
	}
	reRFC3339TS = regexp.MustCompile("^(20[0-9]{2}-[01][0-9]-[0-2][0-9]T[012][0-9]:[0-6][0-9]:[0-6][0-9][A-Z0-9:\\+]{1,3})")
)

func hasLine(t *testing.T, content []byte, wanted string, checkTs bool) bool {
	var (
		line string
		res  [][]string
	)

	for _, line = range strings.Split(string(content), "\n") {
		if line == "" {
			continue
		}

		if checkTs {
			res = reRFC3339TS.FindAllStringSubmatch(line, -1)
			if len(res) == 0 && !strings.Contains(line, " FATAL   :") {
				t.Error("Did not find timestamp in line")
				return false
			}

			line = line[len(res[0][0])+1:]
		}

		if line == wanted {
			return true
		}
	}

	return false
}

func newTestLog(debug, verbose, timestamp bool) (Log, error) {
	var (
		fd  *os.File
		err error
	)

	if fd, err = ioutil.TempFile("/tmp", "ch-logger_test"); err != nil {
		return Log{}, err
	}

	return Log{
		UseDebug:     debug,
		UseVerbose:   verbose,
		UseTimestamp: timestamp,
		TestFd:       fd,
	}, nil
}

func cleanTestLog(fd *os.File) {
	defer os.Remove(fd.Name())
	fd.Close()
}

func runAllTests(t *testing.T, debug, verbose, timestamp bool) {
	var (
		content []byte
		log     Log
		err     error
		key     byte
		value   string
	)

	if log, err = newTestLog(debug, verbose, timestamp); err != nil {
		t.Error("Failed to open test log: " + err.Error())
	}
	defer cleanTestLog(log.TestFd)

	if log.UseDebug != debug {
		t.Error("log.UseDebug != debug")
	}

	if log.UseVerbose != verbose {
		t.Error("log.UseVerbose != verbose")
	}

	if log.UseTimestamp != timestamp {
		t.Error("log.UseTimestamp != timestamp")
	}

	log.Info("test info")
	log.Warn("test warning")
	log.Verbose("test verbose")
	log.Debug("test debug")
	log.Error("test fatal")

	if content, err = ioutil.ReadFile(log.TestFd.Name()); err != nil {
		t.Error("Failed to read test log: " + err.Error())
	}

	for key, value = range TEST_MSGS_NOTS {
		switch key {
		case MSG_VERBOSE:
			{
				if verbose || debug {
					if !hasLine(t, content, value, timestamp) {
						t.Error("Did not find", value, "in log output")
					}
				} else {
					if hasLine(t, content, value, timestamp) {
						t.Error("Found", value, "in log output")
					}
				}
			}
		case MSG_DEBUG:
			{
				if debug {
					if !hasLine(t, content, value, timestamp) {
						t.Error("Did not find", value, "in log output")
					}
				} else {
					if hasLine(t, content, value, timestamp) {
						t.Error("Found", value, "in log output")
					}
				}
			}
		default:
			{
				if !hasLine(t, content, value, timestamp) {
					t.Error("Did not find", value, "in log output")
				}
			}
		}
	}
}

func TestLogger(t *testing.T) {
	runAllTests(t, false, false, false)
	runAllTests(t, true, false, false)
	runAllTests(t, false, true, false)
	runAllTests(t, true, true, false)

	runAllTests(t, false, false, true)
	//runAllTests(t, true, false, true)
	//runAllTests(t, false, true, true)
	//runAllTests(t, true, true, true)
}
