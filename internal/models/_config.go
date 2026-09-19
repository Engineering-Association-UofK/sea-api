package models

import "encoding/json"

type config struct {
	Key   string          `db:"string"`
	Value json.RawMessage `db:"value"`
}
