package event

import (
	"sea-api/internal/errs"
	"sea-api/internal/models"
)

func (s *EventService) LinkForm(req models.EventFormRequest) (int64, error) {
	if _, err := s.EventRepo.GetEventByID(req.EventID); err != nil {
		return 0, errs.New(errs.NotFound, "Event with ID not found", nil)
	}

	if _, err := s.FormRepo.GetFormByID(req.FormID); err != nil {
		return 0, errs.New(errs.NotFound, "Form with ID not found", nil)
	}

	return s.EventRepo.LinkForm(req)
}

func (s *EventService) GetFormID(eventID int64) (*models.EventFormModel, error) {
	if _, err := s.EventRepo.GetEventByID(eventID); err != nil {
		return nil, errs.New(errs.NotFound, "Event with ID not found", nil)
	}

	return s.EventRepo.GetByFormID(eventID)
}
