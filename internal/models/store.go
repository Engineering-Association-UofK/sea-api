package models

type Path string

const (
	IMAGES Path = "images"
	DOCS   Path = "docs"
	ASSETS Path = "assets"
)

type ObjectType int

const (
	ObjBlogPost ObjectType = iota
	ObjNews
	ObjForm
	ObjEvent
	ObjCollaborator
)

type StoreModel struct {
	ID   int64  `db:"id"`
	Fid  string `db:"fid"`
	Size int64  `db:"size"`
	Mime string `db:"mime"`
}

type FileModel struct {
	ID       int64  `db:"id"`
	Key      string `db:"file_key"`
	FileSize int64  `db:"file_size"`
	MimeType string `db:"mime_type"`
}
