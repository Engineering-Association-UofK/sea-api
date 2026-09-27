package electionservice

import (
	"context"
	"sea-api/internal/errs"
	"sea-api/internal/models/electionsmodels"
	"slices"
	"time"
)

func (s *ElectionService) ElectionStatistics() (*electionsmodels.VoteStatistics, error) {
	stats, err := s.repo.GetVotesStatistics()
	if err != nil {
		return nil, err
	}

	cfg, err := s.repo.GetElectionConfig()
	if err != nil {
		return nil, err
	}

	stats.VotePercentage = (float32(stats.NumberOfVoters) / float32(cfg.StudentBase)) * 100
	stats.TicketsStartTime = cfg.TicketStartDate
	stats.StartTime = cfg.StartDate
	stats.EndTime = cfg.EndDate

	return stats, nil
}

func (s *ElectionService) Vote(ctx context.Context, req electionsmodels.VoteRequest) error {
	cfg, err := s.repo.GetElectionConfig()
	if err != nil {
		return err
	}

	// Check if it's voting period
	now := time.Now()
	if now.Before(cfg.StartDate) || now.After(cfg.EndDate) {
		return errs.New(errs.Forbidden, "voting is currently closed", nil)
	}

	// Check the voting count
	maxVotes := cfg.MaxVotes
	if maxVotes <= 0 {
		maxVotes = 30 // Fallback limit
	}
	if len(req.Votes) == 0 || int64(len(req.Votes)) > maxVotes {
		return errs.New(errs.Forbidden, "invalid number of votes submitted", nil)
	}

	// Check for duplicate votes
	uniqueVotes := DeduplicateVotes(req.Votes)

	// Check for wrong IDs
	IDs, err := s.repo.GetCandidateIDsForActiveCycle()
	if err != nil {
		return err
	}
	for _, id := range req.Votes {
		if !slices.Contains(IDs, id) {
			return errs.New(errs.BadRequest, "Wrong candidate ID in vote list", nil)
		}
	}

	// Begin the transaction
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
