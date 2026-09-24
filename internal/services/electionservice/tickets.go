package electionservice

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"sea-api/internal/models/electionsmodels"
	"time"

	"github.com/jmoiron/sqlx"
)

func (s *ElectionService) GetTicket(ctx context.Context, userId int64) (*electionsmodels.TicketResponse, error) {
	cfg, err := s.repo.GetElectionConfig()
	if err != nil {
		return nil, err
	}

	if time.Now().Before(cfg.TicketStartDate) {
		return nil, errors.New("ticket issuance has not started yet")
	}

	hasRecord, err := s.repo.HasRecord(userId, cfg.ActiveCycle)
	if err != nil {
		return nil, err
	}
	if hasRecord {
		return nil, errors.New("user has already received a ticket for this cycle")
	}

	tx, err := s.repo.Transaction(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := s.repo.AddRecord(tx, userId, cfg.ActiveCycle); err != nil {
		return nil, err
	}

	ticket, err := s.GenerateShortTicket(tx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &electionsmodels.TicketResponse{Ticket: ticket}, nil
}

const ticketAlphabet = "23456789ABCDEFGHJKMNPQRSTVWXYZ" // 32 unambiguous characters

// GenerateShortTicket generates an 8-character string
func (s *ElectionService) GenerateShortTicket(tx *sqlx.Tx) (string, error) {
	bytes := make([]byte, 8)
	alphabetLen := big.NewInt(int64(len(ticketAlphabet)))

	// Format as 4-4 for readability (e.g., "A8K9-M2P4")
	var ticket = string(bytes[:4]) + "-" + string(bytes[4:])
	var err error

	for attempts := 0; attempts < 3; attempts++ {
		for i := 0; i < 8; i++ {
			num, err := rand.Int(rand.Reader, alphabetLen)
			if err != nil {
				return "", err
			}
			bytes[i] = ticketAlphabet[num.Int64()]
		}

		err = s.repo.SaveTicket(tx, ticket)
		if err == nil {
			break // Successfully inserted unique ticket
		}
		// If err is unique constraint violation, loop retries with a new code
	}
	if err != nil {
		return "", err // Failed after retries
	}

	return ticket, nil
}
