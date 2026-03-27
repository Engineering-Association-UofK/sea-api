package repositories

import (
	"sea-api/internal/models"

	"github.com/jmoiron/sqlx"
)

type StoreRepository struct {
	DB *sqlx.DB
}

func NewStoreRepository(db *sqlx.DB) *StoreRepository {
	return &StoreRepository{DB: db}
}

// ======== The old file system ========

// Deprecated: Old system used fids and was SeaweedFS reliant
// The new system focuses of S3 compatibility
// Use CreateFile in FileRepository
func (r *StoreRepository) Create(item models.StoreModel) (int64, error) {
	query := `
	INSERT INTO store (fid, size, mime)
	VALUES (:fid, :size, :mime)
	`
	res, err := r.DB.NamedExec(query, &item)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// Deprecated: Old system used fids and was SeaweedFS reliant
// The new system focuses of S3 compatibility
// Use GetFileById in FileRepository
func (r *StoreRepository) GetById(id int64) (*models.StoreModel, error) {
	var item models.StoreModel
	err := r.DB.Get(&item, `SELECT * FROM store WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Deprecated: Old system used fids and was SeaweedFS reliant
// The new system focuses of S3 compatibility
// Use GetFileByKey in FileRepository
func (r *StoreRepository) GetByFid(fid string) (*models.StoreModel, error) {
	var item models.StoreModel
	err := r.DB.Get(&item, `SELECT * FROM store WHERE fid = ?`, fid)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Deprecated: Old system used fids and was SeaweedFS reliant
// The new system focuses of S3 compatibility
// Use GetAllFiles in FileRepository
func (r *StoreRepository) GetAll() ([]models.StoreModel, error) {
	var items []models.StoreModel
	err := r.DB.Select(&items, `SELECT * FROM store`)
	if err != nil {
		return nil, err
	}
	return items, nil
}

// Deprecated: Old system used fids and was SeaweedFS reliant
// The new system focuses of S3 compatibility
// Use UpdateFile in FileRepository
func (r *StoreRepository) Update(item *models.StoreModel) error {
	query := `
	UPDATE store
	SET fid = :fid, size = :size, mime = :mime
	WHERE id = :id
	`
	_, err := r.DB.NamedExec(query, &item)
	return err
}

// Deprecated: Old system used fids and was SeaweedFS reliant
// The new system focuses of S3 compatibility
// Use DeleteFile in FileRepository
func (r *StoreRepository) DeleteStore(id int64) error {
	_, err := r.DB.Exec(`DELETE FROM store WHERE id = ?`, id)
	return err
}

// ======== The new file system ========

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
