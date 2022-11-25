package fetch_source

import (
	"bytes"
	"encoding/json"
)

type Client_message struct {
	Type string
	Data interface{}
}

func Encode_to_bytes(obj interface{}) (error, []byte) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	err := encoder.Encode(obj)
	return err, buf.Bytes()
}
