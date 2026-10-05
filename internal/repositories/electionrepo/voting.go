package electionrepo

import (
	"context"
	"fmt"
	"sea-api/internal/models/electionsmodels"
	"sea-api/internal/models/tables"

	"github.com/jmoiron/sqlx"
)

// Vote performs a bulk insert of candidate selection records within the transaction
func (r *ElectionRepo) Vote(tx *sqlx.Tx, candidateIDs []int64) error {
	if len(candidateIDs) == 0 {
		return nil
	}

	query := fmt.Sprintf(`INSERT INTO %s (candidate_id) VALUES `, tables.Votes)

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

// GetVotesStatistics returns the total number of voters and the number of votes cast in the last 24 hours
func (r *ElectionRepo) GetVotesStatistics() (*electionsmodels.Stats, error) {
	query := fmt.Sprintf(`
		SELECT 
			(SELECT COUNT(*) FROM %s) AS tickets_distributed,
			(SELECT COUNT(*) FROM %s WHERE used = 1) AS number_of_voters,
			(SELECT COUNT(*) FROM %s WHERE created_at >= NOW() - INTERVAL 1 DAY) AS votes_in_last_day
	`, tables.TicketRecords, tables.VoteTickets, tables.Votes)

	var stats = electionsmodels.Stats{}

	err := r.db.Get(&stats, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch vote counts: %w", err)
	}

	return &stats, nil
}

// RemoveAllVotes purges votes after election closing
func (r *ElectionRepo) RemoveAllVotes(tx *sqlx.Tx) error {
	query := fmt.Sprintf(`TRUNCATE TABLE %s`, tables.Votes)
	_, err := tx.Exec(query)
	return err
}

// Transaction opens a database transaction bound to context
func (r *ElectionRepo) Transaction(ctx context.Context) (*sqlx.Tx, error) {
	return r.db.BeginTxx(ctx, nil)
}
