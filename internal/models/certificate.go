package models

import (
	"time"
)

type CertStatus string
type CertType string
type CertVersion string

const (
	CertActive  CertStatus = "ACTIVE"
	CertExpired CertStatus = "EXPIRED"
	CertRevoked CertStatus = "REVOKED"

	CertParticipation CertType = "Participation"
	CertCompletion    CertType = "Completion"

	V0_1 CertVersion = "v0.1"
)

type CertificateModel struct {
	ID          int64       `db:"id"`
	Hash        string      `db:"cert_hash"`
	UserID      int64       `db:"user_id"`
	EventID     int64       `db:"event_id"`
	Type        CertType    `db:"type"`
	CertVersion CertVersion `db:"cert_version"`
	Grade       float64     `db:"grade"`
	IssueDate   time.Time   `db:"issue_date"`
	Status      CertStatus  `db:"status"`
}

type CertificateFileModel struct {
	ID            int64  `db:"id"`
	CertificateID int64  `db:"certificate_id"`
	StoreID       int64  `db:"store_id"`
	Lang          string `db:"lang"`
}

type CertificateVerify struct {
	Valid     bool       `json:"valid"`
	ID        string     `json:"id"`
	NameAr    string     `json:"name_ar"`
	NameEn    string     `json:"name_en"`
	EventName string     `json:"event"`
	Status    CertStatus `json:"status"`
	Grade     string     `json:"grade"`
	Outcomes  []string   `json:"outcomes"`
	EndDate   time.Time  `json:"end_date"`
	IssueDate time.Time  `json:"issue_date"`
}

type CertificateSendEmailData struct {
	EventID int64    `json:"event_id" binding:"required"`
	Cc      []string `json:"cc"`
	Bcc     []string `json:"bcc"`
}

// DTOs

type CertificateResponse struct {
	Hash      string      `db:"cert_hash" json:"hash"`
	EventID   int64       `db:"event_id" json:"event_id"`
	EventName string      `db:"event_name" json:"event_name"`
	Type      CertType    `db:"type" json:"type"`
	Version   CertVersion `db:"cert_version" json:"version"`
	Grade     float64     `db:"grade" json:"grade"`
	IssueDate time.Time   `db:"issue_date" json:"issue_date"`
	Status    CertStatus  `db:"status" json:"status"`
}

type CertificateListResponse struct {
	Current      int64                 `json:"current"`
	Pages        int64                 `json:"pages"`
	Certificates []CertificateResponse `json:"certificates"`
}
