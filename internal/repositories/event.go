package repositories

import (
	"sea-api/internal/models"

	"github.com/jmoiron/sqlx"
)

type EventRepository struct {
	db *sqlx.DB
}

func NewEventRepository(db *sqlx.DB) *EventRepository {
	return &EventRepository{db: db}
}

// ======== CREATE NEW MODELS ========

func (r *EventRepository) CreateEvent(event *models.EventModel) (int64, error) {
	query := `
	INSERT INTO events
	(name, description, event_type, max_participants, outcomes, start_date, end_date)
	VALUES (:name, :description, :event_type, :max_participants, :outcomes, :start_date, :end_date)
	`
	res, err := r.db.NamedExec(query, &event)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *EventRepository) CreateComponent(component *models.EventComponentModel) (int64, error) {
	query := `
	INSERT INTO event_components
	(event_id, name, description, max_score)
	VALUES (:event_id, :name, :description, :max_score)
	`
	res, err := r.db.NamedExec(query, &component)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *EventRepository) MassCreateComponent(components []models.EventComponentModel) error {
	query := `
	INSERT INTO event_components
	(event_id, name, description, max_score)
	VALUES (:event_id, :name, :description, :max_score)
	`
	query, args, err := sqlx.Named(query, components)
	if err != nil {
		return err
	}

	query = r.db.Rebind(query)
	_, err = r.db.Exec(query, args...)
	return err
}

func (r *EventRepository) CreateParticipation(participation *models.EventParticipationModel) (int64, error) {
	query := `
	INSERT INTO event_participation
	(event_id, user_id, grade, status, joined_at, completed)
	VALUES (:event_id, :user_id, :grade, :status, :joined_at, :completed)
	`
	res, err := r.db.NamedExec(query, &participation)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *EventRepository) MassCreateParticipation(participations []models.EventParticipationModel) error {
	query := `
	INSERT INTO event_participation
	(event_id, user_id, grade, status, joined_at, completed)
	VALUES (:event_id, :user_id, :grade, :status, :joined_at, :completed)
	`
	query, args, err := sqlx.Named(query, participations)
	if err != nil {
		return err
	}

	query = r.db.Rebind(query)
	_, err = r.db.Exec(query, args...)
	return err
}

func (r *EventRepository) CreateScore(score *models.ComponentScoreModel) (int64, error) {
	query := `
	INSERT INTO component_scores
	(participation_id, component_id, score)
	VALUES (:participation_id, :component_id, :score)
	`
	res, err := r.db.NamedExec(query, &score)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *EventRepository) MassCreateScore(scores []models.ComponentScoreModel) error {
	query := `
	INSERT INTO component_scores
	(participation_id, component_id, score)
	VALUES (:participation_id, :component_id, :score)
	`
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
	UPDATE events
	SET name = :name, description = :description, event_type = :event_type, max_participants = :max_participants, 
	outcomes = :outcomes, start_date = :start_date, end_date = :end_date
	WHERE id = :id
	`
	_, err := r.db.NamedExec(query, &event)
	return err
}

func (r *EventRepository) UpdateComponent(component *models.EventComponentModel) error {
	query := `
	UPDATE event_components
	SET event_id = :event_id, name = :name, description = :description, max_score = :max_score
	WHERE id = :id
	`
	_, err := r.db.NamedExec(query, &component)
	return err
}

func (r *EventRepository) MassUpdateComponent(components []models.EventComponentModel) error {
	query := `
	UPDATE event_components
	SET event_id = :event_id, name = :name, description = :description, max_score = :max_score
	WHERE id = :id
	`
	_, err := r.db.NamedExec(query, components)
	return err
}

func (r *EventRepository) UpdateParticipation(participation *models.EventParticipationModel) error {
	query := `
	UPDATE event_participation
	SET event_id = :event_id, user_id = :user_id, grade = :grade, status = :status, joined_at = :joined_at, completed = :completed
	WHERE id = :id
	`
	_, err := r.db.NamedExec(query, &participation)
	return err
}

func (r *EventRepository) MassUpdateParticipation(participations []models.EventParticipationModel) error {
	query := `
	UPDATE event_participation
	SET event_id = :event_id, user_id = :user_id, grade = :grade, status = :status, joined_at = :joined_at, completed = :completed
	WHERE id = :id
	`
	_, err := r.db.NamedExec(query, participations)
	return err
}

func (r *EventRepository) UpdateScore(score *models.ComponentScoreModel) error {
	query := `
	UPDATE component_scores
	SET participation_id = :participation_id, component_id = :component_id, score = :score
	WHERE id = :id
	`
	_, err := r.db.NamedExec(query, &score)
	return err
}

func (r *EventRepository) MassUpdateScore(scores []models.ComponentScoreModel) error {
	query := `
	UPDATE component_scores
	SET participation_id = :participation_id, component_id = :component_id, score = :score
	WHERE id = :id
	`
	_, err := r.db.NamedExec(query, scores)
	return err
}

// ======== GET BY ID ========

func (r *EventRepository) GetEventByID(id int64) (*models.EventModel, error) {
	var event models.EventModel
	err := r.db.Get(&event, `SELECT * FROM events WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *EventRepository) GetComponentByID(id int64) (*models.EventComponentModel, error) {
	var component models.EventComponentModel
	err := r.db.Get(&component, `SELECT * FROM event_components WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &component, nil
}

func (r *EventRepository) GetParticipationByID(id int64) (*models.EventParticipationModel, error) {
	var participation models.EventParticipationModel
	err := r.db.Get(&participation, `SELECT * FROM event_participation WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &participation, nil
}

func (r *EventRepository) GetScoreByID(id int64) (*models.ComponentScoreModel, error) {
	var score models.ComponentScoreModel
	err := r.db.Get(&score, `SELECT * FROM component_scores WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &score, nil
}

// ======== GET BY SPECIFIC ID ========

func (r *EventRepository) GetComponentsByEventID(eventID int64) ([]models.EventComponentModel, error) {
	var components []models.EventComponentModel
	err := r.db.Select(&components, `SELECT * FROM event_components WHERE event_id = ?`, eventID)
	if err != nil {
		return nil, err
	}
	return components, nil
}

func (r *EventRepository) GetParticipationByEventID(eventID int64) ([]models.EventParticipationModel, error) {
	var participation []models.EventParticipationModel
	err := r.db.Select(&participation, `SELECT * FROM event_participation WHERE event_id = ?`, eventID)
	if err != nil {
		return nil, err
	}
	return participation, nil
}

func (r *EventRepository) GetScoresByParticipationID(participationID int64) ([]models.ComponentScoreModel, error) {
	var scores []models.ComponentScoreModel
	err := r.db.Select(&scores, `SELECT * FROM component_scores WHERE participation_id = ?`, participationID)
	if err != nil {
		return nil, err
	}
	return scores, nil
}

func (r *EventRepository) GetParticipationByUserID(userID int) ([]models.EventParticipationModel, error) {
	var participation []models.EventParticipationModel
	err := r.db.Select(&participation, `SELECT * FROM event_participation WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	return participation, nil
}

func (r *EventRepository) GetParticipationByUserAndEventIDs(userID int, eventID int64) (*models.EventParticipationModel, error) {
	var participation models.EventParticipationModel
	err := r.db.Select(&participation, `SELECT * FROM event_participation WHERE user_id = ? AND event_id = ?`, userID, eventID)
	if err != nil {
		return nil, err
	}
	return &participation, nil
}

// ======== GET ALL ========

func (r *EventRepository) GetAllEvents() ([]models.EventModel, error) {
	var events []models.EventModel
	err := r.db.Select(&events, `SELECT * FROM events`)
	if err != nil {
		return nil, err
	}
	return events, nil
}

// ======== DELETE ========

func (r *EventRepository) DeleteEvent(id int64) error {
	_, err := r.db.Exec(`DELETE FROM events WHERE id = ?`, id)
	return err
}

func (r *EventRepository) DeleteComponent(id int64) error {
	_, err := r.db.Exec(`DELETE FROM event_components WHERE id = ?`, id)
	return err
}

func (r *EventRepository) MassDeleteComponent(ids []int64) error {
	query, args, err := sqlx.In(`DELETE FROM event_components WHERE id IN (?)`, ids)
	if err != nil {
		return err
	}
	query = r.db.Rebind(query)
	_, err = r.db.Exec(query, args...)
	return err
}

func (r *EventRepository) DeleteParticipation(id int64) error {
	_, err := r.db.Exec(`DELETE FROM event_participation WHERE id = ?`, id)
	return err
}

func (r *EventRepository) MassDeleteParticipation(ids []int64) error {
	query, args, err := sqlx.In(`DELETE FROM event_participation WHERE id IN (?)`, ids)
	if err != nil {
		return err
	}
	query = r.db.Rebind(query)
	_, err = r.db.Exec(query, args...)
	return err
}

func (r *EventRepository) DeleteScore(id int64) error {
	_, err := r.db.Exec(`DELETE FROM component_scores WHERE id = ?`, id)
	return err
}

func (r *EventRepository) MassDeleteScore(ids []int64) error {
	query, args, err := sqlx.In(`DELETE FROM component_scores WHERE id IN (?)`, ids)
	if err != nil {
		return err
	}
	query = r.db.Rebind(query)
	_, err = r.db.Exec(query, args...)
	return err
}
