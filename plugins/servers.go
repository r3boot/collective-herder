package plugins

import (
	"os"
	"strconv"

	"github.com/r3boot/collective-herder/lib/utils"
	"github.com/r3boot/collective-herder/plugins/ping"
)

func NewServers(l utils.Log) *Servers {
	var (
		p   *Servers
		err error
	)
	Log = l

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
	p.runFunc[ping.NAME] = ping.Run
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
