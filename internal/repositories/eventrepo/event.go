package eventrepo

import (
	"fmt"
	"sea-api/internal/models"
	"sea-api/internal/models/eventmodels"
	"strings"

	"github.com/jmoiron/sqlx"
)

type EventRepository struct {
	db *sqlx.DB
}

func NewEventRepository(db *sqlx.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Create(req *eventmodels.Event) (int64, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (name, description, background_id, belonging, require_applying, form_id, max_applications, created_at, start_date, end_date)
	VALUES (:name, :description, :background_id, :belonging, :require_applying, :form_id, :max_applications, :created_at, :start_date, :end_date)
	`, models.TableNewEvents)
	res, err := r.db.NamedExec(query, &req)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *EventRepository) Update(event *eventmodels.Event) error {
	query := fmt.Sprintf(`
	UPDATE %s SET
		name = :name, 
		description = :description, 
		background_id = :background_id, 
		belonging = :belonging, 
		require_applying = :require_applying,
		form_id = :form_id,
		max_applications = :max_applications, 
		start_date = :start_date,
		created_at = :created_at,
		end_date = :end_date
	WHERE id = :id
	`, models.TableNewEvents)
	_, err := r.db.NamedExec(query, event)
	return err
}

func (r *EventRepository) CountEvents() (int64, error) {
	var count int64
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, models.TableNewEvents)
	err := r.db.Get(&count, query)
	return count, err
}

func (r *EventRepository) Get(id int64) (*eventmodels.Event, error) {
	var model eventmodels.Event
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, models.TableNewEvents)

	err := r.db.Get(&model, query, id)
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *EventRepository) GetView(id int64) (*eventmodels.EventsRow, error) {
	var model eventmodels.EventsRow
	query := fmt.Sprintf(`
	SELECT 
		e.id,
		e.name,
		e.description,
		g.file_key AS background_key,
		e.belonging,
		e.require_applying,
		e.form_id,
		e.max_applications,
		e.created_at,
		e.start_date,
		e.end_date
	FROM %s e
	LEFT JOIN %s g ON e.background_id = g.id
	LEFT JOIN %s f ON g.file_id = f.id
	WHERE id = ?
	`, models.TableNewEvents, models.TableGalleryAssets, models.TableFiles)

	err := r.db.Get(model, query, id)
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *EventRepository) GetList(req *eventmodels.EventListRequest) (*[]eventmodels.EventsRow, error) {
	offset := (req.Page - 1) * req.Limit

	// Base query
	query := fmt.Sprintf(`
	SELECT 
		e.id,
		e.name,
		e.description,
		e.background_id,
		f.file_key AS background_key,
		e.belonging,
		e.require_applying,
		e.form_id,
		e.max_applications,
		e.created_at,
		e.start_date,
		e.end_date
	FROM %s e
	LEFT JOIN %s g ON e.background_id = g.id
	LEFT JOIN %s f ON g.file_id = f.id
	`, models.TableNewEvents, models.TableGalleryAssets, models.TableFiles)

	var conditions []string
	var args []interface{}

	// Handle search
	if req.Search != "" {
		conditions = append(conditions, "(e.name LIKE ? OR e.description LIKE ?)")
		searchTerm := "%" + req.Search + "%"
		args = append(args, searchTerm, searchTerm)
	}

	// Handle belonging enum check
	if req.Belonging != "" {
		conditions = append(conditions, "e.belonging = ?")
		args = append(args, req.Belonging)
	}

	// Append WHERE clause if filters exist
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Append ORDER BY, LIMIT, and OFFSET
	query += " ORDER BY e.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, req.Limit, offset)

	var list = []eventmodels.EventsRow{}
	err := r.db.Select(&list, query, args...)
	if err != nil {
		return nil, err
	}

	return &list, nil
}

func (r *EventRepository) Delete(id int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, models.TableNewEvents)
	_, err := r.db.Exec(query, id)
	return err
}
