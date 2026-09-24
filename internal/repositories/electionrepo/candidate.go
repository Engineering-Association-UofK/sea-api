package electionrepo

import (
	"fmt"
	"sea-api/internal/models"
	"sea-api/internal/models/electionsmodels"

	"github.com/jmoiron/sqlx"
)

type ElectionRepo struct {
	db sqlx.DB
}

// CreateCandidate inserts a new candidate for a given cycle
func (r *ElectionRepo) CreateCandidate(req electionsmodels.Candidate) (int64, error) {
	query := fmt.Sprintf(`
		INSERT INTO %s (user_id, cycle, belonging)
		VALUES (:user_id, :cycle, :belonging)
	`, models.TableCandidates)

	res, err := r.db.NamedExec(query, req)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdateCandidate updates belonging or assigned cycle of a candidate
func (r *ElectionRepo) UpdateCandidate(req electionsmodels.Candidate) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET user_id = :user_id, cycle = :cycle, belonging = :belonging
		WHERE id = :id
	`, models.TableCandidates)

	_, err := r.db.NamedExec(query, req)
	return err
}

// GetCandidateList returns candidates joined with user data and profile picture key for a specific cycle
func (r *ElectionRepo) GetCandidateList(cycle int64) ([]electionsmodels.CandidateRaw, error) {
	var candidates []electionsmodels.CandidateRaw

	query := fmt.Sprintf(`
		SELECT 
			c.id,
			u.name_en AS name,
			COALESCE(f.file_key, '') AS pic_key,
			c.cycle,
			c.belonging
		FROM %s c
		JOIN %s u ON c.user_id = u.id
		LEFT JOIN %s f ON u.profile_image_id = f.id
		WHERE c.cycle = ?
		ORDER BY c.id ASC
	`, models.TableCandidates, models.TableUsers, models.TableFiles)

	err := r.db.Select(&candidates, query, cycle)
	if err != nil {
		return nil, err
	}
	return candidates, nil
}

// GetCandidateListForActiveCycle retrieves candidates for the active cycle set in the global config
func (r *ElectionRepo) GetCandidateListForActiveCycle() ([]electionsmodels.CandidateRaw, error) {
	var candidates []electionsmodels.CandidateRaw

	query := fmt.Sprintf(`
		SELECT 
			c.id,
			u.name_en AS name,
			COALESCE(f.file_key, '') AS pic_key,
			c.cycle,
			c.belonging
		FROM %s c
		JOIN %s u ON c.user_id = u.id
		LEFT JOIN %s f ON u.profile_image_id = f.id
		JOIN %s cfg ON cfg.key = 'election_config'
		WHERE c.cycle = CAST(JSON_EXTRACT(cfg.value, '$.active_cycle') AS UNSIGNED)
		ORDER BY c.created_at ASC
	`, models.TableCandidates, models.TableUsers, models.TableFiles, models.TableConfig)

	err := r.db.Select(&candidates, query)
	if err != nil {
		return nil, err
	}
	return candidates, nil
}

// RemoveCandidate deletes a candidate by ID
func (r *ElectionRepo) RemoveCandidate(id int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, models.TableCandidates)
	_, err := r.db.Exec(query, id)
	return err
}
