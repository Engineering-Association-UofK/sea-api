package repositories

import (
	"database/sql"
	"sea-api/internal/models"

	"github.com/jmoiron/sqlx"
)

type VerificationRepo struct {
	db *sqlx.DB
}

func NewVerificationRepo(db *sqlx.DB) *VerificationRepo {
	return &VerificationRepo{
		db: db,
	}
}

func (r *VerificationRepo) Create(verification *models.VerificationCodeModel) error {
	query := `
		INSERT INTO verification_code (code, user_id, created_at)
		VALUES (:code, :user_id, :created_at)
	`
	_, err := r.db.NamedExec(query, verification)
	return err
}

func (r *VerificationRepo) GetByCode(code string) (*models.VerificationCodeModel, error) {
	query := `
		SELECT * FROM verification_code WHERE code = ?
	`
	var verification models.VerificationCodeModel
	err := r.db.Get(&verification, query, code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &verification, nil
}

func (r *VerificationRepo) GetByUserID(user_id int64) (*models.VerificationCodeModel, error) {
	query := `
		SELECT * FROM verification_code WHERE user_id = ?
	`
	var verification models.VerificationCodeModel
	err := r.db.Get(&verification, query, user_id)
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

func (r *VerificationRepo) Delete(id int64) error {
	query := `
		DELETE FROM verification_code WHERE id = ?
	`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *VerificationRepo) Clean() error {
	query := `
	DELETE FROM verification_code
    WHERE created_at < NOW() - INTERVAL 90 MINUTE
	`
	_, err := r.db.Exec(query)
	return err
}
