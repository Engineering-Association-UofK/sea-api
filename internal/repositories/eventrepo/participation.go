package eventrepo

import (
	"fmt"
	"sea-api/internal/models/eventmodels"
	"sea-api/internal/models/tables"

	"github.com/jmoiron/sqlx"
)

func (r *EventRepository) CreateParticipation(tx *sqlx.Tx, participant *eventmodels.Participant) (int64, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (event_id, user_id, joined_at)
	VALUES (:event_id, :user_id, :joined_at)
	`, tables.EventParticipation)

	if tx != nil {
		res, err := tx.NamedExec(query, participant)
		if err != nil {
			return 0, err
		}
		return res.LastInsertId()
	}

	res, err := r.db.NamedExec(query, participant)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *EventRepository) GetParticipation(partID int64) (*eventmodels.Participant, error) {
	var model eventmodels.Participant
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, tables.EventParticipation)
	err := r.db.Get(&model, query, partID)
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *EventRepository) GetParticipationWithEventAndUserIDs(eventID, userID int64) (*eventmodels.Participant, error) {
	var model eventmodels.Participant
	query := fmt.Sprintf(`SELECT * FROM %s WHERE event_id = ? AND user_id = ?`, tables.EventParticipation)
	err := r.db.Get(&model, query, eventID, userID)
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *EventRepository) GetParticipationByEventID(eventID int64, limit, page int64) ([]eventmodels.Participant, error) {
	offset := (page - 1) * limit
	var list = []eventmodels.Participant{}
	query := fmt.Sprintf(`
	SELECT * FROM %s 
	WHERE event_id = ? 
	ORDER BY joined_at DESC LIMIT ? OFFSET ?
	`, tables.EventParticipation)

	err := r.db.Select(&list, query, eventID, limit, offset)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *EventRepository) GetParticipantViews(eventID int64, limit, page int64) ([]eventmodels.ParticipantRaw, error) {
	offset := (page - 1) * limit
	var list = []eventmodels.ParticipantRaw{}
	query := fmt.Sprintf(`
	SELECT
		p.id,
		p.event_id,
		p.user_id,
		u.username,
		f.file_key AS photo_key,
		p.joined_at
	FROM %s p
	JOIN %s u ON p.user_id = u.id
	LEFT JOIN %s f ON u.profile_image_id = f.id
	WHERE p.event_id = ?
	ORDER BY p.joined_at DESC LIMIT ? OFFSET ?
	`, tables.EventParticipation, tables.Users, tables.Files)

	err := r.db.Select(&list, query, eventID, limit, offset)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *EventRepository) CountParticipationByEventID(eventID int64) (int64, error) {
	var count int64
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE event_id = ?`, tables.EventParticipation)
	err := r.db.Get(&count, query, eventID)
	return count, err
}

func (r *EventRepository) IsUserParticipant(eventID, userID int64) (bool, error) {
	var exists bool
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE event_id = ? AND user_id = ?)`, tables.EventParticipation)
	err := r.db.Get(&exists, query, eventID, userID)
	return exists, err
}

func (r *EventRepository) DeleteParticipation(applicationID int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, tables.EventParticipation)
	_, err := r.db.Exec(query, applicationID)
	return err
}
