package models

import (
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
	Data      interface{}      `json:"data" db:"data"`
	CreatedAt time.Time        `json:"created_at" db:"created_at"`
	IsRead    bool             `json:"is_read" db:"is_read"`
}

type NotifyRedirectData struct {
	Title string `json:"title"`
	Path  string `json:"path"`
}

type NotificationRequest struct {
	UserID  int64            `json:"user_id" binding:"required"`
	Title   string           `json:"title" binding:"required"`
	Message string           `json:"message" binding:"required"`
	Type    NotificationType `json:"type" binding:"required"`
	Data    interface{}      `json:"data"`
}

type DemoNotificationRequest struct {
	Title   string `json:"title" binding:"required"`
	Message string `json:"message" binding:"required"`
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

type NotificationListRequest struct {
	ListRequest
}

type NotificationsListResponse struct {
	ListResponse
	List []NotificationResponse `json:"list"`
}
