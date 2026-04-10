package repositories

import (
	"fmt"
	"sea-api/internal/models"

	"github.com/jmoiron/sqlx"
)

type CollaboratorRepo struct {
	DB *sqlx.DB
}

func NewCollaboratorRepo(db *sqlx.DB) *CollaboratorRepo {
	return &CollaboratorRepo{DB: db}
}

func (r *CollaboratorRepo) Create(collab *models.CollaboratorModel) (int64, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (name_ar, name_en, email, signature_id)
	VALUES (:name_ar, :name_en, :email, :signature_id)
	`, models.TableCollaborators)
	res, err := r.DB.NamedExec(query, collab)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *CollaboratorRepo) GetByID(id int64) (*models.CollaboratorModel, error) {
	var collab models.CollaboratorModel
	err := r.DB.Get(&collab, fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, models.TableCollaborators), id)
	if err != nil {
		return nil, err
	}
	return &collab, nil
}

func (r *CollaboratorRepo) GetAll() ([]models.CollaboratorModel, error) {
	var collaborators []models.CollaboratorModel
	err := r.DB.Select(&collaborators, fmt.Sprintf(`SELECT * FROM %s`, models.TableCollaborators))
	if err != nil {
		return nil, err
	}
	return collaborators, nil
}

func (r *CollaboratorRepo) Update(collab *models.CollaboratorModel) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET name_ar = :name_ar, name_en = :name_en, email = :email, signature_id = :signature_id
	WHERE id = :id
	`, models.TableCollaborators)
	_, err := r.DB.NamedExec(query, collab)
	return err
}

func (r *CollaboratorRepo) Delete(id int64) error {
	_, err := r.DB.Exec(fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, models.TableCollaborators), id)
	return err
}
