package models

import (
	"database/sql"
	"time"
)

type EventType string
type ParticipantStatus string
type EventStatus string

const (
	WORKSHOP EventType = "WORKSHOP"
	COURSE   EventType = "COURSE"
	SEMINAR  EventType = "SEMINAR"

	PENDING   ParticipantStatus = "PENDING"
	ACCEPTED  ParticipantStatus = "ACCEPTED"
	REJECTED  ParticipantStatus = "REJECTED"
	COMPLETED ParticipantStatus = "COMPLETED"

	EventUpcoming  EventStatus = "UPCOMING"
	EventOngoing   EventStatus = "ONGOING"
	EventCompleted EventStatus = "COMPLETED"
)

// ====== DataBase Models ======

type EventModel struct {
	ID              int64         `db:"id"`
	Name            string        `db:"name"`
	Description     string        `db:"description"`
	CoordinatorID   int64         `db:"coordinator_id"`
	PresenterID     int64         `db:"presenter_id"`
	EventType       EventType     `db:"event_type"`
	WallpaperID     sql.NullInt64 `db:"wallpaper_id"`
	FormApplication bool          `db:"form_application"`
	MaxParticipants int           `db:"max_participants"`
	StartDate       time.Time     `db:"start_date"`
	EndDate         time.Time     `db:"end_date"`
	Outcomes        string        `db:"outcomes"`
}

type EventComponentModel struct {
	ID          int64   `db:"id"`
	EventID     int64   `db:"event_id"`
	Name        string  `db:"name"`
	Description string  `db:"description"`
	MaxScore    float64 `db:"max_score"`
}

type EventParticipantModel struct {
	ID        int64             `db:"id"`
	EventID   int64             `db:"event_id"`
	UserID    int64             `db:"user_id"`
	Grade     float64           `db:"grade"`
	Status    ParticipantStatus `db:"status"`
	JoinedAt  time.Time         `db:"joined_at"`
	Completed bool              `db:"completed"`
}

type ComponentScoreModel struct {
	ID            int64   `db:"id"`
	ParticipantID int64   `db:"participant_id"`
	ComponentID   int64   `db:"component_id"`
	Score         float64 `db:"score"`
}

type EventFormModel struct {
	ID      int64 `db:"id"`
	FormID  int64 `db:"form_id"`
	EventID int64 `db:"event_id"`
}

type EventApplicationModel struct {
	ID          int64     `db:"id"`
	EventID     int64     `db:"event_id"`
	UserID      int64     `db:"user_id"`
	Status      string    `db:"status"`
	SubmittedAt time.Time `db:"submitted_at"`
}

////////////////
///   VIEW   ///
////////////////

// LIST VIEW

type EventViewListResponse struct {
	Current int64                       `json:"current"`
	Pages   int64                       `json:"pages"`
	Events  []EventViewListItemResponse `json:"events"`
}

type EventViewListItemRow struct {
	ID                int64          `db:"id"`
	Name              string         `db:"name"`
	PresenterID       int64          `db:"presenter_id"`
	EventType         EventType      `db:"event_type"`
	WallpaperFileKey  sql.NullString `db:"wallpaper_file_key"`
	MaxParticipants   int            `db:"max_participants"`
	ParticipantsCount int            `db:"participants_count"`
	Status            EventStatus    `db:"status"`
	StartDate         time.Time      `db:"start_date"`
	EndDate           time.Time      `db:"end_date"`
}

type EventViewListItemResponse struct {
	ID                int64       `json:"id" db:"id"`
	Name              string      `json:"name" db:"name"`
	PresenterID       int64       `json:"presenter_id" db:"presenter_id"`
	EventType         EventType   `json:"event_type" db:"event_type"`
	Wallpaper         string      `json:"wallpaper"`
	MaxParticipants   int         `json:"max_participants" db:"max_participants"`
	ParticipantsCount int         `json:"participants_count" db:"participants_count"`
	Status            EventStatus `json:"status" db:"status"`
	StartDate         time.Time   `json:"start_date" db:"start_date"`
	EndDate           time.Time   `json:"end_date" db:"end_date"`
}

// Pagination & filters

type QueryEventPublicRequest struct {
	Limit  int64       `form:"limit" binding:"required"`
	Page   int64       `form:"page" binding:"required"`
	Type   EventType   `form:"type"`
	Status EventStatus `form:"status"`
}

// DETAILS VIEW

type EventViewDetailsResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	EventType   EventType `json:"event_type"`
	Outcomes    []string  `json:"outcomes"`

	Wallpaper       sql.NullString `json:"wallpaper"`
	FormApplication bool           `json:"form_application"`
	MaxParticipants int            `json:"max_participants"`

	Schedule      EventSchedule      `json:"schedule"`
	Presenters    []PresenterSummary `json:"presenter"`
	GradingScheme []ComponentDTO     `json:"grading_scheme"`
}

type EventDetailsAdminRow struct {
	ID                int64          `db:"id"`
	Name              string         `db:"name"`
	Description       string         `db:"description"`
	EventType         EventType      `db:"event_type"`
	Outcomes          string         `db:"outcomes"`
	WallpaperID       sql.NullInt64  `db:"wallpaper_id"`
	WallpaperFileKey  sql.NullString `db:"wallpaper_file_key"`
	FormApplication   bool           `db:"form_application"`
	FormID            sql.NullInt64  `db:"form_id"`
	MaxParticipants   int            `db:"max_participants"`
	ParticipantsCount int            `db:"participants_count"`
	StartDate         time.Time      `db:"start_date"`
	EndDate           time.Time      `db:"end_date"`
	PresentersJSON    []byte         `db:"presenters_json"`
	CoordinatorsJSON  []byte         `db:"coordinators_json"`
	GradingJSON       []byte         `db:"grading_json"`
	UserRegJSON       []byte         `db:"user_reg_json"`
}

type EventDetailsAdminResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	EventType   EventType `json:"event_type"`
	Outcomes    []string  `json:"outcomes"`

	WallpaperID int64  `json:"wallpaper_id"`
	Wallpaper   string `json:"wallpaper"`

	FormApplication bool  `json:"form_application"`
	FormID          int64 `json:"form_id"`

	MaxParticipants   int           `json:"max_participants"`
	ParticipantsCount int           `json:"participants_count"`
	Schedule          EventSchedule `json:"schedule"`

	Presenters    []PresenterSummary   `json:"presenter"`
	Coordinators  []CoordinatorSummary `json:"coordinator"`
	GradingScheme []ComponentDTO       `json:"grading_scheme"`
}

type EventSchedule struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type PresenterSummary struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type CoordinatorSummary struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ComponentDTO struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	MaxScore    float64 `json:"max_score"`
}

// PARTICIPANTS

type EventParticipantsResponse struct {
	EventID      int64                `json:"event_id"`
	Participants []ParticipantDetails `json:"participants"`
	Page         int64                `json:"page"`
	Current      int64                `json:"current"`
}

type ParticipantDetails struct {
	RegistrationID int64             `json:"registration_id"`
	UserID         int64             `json:"user_id"`
	NameAr         string            `json:"name_ar"`
	NameEn         string            `json:"name_en"`
	JoinedAt       time.Time         `json:"joined_at"`
	Status         ParticipantStatus `json:"status"`
	Grades         []GradeDetail     `json:"grades"`
	Completed      bool              `json:"completed"`
}

type GradeDetail struct {
	ComponentID int64   `json:"component_id"`
	Score       float64 `json:"score"`
}

type ApplicationStatus struct {
	EventID   int64  `json:"event_id" db:"event_id"`
	EventName string `json:"event_name" db:"event_name"`
	Status    string `json:"status" db:"status"`
}

type ApplicationStatusList struct {
	Current      int64               `json:"current"`
	Pages        int64               `json:"pages"`
	Applications []ApplicationStatus `json:"applications"`
}

////////////////
///  CREATE  ///
////////////////

type EventCreateRequest struct {
	Name            string         `json:"name" binding:"required"`
	Description     string         `json:"description" binding:"required"`
	PresenterID     int64          `json:"presenter_id" binding:"required"`
	CoordinatorID   int64          `json:"coordinator_id" binding:"required"`
	EventType       EventType      `json:"event_type" binding:"required"`
	WallpaperID     int64          `json:"wallpaper_id"`
	FormApplication bool           `json:"form_application"`
	MaxParticipants int            `json:"max_participants" binding:"required"`
	StartDate       time.Time      `json:"start_date" binding:"required"`
	EndDate         time.Time      `json:"end_date" binding:"required"`
	Outcomes        []string       `json:"outcomes" binding:"required"`
	Components      []ComponentDTO `json:"components"`
}

////////////////
///  UPDATE  ///
////////////////

type EventUpdateRequest struct {
	ID int64 `json:"id"`
	EventCreateRequest
}

type ParticipantUpdateRequest struct {
	ID        int64             `json:"id" binding:"required"`
	Status    ParticipantStatus `json:"status" binding:"required"`
	Completed bool              `json:"completed"`
	Grades    []GradeDetail     `json:"grades"`
}

type MakeCertificatesForEventRequest struct {
	EventID            int64       `json:"event_id" binding:"required"`
	CertificateType    CertType    `json:"certificate_type" binding:"required"`
	CertificateVersion CertVersion `json:"certificate_version" binding:"required"`
}

type ApplyResponse struct {
	FormRequired bool  `json:"form_required"`
	FormID       int64 `json:"form_id"`
}

type EventFormRequest struct {
	EventID int64 `json:"event_id"`
	FormID  int64 `json:"form_id"`
}
