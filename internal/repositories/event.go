package repositories

import (
	"fmt"
	"sea-api/internal/models"

	"github.com/jmoiron/sqlx"
)

type IEventRepository interface {
	// CREATE
	CreateEvent(event *models.EventModel) (int64, error)
	CreateComponent(component *models.EventComponentModel) (int64, error)
	MassCreateComponent(components []models.EventComponentModel, tx *sqlx.Tx) error
	CreateParticipant(participant *models.EventParticipantModel) (int64, error)
	MassCreateParticipant(participants []models.EventParticipantModel, tx *sqlx.Tx) error
	CreateScore(score *models.ComponentScoreModel) (int64, error)
	MassCreateScore(scores []models.ComponentScoreModel, tx *sqlx.Tx) error

	// UPDATE
	UpdateEvent(event *models.EventModel) error
	UpdateComponent(component *models.EventComponentModel) error
	MassUpdateComponent(components []models.EventComponentModel) error
	UpdateParticipant(participant *models.EventParticipantModel) error
	MassUpdateParticipant(participants []models.EventParticipantModel) error
	UpdateScore(score *models.ComponentScoreModel) error
	MassUpdateScore(scores []models.ComponentScoreModel) error

	// GET One
	GetEventByID(id int64) (*models.EventModel, error)
	GetComponentByID(id int64) (*models.EventComponentModel, error)
	GetParticipantByID(id int64) (*models.EventParticipantModel, error)
	GetParticipantByEventAndUserIDs(eventID int64, user_id int64) (*models.EventParticipantModel, error)
	GetScoreByID(id int64) (*models.ComponentScoreModel, error)
	GetParticipantByUserAndEventIDs(user_id int, eventID int64) (*models.EventParticipantModel, error)

	// GET Many
	GetComponentsByEventID(eventID int64) ([]models.EventComponentModel, error)
	GetParticipantByEventID(eventID int64) ([]models.EventParticipantModel, error)
	GetEligibleParticipantByEventID(eventID int64) ([]models.EventParticipantModel, error)
	GetScoresByParticipantID(participantID int64) ([]models.ComponentScoreModel, error)
	GetParticipantsByEventAndUserIDs(eventID int64, userIDs []int64) ([]models.EventParticipantModel, error)
	GetScoresByParticipantIDs(participantIDs []int64) ([]models.ComponentScoreModel, error)
	GetParticipantByUserID(user_id int) ([]models.EventParticipantModel, error)
	GetAllEvents() ([]models.EventModel, error)

	// DELETE
	DeleteEvent(id int64) error
	DeleteComponent(id int64) error
	MassDeleteComponent(ids []int64) error
	DeleteParticipant(id int64) error
	MassDeleteParticipant(ids []int64) error
	DeleteScore(id int64) error
	MassDeleteScore(ids []int64) error
}

type EventRepository struct {
	db *sqlx.DB
}

func NewEventRepository(db *sqlx.DB) *EventRepository {
	return &EventRepository{db: db}
}

// ======== CREATE NEW MODELS ========

func (r *EventRepository) CreateEvent(event *models.EventModel) (int64, error) {
	query := `
	INSERT INTO event
	(name, description, event_type, max_participants, presenter_id, outcomes, start_date, end_date)
	VALUES (:name, :description, :event_type, :max_participants, :presenter_id, :outcomes, :start_date, :end_date)
	`
	res, err := r.db.NamedExec(query, &event)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *EventRepository) CreateComponent(component *models.EventComponentModel) (int64, error) {
	query := `
	INSERT INTO event_component
	(event_id, name, description, max_score)
	VALUES (:event_id, :name, :description, :max_score)
	`
	res, err := r.db.NamedExec(query, &component)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *EventRepository) MassCreateComponent(components []models.EventComponentModel, tx *sqlx.Tx) error {
	query := `
	INSERT INTO event_component
	(event_id, name, description, max_score)
	VALUES (:event_id, :name, :description, :max_score)
	`
	if tx != nil {
		query = tx.Rebind(query)
		_, err := tx.NamedExec(query, components)
		return err
	}

	query, args, err := sqlx.Named(query, components)
	if err != nil {
		return err
	}

	query = r.db.Rebind(query)
	_, err = r.db.Exec(query, args...)
	return err
}

func (r *EventRepository) CreateParticipant(participant *models.EventParticipantModel) (int64, error) {
	query := `
	INSERT INTO event_participant
	(event_id, user_id, grade, status, joined_at, completed)
	VALUES (:event_id, :user_id, :grade, :status, :joined_at, :completed)
	`
	res, err := r.db.NamedExec(query, &participant)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *EventRepository) MassCreateParticipant(participants []models.EventParticipantModel, tx *sqlx.Tx) error {
	query := `
	INSERT INTO event_participant
	(event_id, user_id, grade, status, joined_at, completed)
	VALUES (:event_id, :user_id, :grade, :status, :joined_at, :completed)
	`
	if tx != nil {
		query = tx.Rebind(query)
		_, err := tx.NamedExec(query, participants)
		return err
	}

	query, args, err := sqlx.Named(query, participants)
	if err != nil {
		return err
	}

	query = r.db.Rebind(query)
	_, err = r.db.Exec(query, args...)
	return err
}

func (r *EventRepository) CreateScore(score *models.ComponentScoreModel) (int64, error) {
	query := `
	INSERT INTO component_score
	(participant_id, component_id, score)
	VALUES (:participant_id, :component_id, :score)
	`
	res, err := r.db.NamedExec(query, &score)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *EventRepository) MassCreateScore(scores []models.ComponentScoreModel, tx *sqlx.Tx) error {
	query := `
	INSERT INTO component_score
	(participant_id, component_id, score)
	VALUES (:participant_id, :component_id, :score)
	`
	if tx != nil {
		query = tx.Rebind(query)
		_, err := tx.NamedExec(query, scores)
		return err
	}

	query, args, err := sqlx.Named(query, scores)
	if err != nil {
		return err
	}

	query = r.db.Rebind(query)
	_, err = r.db.Exec(query, args...)
	return err
}

// ======== UPDATE MODELS ========

func (r *EventRepository) UpdateEvent(event *models.EventModel) error {
	query := `
	UPDATE event
	SET name = :name, description = :description, event_type = :event_type, max_participants = :max_participants,
	presenter_id = :presenter_id, outcomes = :outcomes, start_date = :start_date, end_date = :end_date
	WHERE id = :id
	`
	_, err := r.db.NamedExec(query, &event)
	return err
}

func (r *EventRepository) UpdateComponent(component *models.EventComponentModel) error {
	query := `
	UPDATE event_component
	SET name = :name, description = :description, max_score = :max_score
	WHERE id = :id
	`
	_, err := r.db.NamedExec(query, &component)
	return err
}

func (r *EventRepository) MassUpdateComponent(components []models.EventComponentModel) error {
	if len(components) == 0 {
		return nil
	}

	query := `
	UPDATE event_component
	SET name = :name,
	    description = :description,
	    max_score = :max_score
	WHERE id = :id
	`

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, c := range components {
		if _, err := stmt.Exec(c); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *EventRepository) UpdateParticipant(participant *models.EventParticipantModel) error {
	query := `
	UPDATE event_participant
	SET event_id = :event_id, user_id = :user_id, grade = :grade, status = :status, joined_at = :joined_at, completed = :completed
	WHERE id = :id
	`
	_, err := r.db.NamedExec(query, &participant)
	return err
}

func (r *EventRepository) MassUpdateParticipant(participants []models.EventParticipantModel) error {
	if len(participants) == 0 {
		return nil
	}

	query := `
	UPDATE event_participant
	SET event_id = :event_id,
	    user_id = :user_id,
	    grade = :grade,
	    status = :status,
	    joined_at = :joined_at,
	    completed = :completed
	WHERE id = :id
	`

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, p := range participants {
		if _, err := stmt.Exec(p); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *EventRepository) UpdateScore(score *models.ComponentScoreModel) error {
	query := `
	UPDATE component_score
	SET participant_id = :participant_id, component_id = :component_id, score = :score
	WHERE id = :id
	`
	_, err := r.db.NamedExec(query, &score)
	return err
}

func (r *EventRepository) MassUpdateScore(scores []models.ComponentScoreModel) error {
	if len(scores) == 0 {
		return nil
	}

	query := `
	UPDATE component_score
	SET participant_id = :participant_id,
	    component_id = :component_id,
	    score = :score
	WHERE id = :id
	`

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, s := range scores {
		if _, err := stmt.Exec(s); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// ======== GET BY ID ========

func (r *EventRepository) GetEventByID(id int64) (*models.EventModel, error) {
	var event models.EventModel
	err := r.db.Get(&event, `SELECT * FROM event WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *EventRepository) GetComponentByID(id int64) (*models.EventComponentModel, error) {
	var component models.EventComponentModel
	err := r.db.Get(&component, `SELECT * FROM event_component WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &component, nil
}

func (r *EventRepository) GetParticipantByID(id int64) (*models.EventParticipantModel, error) {
	var participant models.EventParticipantModel
	err := r.db.Get(&participant, `SELECT * FROM event_participant WHERE user_id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &participant, nil
}

func (r *EventRepository) GetParticipantByEventAndUserIDs(eventID int64, user_id int64) (*models.EventParticipantModel, error) {
	var participant models.EventParticipantModel
	err := r.db.Get(&participant, `SELECT * FROM event_participant WHERE event_id = ? AND user_id = ?`, eventID, user_id)
	if err != nil {
		return nil, err
	}
	return &participant, nil
}

func (r *EventRepository) GetScoreByID(id int64) (*models.ComponentScoreModel, error) {
	var score models.ComponentScoreModel
	err := r.db.Get(&score, `SELECT * FROM component_score WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &score, nil
}

// ======== GET BY SPECIFIC ID ========

func (r *EventRepository) GetComponentsByEventID(eventID int64) ([]models.EventComponentModel, error) {
	var components []models.EventComponentModel
	err := r.db.Select(&components, `SELECT * FROM event_component WHERE event_id = ?`, eventID)
	if err != nil {
		return nil, err
	}
	return components, nil
}

func (r *EventRepository) GetParticipantByEventID(eventID int64) ([]models.EventParticipantModel, error) {
	var participant []models.EventParticipantModel
	err := r.db.Select(&participant, `SELECT * FROM event_participant WHERE event_id = ?`, eventID)
	if err != nil {
		return nil, err
	}
	return participant, nil
}

func (r *EventRepository) GetEligibleParticipantByEventID(eventID int64) ([]models.EventParticipantModel, error) {
	query := fmt.Sprintf(`
	SELECT * FROM event_participant
		WHERE event_id = ? 
		AND grade >= 40
		AND completed = true
		AND status = %s
	`, models.COMPLETED)
	var participants []models.EventParticipantModel
	err := r.db.Select(&participants, query, eventID)
	if err != nil {
		return nil, err
	}
	return participants, nil
}

func (r *EventRepository) GetScoresByParticipantID(participantID int64) ([]models.ComponentScoreModel, error) {
	var scores []models.ComponentScoreModel
	err := r.db.Select(&scores, `SELECT * FROM component_score WHERE participant_id = ?`, participantID)
	if err != nil {
		return nil, err
	}
	return scores, nil
}

func (r *EventRepository) GetParticipantsByEventAndUserIDs(eventID int64, userIDs []int64) ([]models.EventParticipantModel, error) {
	if len(userIDs) == 0 {
		return []models.EventParticipantModel{}, nil
	}
	query, args, err := sqlx.In(`SELECT * FROM event_participant WHERE event_id = ? AND user_id IN (?)`, eventID, userIDs)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)
	var participants []models.EventParticipantModel
	err = r.db.Select(&participants, query, args...)
	return participants, err
}

func (r *EventRepository) GetScoresByParticipantIDs(participantIDs []int64) ([]models.ComponentScoreModel, error) {
	if len(participantIDs) == 0 {
		return []models.ComponentScoreModel{}, nil
	}
	query, args, err := sqlx.In(`SELECT * FROM component_score WHERE participant_id IN (?)`, participantIDs)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)
	var scores []models.ComponentScoreModel
	err = r.db.Select(&scores, query, args...)
	return scores, err
}

func (r *EventRepository) GetParticipantByUserID(user_id int) ([]models.EventParticipantModel, error) {
	var participant []models.EventParticipantModel
	err := r.db.Select(&participant, `SELECT * FROM event_participant WHERE user_id = ?`, user_id)
	if err != nil {
		return nil, err
	}
	return participant, nil
}

func (r *EventRepository) GetParticipantByUserAndEventIDs(user_id int, eventID int64) (*models.EventParticipantModel, error) {
	var participant models.EventParticipantModel
	err := r.db.Select(&participant, `SELECT * FROM event_participant WHERE user_id = ? AND event_id = ?`, user_id, eventID)
	if err != nil {
		return nil, err
	}
	return &participant, nil
}

// ======== GET ALL ========

func (r *EventRepository) GetAllEvents() ([]models.EventModel, error) {
	var events []models.EventModel
	err := r.db.Select(&events, `SELECT * FROM event`)
	if err != nil {
		return nil, err
	}
	return events, nil
}

// ======== DELETE ========

func (r *EventRepository) DeleteEvent(id int64) error {
	_, err := r.db.Exec(`DELETE FROM event WHERE id = ?`, id)
	return err
}

func (r *EventRepository) DeleteComponent(id int64) error {
	_, err := r.db.Exec(`DELETE FROM event_component WHERE id = ?`, id)
	return err
}

func (r *EventRepository) MassDeleteComponent(ids []int64) error {
	query, args, err := sqlx.In(`DELETE FROM event_component WHERE id IN (?)`, ids)
	if err != nil {
		return err
	}
	query = r.db.Rebind(query)
	_, err = r.db.Exec(query, args...)
	return err
}

func (r *EventRepository) DeleteParticipant(id int64) error {
	_, err := r.db.Exec(`DELETE FROM event_participant WHERE user_id = ?`, id)
	return err
}

func (r *EventRepository) MassDeleteParticipant(ids []int64) error {
	query, args, err := sqlx.In(`DELETE FROM event_participant WHERE user_id IN (?)`, ids)
	if err != nil {
		return err
	}
	query = r.db.Rebind(query)
	_, err = r.db.Exec(query, args...)
	return err
}

func (r *EventRepository) DeleteScore(id int64) error {
	_, err := r.db.Exec(`DELETE FROM component_score WHERE id = ?`, id)
	return err
}

func (r *EventRepository) MassDeleteScore(ids []int64) error {
	query, args, err := sqlx.In(`DELETE FROM component_score WHERE id IN (?)`, ids)
	if err != nil {
		return err
	}
	query = r.db.Rebind(query)
	_, err = r.db.Exec(query, args...)
	return err
}
