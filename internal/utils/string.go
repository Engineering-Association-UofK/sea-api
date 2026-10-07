package utils

import (
	"encoding/json"
	"strings"
)

type TrimmedString string

func (ts *TrimmedString) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*ts = TrimmedString(strings.TrimSpace(s))
	return nil
}
