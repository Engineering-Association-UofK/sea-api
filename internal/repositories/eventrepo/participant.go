package eventrepo

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sea-api/internal/models"

	"github.com/jmoiron/sqlx"
)

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

func (r *EventRepository) GetParticipantByID(id int64) (*models.EventParticipantModel, error) {
	var participant models.EventParticipantModel
	err := r.db.Get(&participant, fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, models.TableEventParticipants), id)
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

func (r *EventRepository) GetParticipantsCount(eventID int64) (int, error) {
	var count int
	err := r.db.Get(&count, fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE event_id = ?`, models.TableEventParticipants), eventID)
	if err != nil {
		return 0, err
	}
	return count, nil
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
