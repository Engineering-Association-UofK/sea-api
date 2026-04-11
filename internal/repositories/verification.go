package repositories

import (
	"database/sql"
	"fmt"
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
	query := fmt.Sprintf(`
		INSERT INTO %s (code, user_id, created_at)
		VALUES (:code, :user_id, :created_at)
	`, models.TableVerificationCode)
	_, err := r.db.NamedExec(query, verification)
	return err
}

func (r *VerificationRepo) GetByCode(code string) (*models.VerificationCodeModel, error) {
	query := fmt.Sprintf(`
		SELECT * FROM %s WHERE code = ?
	`, models.TableVerificationCode)
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
	query := fmt.Sprintf(`
		SELECT * FROM %s WHERE user_id = ?
	`, models.TableVerificationCode)
	var verification models.VerificationCodeModel
	err := r.db.Get(&verification, query, user_id)
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

func (r *VerificationRepo) Delete(id int64) error {
	query := fmt.Sprintf(`
		DELETE FROM %s WHERE id = ?
	`, models.TableVerificationCode)
	_, err := r.db.Exec(query, id)
	return err
}

func (r *VerificationRepo) Clean() error {
	query := fmt.Sprintf(`
	DELETE FROM %s
    WHERE created_at < NOW() - INTERVAL 90 MINUTE
	`, models.TableVerificationCode)
	_, err := r.db.Exec(query)
	return err
}
