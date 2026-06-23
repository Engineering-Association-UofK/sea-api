package eventrepo

import (
	"fmt"
	"sea-api/internal/models"
)

func (r *EventRepository) LinkForm(req models.EventFormRequest) (int64, error) {
	query := fmt.Sprintf(`
		INSERT INTO %s (form_id, event_id)
		VALUES (:form_id, :event_id)
	`, models.TableEventForms)

	res, err := r.db.NamedExec(query, &req)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *EventRepository) GetByEventID(eventID int64) (*models.EventFormModel, error) {
	var result models.EventFormModel
	err := r.db.Get(&result, fmt.Sprintf(`SELECT * FROM %s Where event_id = ?`, models.TableEventForms), eventID)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *EventRepository) GetByFormID(formID int64) (*models.EventFormModel, error) {
	var result models.EventFormModel
	err := r.db.Get(&result, fmt.Sprintf(`SELECT * FROM %s Where form_id = ?`, models.TableEventForms), formID)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
