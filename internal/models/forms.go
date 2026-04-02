package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

type QuestionType string

const (
	FORM_TEXT      QuestionType = "TEXT"
	FORM_PARAGRAPH QuestionType = "PARAGRAPH"
	FORM_RADIO     QuestionType = "RADIO"
	FORM_CHECKBOX  QuestionType = "CHECKBOX"
	FORM_DROPDOWN  QuestionType = "DROPDOWN"
	FORM_NUMBER    QuestionType = "NUMBER"
)

var AllowedQuestionTypes = map[QuestionType]bool{
	FORM_TEXT:      true,
	FORM_PARAGRAPH: true,
	FORM_RADIO:     true,
	FORM_CHECKBOX:  true,
	FORM_DROPDOWN:  true,
	FORM_NUMBER:    true,
}

type ResponseStatus string

const (
	FORM_DRAFT     ResponseStatus = "DRAFT"
	FORM_SUBMITTED ResponseStatus = "SUBMITTED"
)

// 1. The Form Container
type FormModel struct {
	ID                   int64         `db:"id" json:"id"`
	Title                string        `db:"title" json:"title"`
	Description          string        `db:"description" json:"description"`
	AllowMultipleEntries bool          `db:"allow_multiple" json:"allow_multiple"`
	IsActive             bool          `db:"is_active" json:"is_active"`
	HeaderImageID        sql.NullInt64 `db:"header_image_id" json:"header_image_id"`
	CreatedBy            int64         `db:"created_by" json:"created_by"`
	CreatedAt            time.Time     `db:"created_at" json:"created_at"`
}

type FormPageModel struct {
	ID         int64 `db:"id" json:"id"`
	FormID     int64 `db:"form_id" json:"form_id"`
	PageNumber int   `db:"page_num" json:"page_num"`
}

type FormQuestionModel struct {
	ID           int64            `db:"id" json:"id"`
	FormPageID   int64            `db:"form_page_id" json:"form_page_id"`
	QuestionText string           `db:"question_text" json:"question_text"`
	Type         QuestionType     `db:"type" json:"type"`
	Options      *json.RawMessage `db:"options" json:"options"`
	IsRequired   bool             `db:"is_required" json:"is_required"`
	DisplayOrder int              `db:"display_order" json:"display_order"`
}

type FormResponseModel struct {
	ID          int64          `db:"id" json:"id"`
	FormID      int64          `db:"form_id" json:"form_id"`
	UserID      int64          `db:"user_id" json:"user_id"`
	Status      ResponseStatus `db:"status" json:"status"`
	SubmittedAt time.Time      `db:"submitted_at" json:"submitted_at"`
}

type FormAnswerModel struct {
	ID          int64  `db:"id" json:"id"`
	ResponseID  int64  `db:"response_id" json:"response_id"`
	QuestionID  int64  `db:"question_id" json:"question_id"`
	AnswerValue string `db:"answer_value" json:"answer_value"`
}

// Options

type TextOption struct {
	Text string `json:"text"`
}

type NumberOption struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type Options []string

// Helpers

type FormRow struct {
	FormID        int64         `db:"form_id"`
	Title         string        `db:"title"`
	Description   string        `db:"description"`
	AllowMultiple bool          `db:"allow_multiple"`
	IsActive      bool          `db:"is_active"`
	HeaderImageID sql.NullInt64 `db:"header_image_id"`

	PageID  *int64 `db:"page_id"`
	PageNum *int   `db:"page_num"`

	QuestionID   *int64           `db:"question_id"`
	QuestionText *string          `db:"question_text"`
	Type         *QuestionType    `db:"type"`
	Options      *json.RawMessage `db:"options"`
	IsRequired   *bool            `db:"is_required"`
	DisplayOrder *int             `db:"display_order"`
}

type FormAnalysisRow struct {
	QuestionID   int64          `db:"question_id"`
	QuestionText string         `db:"question_text"`
	Type         QuestionType   `db:"type"`
	AnswerValue  sql.NullString `db:"answer_value"`
	AnswerCount  int            `db:"answer_count"`
}

// Full render DTOs

type FormForEditDTO struct {
	Url       string                  `json:"url"`
	Form      UpdateFormRequest       `json:"form"`
	Pages     []UpdatePageRequest     `json:"pages"`
	Questions []UpdateQuestionRequest `json:"questions"`
}

type FormForUserDTO struct {
	Form      FormDTO           `json:"form"`
	Pages     []FormPageDTO     `json:"pages"`
	Questions []FormQuestionDTO `json:"questions"`
}

// DTOs

type FormDTO struct {
	ID             int64  `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	HeaderImageUrl string `json:"header_image_url"`
}

type FormPageDTO struct {
	ID         int64 `json:"id"`
	FormID     int64 `json:"form_id"`
	PageNumber int   `json:"page_num"`
}

type FormQuestionDTO struct {
	ID           int64           `json:"id"`
	FormPageID   int64           `json:"form_page_id"`
	QuestionText string          `json:"question_text"`
	Type         QuestionType    `json:"type"`
	Options      json.RawMessage `json:"options"`
	IsRequired   bool            `json:"is_required"`
	DisplayOrder int             `json:"display_order"`
}

type FormResponseDTO struct {
	ID          int64           `json:"id"`
	FormID      int64           `json:"form_id"`
	UserID      int64           `json:"user_id"`
	Status      ResponseStatus  `json:"status"`
	SubmittedAt time.Time       `json:"submitted_at"`
	Answers     []FormAnswerDTO `json:"answers"`
}

type FormAnswerDTO struct {
	ID          int64  `json:"id"`
	ResponseID  int64  `json:"response_id"`
	QuestionID  int64  `json:"question_id"`
	AnswerValue string `json:"answer_value"`
}

// Request

type CreateFormRequest struct {
	Title                string `json:"title" binding:"required"`
	Description          string `json:"description" binding:"required"`
	AllowMultipleEntries bool   `json:"allow_multiple" binding:"required"`
	IsActive             bool   `json:"is_active" binding:"required"`
	HeaderImageID        int64  `json:"header_image_id"`
}

type UpdateFormRequest struct {
	ID int64 `json:"id" binding:"required"`
	CreateFormRequest
}

type CreatePageRequest struct {
	FormID     int64 `json:"form_id" binding:"required"`
	PageNumber int   `json:"page_num" binding:"required"`
}

type UpdatePageRequest struct {
	ID int64 `json:"id" binding:"required"`
	CreatePageRequest
}

type CreateQuestionRequest struct {
	FormPageID   int64           `json:"form_page_id" binding:"required"`
	QuestionText string          `json:"question_text" binding:"required"`
	Type         QuestionType    `json:"type" binding:"required"`
	Options      json.RawMessage `json:"options"`
	IsRequired   bool            `json:"is_required" binding:"required"`
	DisplayOrder int             `json:"display_order" binding:"required"`
}

type UpdateQuestionRequest struct {
	ID int64 `json:"id" binding:"required"`
	CreateQuestionRequest
}

type AnswerRequest struct {
	QuestionID  int64  `json:"question_id" binding:"required"`
	AnswerValue string `json:"answer_value" binding:"required"`
}

type SubmitFormRequest struct {
	FormID  int64           `json:"form_id" binding:"required"`
	Answers []AnswerRequest `json:"answers" binding:"required"`
}

type UpdateResponseStatusRequest struct {
	ID     int64          `json:"id" binding:"required"`
	Status ResponseStatus `json:"status" binding:"required"`
}

// Response

type FormSummaryResponse struct {
	ID                   int64     `json:"id"`
	Title                string    `json:"title"`
	Description          string    `json:"description"`
	IsActive             bool      `json:"is_active"`
	AllowMultipleEntries bool      `json:"allow_multiple"`
	CreatedAt            time.Time `json:"created_at"`
}

// Form analysis details

type FormAnalysisResponse struct {
	FormID         int64                 `json:"form_id"`
	Title          string                `json:"title"`
	TotalResponses int                   `json:"total_responses"`
	Questions      []QuestionAnalysisDTO `json:"questions"`
}

type QuestionAnalysisDTO struct {
	QuestionID   int64            `json:"question_id"`
	QuestionText string           `json:"question_text"`
	Type         QuestionType     `json:"type"`
	Answers      []AnswerCountDTO `json:"answers"`
}

type AnswerCountDTO struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}
