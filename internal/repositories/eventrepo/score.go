package eventrepo

import (
	"fmt"
	"sea-api/internal/models"

	"github.com/jmoiron/sqlx"
)

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

func (r *EventRepository) GetScoreByID(id int64) (*models.ComponentScoreModel, error) {
	var score models.ComponentScoreModel
	err := r.db.Get(&score, fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, models.TableComponentScores), id)
	if err != nil {
		return nil, err
	}
	return &score, nil
}

func (r *EventRepository) GetScoresByParticipantID(participantID int64) ([]models.ComponentScoreModel, error) {
	var scores []models.ComponentScoreModel
	err := r.db.Select(&scores, fmt.Sprintf(`SELECT * FROM %s WHERE participant_id = ?`, models.TableComponentScores), participantID)
	if err != nil {
		return nil, err
	}
	return scores, nil
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
