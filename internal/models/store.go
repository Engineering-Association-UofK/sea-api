package models

type ObjectType int

const (
	ObjBlogPost ObjectType = iota
	ObjNews
	ObjForm
	ObjEvent
	ObjCollaborator
)

type FileModel struct {
	ID       int64  `db:"id"`
	Key      string `db:"file_key"`
	FileSize int64  `db:"file_size"`
	MimeType string `db:"mime_type"`
}
