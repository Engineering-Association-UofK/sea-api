package repositories

import (
	"sea-api/internal/models"

	"github.com/jmoiron/sqlx"
)

type IFileRepository interface {
	CreateFile(item models.FileModel) (int64, error)
	GetFileById(id int64) (*models.FileModel, error)
	GetFileByKey(key string) (*models.FileModel, error)
	GetAllFiles() ([]models.FileModel, error)
	UpdateFile(item *models.FileModel) error
	DeleteFile(id int64) error
}

type FileRepository struct {
	DB *sqlx.DB
}

func NewFileRepository(db *sqlx.DB) *FileRepository {
	return &FileRepository{DB: db}
}

func (r *FileRepository) CreateFile(item models.FileModel) (int64, error) {
	query := `
	INSERT INTO files (file_key, file_size, mime_type)
	VALUES (:file_key, :file_size, :mime_type)
	`
	res, err := r.DB.NamedExec(query, &item)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *FileRepository) GetFileById(id int64) (*models.FileModel, error) {
	var item models.FileModel
	err := r.DB.Get(&item, `SELECT * FROM files WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *FileRepository) GetFileByKey(key string) (*models.FileModel, error) {
	var item models.FileModel
	err := r.DB.Get(&item, `SELECT * FROM files WHERE file_key = ?`, key)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *FileRepository) GetAllFiles() ([]models.FileModel, error) {
	var items []models.FileModel
	err := r.DB.Select(&items, `SELECT * FROM files`)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *FileRepository) UpdateFile(item *models.FileModel) error {
	query := `
	UPDATE files
	SET file_key = :file_key, file_size = :file_size, mime_type = :mime_type
	WHERE id = :id
	`
	_, err := r.DB.NamedExec(query, &item)
	return err
}

func (r *FileRepository) DeleteFile(id int64) error {
	_, err := r.DB.Exec(`DELETE FROM files WHERE id = ?`, id)
	return err
}
