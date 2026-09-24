package electionservice

import (
	"errors"
	"sea-api/internal/models/electionsmodels"
	"sea-api/internal/repositories/electionrepo"
	"time"
)

type ElectionService struct {
	repo electionrepo.ElectionRepo
}

func (s *ElectionService) CreateCandidate(req electionsmodels.Candidate) (int64, error) {
	cfg, err := s.repo.GetElectionConfig()
	if err != nil {
		return 0, err
	}
	if time.Now().After(cfg.StartDate) {
		return 0, errors.New("cannot add candidates: election has already started")
	}

	req.Cycle = cfg.ActiveCycle
	return s.repo.CreateCandidate(req)
}

func (s *ElectionService) UpdateCandidate(req electionsmodels.Candidate) error {
	cfg, err := s.repo.GetElectionConfig()
	if err != nil {
		return err
	}
	if time.Now().After(cfg.StartDate) {
		return errors.New("cannot update candidates: election has already started")
	}
	return s.repo.UpdateCandidate(req)
}

func (s *ElectionService) GetCandidateList() ([]electionsmodels.CandidateResponse, error) {
	raws, err := s.repo.GetCandidateListForActiveCycle()
	if err != nil {
		return nil, err
	}

	res := make([]electionsmodels.CandidateResponse, len(raws))
	for i, r := range raws {
		res[i] = electionsmodels.CandidateResponse{
			Placing: int64(i + 1),
			Name:    r.Name,
			Url:     r.PicKey,
		}
	}
	return res, nil
}

func (s *ElectionService) RemoveCandidate(id int64) error {
	cfg, err := s.repo.GetElectionConfig()
	if err != nil {
		return err
	}
	if time.Now().After(cfg.StartDate) {
		return errors.New("cannot remove candidates: election has already started")
	}
	return s.repo.RemoveCandidate(id)
}
