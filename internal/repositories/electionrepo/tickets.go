package electionrepo

import (
	"fmt"
	"sea-api/internal/models"

	"github.com/jmoiron/sqlx"
)

// SaveTicket stores an unlinked, generated ticket UUID into the vote_tickets pool
func (r *ElectionRepo) SaveTicket(tx *sqlx.Tx, ticket string) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (code, used)
		VALUES (?, FALSE)
	`, models.TableVoteTickets)

	_, err := tx.Exec(query, ticket)
	return err
}

// UseTicket verifies that a ticket exists and is unused, then marks it as used
func (r *ElectionRepo) UseTicket(tx *sqlx.Tx, ticket string) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET used = TRUE, used_at = NOW()
		WHERE code = ? AND used = FALSE
	`, models.TableVoteTickets)

	res, err := tx.Exec(query, ticket)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("invalid or already used ticket")
	}

	return nil
}

// RemoveAllTickets purges active vote tickets after election closing
func (r *ElectionRepo) RemoveAllTickets(tx *sqlx.Tx) error {
	query := fmt.Sprintf(`TRUNCATE TABLE %s`, models.TableVoteTickets)
	_, err := tx.Exec(query)
	return err
}
