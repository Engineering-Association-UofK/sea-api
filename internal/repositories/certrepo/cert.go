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
	INSERT INTO %s (template_id, cert_hash, issuer_id, event_id, recipient_user_id, name, title, subtitle, statement, issued_date) 
	VALUES (:template_id, :cert_hash, :issuer_id, :event_id, :recipient_user_id, :name, :title, :subtitle, :statement, :issued_date)
	`, models.TableNewCertificates)

	res, err := r.db.NamedExec(query, model)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (r *CertRepository) GetCertWithHash(hash string) (*certmodels.Certificate, error) {
	var model certmodels.Certificate
	query := fmt.Sprintf(`
	SELECT * FROM %s WHERE cert_hash = ?
	`, models.TableNewCertificates)

	err := r.db.Get(&model, query)
	if err != nil {
		return nil, err
	}
	return &model, nil
}
