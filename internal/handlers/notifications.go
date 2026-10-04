package handlers

import (
	"fmt"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/models/notificationsmodels"
	"sea-api/internal/response"
	"sea-api/internal/services/notificationservice"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	service *notificationservice.NotificationService
}

func NewNotificationHandler(service *notificationservice.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

// BulkCreateForUsers godocs
//
//	@Summary		Create notifications
//	@Description	Create a notification for every user whose ID is in the request
//	@Tags			Notifications
//	@Param			body	body	notificationsmodels.BulkCreateForUsersRequest	true	"Request body"
//	@Produce		json
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/notifications [post]
//
//	@Security		ApiKeyAuth
func (h *NotificationHandler) BulkCreateForUsers(c *gin.Context) {
	var req notificationsmodels.BulkCreateForUsersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.New(errs.BadRequest, "Invalid request body", nil))
		return
	}

	err := h.service.BulkGenerateForUsers(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(201, "Notifications created successfully", 0, c)
}

// BulkCreateForAllUsers godocs
//
//	@Summary		Create a Global Notification
//	@Description	Create a notification for all users
//	@Tags			Notifications
//	@Param			body	body	notificationsmodels.BulkNotificationRequest	true	"Request body"
//	@Produce		json
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/notifications/all [post]
//
//	@Security		ApiKeyAuth
func (h *NotificationHandler) BulkCreateForAllUsers(c *gin.Context) {
	var req notificationsmodels.BulkNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.New(errs.BadRequest, "Invalid request body", nil))
		return
	}

	err := h.service.BulkGenerateForAllUsers(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(201, "Notification created successfully", 0, c)
}

// GetNotifications godocs
//
//	@Summary		Get notifications
//	@Description	Get latest notifications for requester
//	@Tags			Notifications
//	@Param			limit	query	int	true	"Content count limit"
//	@Param			page	query	int	true	"Page number"
//	@Produce		json
//	@Success		200	{object}	notificationsmodels.NotificationsListResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/account/notifications [get]
//
//	@Security		ApiKeyAuth
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	var req notificationsmodels.NotificationListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	value, exists := c.Get("user")
	claims, ok := value.(*models.ManagedClaims)
	if !exists || !ok {
		c.Error(errs.New(errs.Unauthorized, "Unauthorized", nil))
		return
	}

	resp, err := h.service.GetNotificationsByUserID(claims.UserID, req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, resp)
}

// MarkAsRead godocs
//
//	@Summary		Mark one notification as read
//	@Description	Marks the provided notification ID as read, the notification must belong to the requesting user
//	@Tags			Notifications
//	@Param			id	path	int	true	"Notification ID"
//	@Produce		json
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/account/notifications/{id} [post]
//
//	@Security		ApiKeyAuth
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Error(errs.New(errs.BadRequest, "Invalid notification ID", nil))
		return
	}

	value, exists := c.Get("user")
	claims, ok := value.(*models.ManagedClaims)
	if !exists || !ok {
		c.Error(errs.New(errs.Unauthorized, "Unauthorized", nil))
		return
	}

	affected, err := h.service.MarkAsRead(claims.UserID, id)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(200, fmt.Sprintf("%d Notification marked as read", affected), id, c)
}

// MarkAllAsRead godocs
//
//	@Summary		Mark all notification as read
//	@Description	Marks all notifications that belong to the requesting user as read
//	@Tags			Notifications
//	@Produce		json
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/account/notifications [post]
//
//	@Security		ApiKeyAuth
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	value, exists := c.Get("user")
	claims, ok := value.(*models.ManagedClaims)
	if !exists || !ok {
		c.Error(errs.New(errs.Unauthorized, "Unauthorized", nil))
		return
	}

	err := h.service.MarkAllAsRead(claims.UserID)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(200, "All notifications marked as read", claims.UserID, c)
}

// DeleteNotification godocs
//
//	@Summary		Delete notification
//	@Description	Delete one notification with ID, must belong to the requesting user
//	@Tags			Notifications
//	@Param			id	path	int	true	"Notification ID"
//	@Produce		json
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/account/notifications/{id} [delete]
//
//	@Security		ApiKeyAuth
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Error(errs.New(errs.BadRequest, "Invalid notification ID", nil))
		return
	}

	value, exists := c.Get("user")
	claims, ok := value.(*models.ManagedClaims)
	if !exists || !ok {
		c.Error(errs.New(errs.Unauthorized, "Unauthorized", nil))
		return
	}

	affected, err := h.service.Delete(claims.UserID, id)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(200, fmt.Sprintf("%d Notifications deleted", affected), id, c)
}
