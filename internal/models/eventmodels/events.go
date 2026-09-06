package eventmodels

import (
	"database/sql"
	"sea-api/internal/models"
	"time"
)

type Event struct {
	ID           int64  `db:"id"`
	Name         string `db:"name"`
	Description  string `db:"description"`
	BackgroundID int64  `db:"background_id"`

	Belonging models.Secretariat `db:"belonging"`

	RequireApplying bool          `db:"require_applying"`
	MaxApplications sql.NullInt64 `db:"max_applications"`

	StartDate time.Time `db:"start_date"`
	EndDate   time.Time `db:"end_date"`
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

	Accepted  bool      `db:"Accepted"`
	StartedAt time.Time `db:"started_at"`
}
