package services

import (
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/repositories"
	"time"
)

type FeedbackService struct {
	repo *repositories.FeedbackRepository
}

func NewFeedbackService(repo *repositories.FeedbackRepository) *FeedbackService {
	return &FeedbackService{repo: repo}
}

func (s *FeedbackService) Create(Message string, UserID *int64, Type models.FeedbackType) (int64, error) {
	if !models.AllowedFeedbackTypes[Type] {
		return 0, errs.New(errs.NotFound, "feedback type not found", nil)
	}
	feedback := &models.Feedback{
		Message:   Message,
		UserID:    UserID,
		Type:      Type,
		CreatedAt: time.Now(),
	}
	return s.repo.Create(feedback)
}

func (s *FeedbackService) GetByID(id int64) (*models.Feedback, error) {
	return s.repo.GetByID(id)
}

func (s *FeedbackService) GetAll(req *models.ListRequest) ([]models.Feedback, error) {
	return s.repo.GetAll(req)
}

func (s *FeedbackService) GetAllByType(fType models.FeedbackType, req *models.ListRequest) ([]models.Feedback, error) {
	return s.repo.GetByType(fType, req)
}

func (s *FeedbackService) Delete(id int64) error {
	return s.repo.Delete(id)
}
