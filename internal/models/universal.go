package models

import "time"

type ListRequest struct {
	Limit int64 `form:"limit"`
	Page  int64 `form:"page"`
}

type ListResponse struct {
	TotalPages  int64 `json:"total_pages"`
	CurrentPage int64 `json:"current_page"`
	Count       int64 `json:"count"`
}

var AllowedListLimit = map[int64]bool{
	10:  true,
	25:  true,
	50:  true,
	100: true,
}

type Progress struct {
	Total     int       `json:"total"`
	Current   int       `json:"current"`
	ID        int64     `json:"id"`
	Success   bool      `json:"success"`
	Timestamp time.Time `json:"timestamp"`
	Name      string    `json:"name"`
}

type ProgressError struct {
	Error     string    `json:"error"`
	Timestamp time.Time `json:"timestamp"`
}

type Language string

const (
	Arabic  Language = "ar"
	English Language = "en"
)

var AllowedLanguages = map[Language]bool{
	Arabic:  true,
	English: true,
}
