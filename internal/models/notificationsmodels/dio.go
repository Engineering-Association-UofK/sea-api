package notificationsmodels

import (
	"encoding/json"
	"sea-api/internal/models"
	"time"
)

type NotificationRequest struct {
	UserID  int64            `json:"user_id" binding:"required"`
	Title   string           `json:"title" binding:"required"`
	Message string           `json:"message" binding:"required"`
	Type    NotificationType `json:"type" binding:"required"`
	Data    *json.RawMessage `json:"data"`
}

type BulkNotificationRequest struct {
	Title   string           `json:"title" binding:"required"`
	Message string           `json:"message" binding:"required"`
	Type    NotificationType `json:"type" binding:"required"`
	Data    *json.RawMessage `json:"data"`
}

type BulkCreateForUsersRequest struct {
	IDs []int64 `json:"ids" binding:"required"`
	BulkNotificationRequest
}

type NotificationResponse struct {
	ID        int64            `json:"id"`
	Title     string           `json:"title"`
	Message   string           `json:"message"`
	Type      NotificationType `json:"type"`
	Data      *json.RawMessage `json:"data"`
	CreatedAt time.Time        `json:"created_at"`
	IsRead    bool             `json:"is_read"`
}

type NotificationListRequest struct {
	models.ListRequest
}

type NotificationsListResponse struct {
	models.ListResponse
	List []NotificationResponse `json:"list"`
}
