package models

import "encoding/json"

// config
type config struct {
	Key   string          `db:"string"`
	Value json.RawMessage `db:"value"`
}
