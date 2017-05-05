package run

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/r3boot/collective-herder/lib/utils"
)

const (
	NAME        string = "run"
	DESCRIPTION string = "Run a pre-configured command on the collective"
)

type CommandConfig struct {
	Description string
	Command     string
	MaxParams   int
	ParamRegexp *regexp.Regexp
}

type Result struct {
	Node     string
	Uuid     string
	Response map[string]interface{}
	Duration time.Duration
}

var (
	Commands       map[string]CommandConfig
	reTextField    *regexp.Regexp
	reCommandField *regexp.Regexp
)

func runCommand(command string, args ...string) (stdout, stderr string, err error) {
	var (
		stdoutBuf, stderrBuf bytes.Buffer
		proc                 *exec.Cmd
	)

	proc = exec.Command(command, args...)
	proc.Stdout = &stdoutBuf
	proc.Stderr = &stderrBuf

	if err = proc.Run(); err != nil {
		return
	}

	stdout = stdoutBuf.String()
	stderr = stderrBuf.String()

	return
}

func validTextField(value string) bool {
	if reTextField == nil {
		reTextField = regexp.MustCompile("^[a-zA-Z0-9_-. ]{1,256}$")
	}

	return reTextField.MatchString(value)
}

func validCommandField(value string) bool {
	if reCommandField == nil {
		reCommandField = regexp.MustCompile("^a-zA-Z0-9_-. /]{1,256}$")
	}

	return reCommandField.MatchString(value)
}

func loadCommandConfig(content map[string]interface{}) {
	var (
		name        string
		description string
		command     string
		maxParams   int
		paramRegexp string
		ok          bool
	)

	if name, ok = content["name"].(string); !ok {
		fmt.Fprintf(os.Stderr, "Error: failed to parse CommandConfig; name field is invalid")
		os.Exit(2)
	}

	if description, ok = content["description"].(string); !ok {
		fmt.Fprintf(os.Stderr, "Error: failed to parse CommandConfig; description field is invalid")
		os.Exit(2)
	}

	if command, ok = content["command"].(string); !ok {
		fmt.Fprintf(os.Stderr, "Error: failed to parse CommandConfig; command field is invalid")
		os.Exit(2)
	}

	if maxParams, ok = content["max_params"].(int); !ok {
		fmt.Fprintf(os.Stderr, "Error: failed to parse CommandConfig; max_params field is invalid")
		os.Exit(2)
	}

	if paramRegexp, ok = content["param_regexp"].(string); !ok {
		fmt.Fprintf(os.Stderr, "Error: failed to parse CommandConfig; param_regexp field is invalid")
		os.Exit(2)
	}

	if Commands == nil {
		Commands = make(map[string]CommandConfig)
	}

	Commands[name] = CommandConfig{
		Description: description,
		Command:     command,
		MaxParams:   maxParams,
		ParamRegexp: regexp.MustCompile(paramRegexp),
	}
}

func LoadCommands(commandsDir string) {
	var (
		fs      os.FileInfo
		data    []byte
		content map[string]interface{}
		files   []os.FileInfo
		err     error
	)

	if fs, err = os.Stat(commandsDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to stat "+commandsDir+": "+err.Error())
		os.Exit(2)
	}

	if !fs.IsDir() {
		fmt.Fprintf(os.Stderr, "Error: "+commandsDir+": not a directory")
		os.Exit(2)
	}

	if files, err = ioutil.ReadDir(commandsDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to read entries from "+commandsDir+": "+err.Error())
		os.Exit(2)
	}

	for _, fs = range files {
		if data, err = ioutil.ReadFile(commandsDir + "/" + fs.Name()); err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to read "+fs.Name()+": "+err.Error())
			os.Exit(2)
		}

		if content, err = utils.ParseAsYaml(data); err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to parse "+fs.Name()+" as YAML: "+err.Error())
			os.Exit(2)
		}

		loadCommandConfig(content)
	}

}

func ParseArgs(args []string) map[string]interface{} {
	var (
		arguments []string
		opts      map[string]interface{}
		value     string
	)

	opts = make(map[string]interface{})
	opts["help"] = false

	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: Need a command to run\n")
		os.Exit(2)
	}

	if len(args) == 1 {
		switch args[0] {
		case "--help", "-h":
			{
				opts["help"] = true
			}
		default:
			{
				opts["command"] = args[0]
			}
		}
	} else {
		opts["command"] = args[0]
		for _, value = range args[1:] {
			arguments = append(arguments, value)
		}
	}

	opts["arguments"] = arguments

	return opts
}

func PreRun(opts map[string]interface{}) {
	if opts == nil {
		return
	}

	if opts["help"].(bool) {
		fmt.Fprintf(os.Stderr, "Usage of run plugin: ch run [-h, --help] <command> [<param> ...]\n\n")
		fmt.Fprintf(os.Stderr, "Run a pre-configured command on the collective. Parameters\n")
		fmt.Fprintf(os.Stderr, "passed to the command are matched for a regular expression\n")
		fmt.Fprintf(os.Stderr, "before being executed.\n")
		os.Exit(2)
	}
}

func Run(opts map[string]interface{}) interface{} {
	var (
		result    map[string]interface{}
		ok        bool
		command   string
		value     interface{}
		value_s   string
		numValues int
		stdout    string
		stderr    string
		args      []string
	)

	result = make(map[string]interface{})
	result["stdout"] = ""
	result["stderr"] = ""

	if opts == nil {
		result["stderr"] = "opts == nil\n"
		fmt.Printf(result["stderr"].(string))
		return result
	}

	if opts["help"].(bool) {
		result["stderr"] = "opts[help] == true\n"
		fmt.Printf(result["stderr"].(string))
		return result
	}

	if _, ok = Commands[opts["command"].(string)]; !ok {
		result["stderr"] = "Commands[opts[command].(string)]; !ok\n"
		fmt.Printf(result["stderr"].(string))
		return result
	}

	command = Commands[opts["command"].(string)].Command

	numValues = 0
	for _, value = range opts["arguments"].([]interface{}) {
		numValues += 1
		if numValues > Commands[opts["command"].(string)].MaxParams {
			result["stderr"] = "numValues > maxParams\n"
			fmt.Printf(result["stderr"].(string))
			return result
		}

		value_s = value.(string)
		if !Commands[opts["command"].(string)].ParamRegexp.MatchString(value_s) {
			result["stderr"] = "parameter != regexp\n"
			fmt.Printf(result["stderr"].(string))
			return result
		}
		args = append(args, value.(string))
	}

	stdout, stderr, _ = runCommand(command, args...)

	result["stdout"] = stdout
	result["stderr"] = stderr

	return result
}

func Print(startTime time.Time, result interface{}, opts map[string]interface{}) {
	var (
		node      string
		rawStdout string
		rawStderr string
		line      string
	)

	node = result.(Result).Node
	rawStdout = result.(Result).Response["stdout"].(string)
	rawStderr = result.(Result).Response["stderr"].(string)

	if len(rawStderr) > 0 {
		for _, line = range strings.Split(rawStderr, "\n") {
			if len(line) == 0 {
				continue
			}
			fmt.Printf("%-20sstderr: %v\n", node, line)
		}
	}

	if len(rawStdout) > 0 {
		for _, line = range strings.Split(rawStdout, "\n") {
			if len(line) == 0 {
				continue
			}
			fmt.Printf("%-20sstdout: %v\n", node, line)
		}
	}

}

func Summary(opts map[string]interface{}) {
}
