package electionservice

import (
	"context"
	"errors"
	"sea-api/internal/models/electionsmodels"
	"time"
)

func (s *ElectionService) Vote(ctx context.Context, req electionsmodels.VoteRequest) error {
	cfg, err := s.repo.GetElectionConfig()
	if err != nil {
		return err
	}

	now := time.Now()
	if now.Before(cfg.StartDate) || now.After(cfg.EndDate) {
		return errors.New("voting is currently closed")
	}

	maxVotes := cfg.MaxVotes
	if maxVotes <= 0 {
		maxVotes = 30 // Fallback limit
	}

	if len(req.Votes) == 0 || int64(len(req.Votes)) > maxVotes {
		return errors.New("invalid number of votes submitted")
	}

	uniqueVotes := DeduplicateVotes(req.Votes)

	tx, err := s.repo.Transaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.repo.UseTicket(tx, req.Ticket); err != nil {
		return err
	}

	if err := s.repo.Vote(tx, uniqueVotes); err != nil {
		return err
	}

	return tx.Commit()
}

func DeduplicateVotes(candidateIDs []int64) []int64 {
	seen := make(map[int64]bool)
	result := make([]int64, 0, len(candidateIDs))

	for _, id := range candidateIDs {
		if !seen[id] {
			seen[id] = true
			result = append(result, id)
		}
	}
	return result
}
