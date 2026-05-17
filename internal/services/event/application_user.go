package event

import (
	"log/slog"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/utils/valid"
	"time"
)

func (s *EventService) Apply(userID, eventID int64) (*models.ApplyResponse, error) {
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

func (s *EventService) Cancel(userID, eventID int64) error {
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
	if _, err := s.EventRepo.GetEventByID(eventID); err != nil {
		return nil, errs.New(errs.NotFound, "Event not found", nil)
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
