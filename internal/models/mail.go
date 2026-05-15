package models

import "encoding/json"

type EmailType string

const (
	EmailAccCode        EmailType = "acc_code_verification"
	EmailAccLink        EmailType = "acc_link_verification"
	EmailAccPassChanged EmailType = "acc_password_change_success"
	EmailAccPassReset   EmailType = "acc_password_reset"
	EmailAccSuspension  EmailType = "acc_suspension"

	EmailEventAcceptance    EmailType = "event_acceptance"
	EmailEventRejection     EmailType = "event_rejection"
	EmailEventCertificate   EmailType = "event_certificate"
	EmailEventAnnounce      EmailType = "event_announcement"
	EmailEventReminder      EmailType = "event_reminder"
	EmailEventThankyou      EmailType = "event_post_event_thankyou"
	EmailEventSpeakerInvite EmailType = "event_speaker_invitation"

	EmailTechnicalNotify EmailType = "tech_system_notification"
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

type UsersEmails struct {
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
	Year      int
}
