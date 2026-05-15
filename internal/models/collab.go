package models

import (
	"database/sql"
	"mime/multipart"
)

type CollaboratorModel struct {
	ID          int64         `db:"id"`
	NameAr      string        `db:"name_ar"`
	NameEn      string        `db:"name_en"`
	TitleAr     string        `db:"title_ar"`
	TitleEn     string        `db:"title_en"`
	Email       string        `db:"email"`
	SignatureID sql.NullInt64 `db:"signature_id"`
}

type CollaboratorCreateRequest struct {
	NameAr  string `form:"name_ar" binding:"required"`
	NameEn  string `form:"name_en" binding:"required"`
	TitleAr string `form:"title_ar" binding:"required"`
	TitleEn string `form:"title_en" binding:"required"`
	Email   string `form:"email" binding:"email"`

	SignatureFile *multipart.FileHeader `form:"signature_file" binding:"required"`
}

type CollaboratorUpdateRequest struct {
	ID int64 `json:"id"`
	CollaboratorCreateRequest
}

type CollaboratorResponse struct {
	ID           int64  `json:"id"`
	NameAr       string `json:"name_ar"`
	NameEn       string `json:"name_en"`
	Email        string `json:"email"`
	SignatureUrl string `json:"signature_url"`
}
