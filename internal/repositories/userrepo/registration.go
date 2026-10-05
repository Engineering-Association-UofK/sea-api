package userrepo

import (
	"database/sql"
	"fmt"
	"sea-api/internal/models"
	"sea-api/internal/models/tables"

	"github.com/jmoiron/sqlx"
)

func (r *UserRepository) GetAllTempUsers(limit int64, page int64) ([]models.TempUserModel, error) {
	var users []models.TempUserModel
	offset := (page - 1) * limit
	err := r.DB.Select(&users, fmt.Sprintf(`SELECT * FROM %s LIMIT ? OFFSET ?`, tables.TempUsers), limit, offset)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetTempUsersWithNullPasswords() ([]models.TempUserModel, error) {
	var users []models.TempUserModel
	err := r.DB.Select(&users, fmt.Sprintf(`SELECT * FROM %s WHERE password IS NULL OR password = ''`, tables.TempUsers))
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) CreateTempUser(id int64, tx *sqlx.Tx) error {
	user := models.TempUserModel{ID: sql.NullInt64{Int64: id, Valid: true}}
	query := fmt.Sprintf(`
	INSERT INTO %s (
		id, password
	) VALUES (
		:id, 
	)`, tables.TempUsers)
	if tx != nil {
		_, err := tx.NamedExec(query, user)
		return err
	}
	_, err := r.DB.NamedExec(query, user)
	return err
}

func (r *UserRepository) GetTempUser(id int64) (*models.TempUserModel, error) {
	var user models.TempUserModel
	err := r.DB.Get(&user, fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, tables.TempUsers), id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) StartUserRegistration(tx *sqlx.Tx, model *models.RegInitCreate) error {
	query := fmt.Sprintf(`
	INSERT INTO %s (id, email) VALUES (:id, :email)`, tables.Users)

	_, err := tx.NamedExec(query, model)
	return err
}

func (r *UserRepository) DeleteTempUser(id int64, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, tables.TempUsers)
	if tx != nil {
		_, err := tx.Exec(query, id)
		return err
	}
	_, err := r.DB.Exec(query, id)
	return err
}
