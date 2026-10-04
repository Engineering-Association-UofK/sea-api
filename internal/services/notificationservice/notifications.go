package notificationservice

import (
	"encoding/json"
	"fmt"
	"sea-api/internal/errs"
	"sea-api/internal/models/notificationsmodels"
	"sea-api/internal/repositories/notificationsrepo"
	"time"
)

type NotificationService struct {
	repo  *notificationsrepo.NotificationRepository
	queue <-chan notificationsmodels.NotificationRequest
}

func NewNotificationService(repo *notificationsrepo.NotificationRepository) *NotificationService {
	s := &NotificationService{
		repo:  repo,
		queue: make(<-chan notificationsmodels.NotificationRequest),
	}

	for i := 0; i < 5; i++ {
		go s.worker()
	}

	return s
}

func (s *NotificationService) worker() {
	for n := range s.queue {
		s.CreateNotification(&n)
	}
}

// Used only internally for now
func (s *NotificationService) CreateNotification(req *notificationsmodels.NotificationRequest) (int64, error) {
	if !notificationsmodels.AllowedNotificationTypes[req.Type] {
		return 0, errs.New(errs.BadRequest, "Invalid notification type", nil)
	}

	if req.Data != nil {
		err := verifyType(req.Type, req.Data)
		if err != nil {
			return 0, err
		}
	}

	return s.repo.Create(&notificationsmodels.Notification{
		UserID:    req.UserID,
		Title:     req.Title,
		Message:   req.Message,
		Type:      req.Type,
		Data:      req.Data,
		CreatedAt: time.Now(),
		IsRead:    false,
	})
}

func (s *NotificationService) BulkGenerateForUsers(req *notificationsmodels.BulkCreateForUsersRequest) error {
	if !notificationsmodels.AllowedNotificationTypes[req.Type] {
		return errs.New(errs.BadRequest, "Invalid notification type", nil)
	}

	if req.Data != nil {
		err := verifyType(req.Type, req.Data)
		if err != nil {
			return err
		}
	}

	var dataBytes []byte
	var err error
	if req.Data != nil {
		dataBytes, err = json.Marshal(req.Data)
		if err != nil {
			return fmt.Errorf("failed to marshal notification data: %w", err)
		}
	}

	return s.repo.BulkCreateForUsers(&notificationsmodels.NotificationInsertRow{
		IDs:     req.IDs,
		Title:   req.Title,
		Message: req.Message,
		Type:    req.Type,
		Data:    dataBytes,
	})
}

func (s *NotificationService) BulkGenerateForAllUsers(req *notificationsmodels.BulkNotificationRequest) error {
	if !notificationsmodels.AllowedNotificationTypes[req.Type] {
		return errs.New(errs.BadRequest, "Invalid notification type", nil)
	}

	if req.Data != nil {
		err := verifyType(req.Type, req.Data)
		if err != nil {
			return err
		}
	}

	var dataBytes []byte
	var err error
	if req.Data != nil {
		dataBytes, err = json.Marshal(req.Data)
		if err != nil {
			return fmt.Errorf("failed to marshal notification data: %w", err)
		}
	}

	return s.repo.BulkCreateForUsers(&notificationsmodels.NotificationInsertRow{
		IDs:     nil,
		Title:   req.Title,
		Message: req.Message,
		Type:    req.Type,
		Data:    dataBytes,
	})
}

func (s *NotificationService) GetUnreadCount(userID int64) (int, error) {
	return s.repo.GetUnreadCount(userID)
}

func verifyType(dataType notificationsmodels.NotificationType, data *json.RawMessage) error {
	switch dataType {
	case notificationsmodels.NotifyBasic:
		data = nil
		return nil

	case notificationsmodels.NotifyRedirect:
		var d any
		err := json.Unmarshal(*data, &d)
		if err != nil {
			return errs.New(errs.BadRequest, fmt.Sprintf(`data content is invalid, it should be like: %s`, notificationsmodels.NotifyRedirectData{}), nil)
		}
		return nil

		// Add new types here if ever made

	default:
		return errs.New(errs.BadRequest, fmt.Sprintf(`notification Type '%s' is not found, try '%s'`, dataType, notificationsmodels.NotifyBasic), nil)
	}
}
