package electionrepo

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sea-api/internal/models"
	"sea-api/internal/models/electionsmodels"
	"sea-api/internal/utils"

	"github.com/jmoiron/sqlx"
)

type ElectionRepo struct {
	db *sqlx.DB
}

func NewElectionRepo(db *sqlx.DB) *ElectionRepo {
	return &ElectionRepo{db: db}
}

// GetElectionConfig extracts the JSON global config state
func (r *ElectionRepo) GetElectionConfig() (*electionsmodels.ElectionConfig, error) {
	raw, err := utils.GetConfig("election_config", r.db)
	if err != nil {
		slog.Debug("Could not get config", "Error", err)
		return nil, err
	}

	var cfg electionsmodels.ElectionConfig
	err = json.Unmarshal(*raw, &cfg)
	if err != nil {
		slog.Debug("Could not unmarshal config", "Error", err)
		return nil, err
	}
	return &cfg, nil
}

// UpdateElectionConfig serializes and updates the JSON global config state
func (r *ElectionRepo) ResolveElection(tx *sqlx.Tx, cfg *electionsmodels.ElectionConfig) error {
	cfg.CycleDoneState[cfg.ActiveCycle] = true
	return utils.UpdateConfig("election_config", tx, cfg)
}

// UpdateElectionConfig serializes and updates the JSON global config state
func (r *ElectionRepo) UpdateElectionConfig(tx *sqlx.Tx, cfg *electionsmodels.ElectionConfig) error {
	return utils.UpdateConfig("election_config", tx, cfg)
}

// AggregateAndSaveResults computes DENSE_RANK based on vote counts and saves to history
func (r *ElectionRepo) AggregateAndSaveResults(tx *sqlx.Tx, cycle int64) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (name, cycle, place, student_base, number_of_votes, belonging)
		SELECT 
			u.name_en AS name,
			c.cycle,
			DENSE_RANK() OVER (ORDER BY COUNT(v.id) DESC) AS place,
			(SELECT COUNT(*) FROM %s WHERE cycle = ?) AS student_base,
			COUNT(v.id) AS number_of_votes,
			c.belonging
		FROM %s c
		JOIN %s u ON c.user_id = u.id
		LEFT JOIN %s v ON c.id = v.candidate_id
		WHERE c.cycle = ?
		GROUP BY c.id, c.cycle, c.belonging, u.name_en
	`,
		models.TableElectionResults,
		models.TableTicketRecords,
		models.TableCandidates,
		models.TableUsers,
		models.TableVotes,
	)

	_, err := tx.Exec(query, cycle, cycle)
	return err
}
