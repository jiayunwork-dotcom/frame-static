package serialize

import (
	"encoding/json"

	"frame-static/internal/assemble"
)

func ToJSON(res *assemble.Result) ([]byte, error) {
	return json.MarshalIndent(res, "", "  ")
}

func FromJSON(b []byte) (*assemble.Result, error) {
	var r assemble.Result
	if err := json.Unmarshal(b, &r); err != nil {
		return nil, err
	}
	return &r, nil
}
