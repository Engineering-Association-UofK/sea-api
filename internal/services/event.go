package services

import (
	"sea-api/internal/models"
	"sea-api/internal/repositories"
	"sea-api/internal/utils"
	"strings"

	"github.com/jmoiron/sqlx"
)

type EventService struct {
	EventRepo *repositories.EventRepository
	UserRepo  *repositories.UserRepository
}

func NewEventService(db *sqlx.DB) *EventService {
	return &EventService{
		EventRepo: repositories.NewEventRepository(db),
	}
}

// ======== GET ========

func (s *EventService) GetEventByID(id int64) (*models.EventDTO, error) {
	event, err := s.EventRepo.GetEventByID(id)
	if err != nil {
		return nil, err
	}

	components, err := s.EventRepo.GetComponentsByEventID(id)
	if err != nil {
		return nil, err
	}

	participants, err := s.EventRepo.GetParticipationByEventID(id)
	if err != nil {
		return nil, err
	}

	return &models.EventDTO{
		ID:              event.ID,
		Name:            event.Name,
		Description:     event.Description,
		EventType:       event.EventType,
		MaxParticipants: event.MaxParticipants,
		StartDate:       event.StartDate,
		EndDate:         event.EndDate,
		Outcomes:        strings.Split(event.Outcomes, ","),
		Components:      s.componentFromModelToDto(components),
		Participants:    s.participantFromModelToDto(participants),
	}, nil
}

func (s *EventService) GetAllEvents() ([]models.EventListResponse, error) {
	events, err := s.EventRepo.GetAllEvents()
	if err != nil {
		return nil, err
	}
	var eventList []models.EventListResponse
	for _, event := range events {
		eventList = append(eventList, models.EventListResponse{
			ID:              event.ID,
			Name:            event.Name,
			EventType:       event.EventType,
			MaxParticipants: event.MaxParticipants,
			StartDate:       event.StartDate,
			EndDate:         event.EndDate,
		})
	}
	return eventList, nil
}

// ======== CREATE ========

func (s *EventService) CreateEvent(event *models.EventDTO) (int64, error) {
	outcomes := strings.Join(event.Outcomes, ",")

	id, err := s.EventRepo.CreateEvent(&models.EventModel{
		Name:            event.Name,
		Description:     event.Description,
		EventType:       event.EventType,
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

		err = s.EventRepo.MassCreateComponent(components)
		if err != nil {
			return 0, err
		}
	}

	if len(event.Participants) != 0 {
		participants := s.participantFromDtoToModel(event.Participants, id)

		err = s.EventRepo.MassCreateParticipation(participants)
		if err != nil {
			return 0, err
		}

		return event.ID, nil
	}

	return id, nil
}

// ======== UPDATE ========

func (s *EventService) UpdateEvent(event *models.EventDTO) error {
	outcomes := strings.Join(event.Outcomes, ",")
	err := s.EventRepo.UpdateEvent(&models.EventModel{
		ID:              event.ID,
		Name:            event.Name,
		Description:     event.Description,
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

	participants, err := s.EventRepo.GetParticipationByEventID(event.ID)
	if err != nil {
		return err
	}
	if err := syncEntities(
		participants,
		event.Participants,
		func(m models.EventParticipationModel) int64 { return m.ID },
		s.participantFromDtoToModel,
		s.EventRepo.MassCreateParticipation,
		s.EventRepo.MassUpdateParticipation,
		s.EventRepo.MassDeleteParticipation,
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

// ======== HELPERS ========

func (s *EventService) participantFromDtoToModel(dto []models.ParticipantDTO, eventID int64) []models.EventParticipationModel {
	participants := make([]models.EventParticipationModel, len(dto))
	for i, p := range dto {
		participants[i] = models.EventParticipationModel{
			ID:        p.ID,
			EventID:   eventID,
			UserID:    p.UserID,
			Grade:     p.Grade,
			Status:    p.Status,
			JoinedAt:  p.JoinedAt,
			Completed: p.Completed,
		}
	}
	return participants
}

func (s *EventService) participantFromModelToDto(model []models.EventParticipationModel) []models.ParticipantDTO {
	indices := make([]int, len(model))
	for i, p := range model {
		indices[i] = p.UserID
	}

	if len(indices) == 0 {
		return []models.ParticipantDTO{}
	}

	users, err := s.UserRepo.GetAllByIndices(indices)
	if err != nil {
		return nil
	}
	usersMap := utils.NewMpp[int, models.UserModel]()
	for _, u := range users {
		usersMap.Add(u.Index, u)
	}

	participants := make([]models.ParticipantDTO, len(model))
	for i, m := range model {
		participants[i] = models.ParticipantDTO{
			ID:        m.ID,
			UserID:    m.UserID,
			NameAr:    usersMap[m.UserID].NameAr,
			NameEn:    usersMap[m.UserID].NameEn,
			Grade:     m.Grade,
			Status:    m.Status,
			JoinedAt:  m.JoinedAt,
			Completed: m.Completed,
		}
	}
	return participants
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
	createFn func([]Model) error,
	updateFn func([]Model) error,
	deleteFn func([]ID) error,
	eventId int64,
) error {
	m := utils.NewMpp[ID, Model]()
	for _, e := range existing {
		m.Add(getID(e), e)
	}

	var toCreate, toUpdate []Model
	for _, dto := range newDTOs {
		id := getID(dtoToModel([]DTO{dto}, 0)[0]) // temporary conversion to get ID
		if m.Exists(id) {
			toUpdate = append(toUpdate, dtoToModel([]DTO{dto}, eventId)...)
			_ = m.Delete(id)
		} else {
			toCreate = append(toCreate, dtoToModel([]DTO{dto}, eventId)...)
		}
	}

	if len(toCreate) > 0 {
		if err := createFn(toCreate); err != nil {
			return err
		}
	}
	if len(toUpdate) > 0 {
		if err := updateFn(toUpdate); err != nil {
			return err
		}
	}
	if m.Len() > 0 {
		if err := deleteFn(m.Keys()); err != nil {
			return err
		}
	}
	return nil
}
