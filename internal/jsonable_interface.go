package internal

import "encoding/json"

type Jsonable interface {
	JSON() []byte
}

type JsonWrapper struct {
	Data interface{}
}

func (jw *JsonWrapper) JSON() []byte {
	data, err := json.Marshal(jw.Data)
	if err != nil {
		return []byte("{\"error\": \"conveting struct to json\"}")
	}
	return data
}
