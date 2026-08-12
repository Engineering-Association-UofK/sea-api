package event

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/repositories"
	"sea-api/internal/repositories/eventrepo"
	"sea-api/internal/services"
	"sea-api/internal/services/storage"
	"sea-api/internal/utils"
	"sea-api/internal/utils/valid"
	"strings"

	"github.com/jmoiron/sqlx"
)

type EventService struct {
	NotificationService *services.NotificationService

	S3Service  *storage.S3
	Gallery    *services.GalleryService
	EventRepo  *eventrepo.EventRepository
	CollabRepo *repositories.CollaboratorRepo
	FormRepo   *repositories.FormRepository
	UserRepo   *repositories.UserRepository
}

func NewEventService(
	NotificationService *services.NotificationService,
	S3Service *storage.S3,
	Gallery *services.GalleryService,
	EventRepo *eventrepo.EventRepository,
	CollabRepo *repositories.CollaboratorRepo,
	FormRepo *repositories.FormRepository,
	UserRepo *repositories.UserRepository,
) *EventService {
	return &EventService{
		NotificationService: NotificationService,
		S3Service:           S3Service,
		Gallery:             Gallery,
		EventRepo:           EventRepo,
		CollabRepo:          CollabRepo,
		FormRepo:            FormRepo,
		UserRepo:            UserRepo,
	}
}

// ======== GET ========

func (s *EventService) GetViewList(query models.QueryEventPublicRequest) (models.EventViewListResponse, error) {
	total, err := s.EventRepo.GetTotalEvents()
	if err != nil {
		return models.EventViewListResponse{}, err
	}

	limit := models.ListRequest{Page: query.Page, Limit: query.Limit}
	pages := valid.Limit(&limit, total)

	query.Limit = limit.Limit
	query.Page = limit.Page

	eventRows, err := s.EventRepo.GetEventList(query)
	if err != nil {
		return models.EventViewListResponse{}, err
	}

	events := make([]models.EventViewListItemResponse, len(eventRows))
	for i, row := range eventRows {
		image := ""
		if row.WallpaperFileKey.Valid {
			image, err = s.S3Service.GenerateDownloadUrlByKey(context.Background(), row.WallpaperFileKey.String)
			if err != nil {
				return models.EventViewListResponse{}, err
			}
		}
		events[i] = models.EventViewListItemResponse{
			ID:                row.ID,
			Name:              row.Name,
			PresenterID:       row.PresenterID,
			EventType:         row.EventType,
			Wallpaper:         image,
			MaxParticipants:   row.MaxParticipants,
			ParticipantsCount: row.ParticipantsCount,
			Status:            row.Status,
			StartDate:         row.StartDate,
			EndDate:           row.EndDate,
		}
	}

	return models.EventViewListResponse{
		Current: query.Page,
		Pages:   pages,
		Events:  events,
	}, nil
}

func (s *EventService) GetEventViewDetails(ctx context.Context, id int64) (*models.EventViewDetailsResponse, error) {
	resp, err := s.EventRepo.GetEventDetails(id)
	if err != nil {
		return nil, err
	}

	if resp.Wallpaper.Valid {
		resp.Wallpaper.String, err = s.S3Service.GenerateDownloadUrlByKey(ctx, resp.Wallpaper.String)
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func (s *EventService) GetEventByID(id int64) (*models.EventModel, error) {
	return s.EventRepo.GetEventByID(id)
}

func (s *EventService) GetEventDetailsAdmin(ctx context.Context, id int64) (*models.EventDetailsAdminResponse, error) {
	resp, err := s.EventRepo.GetEventDetailsAdmin(id)
	if err != nil {
		return nil, err
	}

	url := ""
	if resp.WallpaperID.Valid {
		url, err = s.S3Service.GenerateDownloadUrlByKey(ctx, resp.WallpaperFileKey.String)
		if err != nil {
			return nil, err
		}
	}

	var presenters = []models.PresenterSummary{}
	var coordinators = []models.CoordinatorSummary{}
	var grading = []models.ComponentDTO{}

	json.Unmarshal(resp.PresentersJSON, &presenters)
	json.Unmarshal(resp.CoordinatorsJSON, &coordinators)
	json.Unmarshal(resp.GradingJSON, &grading)

	return &models.EventDetailsAdminResponse{
		ID:                resp.ID,
		Name:              resp.Name,
		Description:       resp.Description,
		EventType:         resp.EventType,
		Outcomes:          strings.Split(resp.Outcomes, ","),
		WallpaperID:       resp.WallpaperID.Int64,
		Wallpaper:         url,
		FormApplication:   resp.FormApplication,
		FormID:            resp.FormID.Int64,
		MaxParticipants:   resp.MaxParticipants,
		ParticipantsCount: resp.ParticipantsCount,
		Schedule: models.EventSchedule{
			Start: resp.StartDate,
			End:   resp.EndDate,
		},
		Presenters:   presenters,
		Coordinators: coordinators,

		GradingScheme: grading,
	}, nil
}

// ======== CREATE ========

func (s *EventService) CreateEvent(event *models.EventCreateRequest) (int64, error) {
	if err := s.validateEventRequest(event); err != nil {
		return 0, err
	}

	var ID sql.NullInt64
	if event.WallpaperID != 0 {
		ID.Valid = true
		ID.Int64 = event.WallpaperID
	} else {
		ID.Valid = false
	}

	outcomes := strings.Join(event.Outcomes, ",")

	id, err := s.EventRepo.CreateEvent(&models.EventModel{
		Name:            event.Name,
		Description:     event.Description,
		EventType:       event.EventType,
		WallpaperID:     ID,
		PresenterID:     event.PresenterID,
		CoordinatorID:   event.CoordinatorID,
		FormApplication: event.FormApplication,
		MaxParticipants: event.MaxParticipants,
		StartDate:       event.StartDate,
		EndDate:         event.EndDate,
		Outcomes:        outcomes,
	})
	if err != nil {
		return 0, err
	}

	err = s.Gallery.AttachAssetToObject(event.WallpaperID, models.ObjEvent, id)
	if err != nil {
		slog.Error("Failed to attach asset to event object", "Error", err, "Object ID", id, "Asset ID", event.WallpaperID)
	}

	if len(event.Components) != 0 {
		components := s.componentFromDtoToModel(event.Components, id)

		err = s.EventRepo.MassCreateComponent(components, nil)
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}

// // ======== UPDATE ========

func (s *EventService) UpdateEvent(event *models.EventUpdateRequest) error {
	eventModel, err := s.EventRepo.GetEventByID(event.ID)
	if err != nil {
		return err
	}

	if err := s.validateEventRequest(&event.EventCreateRequest); err != nil {
		return err
	}

	var ID sql.NullInt64
	if event.WallpaperID != 0 {
		ID.Valid = true
		ID.Int64 = event.WallpaperID
	} else {
		ID.Valid = false
	}

	outcomes := strings.Join(event.Outcomes, ",")
	err = s.EventRepo.UpdateEvent(&models.EventModel{
		ID:              event.ID,
		Name:            event.Name,
		Description:     event.Description,
		PresenterID:     event.PresenterID,
		CoordinatorID:   event.CoordinatorID,
		EventType:       event.EventType,
		WallpaperID:     ID,
		FormApplication: event.FormApplication,
		MaxParticipants: event.MaxParticipants,
		StartDate:       event.StartDate,
		EndDate:         event.EndDate,
		Outcomes:        outcomes,
	})
	if err != nil {
		return err
	}

	if event.WallpaperID != eventModel.WallpaperID.Int64 {
		if err = s.Gallery.RemoveReference(models.ObjEvent, event.ID); err != nil {
			slog.Error("Failed to remove reference from asset to event object", "Error", err, "Object ID", event.ID, "Asset ID", event.WallpaperID)
		}

		if err = s.Gallery.AttachAssetToObject(event.WallpaperID, models.ObjEvent, event.ID); err != nil && event.WallpaperID != 0 {
			slog.Error("Failed to attach asset to event object", "Error", err, "Object ID", event.ID, "Asset ID", event.WallpaperID)
		}
	}

	components, err := s.EventRepo.GetComponentsByEventID(event.ID)
	if err != nil {
		return err
	}

	if err := syncEntities(
		components,
		event.Components,
		func(m models.EventComponentModel) int64 { return m.ID },
		s.componentFromDtoToModel,
		s.EventRepo.MassCreateComponent,
		s.EventRepo.MassUpdateComponent,
		s.EventRepo.MassDeleteComponent,
		event.ID,
	); err != nil {
		return err
	}

	return nil
}

// ======== DELETE ========

func (s *EventService) DeleteEvent(id int64) error {
	return s.EventRepo.DeleteEvent(id)
}

// ======== HELPERS =======

func (s *EventService) validateEventRequest(event *models.EventCreateRequest) error {
	if _, err := s.CollabRepo.GetByID(event.PresenterID); err != nil {
		return errs.New(errs.NotFound, "presenter not found", nil)
	}
	if _, err := s.CollabRepo.GetByID(event.CoordinatorID); err != nil {
		return errs.New(errs.NotFound, "coordinator not found", nil)
	}
	if event.StartDate.After(event.EndDate) {
		return errs.New(errs.BadRequest, "start date must be before end date", nil)
	}
	if _, err := s.Gallery.GetAssetByID(event.WallpaperID); err != nil && event.WallpaperID != 0 {
		return errs.New(errs.NotFound, "wallpaper image not found", nil)
	}
	return nil
}

func (s *EventService) componentFromDtoToModel(dto []models.ComponentDTO, eventID int64) []models.EventComponentModel {
	components := make([]models.EventComponentModel, len(dto))
	for i, c := range dto {
		components[i] = models.EventComponentModel{
			ID:          c.ID,
			EventID:     eventID,
			Name:        c.Name,
			Description: c.Description,
			MaxScore:    c.MaxScore,
		}
	}
	return components
}

func (s *EventService) componentFromModelToDto(model []models.EventComponentModel) []models.ComponentDTO {
	components := make([]models.ComponentDTO, len(model))
	for i, c := range model {
		components[i] = models.ComponentDTO{
			ID:          c.ID,
			Name:        c.Name,
			Description: c.Description,
			MaxScore:    c.MaxScore,
		}
	}
	return components
}

func syncEntities[ID comparable, Model any, DTO any](
	existing []Model,
	newDTOs []DTO,
	getID func(Model) ID,
	dtoToModel func([]DTO, int64) []Model,
	createFn func([]Model, *sqlx.Tx) error,
	updateFn func([]Model) error,
	deleteFn func([]ID) error,
	eventId int64,
) error {
	m := utils.Mpp[ID, Model]{}
	for _, e := range existing {
		m.Add(getID(e), e)
	}

	var toCreate, toUpdate []Model
	for _, dto := range newDTOs {
		id := getID(dtoToModel([]DTO{dto}, 0)[0])
		if m.Exists(id) {
			toUpdate = append(toUpdate, dtoToModel([]DTO{dto}, eventId)...)
			_ = m.Delete(id)
		} else {
			toCreate = append(toCreate, dtoToModel([]DTO{dto}, eventId)...)
		}
	}

	if len(toCreate) > 0 {
		if err := createFn(toCreate, nil); err != nil {
			return fmt.Errorf("Error Creating: %s", err)
		}
	}
	if len(toUpdate) > 0 {
		if err := updateFn(toUpdate); err != nil {
			return fmt.Errorf("Error Updating: %s", err)
		}
	}
	if m.Len() > 0 {
		if err := deleteFn(m.Keys()); err != nil {
			return fmt.Errorf("Error Deleting: %s", err)
		}
	}
	return nil
}
