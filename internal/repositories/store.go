package repositories

import (
	"fmt"
	"sea-api/internal/models"

	"github.com/jmoiron/sqlx"
)

type FileRepository struct {
	DB *sqlx.DB
}

func NewFileRepository(db *sqlx.DB) *FileRepository {
	return &FileRepository{DB: db}
}

func (r *FileRepository) CreateFile(item models.FileModel) (int64, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (file_key, file_size, mime_type)
	VALUES (:file_key, :file_size, :mime_type)
	`, models.TableFiles)
	res, err := r.DB.NamedExec(query, &item)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *FileRepository) GetFileById(id int64) (*models.FileModel, error) {
	var item models.FileModel
	err := r.DB.Get(&item, fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, models.TableFiles), id)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *FileRepository) GetFileByKey(key string) (*models.FileModel, error) {
	var item models.FileModel
	err := r.DB.Get(&item, fmt.Sprintf(`SELECT * FROM %s WHERE file_key = ?`, models.TableFiles), key)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *FileRepository) GetAllFiles() ([]models.FileModel, error) {
	var items []models.FileModel
	err := r.DB.Select(&items, fmt.Sprintf(`SELECT * FROM %s`, models.TableFiles))
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *FileRepository) UpdateFile(item *models.FileModel) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET file_key = :file_key, file_size = :file_size, mime_type = :mime_type
	WHERE id = :id
	`, models.TableFiles)
	_, err := r.DB.NamedExec(query, &item)
	return err
}

func (r *FileRepository) UpdateID(id int64, fileKey string) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET id = :id
	WHERE file_key = :file_key
	`, models.TableFiles)
	_, err := r.DB.NamedExec(query, map[string]interface{}{
		"id":       id,
		"file_key": fileKey,
	})
	return err
}

func (r *FileRepository) DeleteFile(id int64) error {
	_, err := r.DB.Exec(fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, models.TableFiles), id)
	return err
}
