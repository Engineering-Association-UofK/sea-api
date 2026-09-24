package electionrepo

import (
	"context"
	"fmt"
	"sea-api/internal/models"
	"sea-api/internal/models/electionsmodels"

	"github.com/jmoiron/sqlx"
)

// Vote performs a bulk insert of candidate selection records within the transaction
func (r *ElectionRepo) Vote(tx *sqlx.Tx, candidateIDs []int64) error {
	if len(candidateIDs) == 0 {
		return nil
	}

	query := fmt.Sprintf(`INSERT INTO %s (candidate_id) VALUES `, models.TableVotes)

	// Dynamic bulk-insert query construction
	vals := []interface{}{}
	for i, id := range candidateIDs {
		if i > 0 {
			query += ", "
		}
		query += "(?)"
		vals = append(vals, id)
	}

	_, err := tx.Exec(query, vals...)
	return err
}

// GetVotesStatistics calculates real-time voting metrics across active vote tickets and records
func (r *ElectionRepo) GetVotesStatistics() (*electionsmodels.VoteStatistics, error) {
	var stats electionsmodels.VoteStatistics

	query := fmt.Sprintf(`
		SELECT 
			(SELECT COUNT(*) FROM %s) AS number_of_voters,
			(SELECT COUNT(*) FROM %s WHERE created_at >= NOW() - INTERVAL 1 DAY) AS votes_in_last_day
	`, models.TableTicketRecords, models.TableVotes)

	err := r.db.Get(&stats, query)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// RemoveAllVotes purges votes after election closing
func (r *ElectionRepo) RemoveAllVotes(tx *sqlx.Tx) error {
	query := fmt.Sprintf(`TRUNCATE TABLE %s`, models.TableVotes)
	_, err := tx.Exec(query)
	return err
}

// Transaction opens a database transaction bound to context
func (r *ElectionRepo) Transaction(ctx context.Context) (*sqlx.Tx, error) {
	return r.db.BeginTxx(ctx, nil)
}
