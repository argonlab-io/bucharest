package utils

import (
	"encoding/json"
	"io"
)

func JSONMapper(src interface{}, dest interface{}) error {
	var err error
	var data map[string]interface{}
	body, ok := src.(io.ReadCloser)
	if ok {
		json.NewDecoder(body).Decode(&data)
		body.Close()
		src = data
	}
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dest)
}
