package eventrepo

import (
	"fmt"
	"sea-api/internal/models"
	"sea-api/internal/models/eventmodels"
)

func (r *EventRepository) CreateParticipation(participant *eventmodels.EventParticipant) (int64, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (event_id, user_id, joined_at)
	VALUES (:event_id, :user_id, :joined_at)
	`, models.TableEventParticipation)
	res, err := r.db.NamedExec(query, participant)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *EventRepository) GetParticipation(partID int64) (*eventmodels.EventParticipant, error) {
	var model eventmodels.EventParticipant
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, models.TableEventParticipation)
	err := r.db.Get(&model, query, partID)
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *EventRepository) GetParticipationWithEventAndUserIDs(eventID, userID int64) (*eventmodels.EventParticipant, error) {
	var model eventmodels.EventParticipant
	query := fmt.Sprintf(`SELECT * FROM %s WHERE event_id = ? AND user_id = ?`, models.TableEventParticipation)
	err := r.db.Get(&model, query, eventID, userID)
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *EventRepository) GetParticipationByEventID(eventID int64, limit, page int64) ([]eventmodels.EventParticipant, error) {
	offset := (page - 1) * limit
	var list = []eventmodels.EventParticipant{}
	query := fmt.Sprintf(`
	SELECT * FROM %s 
	WHERE event_id = ? 
	ORDER BY joined_at DESC LIMIT ? OFFSET ?
	`, models.TableEventParticipation)

	err := r.db.Select(&list, query, eventID, limit, offset)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *EventRepository) CountParticipationByEventID(eventID int64) (int64, error) {
	var count int64
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE event_id = ?`, models.TableEventParticipation)
	err := r.db.Get(&count, query, eventID)
	return count, err
}

func (r *EventRepository) IsUserParticipant(eventID, userID int64) (bool, error) {
	var exists bool
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE event_id = ? AND user_id = ?)`, models.TableEventParticipation)
	err := r.db.Get(&exists, query, eventID, userID)
	return exists, err
}

func (r *EventRepository) DeleteParticipation(applicationID int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, models.TableEventParticipation)
	_, err := r.db.Exec(query, applicationID)
	return err
}
