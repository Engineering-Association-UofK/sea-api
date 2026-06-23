package event

import (
	"fmt"
	"log/slog"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/utils/valid"
	"time"
)

func (s *EventService) Apply(userID, eventID int64) (*models.ApplyResponse, error) {
	event, err := s.ValidateRequest(userID, eventID)
	if err != nil {
		return nil, fmt.Errorf("error applying for event (ID: %d, User ID: %d): %w", eventID, userID, err)
	}

	// If the event doesn't need a Form apply directly
	if !event.FormApplication {
		err = s.EventRepo.Apply(userID, eventID)
		if err != nil {
			return nil, err
		}
		return &models.ApplyResponse{
			FormRequired: false,
			FormID:       0,
		}, nil
	}

	eventForm, err := s.EventRepo.GetByEventID(eventID)
	if err != nil {
		return nil, err
	}
	return &models.ApplyResponse{
		FormRequired: true,
		FormID:       eventForm.FormID,
	}, nil
}

func (s *EventService) FormApply(userID, formID, eventID int64) error {
	_, err := s.ValidateRequest(userID, eventID)
	if err != nil {
		return fmt.Errorf("error applying for event with form (ID: %d, User ID: %d): %w", eventID, userID, err)
	}

	if _, err := s.EventRepo.GetByFormID(formID); err == nil {
		return errs.New(errs.Forbidden, fmt.Sprintf("User already applied for form (Form ID: %d, Event ID): %d", formID, eventID), nil)
	}

	return s.EventRepo.Apply(userID, eventID)
}

func (s *EventService) Cancel(userID, eventID int64) error {
	_, err := s.ValidateRequest(userID, eventID)
	if err != nil {
		return fmt.Errorf("error canceling for event application (ID: %d, User ID: %d): %w", eventID, userID, err)
	}

	participant, err := s.EventRepo.GetParticipantByEventAndUserIDs(eventID, userID)
	if err != nil {
		return errs.New(errs.Forbidden, "user is not a participant in that event", nil)
	}

	if participant.Status == models.ACCEPTED {
		return errs.New(errs.BadRequest, "cannot cancel after being accepted", nil)
	}

	return s.EventRepo.Cancel(userID, eventID)
}

func (s *EventService) Status(userID, eventID int64, req models.ListRequest) (*models.ApplicationStatusList, error) {
	_, err := s.ValidateRequest(userID, eventID)
	if err != nil {
		return nil, fmt.Errorf("error applying for event with form (ID: %d, User ID: %d): %w", eventID, userID, err)
	}

	total, err := s.EventRepo.GetTotalApplicationsForUser(userID, eventID)
	if err != nil {
		return nil, errs.New(errs.InternalServerError, "Failed to get total applications for user: "+err.Error(), nil)
	}

	pages := valid.Limit(&req, total)

	applications, err := s.EventRepo.GetApplicationsForUser(userID, eventID, req)
	if err != nil {
		slog.Error("Failed to get application", "error", err.Error())
		return nil, errs.New(errs.NotFound, "Applications not found", nil)
	}

	return &models.ApplicationStatusList{
		Current:      req.Page,
		Pages:        pages,
		Applications: applications,
	}, nil
}

/// Helpers

func (s *EventService) ValidateRequest(userID, eventID int64) (*models.EventModel, error) {
	event, err := s.EventRepo.GetEventByID(eventID)
	if err != nil {
		return nil, errs.New(errs.NotFound, "Event with ID not found", nil)
	}

	participantsCount, err := s.EventRepo.GetParticipantsCount(eventID)
	if err != nil {
		return nil, errs.New(errs.NotFound, "Event with ID not found", nil)
	}

	if participantsCount+1 >= event.MaxParticipants {
		return nil, errs.New(errs.BadRequest, "Event is full", nil)
	}

	if event.StartDate.Before(time.Now()) {
		return nil, errs.New(errs.BadRequest, "Participation time ended", nil)
	}

	return event, nil
}
