package certmodels

import (
	"database/sql"
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

type LayoutPosition struct {
	X        int       `json:"x"`
	Y        int       `json:"y"`
	FontSize int       `json:"font_size"`
	Align    Alignment `json:"align"`
}

type Layout struct {
	Title         *LayoutPosition `json:"title"`
	Subtitle      *LayoutPosition `json:"subtitle"`
	Statement     *LayoutPosition `json:"statement"`
	RecipientName *LayoutPosition `json:"recipient_name"`
	Date          *LayoutPosition `json:"date"`

	CoordOneName *LayoutPosition `json:"coord_one_name"`
	CoordOneRole *LayoutPosition `json:"coord_one_role"`
	CoordOneSign *LayoutPosition `json:"coord_one_sign"`

	CoordTwoName *LayoutPosition `json:"coord_two_name"`
	CoordTwoRole *LayoutPosition `json:"coord_two_role"`
	CoordTwoSign *LayoutPosition `json:"coord_two_sign"`

	QR *LayoutPosition `json:"qr"`
}

type CertificateTemplate struct {
	ID                int64  `db:"id"`
	Name              string `db:"name"`
	Version           string `db:"version"`
	BackgroundImageID int64  `db:"background_image_file_id"`
	LayoutConfig      string `db:"layout_config"`
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

type CertificateCollaborator struct {
	ID            int64          `db:"id"`
	CertID        int64          `db:"cert_id"`
	Name          string         `db:"name"`
	Role          string         `db:"role"`
	SignatureID   sql.NullString `db:"signature_file_id"`
	DisplayOnCert bool           `db:"display_on_cert"`
}
