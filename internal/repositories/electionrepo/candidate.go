package electionrepo

import (
	"sea-api/internal/models/electionsmodels"

	"github.com/jmoiron/sqlx"
)

type ElectionRepo struct {
	db sqlx.DB
}

func (r *ElectionRepo) CreateCandidate(req electionsmodels.Candidate) (int64, error) {
	return 0, nil
}

func (r *ElectionRepo) UpdateCandidate(req electionsmodels.Candidate) error {
	return nil
}

func (r *ElectionRepo) GetCandidateList(req electionsmodels.Candidate) ([]electionsmodels.CandidateRaw, error) {
	return []electionsmodels.CandidateRaw{}, nil
}

func (r *ElectionRepo) RemoveCandidate(req electionsmodels.Candidate) error {
	return nil
}
