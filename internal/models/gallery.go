package models

import (
	"mime/multipart"
	"time"
)

type GalleryAssetModel struct {
	ID         int64     `db:"id" json:"id"`
	FileID     int64     `db:"file_id" json:"file_id"`
	FileName   string    `db:"file_name" json:"file_name"`
	AltText    string    `db:"alt_text" json:"alt_text"`
	UploadedBy int64     `db:"uploaded_by" json:"uploaded_by"`
	Showcase   bool      `db:"showcase" json:"showcase"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type GalleryReferenceModel struct {
	ID         int64      `db:"id" json:"id"`
	AssetID    int64      `db:"asset_id" json:"asset_id"`
	ObjectType ObjectType `db:"object_type" json:"object_type"`
	ObjectID   int64      `db:"object_id" json:"object_id"`
}

type NewGalleryAssetRequest struct {
	FileName string                `form:"file_name" binding:"required"`
	AltText  string                `form:"alt_text" binding:"required"`
	File     *multipart.FileHeader `form:"file" binding:"required"`
}

type GallerySqlModel struct {
	ID             int64     `db:"id"`
	FileKey        string    `db:"file_key"`
	ReferenceTimes int64     `db:"reference_times"`
	FileName       string    `db:"file_name"`
	AltText        string    `db:"alt_text"`
	UploadedBy     int64     `db:"uploaded_by"`
	CreatedAt      time.Time `db:"created_at"`
}

type GalleryAssetResponse struct {
	ID             int64     `json:"id"`
	URL            string    `json:"url"`
	ReferenceTimes int64     `json:"reference_times"`
	FileName       string    `json:"file_name"`
	AltText        string    `json:"alt_text"`
	UploadedBy     int64     `json:"uploaded_by"`
	CreatedAt      time.Time `json:"created_at"`
}

type GalleryListRequest struct {
	// The images
	Images []GalleryAssetResponse `json:"images"`
	// The total number there is in gallery
	Total int64 `json:"total"`
	// The current page
	Current int64 `form:"current"`
	// Total number of pages (determined by total and limit of the page)
	Page int64 `form:"page"`
}
