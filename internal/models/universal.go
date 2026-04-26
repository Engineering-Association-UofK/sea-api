package models

import "time"

type TableName string

const (
	// Users

	TableUsers     TableName = "users"
	TableTempUsers TableName = "users_temp"
	TableUserRoles TableName = "user_roles"

	// Suspensions

	TableSuspensions       TableName = "suspensions"
	TableSuspensionHistory TableName = "suspension_history"

	// CMS

	TablePosts       TableName = "posts"
	TableTeamMembers TableName = "team_members"

	// Gallery

	TableGalleryAssets     TableName = "gallery_assets"
	TableGalleryReferences TableName = "gallery_references"

	// Documents

	TableDocuments         TableName = "documents"
	TableDocumentRelations TableName = "document_relations"
	TableDocumentMetadata  TableName = "document_metadata"

	// certificates

	TableCertificates     TableName = "certificate"
	TableCertificateFiles TableName = "certificate_file"

	// Forms

	TableForms         TableName = "forms"
	TableFormPages     TableName = "form_pages"
	TableFormQuestions TableName = "form_questions"
	TableFormResponses TableName = "form_responses"
	TableFormAnswers   TableName = "form_answers"

	// Events

	TableEvents            TableName = "event"
	TableEventComponents   TableName = "event_component"
	TableEventParticipants TableName = "event_participant"
	TableComponentScores   TableName = "component_score"
	TableCollaborators     TableName = "collaborators"
	TableEventForms        TableName = "event_form"
	TableEventApplications TableName = "event_applications"

	// Bot

	TableBotNodes            TableName = "bot_nodes"
	TableBotEdges            TableName = "bot_edges"
	TableBotNodeTranslations TableName = "bot_node_translations"
	TableBotEdgeTranslations TableName = "bot_edge_translations"
	TableBotUserStates       TableName = "bot_user_states"
	TableBotActions          TableName = "bot_actions"

	// Notifications

	TableNotifications TableName = "notifications"

	// Logs

	TableLogs TableName = "logs"

	// Verification Code

	TableVerificationCode TableName = "verification_code"

	// Files

	TableFiles TableName = "files"

	// Rate Limit

	TableRateLimits TableName = "rate_limits"

	TableFeedback TableName = "feedback"
)

type ListRequest struct {
	Limit int64    `form:"limit"`
	Page  int64    `form:"page"`
	Type  PostType `form:"type"`
}

var AllowedListLimit = map[int64]bool{
	10:  true,
	25:  true,
	50:  true,
	100: true,
}

type Progress struct {
	Total     int       `json:"total"`
	Current   int       `json:"current"`
	ID        int64     `json:"id"`
	Success   bool      `json:"success"`
	Timestamp time.Time `json:"timestamp"`
	Name      string    `json:"name"`
}

type Language string

const (
	Arabic  Language = "ar"
	English Language = "en"
)

var AllowedLanguages = map[Language]bool{
	Arabic:  true,
	English: true,
}
