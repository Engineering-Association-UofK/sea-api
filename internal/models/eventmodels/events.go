package eventmodels

import (
	"database/sql"
	"sea-api/internal/models"
	"time"
)

type Event struct {
	ID           int64  `db:"id" json:"id" binding:"required"`
	Name         string `db:"name" json:"name" binding:"required"`
	Description  string `db:"description" json:"description" binding:"required"`
	BackgroundID int64  `db:"background_id" json:"background_id" binding:"required"`

	Belonging models.Secretariat `db:"belonging" json:"belonging" binding:"required"`

	RequireApplying bool          `db:"require_applying" json:"require_applying"`
	FormID          sql.NullInt64 `db:"form_id" json:"form_id"`
	MaxApplications sql.NullInt64 `db:"max_applications" json:"max_applications"`

	CreatedAt time.Time `db:"created_at" json:"created_at" binding:"required"`
	StartDate time.Time `db:"start_date" json:"start_date" binding:"required"`
	EndDate   time.Time `db:"end_date" json:"end_date" binding:"required"`
}

type EventCoord struct {
	ID      int64  `db:"id"`
	EventID int64  `db:"event_id"`
	Name    string `db:"name"`
	Role    string `db:"role"`
}

type EventParticipant struct {
	ID       int64     `db:"id"`
	EventID  int64     `db:"event_id"`
	UserID   int64     `db:"user_id"`
	JoinedAt time.Time `db:"joined_at"`
}

type EventApplication struct {
	ID      int64 `db:"id"`
	EventID int64 `db:"event_id"`
	UserID  int64 `db:"user_id"`
	FormID  int64 `db:"form_id"`

	Accepted  bool      `db:"Accepted"`
	StartedAt time.Time `db:"started_at"`
}

// Rows

type EventsRow struct {
	ID            int64  `db:"id"`
	Name          string `db:"name"`
	Description   string `db:"description"`
	BackgroundID  int64  `db:"background_id"`
	BackgroundKey string `db:"background_key"`

	Belonging models.Secretariat `db:"belonging"`

	RequireApplying bool          `db:"require_applying"`
	FormID          sql.NullInt64 `db:"form_id"`
	MaxApplications sql.NullInt64 `db:"max_applications"`

	CreatedAt time.Time `db:"created_at"`
	StartDate time.Time `db:"start_date"`
	EndDate   time.Time `db:"end_date"`
}
