package certrepo

import (
	"fmt"
	"sea-api/internal/models"
	"sea-api/internal/models/certmodels"
	"strings"

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

func (r *CertRepository) GetCertsList(req *certmodels.CertListRequest) ([]certmodels.CertResponse, error) {
	certs := []certmodels.CertResponse{}
	conditions := []string{"1=1"}
	params := map[string]interface{}{}

	if req.UserID > 0 {
		conditions = append(conditions, "c.recipient_user_id = :user_id")
		params["user_id"] = req.UserID
	}

	if req.EventID > 0 {
		conditions = append(conditions, "c.event_id = :event_id")
		params["event_id"] = req.EventID
	}

	if !req.IssueDateAfter.IsZero() {
		conditions = append(conditions, "c.issued_date >= :issue_date_after")
		params["issue_date_after"] = req.IssueDateAfter
	}

	if !req.IssueDateBefore.IsZero() {
		conditions = append(conditions, "c.issued_date <= :issue_date_before")
		params["issue_date_before"] = req.IssueDateBefore
	}

	if strings.TrimSpace(req.SearchName) != "" {
		conditions = append(conditions, "c.recipient_name LIKE :search_name")
		params["search_name"] = "%" + strings.TrimSpace(req.SearchName) + "%"
	}

	whereClause := strings.Join(conditions, " AND ")

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	params["limit"] = limit
	params["offset"] = offset

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
		WHERE %s
		ORDER BY c.issued_date DESC
		LIMIT :limit OFFSET :offset
	`, models.TableNewCertificates, models.TableCertTemplate, whereClause)

	nstmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return nil, err
	}
	defer nstmt.Close()

	err = nstmt.Select(&certs, params)
	if err != nil {
		return nil, err
	}

	return certs, nil
}

func (r *CertRepository) GetCertsCount(req *certmodels.CertListRequest) (int64, error) {
	var count int64
	conditions := []string{"1=1"}
	params := map[string]interface{}{}

	if req.UserID > 0 {
		conditions = append(conditions, "c.recipient_user_id = :user_id")
		params["user_id"] = req.UserID
	}

	if req.EventID > 0 {
		conditions = append(conditions, "c.event_id = :event_id")
		params["event_id"] = req.EventID
	}

	if !req.IssueDateAfter.IsZero() {
		conditions = append(conditions, "c.issued_date >= :issue_date_after")
		params["issue_date_after"] = req.IssueDateAfter
	}

	if !req.IssueDateBefore.IsZero() {
		conditions = append(conditions, "c.issued_date <= :issue_date_before")
		params["issue_date_before"] = req.IssueDateBefore
	}

	if strings.TrimSpace(req.SearchName) != "" {
		conditions = append(conditions, "c.recipient_name LIKE :search_name")
		params["search_name"] = "%" + strings.TrimSpace(req.SearchName) + "%"
	}

	whereClause := strings.Join(conditions, " AND ")

	query := fmt.Sprintf(`
		SELECT COUNT(c.id)
		FROM %s c
		JOIN %s t ON c.template_id = t.id
		WHERE %s
	`, models.TableNewCertificates, models.TableCertTemplate, whereClause)

	nstmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return 0, err
	}
	defer nstmt.Close()

	err = nstmt.Get(&count, params)
	if err != nil {
		return 0, err
	}

	return count, nil
}
