package electionrepo

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// Voting is always bulk inserting to the database
// since each ID is a record
func (r *ElectionRepo) Vote(ts sqlx.Tx, candidateIDs []int64) error {
	return nil
}

func (r *ElectionRepo) GetVotesStatistics()

func (r *ElectionRepo) Transaction(ctx context.Context) (*sqlx.Tx, error) {
	return r.db.BeginTxx(ctx, nil)
}
