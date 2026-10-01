package electionservice

import (
	"sea-api/internal/errs"
	"sea-api/internal/models/electionsmodels"
	"time"
)

func (s *ElectionService) Guard(
	beforeElection bool,
	throughElection bool,
	afterElection bool,
) (*electionsmodels.ElectionConfig, error) {
	cfg, err := s.repo.GetElectionConfig()
	if err != nil {
		return nil, err
	}
	if !cfg.EngageElection {
		return nil, errs.New(errs.Forbidden, "Action Prohibited: No active election.", nil)
	}

	now := time.Now()

	switch {
	case now.Before(cfg.StartDate):
		if !beforeElection {
			return nil, errs.New(errs.Forbidden, "Action Prohibited: Election did not start yet.", nil)
		}
	case now.After(cfg.EndDate):
		if !afterElection {
			return nil, errs.New(errs.Forbidden, "Action Prohibited: Election has ended.", nil)
		}
	default:
		if !throughElection {
			return nil, errs.New(errs.Forbidden, "Action Prohibited: Action not allowed during voting period.", nil)
		}
	}

	return cfg, nil
}
