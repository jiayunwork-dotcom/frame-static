package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type InvalidModelError struct {
	Field string
	Msg   string
}

func (e *InvalidModelError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Msg)
	}
	return e.Msg
}

func ParseModel(r io.Reader) (*Model, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var m Model
	if err := dec.Decode(&m); err != nil {
		return nil, &InvalidModelError{Field: "json", Msg: err.Error()}
	}
	return &m, nil
}

func ParseModelBytes(b []byte) (*Model, error) {
	return ParseModel(bytes.NewReader(b))
}
