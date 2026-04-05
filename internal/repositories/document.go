package repositories

import (
	"sea-api/internal/models"

	"github.com/jmoiron/sqlx"
)

type DocumentRepository struct {
	DB *sqlx.DB
}

func NewDocumentRepository(DB *sqlx.DB) *DocumentRepository {
	return &DocumentRepository{DB: DB}
}

func (r *DocumentRepository) Create(doc *models.DocumentModel, tx *sqlx.Tx) (int64, error) {
	query := `
	INSERT INTO documents (doc_hash, file_id, type, created_at)
	VALUES (:doc_hash, :file_id, :type, :created_at)
	`
	if tx != nil {
		res, err := tx.NamedExec(query, doc)
		if err != nil {
			return 0, err
		}
		return res.LastInsertId()
	}
	res, err := r.DB.NamedExec(query, doc)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *DocumentRepository) GetByID(id int64) (*models.DocumentModel, error) {
	var doc models.DocumentModel
	err := r.DB.Get(&doc, `SELECT * FROM documents WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *DocumentRepository) GetByHash(hash string) (*models.DocumentModel, error) {
	var doc models.DocumentModel
	err := r.DB.Get(&doc, `SELECT * FROM documents WHERE doc_hash = ?`, hash)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *DocumentRepository) CreateRelation(rel *models.DocumentRelationModel, tx *sqlx.Tx) (int64, error) {
	query := `
	INSERT INTO document_relations (description, document_id, object_type, object_id)
	VALUES (:description, :document_id, :object_type, :object_id)
	`
	if tx != nil {
		res, err := tx.NamedExec(query, rel)
		if err != nil {
			return 0, err
		}
		return res.LastInsertId()
	}
	res, err := r.DB.NamedExec(query, rel)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *DocumentRepository) GetRelationsByDocumentID(docID int64) ([]models.DocumentRelationModel, error) {
	var relations []models.DocumentRelationModel
	err := r.DB.Select(&relations, `SELECT * FROM document_relations WHERE document_id = ?`, docID)
	if err != nil {
		return nil, err
	}
	return relations, nil
}

func (r *DocumentRepository) GetRelationsByObject(objectType models.ObjectType, objectID int64) ([]models.DocumentRelationModel, error) {
	var relations []models.DocumentRelationModel
	err := r.DB.Select(&relations, `SELECT * FROM document_relations WHERE object_type = ? AND object_id = ?`, objectType, objectID)
	if err != nil {
		return nil, err
	}
	return relations, nil
}

func (r *DocumentRepository) Delete(id int64) error {
	_, err := r.DB.Exec(`DELETE FROM documents WHERE id = ?`, id)
	return err
}

func (r *DocumentRepository) DeleteRelation(id int64) error {
	_, err := r.DB.Exec(`DELETE FROM document_relations WHERE id = ?`, id)
	return err
}

func (r *DocumentRepository) DeleteRelationsByObject(objectType models.ObjectType, objectID int64) error {
	_, err := r.DB.Exec(`DELETE FROM document_relations WHERE object_type = ? AND object_id = ?`, objectType, objectID)
	return err
}
