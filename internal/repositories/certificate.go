package repositories

import (
	"fmt"
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
	query := fmt.Sprintf(`
	INSERT INTO %s (cert_hash, user_id, event_id, grade, issue_date, status)
	VALUES (:cert_hash, :user_id, :event_id, :grade, :issue_date, :status)
	`, models.TableCertificates)
	res, err := r.DB.NamedExec(query, &item)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *CertificateRepository) CreateFile(item models.CertificateFileModel) (int64, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (certificate_id, store_id, lang)
	VALUES (:certificate_id, :store_id, :lang)
	`, models.TableCertificates)
	res, err := r.DB.NamedExec(query, &item)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *CertificateRepository) GetAll() ([]models.CertificateModel, error) {
	var items []models.CertificateModel
	err := r.DB.Select(&items, fmt.Sprintf(`SELECT * FROM %s`, models.TableCertificates))
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *CertificateRepository) GetByID(id int64) (*models.CertificateModel, error) {
	var model models.CertificateModel
	err := r.DB.Get(&model, fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, models.TableCertificates), id)
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *CertificateRepository) GetByHash(hash string) (*models.CertificateModel, error) {
	var item models.CertificateModel
	err := r.DB.Get(&item, fmt.Sprintf(`SELECT * FROM %s WHERE cert_hash = ?`, models.TableCertificates), hash)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *CertificateRepository) GetByUserIDAndEventID(user_id, eventID int64) (*models.CertificateModel, error) {
	var item models.CertificateModel
	err := r.DB.Get(&item, fmt.Sprintf(`SELECT * FROM %s WHERE user_id = ? AND event_id = ?`, models.TableCertificates), user_id, eventID)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *CertificateRepository) GetByEventIDAndUserIDs(eventID int64, userIDs []int64) ([]models.CertificateModel, error) {
	if len(userIDs) == 0 {
		return []models.CertificateModel{}, nil
	}

	query, args, err := sqlx.In(fmt.Sprintf(`SELECT * FROM %s WHERE event_id = ? AND user_id IN (?)`, models.TableCertificates), eventID, userIDs)
	if err != nil {
		return nil, err
	}
	query = r.DB.Rebind(query)
	var items []models.CertificateModel
	err = r.DB.Select(&items, query, args...)
	return items, err
}

func (r *CertificateRepository) GetAllFiles() ([]models.CertificateFileModel, error) {
	var items []models.CertificateFileModel
	err := r.DB.Select(&items, fmt.Sprintf(`SELECT * FROM %s`, models.TableCertificateFiles))
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *CertificateRepository) GetFileById(id int64) (*models.CertificateFileModel, error) {
	var item models.CertificateFileModel
	err := r.DB.Get(&item, fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, models.TableCertificateFiles), id)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *CertificateRepository) GetFilesByCertificateID(certificateID int64) ([]models.CertificateFileModel, error) {
	var items []models.CertificateFileModel
	err := r.DB.Select(&items, fmt.Sprintf(`SELECT * FROM %s WHERE certificate_id = ?`, models.TableCertificateFiles), certificateID)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *CertificateRepository) GetFileByCertificateIDAndLang(certificateID int64, lang string) (*models.CertificateFileModel, error) {
	var item models.CertificateFileModel
	err := r.DB.Get(&item, fmt.Sprintf(`SELECT * FROM %s WHERE certificate_id = ? AND lang = ?`, models.TableCertificateFiles), certificateID, lang)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *CertificateRepository) Update(item *models.CertificateModel) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET cert_hash = :cert_hash, user_id = :user_id, event_id = :event_id, issue_date = :issue_date, status = :status
	WHERE id = :id
	`, models.TableCertificates)
	_, err := r.DB.NamedExec(query, &item)
	return err
}

func (r *CertificateRepository) UpdateFile(id int64, storeID int64) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET store_id = :store_id
	WHERE id = :id	
	`, models.TableCertificateFiles)
	_, err := r.DB.NamedExec(query, map[string]interface{}{
		"store_id": storeID,
		"id":       id,
	})
	return err
}

func (r *CertificateRepository) Delete(id int64) error {
	_, err := r.DB.Exec(`DELETE FROM certificate WHERE id = ?`, id)
	return err
}
