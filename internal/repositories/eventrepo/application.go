package eventrepo

import (
	"fmt"
	"sea-api/internal/models"
	"sea-api/internal/models/eventmodels"
)

func (r *EventRepository) CreateApplication(app *eventmodels.EventApplication) (int64, error) {
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

func (r *EventRepository) GetApplication(applicationID int64) (*eventmodels.EventApplication, error) {
	var app eventmodels.EventApplication
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, models.TableEventApplication)
	err := r.db.Get(&app, query, applicationID)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *EventRepository) GetApplicationByUserAndEvent(eventID, userID int64) (*eventmodels.EventApplication, error) {
	var app eventmodels.EventApplication
	query := fmt.Sprintf(`SELECT * FROM %s WHERE event_id = ? AND user_id = ?`, models.TableEventApplication)
	err := r.db.Get(&app, query, eventID, userID)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *EventRepository) GetApplicationListByEventID(eventID int64, limit, page int64) ([]eventmodels.EventApplication, error) {
	offset := (page - 1) * limit
	var list = []eventmodels.EventApplication{}
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

func (r *EventRepository) UpdateApplicationStatus(id int64, accepted bool) error {
	query := fmt.Sprintf(`UPDATE %s SET Accepted = ? WHERE id = ?`, models.TableEventApplication)
	_, err := r.db.Exec(query, accepted, id)
	return err
}

func (r *EventRepository) CountApplications(eventID int64) (int64, error) {
	var count int64
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE event_id = ?`, models.TableEventApplication)
	err := r.db.Get(&count, query, eventID)
	return count, err
}

func (r *EventRepository) DeleteApplication(id int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, models.TableEventApplication)
	_, err := r.db.Exec(query, id)
	return err
}
