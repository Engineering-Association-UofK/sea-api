package electionservice

import (
	"context"
	"sea-api/internal/models/electionsmodels"
)

func (s *ElectionService) Vote(ctx context.Context, req electionsmodels.VoteRequest) error {
	// Get ElectionConfig to see if election started

	// Make sure no more than 30 votes

	// Make sure no duplicate voting for the same candidate

	// Get Ts

	// UseTicket

	// Bulk Vote using the vote function
	return nil
}
