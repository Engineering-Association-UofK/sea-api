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

func (r *AuthRepository) StartRegistration(model *models.RegistrationStepModel) error {
	model.Step = 1
	query := fmt.Sprintf(`INSERT INTO %s (reg_code, user_id, step) 
		VALUES (:reg_code, :user_id, :step)`, models.TableRegistrationStep)

	_, err := r.db.NamedExec(query, model)
	return err
}

func (r *AuthRepository) IncrementStep(RegCode string) error {
	query := fmt.Sprintf(`UPDATE %s SET step = step + 1 WHERE reg_code = ?`, models.TableRegistrationStep)
	_, err := r.db.Exec(query, RegCode)
	return err
}

func (r *AuthRepository) DecrementStep(RegCode string) error {
	query := fmt.Sprintf(`UPDATE %s SET step = step - 1 WHERE reg_code = ?`, models.TableRegistrationStep)
	_, err := r.db.Exec(query, RegCode)
	return err
}

func (r *AuthRepository) SetStep(RegCode string, stepNumber int64) error {
	query := fmt.Sprintf(`UPDATE %s SET step = step - ? WHERE reg_code = ?`, models.TableRegistrationStep)
	_, err := r.db.Exec(query, stepNumber, RegCode)
	return err
}

func (r *AuthRepository) GetStateWithCode(RegCode string) (*models.RegistrationStepModel, error) {
	var model models.RegistrationStepModel
	query := fmt.Sprintf(`SELECT * FROM %s WHERE reg_code = ?`, models.TableRegistrationStep)
	err := r.db.Get(&model, query, RegCode)
	return &model, err
}

func (r *AuthRepository) GetStateWithID(ID int64) (*models.RegistrationStepModel, error) {
	var model models.RegistrationStepModel
	query := fmt.Sprintf(`SELECT * FROM %s WHERE user_id = ?`, models.TableRegistrationStep)
	err := r.db.Get(&model, query, ID)
	return &model, err
}
