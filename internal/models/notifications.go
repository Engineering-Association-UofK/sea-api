package models

import "time"

type Notification struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Title     string    `json:"title" db:"title"`
	Message   string    `json:"message" db:"message"`
	Type      string    `json:"type" db:"type"`
	Data      interface{} `json:"data" db:"data"` // Using interface{} for JSON flexibility, might need specific handling in code
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	IsRead    bool      `json:"is_read" db:"is_read"`
}