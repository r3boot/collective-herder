package plugins

import (
	"github.com/r3boot/collective-herder/lib/utils"
)

var (
	Log      utils.Logger
	Hostname string
	HostUuid string
)

func NewRequest(msgType string, facts map[string]interface{}, opts map[string]interface{}) *Request {
	return &Request{
		Uuid:    utils.Uuidgen(),
		MsgType: msgType,
		Facts:   facts,
		Opts:    opts,
	}
}

func NewResponse(uuid string, result interface{}) *Response {
	return &Response{
		Uuid:     uuid,
		Node:     Hostname,
		HostUuid: HostUuid,
		Result:   result,
	}
}
