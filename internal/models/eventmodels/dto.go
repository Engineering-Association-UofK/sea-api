package eventmodels

import (
	"sea-api/internal/models"
	"time"
)

// Events

type EventRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description" binding:"required"`
	BackgroundID int64  `json:"background_id" binding:"required"`

	Belonging models.Secretariat `json:"belonging" binding:"required"`

	RequireApplying bool   `json:"require_applying"`
	FormID          *int64 `json:"form_id"`
	MaxApplications *int64 `json:"max_applications"`

	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

type EventUpdateRequest struct {
	ID int64 `json:"id" binding:"required"`
	EventRequest
}

type EventResponse struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	BackgroundID  int64  `json:"background_id"`
	BackgroundURL string `json:"background_url"`

	Belonging models.Secretariat `json:"belonging"`

	RequireApplying bool   `json:"require_applying"`
	FormID          *int64 `json:"form_id"`
	MaxApplications *int64 `json:"max_applications"`

	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type EventListRequest struct {
	models.ListRequest
	Search    string             `form:"search-name"`
	Belonging models.Secretariat `form:"belonging"`
}

type EventCoordResponse struct {
	ID      int64  `json:"id"`
	EventID int64  `json:"event_id"`
	Name    string `json:"name"`
	Role    string `json:"role"`
}

type EventCoordListResponse struct {
	List []EventCoordResponse `json:"list"`
	models.ListResponse
}

type EventListResponse struct {
	List []EventResponse `json:"list"`
	models.ListResponse
}

type EventViewResponse struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	BackgroundURL string `json:"background_url"`

	Belonging models.Secretariat `json:"belonging"`

	RequireApplying bool `json:"require_applying"`

	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type EventListItemResponse struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	BackgroundURL string `json:"background_url"`

	Belonging models.Secretariat `json:"belonging"`

	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type EventViewListResponse struct {
	List []EventListItemResponse `json:"list"`
	models.ListResponse
}

// Participation

type ApplyResponse struct {
	NeedsForm bool   `json:"needs_form"`
	FormID    *int64 `json:"form_id"`
}

type EventApplicationListRequest struct {
	models.ListRequest
}

type EventApplicationResponse struct {
	ID        int64  `json:"id"`
	EventID   int64  `json:"event_id"`
	UserID    int64  `json:"user_id"`
	FormID    int64  `json:"form_id"`
	Accepted  bool   `json:"accepted"`
	StartedAt string `json:"started_at"`
}

type EventApplicationListResponse struct {
	List []EventApplication `json:"list"`
	models.ListResponse
}

// Participation Pagination

type EventParticipantListRequest struct {
	EventID int64 `form:"event_id"`
	models.ListRequest
}

type EventParticipantListResponse struct {
	List []EventParticipant `json:"list"`
	models.ListResponse
}

// Coordinators

type CoordRequest struct {
	Name string `json:"name" binding:"required"`
	Role string `json:"role" binding:"required"`
}

type AddCoordsRequest struct {
	List []CoordRequest `json:"list"`
}
