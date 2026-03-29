package repositories

import (
	"sea-api/internal/models"

	"github.com/jmoiron/sqlx"
)

type GalleryRepository struct {
	db *sqlx.DB
}

func NewGalleryRepository(db *sqlx.DB) *GalleryRepository {
	return &GalleryRepository{db: db}
}

func (r *GalleryRepository) CreateAsset(asset *models.GalleryAssetModel) (int64, error) {
	query := `
	INSERT INTO gallery_assets (file_id, file_name, alt_text, uploaded_by, showcase, created_at)
	VALUES (:file_id, :file_name, :alt_text, :uploaded_by, :showcase, :created_at)
	`
	res, err := r.db.NamedExec(query, asset)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *GalleryRepository) GetAssetByID(id int64) (*models.GalleryAssetModel, error) {
	var asset models.GalleryAssetModel
	err := r.db.Get(&asset, `SELECT * FROM gallery_assets WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *GalleryRepository) GetAllAssets() ([]models.GalleryAssetModel, error) {
	var assets []models.GalleryAssetModel
	err := r.db.Select(&assets, `SELECT * FROM gallery_assets ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	return assets, nil
}

func (r *GalleryRepository) GetAllGallery() ([]models.GallerySqlModel, error) {
	var assets []models.GallerySqlModel
	err := r.db.Select(&assets, `
	SELECT a.id, a.file_id, a.file_name, a.alt_text, a.uploaded_by, a.created_at,
	COUNT(r.asset_id) AS reference_times
	FROM gallery_assets a
	LEFT JOIN gallery_references r ON a.id = r.asset_id
	GROUP BY a.id, a.file_id, a.file_name, a.alt_text, a.uploaded_by, a.created_at
	ORDER BY a.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	return assets, nil
}

func (r *GalleryRepository) GetUnreferencedAssetIDs() ([]models.GalleryAssetModel, error) {
	var models []models.GalleryAssetModel
	query := `
	SELECT a.id, a.file_id, a.file_name, a.alt_text, a.uploaded_by, a.showcase, a.created_at
	FROM gallery_assets a
	LEFT JOIN gallery_references r ON a.id = r.asset_id
	WHERE r.asset_id IS NULL
	`
	err := r.db.Select(&models, query)
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (r *GalleryRepository) UpdateAsset(asset *models.GalleryAssetModel) error {
	query := `
	UPDATE gallery_assets
	SET alt_text = :alt_text, file_name = :file_name
	WHERE id = :id
	`
	_, err := r.db.NamedExec(query, asset)
	return err
}

func (r *GalleryRepository) DeleteAsset(id int64) error {
	_, err := r.db.Exec(`DELETE FROM gallery_assets WHERE id = ?`, id)
	return err
}

// ======== REFERENCES ========

func (r *GalleryRepository) CreateReference(ref *models.GalleryReferenceModel) (int64, error) {
	query := `
	INSERT INTO gallery_references (asset_id, object_type, object_id)
	VALUES (:asset_id, :object_type, :object_id)
	`
	res, err := r.db.NamedExec(query, ref)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *GalleryRepository) GetReferenceByID(id int64) (*models.GalleryReferenceModel, error) {
	var ref models.GalleryReferenceModel
	err := r.db.Get(&ref, `SELECT * FROM gallery_references WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &ref, nil
}

func (r *GalleryRepository) GetReferencesByAssetID(assetID int64) ([]models.GalleryReferenceModel, error) {
	var refs []models.GalleryReferenceModel
	err := r.db.Select(&refs, `SELECT * FROM gallery_references WHERE asset_id = ?`, assetID)
	if err != nil {
		return nil, err
	}
	return refs, nil
}

func (r *GalleryRepository) GetReferencesByObjectType(objectType models.ObjectType) ([]models.GalleryReferenceModel, error) {
	var refs []models.GalleryReferenceModel
	err := r.db.Select(&refs, `SELECT * FROM gallery_references WHERE object_type = ?`, objectType)
	if err != nil {
		return nil, err
	}
	return refs, nil
}

func (r *GalleryRepository) GetReferenceByObject(objectType models.ObjectType, objectID int64) (*models.GalleryReferenceModel, error) {
	var ref models.GalleryReferenceModel
	err := r.db.Get(&ref, `SELECT * FROM gallery_references WHERE object_type = ? AND object_id = ?`, objectType, objectID)
	if err != nil {
		return nil, err
	}
	return &ref, nil
}

func (r *GalleryRepository) DeleteReferencesByAsset(assetID int64) error {
	_, err := r.db.Exec(`DELETE FROM gallery_references WHERE asset_id = ?`, assetID)
	return err
}

func (r *GalleryRepository) DeleteReferencesByObject(objectType models.ObjectType, objectID int64) error {
	_, err := r.db.Exec(`DELETE FROM gallery_references WHERE object_type = ? AND object_id = ?`, objectType, objectID)
	return err
}

func (r *GalleryRepository) DeleteReference(id int64) error {
	_, err := r.db.Exec(`DELETE FROM gallery_references WHERE id = ?`, id)
	return err
}
