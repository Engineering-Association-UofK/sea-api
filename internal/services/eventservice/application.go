package eventservice

import (
	"context"
	"database/sql"
	"errors"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/models/eventmodels"
	"sea-api/internal/utils/valid"
	"time"
)

func (s *EventService) CheckStatus(eventID, userID int64) (*eventmodels.ApplicationStatus, error) {
	var res eventmodels.ApplicationStatus

	_, err := s.repo.GetParticipationWithEventAndUserIDs(eventID, userID)
	if err == nil {
		res.Applied = true
		res.Accepted = true
		return &res, nil
	}

	application, err := s.repo.GetApplicationByUserAndEvent(eventID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &res, nil
		}
		return nil, err
	}

	res.Applied = true
	res.NeedsForm = true
	res.FormID = &application.FormID

	return &res, nil
}

func (s *EventService) Apply(EventID int64, claims models.ManagedClaims) (*eventmodels.ApplyResponse, error) {
	// Start by checking if the application state already started before
	application, err := s.repo.GetApplicationByUserAndEvent(EventID, claims.UserID)
	if err == nil {
		// User already accepted, no need for any further processing
		if application.Accepted {
			return nil, errs.New(errs.Conflict, "User completed application", nil)
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
		return nil, errs.New(errs.Conflict, "User is a participant", nil)
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

		_, err = s.repo.CreateApplication(&eventmodels.Application{
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
	_, err = s.repo.CreateParticipation(nil, &eventmodels.Participant{
		EventID:  EventID,
		UserID:   claims.UserID,
		JoinedAt: time.Now(),
	})
	return &eventmodels.ApplyResponse{NeedsForm: false}, nil
}

func (s *EventService) GetApplications(ctx context.Context, eventID int64, req *eventmodels.ApplicationListRequest) (*eventmodels.ApplicationListResponse, error) {
	count, err := s.repo.CountApplications(eventID)
	if err != nil {
		return nil, err
	}

	pages := valid.Limit(&req.ListRequest, count)

	apps, err := s.repo.GetApplicationViews(eventID, req.Limit, req.Page)
	if err != nil {
		return nil, err
	}

	list := make([]eventmodels.ApplicationResponse, 0, len(apps))
	for _, app := range apps {
		photoURL := ""
		if app.PhotoKey != nil {
			photoURL, err = s.s3.GenerateDownloadUrlByKey(ctx, *app.PhotoKey)
			if err != nil {
				return nil, err
			}
		}

		list = append(list, eventmodels.ApplicationResponse{
			ID:        app.ID,
			EventID:   app.EventID,
			UserID:    app.UserID,
			Username:  app.Username,
			PhotoURL:  photoURL,
			FormID:    app.FormID,
			Accepted:  app.Accepted,
			StartedAt: app.StartedAt.Format(time.RFC3339),
		})
	}

	return &eventmodels.ApplicationListResponse{
		List: list,
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

	tx, err := s.repo.StartTransaction()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err = s.repo.DeleteApplication(tx, applicationID); err != nil {
		return err
	}

	if _, err = s.repo.CreateParticipation(tx, &eventmodels.Participant{
		EventID:  app.EventID,
		UserID:   app.UserID,
		JoinedAt: time.Now(),
	}); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *EventService) RejectApplication(eventID, applicationID int64) error {
	app, err := s.repo.GetApplication(applicationID)
	if err != nil {
		return err
	}

	if app.EventID != eventID {
		return errs.New(errs.Forbidden, "Applicant does not belong to event", nil)
	}

	return s.repo.DeleteApplication(nil, applicationID)
}
