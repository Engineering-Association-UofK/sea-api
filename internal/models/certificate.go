package models

import (
	"time"
)

type CertStatus string

const (
	ACTIVE  CertStatus = "ACTIVE"
	REVOKED CertStatus = "REVOKED"
)

type DefaultCertificateData struct {
	Name        string
	EventName   string
	Grade       float64
	TaskColumns [][]string
	QRCode      string

	CollabName string
	Signature  string

	StartDate string
	EndDate   string
	NowDate   string
}

type CertificateModel struct {
	ID        int64      `db:"id"`
	Hash      string     `db:"cert_hash"`
	UserID    int64      `db:"user_id"`
	EventID   int64      `db:"event_id"`
	Grade     float64    `db:"grade"`
	IssueDate time.Time  `db:"issue_date"`
	Status    CertStatus `db:"status"`
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

type CertificateListResponse struct {
	Hash      string     `db:"cert_hash"`
	UserID    int64      `db:"user_id"`
	EventID   int64      `db:"event_id"`
	Grade     float64    `db:"grade"`
	IssueDate time.Time  `db:"issue_date"`
	Status    CertStatus `db:"status"`
}
