package models

import "time"

type SuspensionModel struct {
	Id        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	AdminID   int64     `db:"admin_id"`
	Reason    string    `db:"reason"`
	StartedAt time.Time `db:"started_at"`
	EndedAt   time.Time `db:"ended_at"`
}

type SuspensionRequest struct {
	UserID   int64  `json:"user_id"`
	Reason   string `json:"reason"`
	Duration int64  `json:"duration"`
}

type SuspensionResponse struct {
	UserID    int64     `json:"user_id"`
	AdminID   int64     `json:"admin_id"`
	Reason    string    `json:"reason"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"`
}
