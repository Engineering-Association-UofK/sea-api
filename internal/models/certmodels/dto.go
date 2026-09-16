package certmodels

import (
	"mime/multipart"
	"sea-api/internal/models"
	"time"
)

type CertResponse struct {
	ID           int64  `json:"id"`
	TemplateName string `json:"template_name"`
	TemplateID   int64  `json:"template_id"`
	IssuerID     int64  `json:"issuer_id"`

	RecipientName  string    `json:"recipient_name"`
	RecipientEmail string    `json:"recipient_email"`
	IssuedDate     time.Time `json:"issued_date"`

	RecipientUserID *int64 `json:"recipient_user_id"`
	EventID         *int64 `json:"event_id"`
}

type CertListResponse struct {
	models.ListResponse
	List []CertResponse `json:"list"`
}

type CertViewResponse struct {
	ID         int64     `json:"id"`
	IssuedDate time.Time `json:"issued_date"`
	EventID    *int64    `json:"event_id"`
}

type CertViewListResponse struct {
	models.ListResponse
	List []CertViewResponse `json:"list"`
}

type UpdateRequest struct {
	ID int64 `json:"id" binding:"required"`

	RecipientEmail  string `json:"recipient_email"`
	RecipientUserID int64  `json:"recipient_user_id"`
	EventID         int64  `json:"event_id"`
}

type CertListRequest struct {
	models.ListRequest
	UserID          int64     `form:"user-id"`
	EventID         int64     `form:"event-id"`
	IssueDateAfter  time.Time `form:"issue-date-after"`
	IssueDateBefore time.Time `form:"issue-date-before"`
	SearchName      string    `form:"search-name"`
}

// Operations

type IssueRequest struct {
	TemplateID     int64  `form:"template_id" binding:"required"`
	RecipientName  string `form:"recipient_name" binding:"required"`
	RecipientEmail string `form:"recipient_email" binding:"required"`

	SignerNameOne string `form:"signer_name_one" binding:"required"`
	SignerRoleOne string `form:"signer_role_one" binding:"required"`

	SignerNameTwo string `form:"signer_name_two" binding:"required"`
	SignerRoleTwo string `form:"signer_role_two" binding:"required"`

	SignerSignatureOne multipart.FileHeader `form:"signer_signature_one" binding:"required"`
	SignerSignatureTwo multipart.FileHeader `form:"signer_signature_two" binding:"required"`

	EventID         *int64 `form:"event_id"`
	RecipientUserID *int64 `form:"recipient_user_id"`
}

type VerifyResponse struct {
	Valid     bool      `json:"valid"`
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	IssueDate time.Time `json:"issue_date"`
	EventID   int64     `json:"event_id"`
}

// Template

type CreateTemplateRequest struct {
	Name     string          `json:"name" binding:"required"`
	Language models.Language `json:"language" binding:"required"`
	Version  string          `json:"version" binding:"required"`
	Layout   Layout          `json:"layout" binding:"required"`
}

type UpdateTemplateRequest struct {
	ID int64 `json:"id" binding:"required"`
	CreateTemplateRequest
}

type TemplateResponse struct {
	ID        int64           `json:"id"`
	Name      string          `json:"name"`
	Version   string          `json:"version"`
	Language  models.Language `json:"language"`
	Layout    Layout          `json:"layout"`
	CreatedAt time.Time       `json:"created_at"`
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
