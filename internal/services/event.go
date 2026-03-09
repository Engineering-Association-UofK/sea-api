package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
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
		UserRepo:  repositories.NewUserRepository(db),
	}
}

// ======== GET ========

func (s *EventService) GetEventByID(id int64) (*models.EventDTO, error) {
	event, err := s.EventRepo.GetEventByID(id)
	if err != nil {
		return nil, errors.New("event not found")
	}

	components, err := s.EventRepo.GetComponentsByEventID(id)
	if err != nil {
		return nil, err
	}
	fmt.Println(components[len(components)-1].ID)

	participants, err := s.EventRepo.GetParticipantByEventID(id)
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
		Participants:    s.participantFromModelToDto(participants, components),
	}, nil
}

func (s *EventService) GetAllEvents() ([]models.EventListResponse, error) {
	events, err := s.EventRepo.GetAllEvents()
	if err != nil {
		return nil, err
	}
	eventList := make([]models.EventListResponse, 0)
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
		participants, _ := s.participantFromDtoToModel(event.Participants, id)

		err = s.EventRepo.MassCreateParticipant(participants)
		if err != nil {
			return 0, err
		}

		// Handle Scores
		createdParticipants, err := s.EventRepo.GetParticipantByEventID(id)
		if err != nil {
			return 0, err
		}
		userPartMap := make(map[int64]int64)
		for _, p := range createdParticipants {
			userPartMap[p.UserID] = p.ID
		}

		scores := s.extractScoresFromDTOs(event.Participants, userPartMap)
		if len(scores) > 0 {
			if err := s.EventRepo.MassCreateScore(scores); err != nil {
				return 0, err
			}
		}
	}

	return id, nil
}

// ======== UPDATE ========

func (s *EventService) UpdateEvent(event *models.EventDTO) error {
	if _, err := s.EventRepo.GetEventByID(event.ID); err != nil {
		return errors.New("event not found")
	}

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
		return errors.Join(fmt.Errorf("Unable to update event: "), err)
	}

	components, err := s.EventRepo.GetComponentsByEventID(event.ID)
	if err != nil {
		return errors.Join(fmt.Errorf("Unable to get components: "), err)
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
		return errors.Join(fmt.Errorf("Unable to sync components: "), err)
	}

	participants, err := s.EventRepo.GetParticipantByEventID(event.ID)
	if err != nil {
		return errors.Join(fmt.Errorf("Unable to get participants: "), err)
	}
	if err := syncEntities(
		participants,
		event.Participants,
		func(m models.EventParticipantModel) int64 { return m.ID },
		func(dtos []models.ParticipantDTO, eid int64) []models.EventParticipantModel {
			p, _ := s.participantFromDtoToModel(dtos, eid)
			return p
		},
		s.EventRepo.MassCreateParticipant,
		s.EventRepo.MassUpdateParticipant,
		s.EventRepo.MassDeleteParticipant,
		event.ID,
	); err != nil {
		return errors.Join(fmt.Errorf("Unable to sync participants: "), err)
	}

	// Sync Scores
	currentParticipants, err := s.EventRepo.GetParticipantByEventID(event.ID)
	if err != nil {
		return errors.Join(fmt.Errorf("Unable to get participants for scores: "), err)
	}
	userPartMap := make(map[int64]int64)
	pIDs := make([]int64, len(currentParticipants))
	for i, p := range currentParticipants {
		userPartMap[p.UserID] = p.ID
		pIDs[i] = p.ID
	}

	existingScores, err := s.EventRepo.GetScoresByParticipantIDs(pIDs)
	if err != nil {
		return errors.Join(fmt.Errorf("Unable to get existing scores: "), err)
	}

	newScores := s.extractScoresFromDTOs(event.Participants, userPartMap)

	if err := syncEntities(
		existingScores,
		newScores,
		func(m models.ComponentScoreModel) int64 { return m.ID },
		func(d []models.ComponentScoreModel, _ int64) []models.ComponentScoreModel { return d },
		s.EventRepo.MassCreateScore,
		s.EventRepo.MassUpdateScore,
		s.EventRepo.MassDeleteScore,
		event.ID,
	); err != nil {
		return errors.Join(fmt.Errorf("Unable to sync scores: "), err)
	}

	return nil
}

// ======== DELETE ========

func (s *EventService) DeleteEvent(id int64) error {
	return s.EventRepo.DeleteEvent(id)
}

// ======== HELPERS ========

func (s *EventService) participantFromDtoToModel(dto []models.ParticipantDTO, eventID int64) ([]models.EventParticipantModel, []models.ComponentScoreModel) {
	if len(dto) == 0 {
		return []models.EventParticipantModel{}, []models.ComponentScoreModel{}
	}

	participants := make([]models.EventParticipantModel, len(dto))
	var ScoreModels []models.ComponentScoreModel
	for i, p := range dto {
		grade := 0.0
		for _, ScoreDTO := range p.Grade {
			ScoreModels = append(ScoreModels, models.ComponentScoreModel{
				ParticipantID: p.ID,
				ID:            ScoreDTO.ID,
				ComponentID:   ScoreDTO.ComponentId,
				Score:         ScoreDTO.Score,
			})
			grade += ScoreDTO.Score
		}
		if len(p.Grade) != 0 {
			grade /= float64(len(p.Grade))
		}
		participants[i] = models.EventParticipantModel{
			ID:        p.ID,
			EventID:   eventID,
			UserID:    p.UserID,
			Grade:     grade,
			Status:    p.Status,
			JoinedAt:  p.JoinedAt,
			Completed: p.Completed,
		}
	}

	return participants, ScoreModels
}

func (s *EventService) extractScoresFromDTOs(dtos []models.ParticipantDTO, userPartMap map[int64]int64) []models.ComponentScoreModel {
	var scores []models.ComponentScoreModel
	for _, p := range dtos {
		pid, ok := userPartMap[p.UserID]
		if !ok {
			continue
		}
		for _, sDto := range p.Grade {
			scores = append(scores, models.ComponentScoreModel{
				ID:            sDto.ID,
				ParticipantID: pid,
				ComponentID:   sDto.ComponentId,
				Score:         sDto.Score,
			})
		}
	}
	return scores
}

func (s *EventService) participantFromModelToDto(model []models.EventParticipantModel, allEventComponents []models.EventComponentModel) []models.ParticipantDTO {
	if len(model) == 0 {
		return []models.ParticipantDTO{}
	}

	participantMap := utils.FromSlice(model, func(p models.EventParticipantModel) int64 { return p.ID })

	userIDs := make([]int64, 0, len(model))
	for _, p := range model {
		userIDs = append(userIDs, p.UserID)
	}
	users, err := s.UserRepo.GetAllByIndices(userIDs)
	if err != nil {
		slog.Error("error getting users", "error", err)
		return []models.ParticipantDTO{}
	}
	usersMap := utils.FromSlice(users, func(u models.UserModel) int64 { return u.Index })

	scores, err := s.EventRepo.GetScoresByParticipantIDs(participantMap.Keys())
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			slog.Error("error getting scores", "error", err)
			return []models.ParticipantDTO{}
		}
	}

	scoresByParticipant := utils.NewMpp[int64, utils.Mpp[int64, models.ComponentScoreModel]]()
	for _, score := range scores {
		participantScores := scoresByParticipant.GetOrCreate(score.ParticipantID, func() utils.Mpp[int64, models.ComponentScoreModel] {
			return utils.NewMpp[int64, models.ComponentScoreModel]()
		})
		participantScores[score.ComponentID] = score
	}

	participants := make([]models.ParticipantDTO, len(model))
	for i, m := range model {
		user, _ := usersMap.Value(m.UserID)

		participantScores, _ := scoresByParticipant.Value(m.ID)
		if participantScores == nil {
			participantScores = utils.NewMpp[int64, models.ComponentScoreModel]()
		}

		grades := make([]models.ComScoreDTO, len(allEventComponents))
		for j, component := range allEventComponents {
			score, _ := participantScores.Value(component.ID)
			grades[j] = models.ComScoreDTO{
				ID:          score.ID,
				Name:        component.Name,
				ComponentId: component.ID,
				Score:       score.Score,
			}
		}

		participants[i] = models.ParticipantDTO{
			ID:        m.ID,
			UserID:    m.UserID,
			NameAr:    user.NameAr,
			NameEn:    user.NameEn,
			Grade:     grades,
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
		id := getID(dtoToModel([]DTO{dto}, 0)[0])
		if m.Exists(id) {
			toUpdate = append(toUpdate, dtoToModel([]DTO{dto}, eventId)...)
			_ = m.Delete(id)
		} else {
			toCreate = append(toCreate, dtoToModel([]DTO{dto}, eventId)...)
		}
	}

	if len(toCreate) > 0 {
		if err := createFn(toCreate); err != nil {
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

func (s *EventService) GetGrades(userMap utils.Mpp[int64, models.UserModel], eventId int64) (utils.Mpp[int64, float64], error) {
	components, err := s.EventRepo.GetComponentsByEventID(eventId)
	if err != nil {
		return utils.NewMpp[int64, float64](), err
	}

	divisor := float64(len(components))
	if divisor == 0 {
		divisor = 1
	}

	userIDs := userMap.Keys()
	participants, err := s.EventRepo.GetParticipantsByEventAndUserIDs(eventId, userIDs)
	if err != nil {
		return utils.NewMpp[int64, float64](), err
	}

	partMap := make(map[int64]models.EventParticipantModel)
	for _, p := range participants {
		partMap[p.UserID] = p
	}

	participantIDs := make([]int64, 0, len(participants))
	for _, uid := range userIDs {
		if p, ok := partMap[uid]; ok {
			participantIDs = append(participantIDs, p.ID)
		} else {
			user, _ := userMap.Value(uid)
			return utils.NewMpp[int64, float64](), fmt.Errorf("user %s is not a participant", user.Username)
		}
	}

	scores, err := s.EventRepo.GetScoresByParticipantIDs(participantIDs)
	if err != nil {
		return utils.NewMpp[int64, float64](), err
	}

	scoresMap := make(map[int64][]float64)
	for _, s := range scores {
		scoresMap[s.ParticipantID] = append(scoresMap[s.ParticipantID], s.Score)
	}

	result := utils.NewMpp[int64, float64]()
	for _, p := range participants {
		pScores := scoresMap[p.ID]
		var sum float64
		for _, val := range pScores {
			sum += val
		}
		result.Add(p.UserID, sum/divisor)
	}

	return result, nil
}
