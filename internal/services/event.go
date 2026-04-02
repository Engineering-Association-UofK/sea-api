package services

import (
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"sea-api/internal/config"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/repositories"
	"sea-api/internal/utils"
	"sea-api/internal/utils/sheets"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type EventService struct {
	EventRepo repositories.IEventRepository
	UserRepo  repositories.IUserRepository
}

func NewEventService(EventRepo repositories.IEventRepository, UserRepo repositories.IUserRepository) *EventService {
	return &EventService{
		EventRepo: EventRepo,
		UserRepo:  UserRepo,
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

	participants, err := s.EventRepo.GetParticipantByEventID(id)
	if err != nil {
		return nil, err
	}

	return &models.EventDTO{
		ID:              event.ID,
		Name:            event.Name,
		Description:     event.Description,
		PresenterID:     event.PresenterID,
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
			PresenterID:     event.PresenterID,
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
		PresenterID:     event.PresenterID,
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

	if len(event.Participants) != 0 {
		participants, _ := s.participantFromDtoToModel(event.Participants, id)

		err = s.EventRepo.MassCreateParticipant(participants, nil)
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
			if err := s.EventRepo.MassCreateScore(scores, nil); err != nil {
				return 0, err
			}
		}
	}

	return id, nil
}

// ======== UPDATE ========

func (s *EventService) UpdateEvent(event *models.EventDTO) error {
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
		return errs.New(errs.InternalServerError, err.Error(), nil)
	}

	components, err := s.EventRepo.GetComponentsByEventID(event.ID)
	if err != nil {
		return errs.New(errs.InternalServerError, err.Error(), nil)
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
		return errs.New(errs.InternalServerError, err.Error(), nil)
	}

	participants, err := s.EventRepo.GetParticipantByEventID(event.ID)
	if err != nil {
		return errs.New(errs.InternalServerError, err.Error(), nil)
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
		return errs.New(errs.InternalServerError, err.Error(), nil)
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

// ======== SPECIAL ========

func (s *EventService) ImportUsers(eventID int64, file io.Reader) error {
	users, err := sheets.ParseExcelToStructs[models.EventUsersImport](file)
	if err != nil {
		return err
	}

	ids := utils.ExtractField(users, func(u models.EventUsersImport) int64 {
		index, _ := strconv.ParseInt(u.Index, 10, 64)
		return index
	})
	existing, err := s.UserRepo.GetAllByIndices(ids)
	if err != nil {
		return err
	}

	existingMap := utils.FromSlice(existing, func(u models.UserModel) int64 { return u.ID })

	tx, err := s.UserRepo.GetTransaction()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var participants []models.EventParticipantModel
	for _, u := range users {
		index, err := strconv.ParseInt(u.Index, 10, 64)
		if err != nil {
			return err
		}
		if _, ok := existingMap[index]; !ok {
			username := sha512.Sum512([]byte(u.NameEn + "|" + config.App.SecretSalt))
			p, _ := generatePasscode(8)
			pass, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
			err = s.UserRepo.Create(&models.UserModel{
				ID:         index,
				UniID:      0,
				Username:   hex.EncodeToString(username[:]),
				NameEn:     u.NameEn,
				NameAr:     u.NameAr,
				Email:      u.Email,
				Phone:      "",
				Department: "",
				Verified:   false,
				Password:   string(pass),
				Status:     models.STATUS_INACTIVE,
				Gender:     models.MALE,
			}, tx)
			if err != nil {
				return err
			}

			err = s.UserRepo.DeleteTempUser(index, tx)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return err
			}
		}

		grade := 0.0
		if g, err := strconv.ParseFloat(u.Grade, 64); err == nil {
			grade = g
		}

		participants = append(participants, models.EventParticipantModel{
			EventID:   eventID,
			UserID:    index,
			Status:    models.COMPLETED,
			Grade:     grade,
			JoinedAt:  time.Now(),
			Completed: true,
		})
	}

	if len(participants) > 0 {
		err = s.EventRepo.MassCreateParticipant(participants, tx)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// ======== HELPERS ========

func (s *EventService) participantFromDtoToModel(dto []models.ParticipantDTO, eventID int64) ([]models.EventParticipantModel, []models.ComponentScoreModel) {
	if len(dto) == 0 {
		return []models.EventParticipantModel{}, []models.ComponentScoreModel{}
	}
	components, err := s.EventRepo.GetComponentsByEventID(eventID)
	if err != nil {
		slog.Error("error getting components", "error", err)
		return []models.EventParticipantModel{}, []models.ComponentScoreModel{}
	}

	participants := make([]models.EventParticipantModel, len(dto))
	var ScoreModels []models.ComponentScoreModel
	for i, p := range dto {
		gradeMap := utils.FromSlice(p.Grade, func(dto models.ComScoreDTO) int64 { return dto.ComponentId })
		grade := 0.0
		maxGrade := 0.0
		for _, Component := range components {
			score := gradeMap.GetOrCreate(Component.ID, func() models.ComScoreDTO {
				return models.ComScoreDTO{
					ID:          0,
					Name:        Component.Name,
					ComponentId: Component.ID,
					Score:       0,
				}
			})
			if score.Score > Component.MaxScore {
				score.Score = Component.MaxScore
			}
			ScoreModels = append(ScoreModels, models.ComponentScoreModel{
				ParticipantID: p.ID,
				ID:            score.ID,
				ComponentID:   score.ComponentId,
				Score:         score.Score,
			})
			maxGrade += Component.MaxScore
			grade += score.Score
		}
		if len(p.Grade) != 0 {
			grade = grade / maxGrade * 100
		}
		t := time.Now()
		if !p.JoinedAt.IsZero() {
			t = p.JoinedAt
		}
		s := models.ACCEPTED
		if p.Status != "" {
			s = p.Status
		}
		participants[i] = models.EventParticipantModel{
			ID:        p.ID,
			EventID:   eventID,
			UserID:    p.UserID,
			Grade:     grade,
			Status:    s,
			JoinedAt:  t,
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
	usersMap := utils.FromSlice(users, func(u models.UserModel) int64 { return u.ID })

	scores, err := s.EventRepo.GetScoresByParticipantIDs(participantMap.Keys())
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			slog.Error("error getting scores", "error", err)
			return []models.ParticipantDTO{}
		}
	}

	scoresByParticipant := utils.Mpp[int64, utils.Mpp[int64, models.ComponentScoreModel]]{}
	for _, score := range scores {
		participantScores := scoresByParticipant.GetOrCreate(score.ParticipantID, func() utils.Mpp[int64, models.ComponentScoreModel] {
			return utils.Mpp[int64, models.ComponentScoreModel]{}
		})
		participantScores[score.ComponentID] = score
	}

	participants := make([]models.ParticipantDTO, len(model))
	for i, m := range model {
		user, _ := usersMap.Value(m.UserID)

		participantScores, _ := scoresByParticipant.Value(m.ID)
		if participantScores == nil {
			participantScores = utils.Mpp[int64, models.ComponentScoreModel]{}
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
