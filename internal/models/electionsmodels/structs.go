package electionsmodels

import (
	"sea-api/internal/models"
)

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
