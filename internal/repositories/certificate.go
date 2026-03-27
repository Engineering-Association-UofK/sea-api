package repositories

import (
	"sea-api/internal/models"

	"github.com/jmoiron/sqlx"
)

type CertificateRepository struct {
	DB *sqlx.DB
}

func NewCertificateRepository(db *sqlx.DB) *CertificateRepository {
	return &CertificateRepository{DB: db}
}

func (r *CertificateRepository) Create(item models.CertificateModel) (int64, error) {
	query := `
	INSERT INTO certificate (cert_hash, user_id, event_id, grade, issue_date, status)
	VALUES (:cert_hash, :user_id, :event_id, :grade, :issue_date, :status)
	`
	res, err := r.DB.NamedExec(query, &item)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *CertificateRepository) CreateFile(item models.CertificateFileModel) (int64, error) {
	query := `
	INSERT INTO certificate_file (certificate_id, store_id, lang)
	VALUES (:certificate_id, :store_id, :lang)
	`
	res, err := r.DB.NamedExec(query, &item)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *CertificateRepository) GetByID(id int64) (*models.CertificateModel, error) {
	var model models.CertificateModel
	err := r.DB.Get(&model, `SELECT * FROM certificate WHERE user_id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *CertificateRepository) GetByHash(hash string) (*models.CertificateModel, error) {
	var item models.CertificateModel
	err := r.DB.Get(&item, `SELECT * FROM certificate WHERE cert_hash = ?`, hash)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *CertificateRepository) GetByUserIDAndEventID(user_id, eventID int64) (*models.CertificateModel, error) {
	var item models.CertificateModel
	err := r.DB.Get(&item, `SELECT * FROM certificate WHERE user_id = ? AND event_id = ?`, user_id, eventID)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *CertificateRepository) GetByEventIDAndUserIDs(eventID int64, userIDs []int64) ([]models.CertificateModel, error) {
	if len(userIDs) == 0 {
		return []models.CertificateModel{}, nil
	}

	query, args, err := sqlx.In(`SELECT * FROM certificate WHERE event_id = ? AND user_id IN (?)`, eventID, userIDs)
	if err != nil {
		return nil, err
	}
	query = r.DB.Rebind(query)
	var items []models.CertificateModel
	err = r.DB.Select(&items, query, args...)
	return items, err
}

func (r *CertificateRepository) GetFileById(id int64) (*models.CertificateFileModel, error) {
	var item models.CertificateFileModel
	err := r.DB.Get(&item, `SELECT * FROM certificate_file WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *CertificateRepository) GetFilesByCertificateID(certificateID int64) ([]models.CertificateFileModel, error) {
	var items []models.CertificateFileModel
	err := r.DB.Select(&items, `SELECT * FROM certificate_file WHERE certificate_id = ?`, certificateID)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *CertificateRepository) GetFileByCertificateIDAndLang(certificateID int64, lang string) (*models.CertificateFileModel, error) {
	var item models.CertificateFileModel
	err := r.DB.Get(&item, `SELECT * FROM certificate_file WHERE certificate_id = ? AND lang = ?`, certificateID, lang)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *CertificateRepository) Update(item *models.CertificateModel) error {
	query := `
	UPDATE certificate
	SET cert_hash = :cert_hash, user_id = :user_id, event_id = :event_id, issue_date = :issue_date, status = :status
	WHERE id = :id
	`
	_, err := r.DB.NamedExec(query, &item)
	return err
}

func (r *CertificateRepository) Delete(id int64) error {
	_, err := r.DB.Exec(`DELETE FROM certificate WHERE user_id = ?`, id)
	return err
}
