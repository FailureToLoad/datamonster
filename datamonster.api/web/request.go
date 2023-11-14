package web

import (
	"encoding/json"
	"io"
)

func DecodeJson(rc io.ReadCloser, data interface{}) error {
	defer rc.Close()
	decoder := json.NewDecoder(rc)
	return decoder.Decode(data)
}
