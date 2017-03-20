package plugins

import (
	"encoding/json"
)

func (req *Request) ToJSON() ([]byte, error) {
	var (
		data []byte
		err  error
	)

	data, err = json.Marshal(req)
	return data, err
}

func (resp *Response) ToJSON() ([]byte, error) {
	var (
		data []byte
		err  error
	)

	data, err = json.Marshal(resp)
	return data, err
}
