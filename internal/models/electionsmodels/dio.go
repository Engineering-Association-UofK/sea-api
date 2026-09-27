package electionsmodels

import "sea-api/internal/models"

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
