package electionservice

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"sea-api/internal/errs"
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
		return nil, errs.New(errs.Forbidden, "ticket issuance has not started yet", nil)
	}

	hasRecord, err := s.repo.HasRecord(userId, cfg.ActiveCycle)
	if err != nil {
		return nil, err
	}
	if hasRecord {
		return nil, errs.New(errs.Forbidden, "user has already received a ticket for this cycle", nil)
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

const ticketAlphabet = "23456789ABCDEFGHJKMNPQRSTVWXYZ" // 30 unambiguous characters

// GenerateShortTicket generates an 8-character string formatted as XXXX-XXXX
func (s *ElectionService) GenerateShortTicket(tx *sqlx.Tx) (string, error) {
	bytes := make([]byte, 8)
	alphabetLen := big.NewInt(int64(len(ticketAlphabet)))

	var ticket string
	var err error

	for attempts := 0; attempts < 3; attempts++ {
		// Populate random bytes
		for i := 0; i < 8; i++ {
			num, randErr := rand.Int(rand.Reader, alphabetLen)
			if randErr != nil {
				return "", randErr
			}
			bytes[i] = ticketAlphabet[num.Int64()]
		}

		ticket = string(bytes[:4]) + "-" + string(bytes[4:])

		// Attempt insertion
		err = s.repo.SaveTicket(tx, ticket)
		if err == nil {
			fmt.Println(ticket)
			return ticket, nil
		}
		// If duplicate key collision occurs, the loop retries with fresh random bytes
	}

	return "", fmt.Errorf("failed to generate unique ticket after retries: %w", err)
}
