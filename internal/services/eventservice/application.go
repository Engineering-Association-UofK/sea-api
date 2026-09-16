package eventservice

import (
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/models/eventmodels"
	"sea-api/internal/utils/valid"
	"time"
)

func (s *EventService) Apply(EventID int64, claims models.ManagedClaims) (*eventmodels.ApplyResponse, error) {
	// Start by checking if the application state already started before
	application, err := s.repo.GetApplicationByUserAndEvent(EventID, claims.UserID)
	if err == nil {
		// User already accepted, no need for any further processing
		if application.Accepted {
			return nil, errs.New(errs.Conflict, "Already a participant", nil)
		}

		// User still in application process, send back form ID
		return &eventmodels.ApplyResponse{
			NeedsForm: true,
			FormID:    &application.FormID,
		}, nil
	}

	// Then check for participation in case the event does not require a form
	_, err = s.repo.GetParticipationWithEventAndUserIDs(EventID, claims.UserID)
	if err == nil {
		return nil, errs.New(errs.Conflict, "Already a participant", nil)
	}

	event, err := s.repo.Get(EventID)
	if err != nil {
		return nil, err
	}

	if !event.RequireApplying {
		return nil, errs.New(errs.Forbidden, "This event is not joinable", nil)
	}

	// If form ID is valid, the event requires form application to join
	if event.FormID.Valid {
		_, err := s.formsRepo.GetFormByID(event.FormID.Int64)
		if err != nil {
			return nil, err
		}

		_, err = s.repo.CreateApplication(&eventmodels.EventApplication{
			EventID:   EventID,
			UserID:    claims.UserID,
			FormID:    event.FormID.Int64,
			Accepted:  false,
			StartedAt: time.Now(),
		})

		return &eventmodels.ApplyResponse{
			NeedsForm: true,
			FormID:    &event.FormID.Int64,
		}, nil
	}

	// Otherwise, make instantly as a participant
	_, err = s.repo.CreateParticipation(&eventmodels.EventParticipant{
		EventID:  EventID,
		UserID:   claims.UserID,
		JoinedAt: time.Now(),
	})
	return &eventmodels.ApplyResponse{NeedsForm: false}, nil
}

func (s *EventService) ProcessApplication(eventID, userID int64, accept bool) error {
	app, err := s.repo.GetApplicationByUserAndEvent(eventID, userID)
	if err != nil {
		return errs.New(errs.NotFound, "application not found", nil)
	}

	if err := s.repo.UpdateApplicationStatus(app.ID, accept); err != nil {
		return err
	}

	if accept {
		_, err := s.repo.CreateParticipation(&eventmodels.EventParticipant{
			EventID:  eventID,
			UserID:   userID,
			JoinedAt: time.Now(),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *EventService) GetApplications(eventID int64, req *eventmodels.EventApplicationListRequest) (*eventmodels.EventApplicationListResponse, error) {
	count, err := s.repo.CountApplications(eventID)
	if err != nil {
		return nil, err
	}

	pages := valid.Limit(&req.ListRequest, count)

	apps, err := s.repo.GetApplicationListByEventID(eventID, req.Limit, req.Page)
	if err != nil {
		return nil, err
	}

	return &eventmodels.EventApplicationListResponse{
		List: apps,
		ListResponse: models.ListResponse{
			TotalPages:  pages,
			CurrentPage: req.Page,
			Count:       count,
		},
	}, nil
}

func (s *EventService) AcceptApplication(eventID, applicationID int64) error {
	app, err := s.repo.GetApplication(applicationID)
	if err != nil {
		return err
	}

	if app.EventID != eventID {
		return errs.New(errs.Forbidden, "Applicant does not belong to event", nil)
	}

	if app.Accepted {
		return errs.New(errs.Conflict, "Applicant Already accepted", nil)
	}

	return s.repo.UpdateApplicationStatus(applicationID, true)
}

func (s *EventService) RejectApplication(eventID, applicationID int64) error {
	app, err := s.repo.GetApplication(applicationID)
	if err != nil {
		return err
	}

	if app.EventID != eventID {
		return errs.New(errs.Forbidden, "Applicant does not belong to event", nil)
	}

	return s.repo.DeleteApplication(applicationID)
}
