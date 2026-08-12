package eventrepo

import (
	"fmt"
	"sea-api/internal/models"
	"time"
)

func (r *EventRepository) Apply(userID, eventID int64) error {
	query := fmt.Sprintf(`
	INSERT INTO %s (event_id, user_id, status, joined_at, grade, completed)
	VALUES (?, ?, ?, ?, 0, 0)
	`, models.TableEventParticipants)
	_, err := r.db.Exec(query, eventID, userID, models.PENDING, time.Now())
	return err
}

func (r *EventRepository) Cancel(userID, eventID int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE user_id = ? AND event_id = ?`, models.TableEventParticipants)
	_, err := r.db.Exec(query, userID, eventID)
	return err
}

func (r *EventRepository) GetTotalApplicationsForUser(userID int64) (int64, error) {
	var total int64
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE user_id = ?`, models.TableEventParticipants)
	err := r.db.Get(&total, query, userID)
	return total, err
}

func (r *EventRepository) GetApplicationsForUser(userID int64, req models.ListRequest) ([]models.ApplicationStatus, error) {
	var applications = []models.ApplicationStatus{}
	query := fmt.Sprintf(`
	SELECT 
		ep.event_id,
		e.name as event_name,
		ep.status
	FROM %s ep
	JOIN %s e ON ep.event_id = e.id
	WHERE ep.user_id = ?
	ORDER BY ep.joined_at DESC LIMIT ? OFFSET ?
	`, models.TableEventParticipants, models.TableEvents)

	err := r.db.Select(&applications, query, userID, req.Limit, (req.Page-1)*req.Limit)
	return applications, err
}

func (r *EventRepository) GetApplicationStatus(userID, eventID int64) (*models.ApplicationStatus, error) {
	var application models.ApplicationStatus

	query := fmt.Sprintf(`
	SELECT 
		ep.event_id,
		e.name as event_name,
		ep.status
	FROM %s ep
	JOIN %s e ON ep.event_id = e.id
	WHERE ep.user_id = ?
	AND ep.event_id = ?
	`, models.TableEventParticipants, models.TableEvents)

	err := r.db.Get(&application, query, userID, eventID)
	return &application, err
}
