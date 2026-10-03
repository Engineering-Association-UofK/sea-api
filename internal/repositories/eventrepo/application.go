package eventrepo

import (
	"fmt"
	"sea-api/internal/models"
	"sea-api/internal/models/eventmodels"

	"github.com/jmoiron/sqlx"
)

func (r *EventRepository) CreateApplication(app *eventmodels.Application) (int64, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (event_id, user_id, form_id, Accepted, started_at)
	VALUES (:event_id, :user_id, :form_id, :Accepted, :started_at)
	`, models.TableEventApplication)
	res, err := r.db.NamedExec(query, app)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *EventRepository) GetApplication(applicationID int64) (*eventmodels.Application, error) {
	var app eventmodels.Application
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, models.TableEventApplication)
	err := r.db.Get(&app, query, applicationID)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *EventRepository) GetApplicationByUserAndEvent(eventID, userID int64) (*eventmodels.Application, error) {
	var app eventmodels.Application
	query := fmt.Sprintf(`SELECT * FROM %s WHERE event_id = ? AND user_id = ?`, models.TableEventApplication)
	err := r.db.Get(&app, query, eventID, userID)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *EventRepository) GetApplicationListByEventID(eventID int64, limit, page int64) ([]eventmodels.Application, error) {
	offset := (page - 1) * limit
	var list = []eventmodels.Application{}
	query := fmt.Sprintf(`
	SELECT * FROM %s 
	WHERE event_id = ? 
	ORDER BY started_at DESC LIMIT ? OFFSET ?
	`, models.TableEventApplication)

	err := r.db.Select(&list, query, eventID, limit, offset)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *EventRepository) GetApplicationViews(eventID int64, limit, page int64) ([]eventmodels.ApplicationRaw, error) {
	offset := (page - 1) * limit
	var list = []eventmodels.ApplicationRaw{}
	query := fmt.Sprintf(`
	SELECT
		a.id,
		a.event_id,
		a.form_id,
		a.user_id,
		u.username,
		f.file_key AS photo_key,
		a.Accepted,
		a.started_at
	FROM %s a
	JOIN %s u ON a.user_id = u.id
	LEFT JOIN %s f ON u.profile_image_id = f.id
	WHERE a.event_id = ?
	ORDER BY a.started_at DESC LIMIT ? OFFSET ?
	`, models.TableEventApplication, models.TableUsers, models.TableFiles)

	err := r.db.Select(&list, query, eventID, limit, offset)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *EventRepository) CountApplications(eventID int64) (int64, error) {
	var count int64
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE event_id = ?`, models.TableEventApplication)
	err := r.db.Get(&count, query, eventID)
	return count, err
}

func (r *EventRepository) DeleteApplication(tx *sqlx.Tx, id int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, models.TableEventApplication)
	_, err := tx.Exec(query, id)
	return err
}
