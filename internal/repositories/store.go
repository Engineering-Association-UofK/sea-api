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

func (r *StoreRepository) GetById(id int64) (*models.StoreModel, error) {
	var item models.StoreModel
	err := r.DB.Get(&item, `SELECT * FROM store WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *StoreRepository) GetByFid(fid string) (*models.StoreModel, error) {
	var item models.StoreModel
	err := r.DB.Get(&item, `SELECT * FROM store WHERE fid = ?`, fid)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *StoreRepository) GetAll() ([]models.StoreModel, error) {
	var items []models.StoreModel
	err := r.DB.Select(&items, `SELECT * FROM store`)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *StoreRepository) Update(item *models.StoreModel) error {
	query := `
	UPDATE store
	SET fid = :fid, size = :size, mime = :mime
	WHERE id = :id
	`
	_, err := r.DB.NamedExec(query, &item)
	return err
}

func (r *StoreRepository) DeleteStore(id int64) error {
	_, err := r.DB.Exec(`DELETE FROM store WHERE id = ?`, id)
	return err
}
