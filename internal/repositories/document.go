package repositories

import (
	"fmt"
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
	query := fmt.Sprintf(`
	INSERT INTO %s (doc_hash, file_id, type, created_at)
	VALUES (:doc_hash, :file_id, :type, :created_at)
	`, models.TableDocuments)
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
	err := r.DB.Get(&doc, fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, models.TableDocuments), id)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *DocumentRepository) GetByHash(hash string) (*models.DocumentModel, error) {
	var doc models.DocumentModel
	err := r.DB.Get(&doc, fmt.Sprintf(`SELECT * FROM %s WHERE doc_hash = ?`, models.TableDocuments), hash)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *DocumentRepository) CreateRelation(rel *models.DocumentRelationModel, tx *sqlx.Tx) (int64, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (description, document_id, object_type, object_id)
	VALUES (:description, :document_id, :object_type, :object_id)
	`, models.TableDocumentRelations)
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
	err := r.DB.Select(&relations, fmt.Sprintf(`SELECT * FROM %s WHERE document_id = ?`, models.TableDocumentRelations), docID)
	if err != nil {
		return nil, err
	}
	return relations, nil
}

func (r *DocumentRepository) GetRelationsByObject(objectType models.ObjectType, objectID int64) ([]models.DocumentRelationModel, error) {
	var relations []models.DocumentRelationModel
	err := r.DB.Select(&relations, fmt.Sprintf(`SELECT * FROM %s WHERE object_type = ? AND object_id = ?`, models.TableDocumentRelations), objectType, objectID)
	if err != nil {
		return nil, err
	}
	return relations, nil
}

func (r *DocumentRepository) Delete(id int64) error {
	_, err := r.DB.Exec(fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, models.TableDocuments), id)
	return err
}

func (r *DocumentRepository) DeleteRelation(id int64) error {
	_, err := r.DB.Exec(fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, models.TableDocumentRelations), id)
	return err
}

func (r *DocumentRepository) DeleteRelationsByObject(objectType models.ObjectType, objectID int64) error {
	_, err := r.DB.Exec(fmt.Sprintf(`DELETE FROM %s WHERE object_type = ? AND object_id = ?`, models.TableDocumentRelations), objectType, objectID)
	return err
}

// =========================================
// ==========  document metadata  ==========
// =========================================

func (d *DocumentRepository) CreateMetadata(item *models.DocumentMetadataModel, tx *sqlx.Tx) (int64, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (document_id, d_key, d_value)
	VALUES (:document_id, :d_key, :d_value)
	`, models.TableDocumentMetadata)
	if tx != nil {
		res, err := tx.NamedExec(query, item)
		if err != nil {
			return 0, err
		}
		return res.LastInsertId()
	}
	res, err := d.DB.NamedExec(query, item)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (d *DocumentRepository) GetMetadataByDocumentID(documentID int64) ([]models.DocumentMetadataModel, error) {
	items := []models.DocumentMetadataModel{}
	err := d.DB.Select(&items, fmt.Sprintf(`SELECT * FROM %s WHERE document_id = ?`, models.TableDocumentMetadata), documentID)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (d *DocumentRepository) GetMetadataByID(id int64) (*models.DocumentMetadataModel, error) {
	var item models.DocumentMetadataModel
	err := d.DB.Get(&item, fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, models.TableDocumentMetadata), id)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (d *DocumentRepository) UpdateMetadata(item *models.DocumentMetadataModel, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET document_id = :document_id, d_key = :d_key, d_value = :d_value
	WHERE id = :id
	`, models.TableDocumentMetadata)
	if tx != nil {
		_, err := tx.NamedExec(query, item)
		return err
	}
	_, err := d.DB.NamedExec(query, item)
	return err
}

func (d *DocumentRepository) DeleteMetadata(id int64, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, models.TableDocumentMetadata)
	if tx != nil {
		_, err := tx.Exec(query, id)
		return err
	}
	_, err := d.DB.Exec(query, id)
	return err
}
