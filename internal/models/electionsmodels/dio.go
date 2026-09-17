package electionsmodels

type CandidateRequest struct {
	UserID int64 `json:"user_id"`
}

type Voterequest struct {
	Ticket string  `json:"ticket"`
	Votes  []int64 `json:"votes"`
}

type CandidateResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}
