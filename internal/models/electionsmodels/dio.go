package electionsmodels

type TicketResponse struct {
	Ticket string `json:"ticket"`
}
type VoteRequest struct {
	Ticket string  `json:"ticket"`
	Votes  []int64 `json:"votes"`
}

type CandidateResponse struct {
	Placing int64  `json:"placing"`
	Name    string `json:"name"`
	Url     string `json:"url"`
}
