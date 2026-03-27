package utils

import (
	"encoding/json"
	"log/slog"
	"sea-api/internal/models"
)

func ParseProgressStruct(total, current int, id int64, success bool, name string, progressChan chan string) {
	s, err := parseToJsonString(models.Progress{
		Total:   total,
		Current: current,
		ID:      id,
		Success: success,
		Name:    name,
	})
	if err != nil {
		slog.Error("Error parsing progress to JSON string", "error", err)
	} else {
		progressChan <- s
	}
}

func parseToJsonString(data any) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
