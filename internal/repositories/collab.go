package repositories

import (
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
	query := `
	INSERT INTO collaborators (name_ar, name_en, email, signature_id)
	VALUES (:name_ar, :name_en, :email, :signature_id)
	`
	res, err := r.DB.NamedExec(query, collab)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *CollaboratorRepo) GetByID(id int64) (*models.CollaboratorModel, error) {
	var collab models.CollaboratorModel
	err := r.DB.Get(&collab, `SELECT * FROM collaborators WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &collab, nil
}

func (r *CollaboratorRepo) GetAll() ([]models.CollaboratorModel, error) {
	var collaborators []models.CollaboratorModel
	err := r.DB.Select(&collaborators, `SELECT * FROM collaborators`)
	if err != nil {
		return nil, err
	}
	return collaborators, nil
}

func (r *CollaboratorRepo) Update(collab *models.CollaboratorModel) error {
	query := `
	UPDATE collaborators
	SET name_ar = :name_ar, name_en = :name_en, email = :email, signature_id = :signature_id
	WHERE id = :id
	`
	_, err := r.DB.NamedExec(query, collab)
	return err
}

func (r *CollaboratorRepo) Delete(id int64) error {
	_, err := r.DB.Exec(`DELETE FROM collaborators WHERE id = ?`, id)
	return err
}
