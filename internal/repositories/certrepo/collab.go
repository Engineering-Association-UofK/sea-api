package certrepo

import (
	"fmt"
	"sea-api/internal/models"
	"sea-api/internal/models/certmodels"
)

func (r *CertRepository) CreateCollaborator(model *certmodels.CertificateCollaborator) (int64, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (cert_id, name, role, signature_file_id, display_on_cert)
	VALUES (:cert_id, :name, :role, :signature_file_id, :display_on_cert)
	`, models.TableCertCollaborator)

	res, err := r.db.NamedExec(query, model)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *CertRepository) GetCollaboratorByID(id int64) (*certmodels.CertificateCollaborator, error) {
	var model certmodels.CertificateCollaborator
	query := fmt.Sprintf(`
	SELECT id, cert_id, name, role, signature_file_id, display_on_cert 
	FROM %s WHERE id = ?
	`, models.TableCertificateCollaborators)

	err := r.db.Get(&model, query, id)
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *CertRepository) ListCollaboratorsByCertID(certID int64) ([]certmodels.CertificateCollaborator, error) {
	var collaborators []certmodels.CertificateCollaborator
	query := fmt.Sprintf(`
	SELECT id, cert_id, name, role, signature_file_id, display_on_cert 
	FROM %s WHERE cert_id = ? ORDER BY id ASC
	`, models.TableCertificateCollaborators)

	err := r.db.Select(&collaborators, query, certID)
	if err != nil {
		return nil, err
	}
	return collaborators, nil
}

func (r *CertRepository) UpdateCollaborator(model *certmodels.CertificateCollaborator) error {
	query := fmt.Sprintf(`
	UPDATE %s 
	SET cert_id = :cert_id, name = :name, role = :role, signature_file_id = :signature_file_id, display_on_cert = :display_on_cert
	WHERE id = :id
	`, models.TableCertificateCollaborators)

	_, err := r.db.NamedExec(query, model)
	return err
}

func (r *CertRepository) DeleteCollaborator(id int64) error {
	query := fmt.Sprintf(`
	DELETE FROM %s WHERE id = ?
	`, models.TableCertificateCollaborators)

	_, err := r.db.Exec(query, id)
	return err
}
