package electionsmodels

import (
	"sea-api/internal/models"
	"time"
)

type ElectionConfig struct {
	ActiveCycle int64 `db:"active_cycle"`
	MaxVotes    int64 `db:"max_votes"`
	StudentBase int64 `db:"student_base"`

	TicketStartDate time.Time `db:"ticket_start_date"`
	StartDate       time.Time `db:"start_date"`
	EndDate         time.Time `db:"end_date"`
}

type CandidateRaw struct {
	ID     int64  `db:"id"`
	Name   string `db:"name"`
	PicKey string `db:"pic_key"`

	Cycle int64 `db:"cycle"`

	Belonging models.Department `db:"belonging"`
}

type VoteStatistics struct {
	NumberOfVoters   int64   `db:"number_of_voters" json:"number_of_voters"`
	VotesInLastDay   int64   `db:"votes_in_last_day" json:"votes_in_last_day"`
	VotePercentage   float32 `db:"vote_percentage" json:"vote_percentage"`
	SecondsRemaining int64   `db:"time_remaining" json:"time_remaining"`
}
