package handlers

import (
	"fmt"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/response"
	"sea-api/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	service *services.NotificationService
}

func NewNotificationHandler(service *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

// CreateDemoNotifications godocs
//
//	@Summary		Create demo notification
//	@Description	Create a demo notification for the account that sends the request
//	@Tags			Notifications
//	@Param			body	body	models.DemoNotificationRequest	true	"Request body"
//	@Produce		json
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/account/notifications/demo [post]
//
//	@Security		ApiKeyAuth
func (h *NotificationHandler) CreateDemoNotifications(c *gin.Context) {
	var req models.DemoNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.New(errs.BadRequest, "Invalid request body", nil))
		return
	}

	value, exists := c.Get("user")
	claims, ok := value.(*models.ManagedClaims)
	if !exists || !ok {
		c.Error(errs.New(errs.Unauthorized, "Unauthorized", nil))
		return
	}

	id, err := h.service.CreateDemoNotifications(claims.UserID, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(201, "Notification created successfully", id, c)
}

// GetNotifications godocs
//
//	@Summary		Get notifications
//	@Description	Get latest notifications for requester
//	@Tags			Notifications
//	@Param			limit	query	int	true	"Content count limit"
//	@Param			page	query	int	true	"Page number"
//	@Produce		json
//	@Success		200	{object}	models.NotificationsListResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/account/notifications [get]
//
//	@Security		ApiKeyAuth
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	var req models.ListRequest
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
//	@Router			/account/notifications/{id} [get]
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
