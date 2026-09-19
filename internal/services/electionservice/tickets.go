package electionservice

import (
	"context"
	"sea-api/internal/models/electionsmodels"
)

func (s *ElectionService) GetTicket(ctx context.Context, userId int64) (*electionsmodels.TicketResponse, error) {
	// Get ElectionConfig to make sure TicketStartDate passed

	// Get Ts

	// Add record

	// Generate Ticket

	// Save ticket

	//return ticket
	return nil, nil
}

func (s *ElectionService) Reset() error {
	// Use ElectionConfig to make sure elections ended
	return nil
}
