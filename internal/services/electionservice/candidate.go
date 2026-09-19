package electionservice

import (
	"sea-api/internal/models/electionsmodels"
	"sea-api/internal/repositories/electionrepo"
)

type ElectionService struct {
	repo electionrepo.ElectionRepo
}

func (s *ElectionService) SetActiveCycle(cycle int64) error {
	return nil
}

func (s *ElectionService) CreateCandidate(req electionsmodels.Candidate) (int64, error) {
	// use ElectionConfig to get start date to prevent changes after starting
	return 0, nil
}

func (s *ElectionService) UpdateCandidate(req electionsmodels.Candidate) error {
	// use ElectionConfig to get start date to prevent changes after starting
	return nil
}

// The placing is the placing on the list as a static simple one
// to three digits for simple media sharing
func (s *ElectionService) GetCandidateList() ([]electionsmodels.CandidateResponse, error) {
	// use ElectionConfig to get active cycle
	return []electionsmodels.CandidateResponse{}, nil
}

func (s *ElectionService) RemoveCandidate(id int64) error {
	return nil
}
