package electionrepo

import (
	"fmt"
	"sea-api/internal/models"
	"sea-api/internal/models/electionsmodels"

	"github.com/jmoiron/sqlx"
)

// GetRawVoteResults aggregates the votes for all candidates
func (r *ElectionRepo) GetRawVoteResults(tx *sqlx.Tx, cycle int64, studentBase int64) ([]electionsmodels.Result, error) {
	var results []electionsmodels.Result
	query := fmt.Sprintf(`
		SELECT 
			u.name_en AS name,
			c.cycle,
			? AS student_base,
			COUNT(v.id) AS number_of_votes,
			c.belonging
		FROM %s c
		JOIN %s u ON c.user_id = u.id
		LEFT JOIN %s v ON c.id = v.candidate_id
		WHERE c.cycle = ?
		GROUP BY c.id, c.cycle, c.belonging, u.name_en
		ORDER BY number_of_votes DESC, c.id ASC
	`, models.TableCandidates, models.TableUsers, models.TableVotes)

	err := tx.Select(&results, query, studentBase, cycle)
	return results, err
}

// SaveFinalResults inserts the algorithmically placed results in bulk
func (r *ElectionRepo) SaveFinalResults(tx *sqlx.Tx, results []electionsmodels.Result) error {
	if len(results) == 0 {
		return nil
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (name, cycle, place, student_base, number_of_votes, belonging)
		VALUES (:name, :cycle, :place, :student_base, :number_of_votes, :belonging)
	`, models.TableElectionResults)

	_, err := tx.NamedExec(query, results)
	return err
}
