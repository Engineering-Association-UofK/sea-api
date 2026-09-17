package electionsmodels

import "sea-api/internal/models"

type Candidate struct {
	ID     int64 `db:"id"`
	UserID int64 `db:"user_id"`
	Cycle  int64 `db:"cycle"`

	Belonging models.Department `db:"belonging"`
}

type TicketRecord struct {
	UserID int64 `db:"user_id"`
}

type VoteTicket struct {
	Code string `db:"code"`
	Used bool   `db:"used"`
}

type Vote struct {
	CandidateID int64 `db:"candidate_id"`
}
