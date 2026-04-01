package models

import "database/sql"

type CollaboratorModel struct {
	ID          int64         `db:"id"`
	NameAr      string        `db:"name_ar"`
	NameEn      string        `db:"name_en"`
	Email       string        `db:"email"`
	SignatureID sql.NullInt64 `db:"signature_id"`
}

type CollaboratorCreateRequest struct {
	NameAr string `json:"name_ar"`
	NameEn string `json:"name_en"`
	Email  string `json:"email"`
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
