package certrepo

import (
	"fmt"
	"sea-api/internal/models"
	"sea-api/internal/models/certmodels"
)

func (r *CertRepository) CreateTemplate(model *certmodels.CertificateTemplate) (int64, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (name, version, background_image_file_id, layout_config)
	VALUES (:name, :version, :background_image_file_id, :layout_config)
	`, models.TableCertTemplate)

	res, err := r.db.NamedExec(query, model)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *CertRepository) GetTemplateByID(id int64) (*certmodels.CertificateTemplate, error) {
	var model certmodels.CertificateTemplate
	query := fmt.Sprintf(`
	SELECT id, name, version, background_image_file_id, layout_config 
	FROM %s WHERE id = ?
	`, models.TableCertTemplate)

	err := r.db.Get(&model, query, id)
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *CertRepository) ListTemplates() ([]certmodels.CertificateTemplate, error) {
	var templates = []certmodels.CertificateTemplate{}
	query := fmt.Sprintf(`
	SELECT id, name, version, background_image_file_id, layout_config 
	FROM %s ORDER BY id DESC
	`, models.TableCertTemplate)

	err := r.db.Select(&templates, query)
	if err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *CertRepository) UpdateTemplate(model *certmodels.CertificateTemplate) error {
	query := fmt.Sprintf(`
	UPDATE %s 
	SET name = :name, version = :version, background_image_file_id = :background_image_file_id, layout_config = :layout_config
	WHERE id = :id
	`, models.TableCertTemplate)

	_, err := r.db.NamedExec(query, model)
	return err
}

func (r *CertRepository) DeleteTemplate(id int64) error {
	query := fmt.Sprintf(`
	DELETE FROM %s WHERE id = ?
	`, models.TableCertTemplate)

	_, err := r.db.Exec(query, id)
	return err
}
