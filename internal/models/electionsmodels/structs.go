package electionsmodels

import (
	"sea-api/internal/models"
	"time"
)

type ElectionConfig struct {
	ActiveCycle int64 `json:"active_cycle"`
	MaxVotes    int64 `json:"max_votes"`
	StudentBase int64 `json:"student_base"`

	CycleDoneState map[int64]bool `json:"cycle_done_state"`

	TicketStartDate time.Time `json:"ticket_start_date"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
}

type CandidateRaw struct {
	ID     int64  `db:"id"`
	NameAr string `db:"name_ar"`
	NameEn string `db:"name_en"`
	PicKey string `db:"pic_key"`

	Cycle int64 `db:"cycle"`

	Belonging models.Department `db:"belonging"`
}

type Stats struct {
	TicketsDistributed int64 `db:"tickets_distributed" json:"tickets_distributed"`
	NumberOfVoters     int64 `db:"number_of_voters" json:"number_of_voters"`
	VotesInLastDay     int64 `db:"votes_in_last_day" json:"votes_in_last_day"`
}

type PublicStats struct {
	VotePercentage   float32   `json:"vote_percentage"`
	TicketsStartTime time.Time `json:"tickets_start_time"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`

	CurrentCycle int64 `json:"current_cycle"`
	VoteLimit    int64 `json:"vote_limit"`
}

type PrivateStats struct {
	Stats

	VotePercentage float32 `json:"vote_percentage"`
	StudentBody    int64   `json:"student_body"`

	TicketsStartTime time.Time `json:"tickets_start_time"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
}
