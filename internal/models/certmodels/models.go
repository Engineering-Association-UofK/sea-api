package certmodels

import (
	"database/sql"
	"encoding/json"
	"sea-api/internal/models"
	"time"
)

type Alignment string

const (
	ALIGN_LEFT   Alignment = "left"
	ALIGN_Right  Alignment = "right"
	ALIGN_CENTER Alignment = "center"
	ALIGN_TOP    Alignment = "top"
	ALIGN_BOTTOM Alignment = "bottom"
)

type Layout struct {
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle"`
	Statement string `json:"statement"`
}

type CertificateTemplate struct {
	ID           int64           `db:"id"`
	Name         string          `db:"name"`
	Language     models.Language `db:"language"`
	Version      string          `db:"version"`
	LayoutConfig json.RawMessage `db:"layout_config"`
	CreatedAt    time.Time       `db:"created_at"`
}

type Certificate struct {
	ID         int64  `db:"id"`
	Hash       string `db:"cert_hash"`
	TemplateID int64  `db:"template_id"`
	IssuerID   int64  `db:"issuer_id"`

	EventID         sql.NullInt64 `db:"event_id"`
	RecipientUserID sql.NullInt64 `db:"recipient_user_id"`

	RecipientName string         `db:"name"`
	Title         sql.NullString `db:"title"`
	Subtitle      sql.NullString `db:"subtitle"`
	Statement     sql.NullString `db:"statement"`
	IssuedDate    time.Time      `db:"issued_date"`
}
