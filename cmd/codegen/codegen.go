package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"text/template"
)

type TemplateVars struct {
	Name    string
	Plugins []string
}

var TEMPLATE string = `package {{.Name}}
/*
 * WARNING: Modify at your own risk, this file is auto-generated ...
 */

import (
	"fmt"
	"os"
	"time"

{{range .Plugins}}
	"github.com/r3boot/collective-herder/plugins/{{.}}"
{{end}}
)

func (p *Agents) LoadAllAgents() {
{{range .Plugins}}
	// Glue code for {{.}} plugin
	p.argsFunc[{{.}}.NAME] = {{.}}.ParseArgs
	p.preRunFunc[{{.}}.NAME] = {{.}}.PreRun
	p.printFunc[{{.}}.NAME] = {{.}}.Print
	p.summaryFunc[{{.}}.NAME] = {{.}}.Summary
	p.Meta[{{.}}.NAME] = {{.}}.DESCRIPTION

{{end}}
}

func (p *Agents) Print(plugin, uuid string, startTime time.Time, response interface{}, opts map[string]interface{}) {
	var (
		node           string
		hostUuid       string
		responseResult interface{}
	)

	node = response.(Response).Node
	hostUuid = response.(Response).HostUuid
	responseResult = response.(Response).Result

	switch plugin {
{{range .Plugins}}
	case {{.}}.NAME:
		{
			result := {{.}}.Result{
				Node: node,
				Uuid: hostUuid,
				Duration: time.Since(startTime),
				Response: responseResult.(map[string]interface{}),
			}
			p.printFunc[plugin](startTime, result, opts)
		}
{{end}}
	default:
		{
			fmt.Fprintf(os.Stderr, "Error: Print: unknown plugin: %s", plugin)
			os.Exit(1)
		}
	}
}

func (p *Agents) Summary(plugin string, opts map[string]interface{}) {
	switch plugin {
{{range .Plugins}}
	case {{.}}.NAME:
		{
			p.summaryFunc[plugin](opts)
		}
{{end}}
	default:
		{
			fmt.Fprintf(os.Stderr, "Error: Summary: unknown plugin: %s", plugin)
			os.Exit(1)
		}
	}
}
`

/*
 * This program assumes it's running from the ./build directory, relative to
 * the source root directory.
 */
func main() {
	var (
		fullPath       string
		baseDir        string
		pluginDir      string
		agentsDestFile string
		vars           TemplateVars
		tmpl           *template.Template
		fd             *os.File
		fs             os.FileInfo
		entries        []os.FileInfo
		plugins        []string
		err            error
	)

	if fullPath, err = filepath.Abs(os.Args[0]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to get absolute path for self: "+err.Error())
		os.Exit(1)
	}

	baseDir = filepath.Dir(filepath.Dir(fullPath))
	pluginDir = baseDir + "/plugins"
	agentsDestFile = baseDir + "/plugins/agents_gen.go"

	if fs, err = os.Stat(pluginDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to stat "+pluginDir+": "+err.Error())
		os.Exit(1)
	}

	if !fs.IsDir() {
		fmt.Fprintf(os.Stderr, "Error: "+pluginDir+": not a directory")
		os.Exit(1)
	}

	if entries, err = ioutil.ReadDir(pluginDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to read "+pluginDir+": "+err.Error())
		os.Exit(1)
	}

	for _, fs = range entries {
		if !fs.IsDir() {
			continue
		}

		plugins = append(plugins, fs.Name())
	}

	if tmpl, err = template.New("codegen").Parse(TEMPLATE); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to load template: "+err.Error())
		os.Exit(1)
	}

	vars = TemplateVars{
		Name:    "plugins",
		Plugins: plugins,
	}

	fmt.Printf("Generating plugins/%s: ", path.Base(agentsDestFile))

	if _, err = os.Stat(agentsDestFile); err == nil {
		if err = os.Remove(agentsDestFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to remove "+agentsDestFile+": "+err.Error())
			os.Exit(1)
		}
	}

	if fd, err = os.OpenFile(agentsDestFile, os.O_RDWR|os.O_CREATE, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to open "+agentsDestFile+": "+err.Error())
		os.Exit(1)
	}
	defer fd.Close()

	if err = tmpl.Execute(fd, vars); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to execute template: "+err.Error())
		os.Exit(1)
	}

	fmt.Printf("ok\n")

}
