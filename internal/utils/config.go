package utils

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sea-api/internal/models"
	"strings"

	"github.com/jmoiron/sqlx"
)

// GetConfig returns the JSON global config state
func GetConfig(name string, db *sqlx.DB) (*json.RawMessage, error) {
	var raw json.RawMessage
	query := fmt.Sprintf("SELECT value FROM %s WHERE `key` = ?", models.TableConfig)

	err := db.Get(&raw, query, name)
	if err != nil {
		return nil, err
	}
	return &raw, nil
}

// UpdateConfig serializes and updates the JSON global config state
func UpdateConfig(name string, db *sqlx.DB, cfg any) error {
	val, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE %s SET value = ? WHERE `key` = ?", models.TableConfig)
	_, err = db.Exec(query, val, name)
	return err
}

// ConfigExists checks if a config key exists and has a non-empty JSON value
func ConfigExists(name string, db *sqlx.DB) (bool, error) {
	var raw *json.RawMessage
	query := fmt.Sprintf("SELECT value FROM %s WHERE `key` = ?", models.TableConfig)

	err := db.Get(&raw, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil // Key does not exist
		}
		return false, err // DB error
	}

	// Check if pointer is nil, or if JSON string is empty/null/whitespace
	if raw == nil {
		return false, nil
	}

	trimmed := strings.TrimSpace(string(*raw))
	if trimmed == "" || trimmed == "null" || trimmed == "{}" {
		return false, nil
	}

	return true, nil
}
