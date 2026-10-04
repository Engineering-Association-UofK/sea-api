package notificationsmodels

import (
	"encoding/json"
	"time"
)

type NotificationType string

const (
	NotifyBasic    NotificationType = "basic"
	NotifyRedirect NotificationType = "redirect"
)

var AllowedNotificationTypes = map[NotificationType]bool{
	NotifyBasic:    true,
	NotifyRedirect: true,
}

type Notification struct {
	ID        int64            `json:"id" db:"id"`
	UserID    int64            `json:"user_id" db:"user_id"`
	Title     string           `json:"title" db:"title"`
	Message   string           `json:"message" db:"message"`
	Type      NotificationType `json:"type" db:"type"`
	Data      *json.RawMessage `json:"data" db:"data"`
	CreatedAt time.Time        `json:"created_at" db:"created_at"`
	IsRead    bool             `json:"is_read" db:"is_read"`
}

type NotifyRedirectData struct {
	Name string `json:"name" bending:"required"`
	Path string `json:"path" bending:"required"`
}

type NotificationInsertRow struct {
	IDs     []int64
	Title   string
	Message string
	Type    NotificationType
	Data    []byte
}
