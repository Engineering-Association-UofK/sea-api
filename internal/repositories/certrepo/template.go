package certrepo

import (
	"fmt"
	"sea-api/internal/models"
	"sea-api/internal/models/certmodels"
	"sea-api/internal/models/tables"
)

func (r *CertRepository) CreateTemplate(model *certmodels.CertificateTemplate) (int64, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (name, language, version, layout_config, created_at)
	VALUES (:name, :language, :version, :layout_config, :created_at)
	`, tables.CertTemplate)

	res, err := r.db.NamedExec(query, model)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *CertRepository) GetTemplateByID(id int64) (*certmodels.CertificateTemplate, error) {
	var model certmodels.CertificateTemplate
	query := fmt.Sprintf(`
	SELECT * FROM %s WHERE id = ?
	`, tables.CertTemplate)

	err := r.db.Get(&model, query, id)
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *CertRepository) GetCount() int64 {
	var count int64
	err := r.db.Get(&count, fmt.Sprintf(`SELECT COUNT(*) FROM %s`, tables.CertTemplate))
	if err != nil {
		return 0
	}
	return count
}

func (r *CertRepository) ListTemplates(req *models.ListRequest) ([]certmodels.CertificateTemplate, error) {
	var templates = []certmodels.CertificateTemplate{}
	query := fmt.Sprintf(`
	SELECT * FROM %s
	ORDER BY created_at DESC
	LIMIT ? OFFSET ?
	`, tables.CertTemplate)

	offset := (req.Page - 1) * req.Limit
	err := r.db.Select(&templates, query, req.Limit, offset)
	if err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *CertRepository) UpdateTemplate(model *certmodels.CertificateTemplate) error {
	query := fmt.Sprintf(`
	UPDATE %s 
	SET name = :name, language = :language, version = :version, layout_config = :layout_config
	WHERE id = :id
	`, tables.CertTemplate)

	_, err := r.db.NamedExec(query, model)
	return err
}

func (r *CertRepository) DeleteTemplate(id int64) error {
	query := fmt.Sprintf(`
	DELETE FROM %s WHERE id = ?
	`, tables.CertTemplate)

	_, err := r.db.Exec(query, id)
	return err
}
