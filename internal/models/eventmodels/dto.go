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

type EventPrivateListResponse struct {
	ID            int64              `json:"id"`
	Name          string             `json:"name"`
	BackgroundURL string             `json:"background_url"`
	Belonging     models.Secretariat `json:"belonging"`

	RequireApplying bool   `json:"require_applying"`
	FormID          *int64 `json:"form_id"`

	CreatedAt time.Time `json:"created_at"`
}

type EventListRequest struct {
	models.ListRequest
	Search    string             `form:"search-name"`
	Belonging models.Secretariat `form:"belonging"`
}

type CoordinatorResponse struct {
	ID      int64  `json:"id"`
	EventID int64  `json:"event_id"`
	Name    string `json:"name"`
	Role    string `json:"role"`
}

type CoordinatorListResponse struct {
	List []CoordinatorResponse `json:"list"`
	models.ListResponse
}

type EventListResponse struct {
	List []EventPrivateListResponse `json:"list"`
	models.ListResponse
}

type EventViewResponse struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	BackgroundURL   string    `json:"background_url"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	RequireApplying bool      `json:"require_applying"`

	Belonging models.Secretariat `json:"belonging"`
	Coords    []CoordResponse    `json:"coords"`
}

type EventViewListResponse struct {
	List []EventViewResponse `json:"list"`
	models.ListResponse
}

// Application

type ApplicationStatus struct {
	Applied   bool   `json:"applied"`
	Accepted  bool   `json:"accepted"`
	NeedsForm bool   `json:"needs_form"`
	FormID    *int64 `json:"form_id,omitempty"`
}

type ApplyResponse struct {
	NeedsForm bool   `json:"needs_form"`
	FormID    *int64 `json:"form_id"`
}

type ApplicationListRequest struct {
	models.ListRequest
}

type ApplicationResponse struct {
	ID        int64  `json:"id"`
	EventID   int64  `json:"event_id"`
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	PhotoURL  string `json:"photo_url"`
	FormID    int64  `json:"form_id"`
	Accepted  bool   `json:"accepted"`
	StartedAt string `json:"started_at"`
}

type ApplicationListResponse struct {
	List []ApplicationResponse `json:"list"`
	models.ListResponse
}

// Participation

type ParticipantResponse struct {
	ID       int64     `json:"id"`
	EventID  int64     `json:"event_id"`
	UserID   int64     `json:"user_id"`
	Username string    `json:"username"`
	PhotoURL string    `json:"photo_url"`
	JoinedAt time.Time `json:"joined_at"`
}

type ParticipantListRequest struct {
	EventID int64 `form:"event_id"`
	models.ListRequest
}

type ParticipantListResponse struct {
	List []ParticipantResponse `json:"list"`
	models.ListResponse
}

// Coordinators

type CoordRequest struct {
	Name string `json:"name" binding:"required"`
	Role string `json:"role" binding:"required"`
}

type CoordResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type AddCoordsRequest struct {
	List []CoordRequest `json:"list"`
}
