package repositories

import (
	"encoding/json"
	"fmt"
	"log/slog"
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
	query := fmt.Sprintf(`
	INSERT INTO %s
	(name, description, event_type, max_participants, presenter_id, outcomes, start_date, end_date)
	VALUES (:name, :description, :event_type, :max_participants, :presenter_id, :outcomes, :start_date, :end_date)
	`, models.TableEvents)
	res, err := r.db.NamedExec(query, &event)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *EventRepository) CreateComponent(component *models.EventComponentModel) (int64, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s
	(event_id, name, description, max_score)
	VALUES (:event_id, :name, :description, :max_score)
	`, models.TableEventComponents)
	res, err := r.db.NamedExec(query, &component)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *EventRepository) MassCreateComponent(components []models.EventComponentModel, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`
	INSERT INTO %s
	(event_id, name, description, max_score)
	VALUES (:event_id, :name, :description, :max_score)
	`, models.TableEventComponents)
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
	query := fmt.Sprintf(`
	INSERT INTO %s
	(event_id, user_id, grade, status, joined_at, completed)
	VALUES (:event_id, :user_id, :grade, :status, :joined_at, :completed)
	`, models.TableEventParticipants)
	res, err := r.db.NamedExec(query, &participant)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *EventRepository) MassCreateParticipant(participants []models.EventParticipantModel, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`
	INSERT INTO %s
	(event_id, user_id, grade, status, joined_at, completed)
	VALUES (:event_id, :user_id, :grade, :status, :joined_at, :completed)
	`, models.TableEventParticipants)
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
	query := fmt.Sprintf(`
	INSERT INTO %s
	(participant_id, component_id, score)
	VALUES (:participant_id, :component_id, :score)
	`, models.TableComponentScores)
	res, err := r.db.NamedExec(query, &score)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *EventRepository) MassCreateScore(scores []models.ComponentScoreModel, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`
	INSERT INTO %s
	(participant_id, component_id, score)
	VALUES (:participant_id, :component_id, :score)
	`, models.TableComponentScores)
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
	query := fmt.Sprintf(`
	UPDATE %s
	SET name = :name, description = :description, event_type = :event_type, max_participants = :max_participants,
	presenter_id = :presenter_id, outcomes = :outcomes, start_date = :start_date, end_date = :end_date
	WHERE id = :id
	`, models.TableEvents)
	_, err := r.db.NamedExec(query, &event)
	return err
}

func (r *EventRepository) UpdateComponent(component *models.EventComponentModel) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET name = :name, description = :description, max_score = :max_score
	WHERE id = :id
	`, models.TableEventComponents)
	_, err := r.db.NamedExec(query, &component)
	return err
}

func (r *EventRepository) MassUpdateComponent(components []models.EventComponentModel) error {
	if len(components) == 0 {
		return nil
	}

	query := fmt.Sprintf(`
	UPDATE %s
	SET name = :name,
	    description = :description,
	    max_score = :max_score
	WHERE id = :id
	`, models.TableEventComponents)

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
	query := fmt.Sprintf(`
	UPDATE %s
	SET event_id = :event_id, user_id = :user_id, grade = :grade, status = :status, joined_at = :joined_at, completed = :completed
	WHERE id = :id
	`, models.TableEventParticipants)
	_, err := r.db.NamedExec(query, &participant)
	return err
}

func (r *EventRepository) MassUpdateParticipant(participants []models.EventParticipantModel) error {
	if len(participants) == 0 {
		return nil
	}

	query := fmt.Sprintf(`
	UPDATE %s
	SET event_id = :event_id,
	    user_id = :user_id,
	    grade = :grade,
	    status = :status,
	    joined_at = :joined_at,
	    completed = :completed
	WHERE id = :id
	`, models.TableEventParticipants)

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
	query := fmt.Sprintf(`
	UPDATE %s
	SET participant_id = :participant_id, component_id = :component_id, score = :score
	WHERE id = :id
	`, models.TableComponentScores)
	_, err := r.db.NamedExec(query, &score)
	return err
}

func (r *EventRepository) MassUpdateScore(scores []models.ComponentScoreModel) error {
	if len(scores) == 0 {
		return nil
	}

	query := fmt.Sprintf(`
	UPDATE %s
	SET participant_id = :participant_id,
	    component_id = :component_id,
	    score = :score
	WHERE id = :id
	`, models.TableComponentScores)

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

func (r *EventRepository) GetEventDetails(eventID int64, userID *int64) (*models.EventViewDetailsResponse, error) {
	var resp models.EventViewDetailsResponse
	var presentersJSON, gradingJSON, userRegJSON, outcomesJSON []byte

	query := fmt.Sprintf(`
	SELECT 
		e.id, 
		e.name, 
		e.description, 
		e.event_type, 
		e.start_date, 
		e.end_date,
		-- Outcomes: Assuming stored as JSON or simple string in TEXT
		e.outcomes, 
		-- Aggregate Presenters (handling the slice in your struct)
		(SELECT JSON_ARRAYAGG(JSON_OBJECT('id', c.id, 'name', c.name_ar))
		FROM %s c WHERE c.id = e.presenter_id) AS presenters_json,
		-- Aggregate Grading Scheme
		(SELECT JSON_ARRAYAGG(JSON_OBJECT(
			'id', ec.id, 
			'name', ec.name, 
			'description', ec.description, 
			'max_score', ec.max_score
		))
		FROM %s ec WHERE ec.event_id = e.id) AS grading_json
	FROM 
		%s e
	`, models.TableCollaborators, models.TableEventComponents, models.TableEvents)

	if userID != nil {
		query += fmt.Sprintf(`
		-- User Registration Status (Optional based on logged-in user)
		JSON_OBJECT(
			'is_registered', IF(ep.id IS NOT NULL, 1, 0),
			'status', ep.status
		) AS user_reg_json

		LEFT JOIN %s ep ON e.id = ep.event_id AND ep.user_id = %d
		`, models.TableEventParticipants, userID)
	}

	query += ` WHERE e.id = ?`

	if userID != nil {
		err := r.db.QueryRow(query, eventID).Scan(
			&resp.ID,
			&resp.Name,
			&resp.Description,
			&resp.EventType,
			&resp.Schedule.Start,
			&resp.Schedule.End,
			&outcomesJSON,
			&presentersJSON,
			&gradingJSON,
			&userRegJSON,
		)
		if err != nil {
			return nil, err
		}
	} else {
		err := r.db.QueryRow(query, eventID).Scan(
			&resp.ID,
			&resp.Name,
			&resp.Description,
			&resp.EventType,
			&resp.Schedule.Start,
			&resp.Schedule.End,
			&outcomesJSON,
			&presentersJSON,
			&gradingJSON,
		)
		if err != nil {
			return nil, err
		}
	}

	json.Unmarshal(outcomesJSON, &resp.Outcomes)
	json.Unmarshal(presentersJSON, &resp.Presenters)
	json.Unmarshal(gradingJSON, &resp.GradingScheme)

	var regStatus models.UserRegStatus
	if err := json.Unmarshal(userRegJSON, &regStatus); err == nil && regStatus.IsRegistered {
		resp.UserRegistration = &regStatus
	}

	return &resp, nil
}

func (r *EventRepository) GetEventByID(id int64) (*models.EventModel, error) {
	var event models.EventModel
	err := r.db.Get(&event, fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, models.TableEvents), id)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *EventRepository) GetComponentByID(id int64) (*models.EventComponentModel, error) {
	var component models.EventComponentModel
	err := r.db.Get(&component, fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, models.TableEventComponents), id)
	if err != nil {
		return nil, err
	}
	return &component, nil
}

func (r *EventRepository) GetParticipantByID(id int64) (*models.EventParticipantModel, error) {
	var participant models.EventParticipantModel
	err := r.db.Get(&participant, fmt.Sprintf(`SELECT * FROM %s WHERE user_id = ?`, models.TableEventParticipants), id)
	if err != nil {
		return nil, err
	}
	return &participant, nil
}

func (r *EventRepository) GetParticipantByEventAndUserIDs(eventID int64, user_id int64) (*models.EventParticipantModel, error) {
	var participant models.EventParticipantModel
	err := r.db.Get(&participant, fmt.Sprintf(`SELECT * FROM %s WHERE event_id = ? AND user_id = ?`, models.TableEventParticipants), eventID, user_id)
	if err != nil {
		return nil, err
	}
	return &participant, nil
}

func (r *EventRepository) GetScoreByID(id int64) (*models.ComponentScoreModel, error) {
	var score models.ComponentScoreModel
	err := r.db.Get(&score, fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, models.TableComponentScores), id)
	if err != nil {
		return nil, err
	}
	return &score, nil
}

// ======== GET BY SPECIFIC ID ========

func (r *EventRepository) GetComponentsByEventID(eventID int64) ([]models.EventComponentModel, error) {
	var components []models.EventComponentModel
	err := r.db.Select(&components, fmt.Sprintf(`SELECT * FROM %s WHERE event_id = ?`, models.TableEventComponents), eventID)
	if err != nil {
		return nil, err
	}
	return components, nil
}

func (r *EventRepository) GetParticipantByEventID(eventID int64) ([]models.EventParticipantModel, error) {
	var participant []models.EventParticipantModel
	err := r.db.Select(&participant, fmt.Sprintf(`SELECT * FROM %s WHERE event_id = ?`, models.TableEventParticipants), eventID)
	if err != nil {
		return nil, err
	}
	return participant, nil
}

func (r *EventRepository) GetEligibleParticipantByEventID(eventID int64) ([]models.EventParticipantModel, error) {
	query := fmt.Sprintf(`
	SELECT * FROM %s
		WHERE event_id = ? 
		AND grade >= 40
		AND completed = true
		AND status = %s
	`, models.COMPLETED, models.TableEventParticipants)
	var participants []models.EventParticipantModel
	err := r.db.Select(&participants, query, eventID)
	if err != nil {
		return nil, err
	}
	return participants, nil
}

func (r *EventRepository) GetScoresByParticipantID(participantID int64) ([]models.ComponentScoreModel, error) {
	var scores []models.ComponentScoreModel
	err := r.db.Select(&scores, fmt.Sprintf(`SELECT * FROM %s WHERE participant_id = ?`, models.TableComponentScores), participantID)
	if err != nil {
		return nil, err
	}
	return scores, nil
}

func (r *EventRepository) GetParticipantsByEventAndUserIDs(eventID int64, userIDs []int64) ([]models.EventParticipantModel, error) {
	if len(userIDs) == 0 {
		return []models.EventParticipantModel{}, nil
	}
	query, args, err := sqlx.In(fmt.Sprintf(`SELECT * FROM %s WHERE event_id = ? AND user_id IN (?)`, models.TableEventParticipants), eventID, userIDs)
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
	query, args, err := sqlx.In(fmt.Sprintf(`SELECT * FROM %s WHERE participant_id IN (?)`, models.TableComponentScores), participantIDs)
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
	err := r.db.Select(&participant, fmt.Sprintf(`SELECT * FROM %s WHERE user_id = ?`, models.TableEventParticipants), user_id)
	if err != nil {
		return nil, err
	}
	return participant, nil
}

func (r *EventRepository) GetParticipantByUserAndEventIDs(user_id int, eventID int64) (*models.EventParticipantModel, error) {
	var participant models.EventParticipantModel
	err := r.db.Select(&participant, fmt.Sprintf(`SELECT * FROM %s WHERE user_id = ? AND event_id = ?`, models.TableEventParticipants), user_id, eventID)
	if err != nil {
		return nil, err
	}
	return &participant, nil
}

// ======== GET ALL ========

func (r *EventRepository) GetEventList(req models.QueryEventPublicRequest) ([]models.EventViewListItemResponse, error) {
	var events = []models.EventViewListItemResponse{}
	query := fmt.Sprintf(`
	SELECT 
		e.id,
		e.name,
		e.presenter_id,
		e.event_type,
		e.max_participants,
		e.start_date,
		e.end_date,
		(SELECT COUNT(*) FROM %s p WHERE p.event_id = e.id) as participants_count,
		CASE 
			WHEN NOW() < e.start_date THEN 'UPCOMING'
			WHEN NOW() BETWEEN e.start_date AND e.end_date THEN 'ONGOING'
			ELSE 'COMPLETED'
		END as status
	FROM %s e
	WHERE 1=1
	`, models.TableEventParticipants, models.TableEvents)

	var args []interface{}
	if req.Type != "" {
		query += " AND e.event_type = ?"
		args = append(args, req.Type)
	}

	if req.Status != "" {
		query += ` AND (
			CASE 
				WHEN NOW() < e.start_date THEN 'UPCOMING'
				WHEN NOW() BETWEEN e.start_date AND e.end_date THEN 'ONGOING'
				ELSE 'COMPLETED'
			END
		) = ?`
		args = append(args, req.Status)
	}

	query += " ORDER BY e.start_date DESC LIMIT ? OFFSET ?"
	args = append(args, req.Limit, (req.Page-1)*req.Limit)

	err := r.db.Select(&events, query, args...)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (r *EventRepository) GetEventParticipantsList(req models.ListRequest, eventID int64) ([]models.ParticipantDetails, error) {
	// 1. Fixed the semicolon and the table alias (ep vs p)
	query := fmt.Sprintf(`
    SELECT 
        ep.id AS registration_id,
        ep.user_id,
        u.name_ar,
        u.name_en,
        ep.joined_at,
        ep.status,
        ep.completed,
        COALESCE(
            (SELECT JSON_ARRAYAGG(
                JSON_OBJECT(
                    'component_id', cs.component_id,
                    'score', cs.score
                )
            )
            FROM %s cs
            WHERE cs.participant_id = ep.id), 
            JSON_ARRAY()
        ) AS grades_json
    FROM 
        %s ep
    JOIN 
        %s u ON ep.user_id = u.id
    WHERE 
        ep.event_id = ?
    ORDER BY ep.joined_at DESC 
    LIMIT ? OFFSET ?
    `, models.TableComponentScores, models.TableEventParticipants, models.TableUsers)

	rows, err := r.db.Query(query, eventID, req.Limit, (req.Page-1)*req.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants = []models.ParticipantDetails{}
	for rows.Next() {
		var p models.ParticipantDetails
		var gradesRaw []byte

		err := rows.Scan(
			&p.RegistrationID,
			&p.UserID,
			&p.NameAr,
			&p.NameEn,
			&p.JoinedAt,
			&p.Status,
			&p.Completed,
			&gradesRaw,
		)
		if err != nil {
			return nil, err
		}

		if len(gradesRaw) > 0 {
			if err := json.Unmarshal(gradesRaw, &p.Grades); err != nil {
				slog.Error("failed to unmarshal grades", "error", err)
			}
		}

		participants = append(participants, p)
	}

	return participants, nil
}

// func (r *EventRepository) GetAllEvents(req models.ListRequest) ([]models.EventModel, error) {
// 	var events []models.EventModel
// 	query := fmt.Sprintf(`SELECT * FROM %s ORDER BY start_date DESC LIMIT ? OFFSET ?`, models.TableEvents)
// 	err := r.db.Select(&events, query, req.Limit, (req.Page-1)*req.Limit)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return events, nil
// }

func (r *EventRepository) GetTotalEvents() (int64, error) {
	var total int64
	err := r.db.Get(&total, fmt.Sprintf(`SELECT COUNT(*) FROM %s`, models.TableEvents))
	if err != nil {
		return 0, err
	}
	return total, nil
}

// ======== DELETE ========

func (r *EventRepository) DeleteEvent(id int64) error {
	_, err := r.db.Exec(fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, models.TableEvents), id)
	return err
}

func (r *EventRepository) DeleteComponent(id int64) error {
	_, err := r.db.Exec(fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, models.TableEventComponents), id)
	return err
}

func (r *EventRepository) MassDeleteComponent(ids []int64) error {
	query, args, err := sqlx.In(fmt.Sprintf(`DELETE FROM %s WHERE id IN (?)`, models.TableEventComponents), ids)
	if err != nil {
		return err
	}
	query = r.db.Rebind(query)
	_, err = r.db.Exec(query, args...)
	return err
}

func (r *EventRepository) DeleteParticipant(id int64) error {
	_, err := r.db.Exec(fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, models.TableEventParticipants), id)
	return err
}

func (r *EventRepository) MassDeleteParticipant(ids []int64) error {
	query, args, err := sqlx.In(fmt.Sprintf(`DELETE FROM %s WHERE id IN (?)`, models.TableEventParticipants), ids)
	if err != nil {
		return err
	}
	query = r.db.Rebind(query)
	_, err = r.db.Exec(query, args...)
	return err
}

func (r *EventRepository) DeleteScore(id int64) error {
	_, err := r.db.Exec(fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, models.TableComponentScores), id)
	return err
}

func (r *EventRepository) MassDeleteScore(ids []int64) error {
	query, args, err := sqlx.In(fmt.Sprintf(`DELETE FROM %s WHERE id IN (?)`, models.TableComponentScores), ids)
	if err != nil {
		return err
	}
	query = r.db.Rebind(query)
	_, err = r.db.Exec(query, args...)
	return err
}
