package repositories

import (
	"fmt"
	"sea-api/internal/models"

	"github.com/jmoiron/sqlx"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) StartRegistration(userID int64) error {
	query := fmt.Sprintf(`INSERT INTO %s (user_id, step) VALUES (?, ?)`, models.TableRegistrationStep)
	_, err := r.db.Exec(query, userID, 1)
	return err
}

func (r *AuthRepository) IncrementStep(userID int64) error {
	query := fmt.Sprintf(`UPDATE %s SET step = step + 1 WHERE user_id = ?`, models.TableRegistrationStep)
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *AuthRepository) GetState(userID int64) (*models.RegistrationStepModel, error) {
	var model *models.RegistrationStepModel
	query := fmt.Sprintf(`SELECT * FROM %s WHERE user_id = ?`, models.TableRegistrationStep)
	err := r.db.Get(model, query, userID)
	return model, err
}
