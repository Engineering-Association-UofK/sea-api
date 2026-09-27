package electionrepo

import (
	"database/sql"
	"errors"
	"fmt"
	"sea-api/internal/models"
	"sea-api/internal/models/electionsmodels"

	"github.com/jmoiron/sqlx"
)

// AddRecord registers that a user has requested/issued a ticket for a specific election cycle
func (r *ElectionRepo) AddRecord(tx *sqlx.Tx, userID int64, cycle int64) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (user_id, cycle)
		VALUES (?, ?)
	`, models.TableTicketRecords)

	_, err := tx.Exec(query, userID, cycle)
	return err
}

// HasRecord checks if a user has already requested a ticket for a specific election cycle
func (r *ElectionRepo) HasRecord(userID int64, cycle int64) (bool, error) {
	var exists bool
	query := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT 1 
			FROM %s 
			WHERE user_id = ? AND cycle = ?
		)
	`, models.TableTicketRecords)

	err := r.db.Get(&exists, query, userID, cycle)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// GetRecordByUserIDAndCycle retrieves the ticket record entry for a user in a specific cycle
func (r *ElectionRepo) GetRecordByUserIDAndCycle(userID int64, cycle int64) (*electionsmodels.TicketRecord, error) {
	var record electionsmodels.TicketRecord
	query := fmt.Sprintf(`
		SELECT id, user_id, cycle, created_at 
		FROM %s 
		WHERE user_id = ? AND cycle = ?
	`, models.TableTicketRecords)

	err := r.db.Get(&record, query, userID, cycle)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &record, nil
}
