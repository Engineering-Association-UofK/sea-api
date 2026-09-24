package electionservice

import (
	"context"
	"errors"
	"sort"
	"time"

	"sea-api/internal/models"
	"sea-api/internal/models/electionsmodels"
)

func (s *ElectionService) Reset(ctx context.Context) error {
	cfg, err := s.repo.GetElectionConfig()
	if err != nil {
		return err
	}

	if time.Now().Before(cfg.EndDate) {
		return errors.New("cannot reset: election is still active")
	}

	tx, err := s.repo.Transaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get raw aggregated votes
	rawResults, err := s.repo.GetRawVoteResults(tx, cfg.ActiveCycle, cfg.StudentBase)
	if err != nil {
		return err
	}

	// Resolve Council of 30 Quotas and Placements
	finalResults := ResolveCouncilOfThirty(rawResults)

	// Save historical data
	if err := s.repo.SaveFinalResults(tx, finalResults); err != nil {
		return err
	}

	// Clean up tables
	if err := s.repo.RemoveAllVotes(tx); err != nil {
		return err
	}
	if err := s.repo.RemoveAllTickets(tx); err != nil {
		return err
	}

	return tx.Commit()
}

// ResolveCouncilOfThirty enforces department quotas and calculates DENSE_RANK placements
func ResolveCouncilOfThirty(results []electionsmodels.Result) []electionsmodels.Result {
	if len(results) <= 30 {
		return assignPlaces(results) // No bumping needed if <= 30 candidates ran
	}

	top30 := results[:30]
	others := results[30:]

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
	finalResults := append(top30, others...)
	return assignPlaces(finalResults)
}

// assignPlaces iterates through the resolved list and maps placement numbers
func assignPlaces(results []electionsmodels.Result) []electionsmodels.Result {
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
