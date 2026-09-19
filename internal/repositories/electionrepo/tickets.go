package electionrepo

import "github.com/jmoiron/sqlx"

func (r *ElectionRepo) AddRecord(ts sqlx.Tx, userID int64) error {
	return nil
}

func (r *ElectionRepo) SaveTicket(ts sqlx.Tx, ticket string) error {
	return nil
}

func (r *ElectionRepo) UseTicket(ts sqlx.Tx, ticket string) error {
	return nil
}

func (r *ElectionRepo) RemoveAllRecords(ts sqlx.Tx) error {
	return nil
}

func (r *ElectionRepo) RemoveAllTickets(ts sqlx.Tx) error {
	return nil
}
