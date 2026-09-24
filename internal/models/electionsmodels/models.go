package electionsmodels

import (
	"sea-api/internal/models"
	"time"
)

// candidates
type Candidate struct {
	ID        int64             `db:"id" json:"id"`
	UserID    int64             `db:"user_id" json:"user_id"`
	Cycle     int64             `db:"cycle" json:"cycle"`
	Belonging models.Department `db:"belonging" json:"belonging"`
	CreatedAt time.Time         `db:"created_at"`
}

// election_results
type Result struct {
	ID    int64  `db:"id" json:"id"`
	Name  string `db:"name"`
	Cycle int64  `db:"cycle"`
	Place int64  `db:"place"`

	NumberOfVotes int64 `db:"number_of_votes"`

	Belonging models.Department `db:"belonging"`
	CreatedAt time.Time         `db:"created_at"`
}

// ticket_records
type TicketRecord struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `db:"user_id"`
	Cycle     int64     `db:"cycle"`
	CreatedAt time.Time `db:"created_at"`
}

// vote_tickets
type VoteTicket struct {
	Code      string    `db:"code"`
	Used      bool      `db:"used"`
	UsedAt    time.Time `db:"used_at"`
	CreatedAt time.Time `db:"created_at"`
}

// votes
type Vote struct {
	ID          int64     `db:"id" json:"id"`
	CandidateID int64     `db:"candidate_id"`
	CreatedAt   time.Time `db:"created_at"`
}
