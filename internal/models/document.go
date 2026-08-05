package models

import (
	"mime/multipart"
	"time"
)

type DocumentType string

const (
	DocGratitude DocumentType = "gratitude"
)

type DocumentModel struct {
	ID       int64        `db:"id"`
	DocHash  string       `db:"doc_hash"`
	FileID   int64        `db:"file_id"`
	Type     DocumentType `db:"type"`
	CreateAt time.Time    `db:"created_at"`
}

type DocumentRelationModel struct {
	ID          int64      `db:"id"`
	Description string     `db:"description"`
	DocumentID  int64      `db:"document_id"`
	ObjectType  ObjectType `db:"object_type"`
	ObjectID    int64      `db:"object_id"`
}

type DocumentMetadataModel struct {
	ID         int64  `db:"id"`
	DocumentID int64  `db:"document_id"`
	Key        string `db:"d_key"`
	Value      string `db:"d_value"`
}

type DocumentMetadata struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type DocumentVerifyResponse struct {
	Valid     bool         `json:"valid"`
	Type      DocumentType `json:"type"`
	CreatedAt time.Time    `json:"created_at"`

	Details []DocumentMetadata `json:"details"`
}

type SignPdfRequest struct {
	EventID int64        `form:"event_id" binding:"required"`
	Type    DocumentType `form:"type" binding:"required"`

	Metadata string `form:"metadata" binding:"required"`

	QrX float64 `form:"qr_x"`
	QrY float64 `form:"qr_y"`
	QrS float64 `form:"qr_s" binding:"required"`

	File *multipart.FileHeader `form:"file" binding:"required"`
}
