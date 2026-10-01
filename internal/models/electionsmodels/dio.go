package electionsmodels

import (
	"sea-api/internal/models"
	"time"
)

type CheckElection struct {
	Live  bool  `json:"live"`
	Cycle int64 `json:"cycle"`
}

type CandidateCreateRequest struct {
	UserID    int64             `json:"user_id"`
	Belonging models.Department `json:"belonging"`
}

type CandidateUpdateRequest struct {
	Belonging models.Department `json:"belonging"`
}

type TicketResponse struct {
	Ticket string `json:"ticket"`
}
type VoteRequest struct {
	Ticket string  `json:"ticket"`
	Votes  []int64 `json:"votes"`
}

type CandidateResponse struct {
	Placing   int64             `json:"placing"`
	ID        int64             `json:"id"`
	NameAr    string            `json:"name_ar"`
	NameEn    string            `json:"name_en"`
	Belonging models.Department `json:"belonging"`
	Url       string            `json:"url"`
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

// election_results
type ResultResponse struct {
	Result
	Belonging models.Department `db:"belonging" json:"belonging"`
	NameAr    string            `db:"name_ar" json:"name_ar"`
	NameEn    string            `db:"name_en" json:"name_en"`
}
