package models

import "encoding/json"

type EmailType string

const (
	ANNOUNCEMENT EmailType = "ANNOUNCEMENT"
	NOTIFICATION EmailType = "NOTIFICATION"
	INVITATION   EmailType = "INVITATION"
	WELCOME      EmailType = "WELCOME"
	RESET_PASS   EmailType = "RESET_PASS"
	TECHNICAL    EmailType = "TECHNICAL"
	CERTIFICATE  EmailType = "CERTIFICATE"
)

type Email struct {
	To      []string `json:"to" binding:"required"`
	Cc      []string `json:"cc"`
	Bcc     []string `json:"bcc"`
	Subject string   `json:"subject" binding:"required"`
	ReplyTo string   `json:"reply_to"`
	HTML    string   `json:"html"`
	Text    string   `json:"text"`
}

type UserEmails struct {
	UserIDs []int64         `json:"user_ids" binding:"required"`
	Subject string          `json:"subject" binding:"required"`
	Type    EmailType       `json:"type" binding:"required"`
	Preview bool            `json:"preview"`
	Data    json.RawMessage `json:"data" binding:"required"`
}

type TechnicalEmail struct {
	ServiceName string `json:"service_name"`
	Status      string `json:"status"`
	Message     string `json:"message"`
	ActionURL   string `json:"action_url"`
	ActionText  string `json:"action_text"`
	Lang        string `json:"lang" binding:"required"`
}

type TechnicalEmailTemplate struct {
	TechnicalEmail
	Username string `json:"username"`
	Year     int
}

type CertificateEmailData struct {
	Username  string
	EventName string
	EventType string
	CertURL   string
	Year      int
}
