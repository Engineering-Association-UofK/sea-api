package certrepo

import (
	"fmt"
	"sea-api/internal/models"
	"sea-api/internal/models/certmodels"

	"github.com/jmoiron/sqlx"
)

type CertRepository struct {
	db *sqlx.DB
}

func NewCertRepository(db *sqlx.DB) *CertRepository {
	return &CertRepository{db: db}
}

func (r *CertRepository) CreateCert(model *certmodels.Certificate) (int64, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (template_id, cert_hash, file_key, issuer_id, event_id, recipient_user_id, recipient_name, recipient_email, issued_date) 
	VALUES (:template_id, :cert_hash, :file_key, :issuer_id, :event_id, :recipient_user_id, :recipient_name, :recipient_email, :issued_date)
	`, models.TableNewCertificates)

	res, err := r.db.NamedExec(query, model)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (r *CertRepository) Update(model *certmodels.Certificate) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET template_id = :template_id, 
	    cert_hash = :cert_hash, 
	    file_key = :file_key, 
	    issuer_id = :issuer_id, 
	    event_id = :event_id,
	    recipient_user_id = :recipient_user_id, 
	    recipient_name = :recipient_name, 
	    recipient_email = :recipient_email, 
	    issued_date = :issued_date
	WHERE id = :id
	`, models.TableNewCertificates)
	_, err := r.db.NamedExec(query, model)
	return err
}

func (r *CertRepository) GetCertWithHash(hash string) (*certmodels.Certificate, error) {
	var model certmodels.Certificate
	query := fmt.Sprintf(`
	SELECT * FROM %s WHERE cert_hash = ?
	`, models.TableNewCertificates)

	err := r.db.Get(&model, query, hash)
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *CertRepository) GetCertWithID(id int64) (*certmodels.Certificate, error) {
	var model certmodels.Certificate
	query := fmt.Sprintf(`
	SELECT * FROM %s WHERE id = ?
	`, models.TableNewCertificates)

	err := r.db.Get(&model, query, id)
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *CertRepository) GetCertsList(req models.ListRequest) ([]certmodels.CertResponse, error) {
	var certs = []certmodels.CertResponse{}
	query := fmt.Sprintf(`
		SELECT 
		c.id,
		c.template_id,
		c.issuer_id,
		c.recipient_name,
		c.recipient_email,
		c.issued_date,
		c.recipient_user_id,
		c.event_id,
		t.name AS template_name
	FROM %s c
	JOIN %s t ON c.template_id = t.id
	`, models.TableNewCertificates, models.TableCertTemplate)

	err := r.db.Select(&certs, query)
	if err != nil {
		return nil, err
	}
	return certs, nil
}
