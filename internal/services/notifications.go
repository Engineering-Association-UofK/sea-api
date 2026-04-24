package services

import (
	"fmt"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/repositories"
	"sea-api/internal/utils/valid"
	"time"
)

type NotificationService struct {
	repo *repositories.NotificationRepository
}

func NewNotificationService(repo *repositories.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) CreateNotification(notification *models.NotificationRequest) (int64, error) {
	if !models.AllowedNotificationTypes[notification.Type] {
		return 0, errs.New(errs.BadRequest, "Invalid notification type", nil)
	}

	return s.repo.Create(&models.Notification{
		UserID:    notification.UserID,
		Title:     notification.Title,
		Message:   notification.Message,
		Type:      notification.Type,
		Data:      notification.Data,
		CreatedAt: time.Now(),
		IsRead:    false,
	})
}

func (s *NotificationService) CreateDemoNotifications(userId int64, req *models.DemoNotificationRequest) (int64, error) {
	return s.CreateNotification(&models.NotificationRequest{
		UserID:  userId,
		Title:   req.Title,
		Message: req.Message,
		Type:    models.NotifyBasic,
		Data:    nil,
	})
}

func (s *NotificationService) GetNotificationsByUserID(userID int64, limit models.ListRequest) (*models.NotificationsListResponse, error) {
	total, err := s.repo.GetTotalWithUserID(userID)
	if err != nil {
		return nil, err
	}

	valid.Limit(&limit, total)

	numPages := total / limit.Limit
	if total%limit.Limit != 0 {
		numPages++
	}

	if limit.Page >= numPages {
		limit.Page = numPages
	}

	responses, err := s.repo.GetByUserIDWithLimit(userID, limit)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%d notifications found\n", len(responses))

	var notifications = []models.NotificationResponse{}
	for _, r := range responses {
		for h, v := range models.AllowedNotificationTypes {
			if r.Type == h && v {
				switch h {
				case models.NotifyEvent:
					r.Data = r.Data.(models.NotifyEventData)
				case models.NotifyCertificate:
					r.Data = r.Data.(models.NotifyCertificateData)
				default:
					r.Data = nil
				}
			} else {
				r.Data = nil
			}
		}
		notifications = append(
			notifications,
			models.NotificationResponse{
				ID:        r.ID,
				Title:     r.Title,
				Message:   r.Message,
				Type:      r.Type,
				Data:      r.Data,
				CreatedAt: r.CreatedAt,
				IsRead:    r.IsRead,
			},
		)
	}

	return &models.NotificationsListResponse{
		Notifications: notifications,
		Pages:         numPages,
		Total:         total,
	}, nil
}

func (s *NotificationService) MarkAsRead(userId, id int64) (int64, error) {
	return s.repo.MarkAsRead(userId, id)
}

func (s *NotificationService) MarkAllAsRead(userID int64) error {
	return s.repo.MarkAllAsRead(userID)
}

func (s *NotificationService) Delete(userId, id int64) (int64, error) {
	return s.repo.Delete(userId, id)
}

func (s *NotificationService) GetUnreadCount(userID int64) (int, error) {
	return s.repo.GetUnreadCount(userID)
}
