package fetch_source

import (
	"bytes"
	"encoding/json"
)

type Client_message struct {
	Type string
	Data interface{}
}

type Generic_data_single struct {
	Data map[string]interface{} `bson:"-"`
}

type Generic_data_multiple struct {
	Data []map[string]interface{} `bson:"-"`
}

func Encode_to_bytes(obj interface{}) (error, []byte) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	err := encoder.Encode(obj)
	return err, buf.Bytes()
}
