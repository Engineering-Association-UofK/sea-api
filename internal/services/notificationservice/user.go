package notificationservice

import (
	"sea-api/internal/models"
	"sea-api/internal/models/notificationsmodels"
	"sea-api/internal/utils/valid"
)

func (s *NotificationService) GetNotificationsByUserID(userID int64, req notificationsmodels.NotificationListRequest) (*notificationsmodels.NotificationsListResponse, error) {
	total, err := s.repo.GetTotalWithUserID(userID)
	if err != nil {
		return nil, err
	}

	pages := valid.Limit(&req.ListRequest, total)

	responses, err := s.repo.GetByUserIDWithLimit(userID, req.ListRequest)
	if err != nil {
		return nil, err
	}

	var notifications = []notificationsmodels.NotificationResponse{}
	for _, r := range responses {
		notifications = append(
			notifications,
			notificationsmodels.NotificationResponse{
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

	return &notificationsmodels.NotificationsListResponse{
		List: notifications,
		ListResponse: models.ListResponse{
			TotalPages:  pages,
			CurrentPage: req.Page,
			Count:       total,
		},
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
