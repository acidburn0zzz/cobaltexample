package main

import (
	"encoding/json"
	"io"
)

type JSONEncoder struct{}

func (enc JSONEncoder) Encode(w io.Writer, val interface{}) error {
	return json.NewEncoder(w).Encode(val)
}

func (enc JSONEncoder) Decode(r io.Reader, val interface{}) error {
	return json.NewDecoder(r).Decode(val)
}

func (enc JSONEncoder) ContentType() string {
	return "application/json;charset=UTF-8"
}
