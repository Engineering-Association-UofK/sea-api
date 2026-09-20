package electionsmodels

import (
	"sea-api/internal/models"
	"time"
)

// election_config
type ElectionConfig struct {
	ActiveCycle int64 `db:"active_cycle"`
	MaxVotes    int64 `db:"max_votes"`

	TicketStartDate time.Time `db:"ticket_start_date"`
	StartDate       time.Time `db:"start_date"`
	EndDate         time.Time `db:"end_date"`
}

// candidates
type Candidate struct {
	ID        int64             `db:"id" json:"id"`
	UserID    int64             `db:"user_id" json:"user_id"`
	Cycle     int64             `db:"cycle" json:"cycle"`
	Belonging models.Department `db:"belonging" json:"belonging"`
}

// results
type Result struct {
	Name  string `db:"name"`
	Cycle int64  `db:"cycle"`
	Place int64  `db:"place"`

	StudentBase   int64 `db:"student_base"`
	NumberOfVotes int64 `db:"number_of_votes"`

	Belonging models.Department `db:"belonging"`
}

// ticket_records
type TicketRecord struct {
	UserID int64 `db:"user_id"`
}

// vote_tickets
type VoteTicket struct {
	Code   string    `db:"code"`
	Used   bool      `db:"used"`
	UsedAt time.Time `db:"used_at"`
}

// votes
type Vote struct {
	CandidateID int64 `db:"candidate_id"`
}
