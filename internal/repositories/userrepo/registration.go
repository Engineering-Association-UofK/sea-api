package userrepo

import (
	"fmt"
	"sea-api/internal/models/authmodels"
	"sea-api/internal/models/tables"

	"github.com/jmoiron/sqlx"
)

func (r *UserRepository) CreatePasscode(userID int64, passcode string) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, passcode) VALUES (?, ?)`, tables.Passcodes)
	_, err := r.DB.Exec(query, userID, passcode)
	return err
}

func (r *UserRepository) GetPasscode(userID int64) (string, error) {
	query := fmt.Sprintf(`SELECT passcode FROM %s WHERE id = ?`, tables.Passcodes)
	var passcode string
	err := r.DB.Select(passcode, query, userID)
	return passcode, err
}

func (r *UserRepository) DeletePasscode(id int64, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, tables.Passcodes)
	_, err := tx.Exec(query, id)
	return err
}

func (r *UserRepository) StartUserRegistration(tx *sqlx.Tx, model *authmodels.RegInitCreate) error {
	query := fmt.Sprintf(`
	INSERT INTO %s (id, email) VALUES (:id, :email)`, tables.Users)

	_, err := tx.NamedExec(query, model)
	return err
}
