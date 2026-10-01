package electionservice

import (
	"context"
	"sea-api/internal/models"
	"sea-api/internal/models/electionsmodels"
	"sea-api/internal/repositories"
	"sea-api/internal/repositories/electionrepo"
	"sea-api/internal/services/storage"
	"sync"
)

type ElectionService struct {
	repo      *electionrepo.ElectionRepo
	cmsRepo   *repositories.CmsRepository
	s3        *storage.S3
	resolveMu sync.Mutex
}

func NewElectionService(repo *electionrepo.ElectionRepo, cmsRepo *repositories.CmsRepository, s3 *storage.S3) *ElectionService {
	return &ElectionService{
		repo:    repo,
		cmsRepo: cmsRepo,
		s3:      s3,
	}
}

func (s *ElectionService) CreateCandidate(req *electionsmodels.CandidateCreateRequest) (int64, error) {
	cfg, err := s.Guard(true, false, false)
	if err != nil {
		return 0, err
	}

	return s.repo.CreateCandidate(electionsmodels.Candidate{
		UserID:    req.UserID,
		Cycle:     cfg.ActiveCycle,
		Belonging: req.Belonging,
	})
}

func (s *ElectionService) UpdateCandidate(id int64, belonging models.Department) error {
	_, err := s.Guard(true, false, false)
	if err != nil {
		return err
	}
	return s.repo.UpdateCandidateBelonging(id, belonging)
}

func (s *ElectionService) GetCandidateList(ctx context.Context) ([]electionsmodels.CandidateResponse, error) {
	_, err := s.Guard(true, true, true)
	if err != nil {
		return nil, err
	}

	raws, err := s.repo.GetCandidateListForActiveCycle()
	if err != nil {
		return nil, err
	}

	res := make([]electionsmodels.CandidateResponse, len(raws))
	for i, r := range raws {
		var url string
		if r.PicKey != "" {
			url, err = s.s3.GenerateDownloadUrlByKey(ctx, r.PicKey)
			if err != nil {
				return nil, err
			}
		}

		res[i] = electionsmodels.CandidateResponse{
			ID:        r.ID,
			Placing:   int64(i + 1),
			NameAr:    r.NameAr,
			NameEn:    r.NameEn,
			Belonging: r.Belonging,
			Url:       url,
		}
	}
	return res, nil
}

func (s *ElectionService) RemoveCandidate(id int64) error {
	_, err := s.Guard(true, false, false)
	if err != nil {
		return err
	}
	return s.repo.RemoveCandidate(id)
}
