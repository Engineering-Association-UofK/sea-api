package tables

type Name string

const (

	// Documents

	// Bot

	// Elections

	Candidates      Name = "election_candidates"
	TicketRecords   Name = "election_ticket_records"
	VoteTickets     Name = "election_vote_tickets"
	Votes           Name = "election_votes"
	ElectionResults Name = "election_results"

	// TODO: Remove table
	VerificationCode Name = "verification_code"
)
