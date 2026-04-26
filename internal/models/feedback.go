package models

import "time"

type Feedback struct {
	ID        int          `db:"id" json:"id"`
	Message   string       `db:"message" json:"message"`
	UserID    *int64       `db:"user_id" json:"user_id"`
	Type      FeedbackType `db:"type" json:"type"`
	CreatedAt time.Time    `db:"created_at" json:"created_at"`
}
