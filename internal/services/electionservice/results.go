package electionservice

import (
	"context"
	"fmt"
	"sort"
	"time"

	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/models/electionsmodels"
)

func (s *ElectionService) GetResults(ctx context.Context, cycle int) ([]electionsmodels.ResultResponse, error) {
	cfg, err := s.repo.GetElectionConfig()
	if err != nil {
		return nil, err
	}

	if cfg.ActiveCycle > int64(cycle) {
		return s.repo.GetResultsViewByCycle(cycle)
	}
	if cfg.ActiveCycle < int64(cycle) {
		return nil, errs.New(errs.BadRequest, "Unreachable cycle number", nil)
	}

	if time.Now().Before(cfg.EndDate) {
		return nil, errs.New(errs.Forbidden, "cannot resolve: election is still active", nil)
	}
	if !cfg.CycleDoneState[cfg.ActiveCycle] {
		if err := s.ResolveAndReset(ctx); err != nil {
			// Handles potential concurrent attempt race conditions
			cfgCheck, checkErr := s.repo.GetElectionConfig()
			if checkErr == nil && cfgCheck.CycleDoneState[cfgCheck.ActiveCycle] {
				// It resolved successfully by another request, proceed silently
			} else {
				return nil, fmt.Errorf("auto-resolving election results failed: %w", err)
			}
		}
	}

	return s.repo.GetResultsViewByCycle(cycle)
}

func (s *ElectionService) ResolveAndReset(ctx context.Context) error {
	s.resolveMu.Lock()
	defer s.resolveMu.Unlock()

	cfg, err := s.repo.GetElectionConfig()
	if err != nil {
		return err
	}

	if cfg.CycleDoneState[cfg.ActiveCycle] {
		return errs.New(errs.Forbidden, "Cycle already resolved", nil)
	}

	if time.Now().Before(cfg.EndDate) {
		return errs.New(errs.Forbidden, "cannot resolve: election is still active", nil)
	}

	tx, err := s.repo.Transaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get raw aggregated votes
	rawResults, err := s.repo.GetRawVoteResults(tx, cfg.ActiveCycle)
	if err != nil {
		return err
	}

	// Resolve Council of 30 Quotas and Placements
	top30, finalResults := ResolveCouncilOfThirty(rawResults)

	// Save historical data
	var res = []electionsmodels.Result{}
	for _, r := range finalResults {
		res = append(res, r.Result)
	}
	if err := s.repo.SaveFinalResults(tx, res); err != nil {
		return err
	}

	// Clean up tables
	if err := s.repo.RemoveAllVotes(tx); err != nil {
		return err
	}
	if err := s.repo.RemoveAllTickets(tx); err != nil {
		return err
	}

	// Create the new list of the thirty council members to be changed later on
	var replacement []models.TeamMemberModel
	for i, t := range top30 {
		replacement = append(replacement, models.TeamMemberModel{
			ID:           int64(i + 1),
			UserID:       t.UserID,
			Role:         fmt.Sprintf("Candidate Num: %d", i+1),
			Bio:          fmt.Sprintf("With a total number of votes: %d", t.NumberOfVotes),
			DisplayOrder: i + 1,
			IsActive:     true,
			CreatedAt:    time.Now(),
		})
	}
	if err := s.cmsRepo.ReplaceTeamMembers(tx, replacement); err != nil {
		return err
	}

	if err := s.repo.ResolveElection(tx, cfg); err != nil {
		return err
	}

	return tx.Commit()
}

// ResolveCouncilOfThirty enforces department quotas and calculates DENSE_RANK placements
func ResolveCouncilOfThirty(results []electionsmodels.ResultsRaw) ([]electionsmodels.ResultsRaw, []electionsmodels.ResultsRaw) {
	if len(results) <= 30 {
		res := assignPlaces(results)
		return res, res // No bumping needed if <= 30 candidates ran
	}

	top30 := make([]electionsmodels.ResultsRaw, 30)
	copy(top30, results[:30])

	others := make([]electionsmodels.ResultsRaw, len(results)-30)
	copy(others, results[30:])

	deptCounts := make(map[models.Department]int)
	for _, r := range top30 {
		deptCounts[r.Belonging]++
	}

	allDepts := []models.Department{
		models.DEP_MECHANICAL, models.DEP_CIVIL, models.DEP_ELECTRICAL,
		models.DEP_CHEMICAL, models.DEP_PETROLEUM, models.DEP_AGRICULTURAL,
		models.DEP_MINING, models.DEP_SURVEYING,
	}

	for _, dept := range allDepts {
		if deptCounts[dept] == 0 {
			// Find the highest voted candidate from this missing department outside top 30
			bestIdx := -1
			for i, o := range others {
				if o.Belonging == dept {
					bestIdx = i
					break // 'others' is already ordered by votes DESC from SQL
				}
			}

			if bestIdx != -1 {
				// Find candidate to demote (lowest votes in top30 whose dept has > 1 rep)
				replaceIdx := -1
				for i := len(top30) - 1; i >= 0; i-- {
					if deptCounts[top30[i].Belonging] > 1 {
						replaceIdx = i
						break
					}
				}

				if replaceIdx != -1 {
					// Update counts and swap them
					deptCounts[top30[replaceIdx].Belonging]--
					deptCounts[dept]++

					promoted := others[bestIdx]
					demoted := top30[replaceIdx]

					top30[replaceIdx] = promoted
					others[bestIdx] = demoted

					// Re-sort the slices to maintain integrity for subsequent loops
					sort.SliceStable(top30, func(i, j int) bool {
						return top30[i].NumberOfVotes > top30[j].NumberOfVotes
					})
					sort.SliceStable(others, func(i, j int) bool {
						return others[i].NumberOfVotes > others[j].NumberOfVotes
					})
				}
			}
		}
	}

	// Recombine and assign numerical placements (Dense Rank)
	finalResults := append(append([]electionsmodels.ResultsRaw{}, top30...), others...)
	return assignPlaces(top30), assignPlaces(finalResults)
}

// assignPlaces iterates through the resolved list and maps placement numbers
func assignPlaces(results []electionsmodels.ResultsRaw) []electionsmodels.ResultsRaw {
	if len(results) == 0 {
		return results
	}

	place := int64(1)
	results[0].Place = place

	for i := 1; i < len(results); i++ {
		// Standard Dense Rank tie logic
		isTie := results[i].NumberOfVotes == results[i-1].NumberOfVotes

		// Hard boundary check: Someone bumped OUT of top 30 cannot tie with someone bumped IN
		if i == 30 {
			isTie = false
		}

		if !isTie {
			place++
		}

		results[i].Place = place
	}

	return results
}
