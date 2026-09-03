package certmodels

import (
	"mime/multipart"
	"sea-api/internal/models"
	"time"
)

type IssueRequest struct {
	TemplateID    int64  `form:"template_id" binding:"required"`
	RecipientName string `form:"recipient_name" binding:"required"`

	SignerNameOne string `form:"signer_name_one" binding:"required"`
	SignerRoleOne string `form:"signer_role_one" binding:"required"`

	SignerNameTwo string `form:"signer_name_two" binding:"required"`
	SignerRoleTwo string `form:"signer_role_two" binding:"required"`

	SignerSignatureOne multipart.FileHeader `form:"signer_signature_one" binding:"required"`
	SignerSignatureTwo multipart.FileHeader `form:"signer_signature_two" binding:"required"`

	EventID         *int64 `form:"event_id"`
	RecipientUserID *int64 `form:"recipient_user_id"`
}

type CreateTemplateRequest struct {
	Name     string          `json:"name" binding:"required"`
	Language models.Language `json:"language", binding:"required"`
	Version  string          `json:"version" binding:"required"`
	Layout   Layout          `json:"layout" binding:"required"`
}

type UpdateTemplateRequest struct {
	ID int64 `json:"id" binding:"required"`
	CreateTemplateRequest
}

type TemplateResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Version   string    `json:"version"`
	Layout    Layout    `json:"layout"`
	CreatedAt time.Time `json:"created_at"`
}

type TemplateListResponse struct {
	Pages   int64 `json:"pages"`
	Current int64 `json:"current"`
	Total   int64 `json:"total"`

	List []TemplateResponse `json:"list"`
}

// Test

type TestImageResponse struct {
	Url string `json:"url"`
}
