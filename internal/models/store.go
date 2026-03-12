package models

type Path string

const (
	IMAGES Path = "images"
	DOCS   Path = "docs"
	ASSETS Path = "assets"
)

type StoreModel struct {
	ID   int64  `db:"id"`
	Fid  string `db:"fid"`
	Size int64  `db:"size"`
	Mime string `db:"mime"`
}
