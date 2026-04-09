package models

import (
	"mime/multipart"
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

type Progress struct {
	Total     int       `json:"total"`
	Current   int       `json:"current"`
	ID        int64     `json:"id"`
	Success   bool      `json:"success"`
	Timestamp time.Time `json:"timestamp"`
	Name      string    `json:"name"`
}

type CertificateSendEmailData struct {
	EventID int64    `json:"event_id" binding:"required"`
	Cc      []string `json:"cc"`
	Bcc     []string `json:"bcc"`
}

type SignPdfRequest struct {
	EventID int64        `form:"event_id" binding:"required"`
	Type    DocumentType `form:"type" binding:"required"`

	Metadata map[string]string `form:"metadata"`

	QrX float64 `form:"qr_x"`
	QrY float64 `form:"qr_y"`
	QrS float64 `form:"qr_s" binding:"required"`

	File *multipart.FileHeader `form:"file" binding:"required"`
}
