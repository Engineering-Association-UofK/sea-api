package event

import (
	"fmt"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/repositories"
	"sea-api/internal/repositories/eventrepo"
	"sea-api/internal/services"
	"sea-api/internal/utils"
	"sea-api/internal/utils/valid"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type EventService struct {
	NotificationService *services.NotificationService

	EventRepo  *eventrepo.EventRepository
	CollabRepo *repositories.CollaboratorRepo
	FormRepo   *repositories.FormRepository
	UserRepo   *repositories.UserRepository
}

func NewEventService(
	NotificationService *services.NotificationService,
	EventRepo *eventrepo.EventRepository,
	CollabRepo *repositories.CollaboratorRepo,
	FormRepo *repositories.FormRepository,
	UserRepo *repositories.UserRepository,
) *EventService {
	return &EventService{
		NotificationService: NotificationService,
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

	events, err := s.EventRepo.GetEventList(query)
	if err != nil {
		return models.EventViewListResponse{}, err
	}

	return models.EventViewListResponse{
		Current: query.Page,
		Pages:   pages,
		Events:  events,
	}, nil
}

func (s *EventService) GetEventViewDetails(id int64) (*models.EventViewDetailsResponse, error) {
	resp, err := s.EventRepo.GetEventDetails(id)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *EventService) GetEventByID(id int64) (*models.EventModel, error) {
	return s.EventRepo.GetEventByID(id)
}

// ======== CREATE ========

func (s *EventService) CreateEvent(event *models.EventCreateRequest) (int64, error) {
	if _, err := s.CollabRepo.GetByID(event.PresenterID); err != nil {
		return 0, errs.New(errs.NotFound, "presenter not found", nil)
	}
	if event.StartDate.After(event.EndDate) {
		return 0, errs.New(errs.BadRequest, "start date must be before end date", nil)
	}
	if event.EndDate.After(time.Now()) {
		return 0, errs.New(errs.BadRequest, "end date must be in the future", nil)
	}

	outcomes := strings.Join(event.Outcomes, ",")

	id, err := s.EventRepo.CreateEvent(&models.EventModel{
		Name:            event.Name,
		Description:     event.Description,
		EventType:       event.EventType,
		PresenterID:     event.PresenterID,
		FormApplication: event.FormApplication,
		MaxParticipants: event.MaxParticipants,
		StartDate:       event.StartDate,
		EndDate:         event.EndDate,
		Outcomes:        outcomes,
	})
	if err != nil {
		return 0, err
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
	if _, err := s.EventRepo.GetEventByID(event.ID); err != nil {
		return err
	}

	outcomes := strings.Join(event.Outcomes, ",")
	err := s.EventRepo.UpdateEvent(&models.EventModel{
		ID:              event.ID,
		Name:            event.Name,
		Description:     event.Description,
		PresenterID:     event.PresenterID,
		EventType:       event.EventType,
		MaxParticipants: event.MaxParticipants,
		StartDate:       event.StartDate,
		EndDate:         event.EndDate,
		Outcomes:        outcomes,
	})
	if err != nil {
		return err
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
