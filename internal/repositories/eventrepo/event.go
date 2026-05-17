package eventrepo

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
	(name, description, event_type, form_application, max_participants, presenter_id, outcomes, start_date, end_date)
	VALUES (:name, :description, :event_type, :form_application, :max_participants, :presenter_id, :outcomes, :start_date, :end_date)
	`, models.TableEvents)
	res, err := r.db.NamedExec(query, &event)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// ======== UPDATE MODELS ========

func (r *EventRepository) UpdateEvent(event *models.EventModel) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET name = :name, description = :description, event_type = :event_type, form_application = :form_application, max_participants = :max_participants,
	presenter_id = :presenter_id, outcomes = :outcomes, start_date = :start_date, end_date = :end_date
	WHERE id = :id
	`, models.TableEvents)
	_, err := r.db.NamedExec(query, &event)
	return err
}

// ======== GET BY ID ========

func (r *EventRepository) GetEventDetails(eventID int64) (*models.EventViewDetailsResponse, error) {
	var resp models.EventViewDetailsResponse
	var presentersJSON, gradingJSON, userRegJSON, outcomesJSON []byte

	query := fmt.Sprintf(`
	SELECT 
		e.id, 
		e.name, 
		e.description, 
		e.form_application, 
		e.max_participants, 
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
	FROM %s e 
	WHERE e.id = ?
	`, models.TableCollaborators, models.TableEventComponents, models.TableEvents)

	slog.Debug("User ID not present, Starting scan")
	err := r.db.QueryRow(query, eventID).Scan(
		&resp.ID,
		&resp.Name,
		&resp.Description,
		&resp.FormApplication,
		&resp.MaxParticipants,
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

	resp.Outcomes = []string{}
	resp.Presenters = []models.PresenterSummary{}
	resp.GradingScheme = []models.ComponentDTO{}

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

// ======== APPLICATIONS ========
// ======== DELETE ========

func (r *EventRepository) DeleteEvent(id int64) error {
	_, err := r.db.Exec(fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, models.TableEvents), id)
	return err
}
