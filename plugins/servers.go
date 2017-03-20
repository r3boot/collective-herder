package plugins

import (
	"os"
	"strconv"

	sysFacts "github.com/r3boot/collective-herder/lib/facts"
	"github.com/r3boot/collective-herder/lib/utils"
	"github.com/r3boot/collective-herder/plugins/facts"
	"github.com/r3boot/collective-herder/plugins/ping"
)

var (
	Facts *sysFacts.Facts
)

func NewServers(l utils.Log, f *sysFacts.Facts) *Servers {
	var (
		p   *Servers
		err error
	)
	Log = l
	Facts = f

	p = &Servers{}
	p.LoadAllServers()

	if Hostname, err = os.Hostname(); err != nil {
		Log.Error("NewServers: Failed to get hostname: " + err.Error())
		os.Exit(1)
	}

	HostUuid = utils.Uuidgen()

	return p
}

func (p *Servers) LoadAllServers() {
	p.runFunc = make(map[string]func(map[string]interface{}) interface{})

	// Ping server
	p.runFunc[ping.NAME] = ping.Run

	// Facts server
	facts.LoadFacts(Facts)
	p.runFunc[facts.NAME] = facts.Run
}

func (p *Servers) NumServersAsString() string {
	return strconv.Itoa(len(p.runFunc))
}

func (p *Servers) HasServer(name string) bool {
	var (
		key string
	)

	for key, _ = range p.runFunc {
		if key == name {
			return true
		}
	}

	return false
}

func (p *Servers) RunServer(name string, opts map[string]interface{}) interface{} {
	if !p.HasServer(name) {
		Log.Warn("RunServer: No such plugin: " + name)
		return nil
	}

	return p.runFunc[name](opts)
}
