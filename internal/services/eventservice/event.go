package eventservice

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/models/eventmodels"
	"sea-api/internal/repositories"
	"sea-api/internal/repositories/eventrepo"
	"sea-api/internal/services"
	"sea-api/internal/services/storage"
	"sea-api/internal/utils/valid"
)

type EventService struct {
	repo           *eventrepo.EventRepository
	formsRepo      *repositories.FormRepository
	s3             *storage.S3
	GalleryService *services.GalleryService
}

func NewEventService(repo *eventrepo.EventRepository, formsRepo *repositories.FormRepository, s3 *storage.S3, galleryService *services.GalleryService) *EventService {
	return &EventService{
		repo:           repo,
		formsRepo:      formsRepo,
		s3:             s3,
		GalleryService: galleryService,
	}
}

func (s *EventService) Create(req *eventmodels.EventRequest) (int64, error) {
	var maxApps sql.NullInt64
	if req.MaxApplications != nil {
		maxApps = sql.NullInt64{Int64: *req.MaxApplications, Valid: true}
	}

	if _, err := s.GalleryService.GetAssetByID(req.BackgroundID); err != nil {
		return 0, errs.New(errs.BadRequest, "invalid image ID provided", nil)
	}

	event := &eventmodels.Event{
		Name:            req.Name,
		Description:     req.Description,
		Belonging:       req.Belonging,
		RequireApplying: req.RequireApplying,
		MaxApplications: maxApps,
		CreatedAt:       time.Now(),
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
	}

	id, err := s.repo.Create(event)
	if err != nil {
		return 0, err
	}

	err = s.GalleryService.AttachAssetToObject(req.BackgroundID, models.ObjEvent, id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *EventService) Update(req *eventmodels.EventUpdateRequest) error {
	event, err := s.repo.Get(req.ID)
	if err != nil {
		slog.Debug("Event not fetched from DB", "Error", err)
		return errs.New(errs.NotFound, "event not found", nil)
	}

	formID := sql.NullInt64{}
	if req.FormID != nil {
		formID.Int64 = *req.FormID
		formID.Valid = *req.FormID != 0
	}

	max := sql.NullInt64{}
	if req.MaxApplications != nil {
		max.Int64 = *req.MaxApplications
		max.Valid = *req.MaxApplications != 0
	}

	if req.BackgroundID != event.BackgroundID {
		if _, err := s.GalleryService.GetAssetByID(req.BackgroundID); err != nil {
			return errs.New(errs.BadRequest, "invalid image ID provided", nil)
		}
		s.GalleryService.RemoveReference(models.ObjPost, event.ID)
		s.GalleryService.AttachAssetToObject(req.BackgroundID, models.ObjPost, req.ID)
	}

	if err := s.repo.Update(&eventmodels.Event{
		ID:              req.ID,
		Name:            req.Name,
		Description:     req.Description,
		BackgroundID:    req.BackgroundID,
		Belonging:       req.Belonging,
		RequireApplying: req.RequireApplying,
		FormID:          formID,
		MaxApplications: max,
		CreatedAt:       event.CreatedAt,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
	}); err != nil {
		return err
	}

	return nil
}

func (s *EventService) Get(id int64) (*eventmodels.EventUpdateRequest, error) {
	event, err := s.repo.Get(id)
	if err != nil {
		return nil, errs.New(errs.NotFound, "event not found", nil)
	}
	return &eventmodels.EventUpdateRequest{
		ID: event.ID,
		EventRequest: eventmodels.EventRequest{
			Name:            event.Name,
			Description:     event.Description,
			BackgroundID:    event.BackgroundID,
			Belonging:       event.Belonging,
			RequireApplying: event.RequireApplying,
			FormID:          &event.FormID.Int64,
			MaxApplications: &event.MaxApplications.Int64,
			StartDate:       event.StartDate,
			EndDate:         event.EndDate,
		},
	}, nil
}

func (s *EventService) GetList(ctx context.Context, req *eventmodels.EventListRequest) (*eventmodels.EventListResponse, error) {
	count, err := s.repo.CountEvents()
	if err != nil {
		return nil, err
	}

	totalPages := valid.Limit(&req.ListRequest, count)

	events, err := s.repo.GetList(req)
	if err != nil {
		return nil, err
	}

	var list = []eventmodels.EventResponse{}
	for _, e := range *events {
		url, err := s.s3.GenerateDownloadUrlByKey(ctx, e.BackgroundKey)
		if err != nil {
			return nil, err
		}

		var formId int64
		if e.FormID.Valid {
			formId = e.FormID.Int64
		}

		var max int64
		if e.MaxApplications.Valid {
			max = e.MaxApplications.Int64
		}

		list = append(list, eventmodels.EventResponse{
			ID:              e.ID,
			Name:            e.Name,
			Description:     e.Description,
			BackgroundID:    e.BackgroundID,
			BackgroundURL:   url,
			Belonging:       e.Belonging,
			RequireApplying: e.RequireApplying,
			FormID:          &formId,
			MaxApplications: &max,
			StartDate:       e.StartDate,
			EndDate:         e.EndDate,
		})
	}

	return &eventmodels.EventListResponse{
		List: list,
		ListResponse: models.ListResponse{
			TotalPages:  totalPages,
			CurrentPage: req.Page,
			Count:       count,
		},
	}, nil
}

func (s *EventService) Delete(id int64) error {
	_, err := s.repo.Get(id)
	if err != nil {
		return errs.New(errs.NotFound, "event not found", nil)
	}

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	return nil
}
