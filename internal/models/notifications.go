package models

import "time"

type NotificationType string

const (
	NotifyBasic       NotificationType = "basic"
	NotifyEvent       NotificationType = "event"
	NotifyCertificate NotificationType = "certificate"
	NotifyApplication NotificationType = "application"
)

var AllowedNotificationTypes = map[NotificationType]bool{
	NotifyBasic:       true,
	NotifyEvent:       true,
	NotifyCertificate: true,
	NotifyApplication: true,
}

type Notification struct {
	ID        int64            `json:"id" db:"id"`
	UserID    int64            `json:"user_id" db:"user_id"`
	Title     string           `json:"title" db:"title"`
	Message   string           `json:"message" db:"message"`
	Type      NotificationType `json:"type" db:"type"`
	Data      interface{}      `json:"data" db:"data"`
	CreatedAt time.Time        `json:"created_at" db:"created_at"`
	IsRead    bool             `json:"is_read" db:"is_read"`
}

type NotifyEventData struct {
	EventID   int64  `json:"event_id"`
	EventName string `json:"event_name"`
	Action    string `json:"action"`
}

type NotifyCertificateData struct {
	EventID         int64  `json:"event_id"`
	CertificateHash string `json:"certificate_hash"`
}

type NotificationRequest struct {
	UserID  int64            `json:"user_id" binding:"required"`
	Title   string           `json:"title" binding:"required"`
	Message string           `json:"message" binding:"required"`
	Type    NotificationType `json:"type" binding:"required"`
	Data    interface{}      `json:"data"`
}

type NotificationResponse struct {
	ID        int64            `json:"id"`
	Title     string           `json:"title"`
	Message   string           `json:"message"`
	Type      NotificationType `json:"type"`
	Data      interface{}      `json:"data"`
	CreatedAt time.Time        `json:"created_at"`
	IsRead    bool             `json:"is_read"`
}

type NotificationsListResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	Pages         int                    `json:"pages"`
	Total         int                    `json:"total"`
}
