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

func (h *NotificationHandler) Delete(c *gin.Context) {
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
