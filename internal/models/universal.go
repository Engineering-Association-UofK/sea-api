package models

type TableName string

const (
	// Users

	TableUsers       TableName = "users"
	TableTempUsers   TableName = "temp_users"
	TableUserRoles   TableName = "user_roles"
	TableSuspensions TableName = "suspensions"
	TableRateLimits  TableName = "rate_limits"

	// CMS

	TableBlogPosts   TableName = "blog_posts"
	TableTeamMembers TableName = "team_members"

	// Files

	TableFiles TableName = "files"

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
)

type ListRequest struct {
	Limit int `json:"limit" binding:"required"`
	Page  int `json:"page" binding:"required"`
}
