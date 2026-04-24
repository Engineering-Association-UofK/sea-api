package models

import (
	"database/sql"
	"time"
)

type PostType string

const (
	PostNews       PostType = "NEWS"
	PostIssues     PostType = "ISSUES"
	PostBlog       PostType = "BLOG"
	PostPhilosophy PostType = "DONATIONS"
)

var AllowedPostTypes = map[PostType]bool{
	PostNews:       true,
	PostIssues:     true,
	PostBlog:       true,
	PostPhilosophy: true,
}

// Models

type PostModel struct {
	ID           int64          `db:"id" json:"id"`
	CoverImageID int64          `db:"cover_image_id" json:"cover_image_id"`
	Title        string         `db:"title" json:"title"`
	Slug         string         `db:"slug" json:"slug"`
	Summary      sql.NullString `db:"summary" json:"summary"`
	Content      string         `db:"content" json:"content"`
	PostType     PostType       `db:"post_type" json:"post_type"`
	AuthorID     int64          `db:"author_id" json:"author_id"`
	IsPublished  bool           `db:"is_published" json:"is_published"`
	CreatedAt    time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at" json:"updated_at"`
}

type TeamMemberModel struct {
	ID           int64          `db:"id" json:"id"`
	UserID       int64          `db:"user_id" json:"user_id"`
	Role         string         `db:"role" json:"role"`
	Bio          string         `db:"bio" json:"bio"`
	Link         sql.NullString `db:"link" json:"link"`
	DisplayOrder int            `db:"display_order" json:"display_order"`
	IsActive     bool           `db:"is_active" json:"is_active"`
	CreatedAt    time.Time      `db:"created_at" json:"created_at"`
}

//////////////
//// ROWS ////
//////////////

// Posts

type PostAdminViewRow struct {
	ID           int64          `db:"id"`
	CoverImageID int64          `db:"cover_image_id"`
	ImageFileKey string         `db:"image_file_key"`
	Title        string         `db:"title"`
	Slug         string         `db:"slug"`
	Summary      sql.NullString `db:"summary"`
	Content      string         `db:"content"`
	PostType     PostType       `db:"post_type"`
	AuthorID     int64          `db:"author_id"`
	AuthorName   string         `db:"author_name"`
	IsPublished  bool           `db:"is_published"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
}

type PostAdminListViewRow struct {
	ID           int64          `db:"id"`
	ImageFileKey string         `db:"image_file_key"`
	Title        string         `db:"title"`
	Summary      sql.NullString `db:"summary"`
	PostType     PostType       `db:"post_type"`
	AuthorName   string         `db:"author_name"`
	IsPublished  bool           `db:"is_published"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
}

type PostViewRow struct {
	ImageFileKey string         `db:"image_file_key"`
	Title        string         `db:"title"`
	Slug         string         `db:"slug"`
	Summary      sql.NullString `db:"summary"`
	Content      string         `db:"content"`
	AuthorName   string         `db:"author_name"`
	UpdatedAt    time.Time      `db:"updated_at"`
}

type PostListViewRaw struct {
	ImageFileKey string         `db:"image_file_key"`
	Title        string         `db:"title"`
	Slug         string         `db:"slug"`
	Summary      sql.NullString `db:"summary"`
	AuthorName   string         `db:"author_name"`
	UpdatedAt    time.Time      `db:"updated_at"`
}

// Team

type TeamMemberRow struct {
	ID            int64          `db:"id"`
	UserID        int64          `db:"user_id"`
	NameAr        string         `db:"name_ar"`
	NameEn        string         `db:"name_en"`
	Role          string         `db:"role"`
	Bio           string         `db:"bio"`
	Link          sql.NullString `db:"link"`
	DisplayOrder  int            `db:"display_order"`
	IsActive      bool           `db:"is_active"`
	CreatedAt     time.Time      `db:"created_at"`
	ProfilePicKey sql.NullString `db:"profile_pic_key"`
}

//////////////
///  DTOs  ///
//////////////

// For Admin

type PostRequest struct {
	CoverImageID int64    `json:"cover_image_id"`
	Title        string   `json:"title" binding:"required,min=3,max=255"`
	Slug         string   `json:"slug"`
	Summary      string   `json:"summary"`
	Type         PostType `json:"type"`
	Content      string   `json:"content" binding:"required"`
	IsPublished  bool     `json:"is_published"`
}

type PostUpdateRequest struct {
	ID int64 `db:"id" json:"id"`
	PostRequest
}

type PostAdminViewResponse struct {
	ID           int64     `json:"id"`
	CoverImageID int64     `json:"cover_image_id"`
	ImageUrl     string    `json:"image_url"`
	Title        string    `json:"title"`
	Slug         string    `json:"slug"`
	Summary      string    `json:"summary"`
	Content      string    `json:"content"`
	PostType     PostType  `json:"post_type"`
	AuthorID     int64     `json:"author_id"`
	AuthorName   string    `json:"author_name"`
	IsPublished  bool      `json:"is_published"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PostAdminListViewResponse struct {
	ID          int64     `json:"id"`
	ImageUrl    string    `json:"image_url"`
	Title       string    `json:"title"`
	Summary     string    `json:"summary"`
	PostType    PostType  `json:"post_type"`
	AuthorName  string    `json:"author_name"`
	IsPublished bool      `json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type BatchPostAdminListViewResponse struct {
	Posts   []PostAdminListViewResponse `json:"posts"`
	Current int64                       `json:"current"`
	Pages   int64                       `json:"pages"`
}

// For View

type PostViewResponse struct {
	ImageUrl   string    `json:"image_url"`
	Title      string    `json:"title"`
	Slug       string    `json:"slug"`
	Summary    string    `json:"summary"`
	Content    string    `json:"content"`
	AuthorName string    `json:"author_name"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type PostViewListResponse struct {
	ImageUrl   string    `json:"image_url"`
	Title      string    `json:"title"`
	Slug       string    `json:"slug"`
	Summary    string    `json:"summary"`
	AuthorName string    `json:"author_name"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type BatchPostListViewResponse struct {
	Posts   []PostViewListResponse `json:"posts"`
	Current int64                  `json:"current"`
	Pages   int64                  `json:"pages"`
}

//////////////
///  TEAM  ///
//////////////

type TeamMemberRequest struct {
	UserID       int64  `json:"user_id" binding:"required"`
	Role         string `json:"role" binding:"required"`
	Bio          string `json:"bio" binding:"required"`
	Link         string `json:"link"`
	DisplayOrder int    `json:"display_order" binding:"required"`
	IsActive     bool   `json:"is_active" binding:"required"`
}

type TeamMemberUpdateRequest struct {
	ID int64 `json:"id" binding:"required"`
	TeamMemberRequest
}

type TeamMemberResponse struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	NameAr       string    `json:"name_ar"`
	NameEn       string    `json:"name_en"`
	Role         string    `json:"role"`
	Bio          string    `json:"bio"`
	Link         string    `json:"link"`
	ProfilePic   string    `json:"profile_pic"`
	DisplayOrder int       `json:"display_order"`
	CreatedAt    time.Time `json:"created_at"`
}

type TeamMemberViewResponse struct {
	UserID       int64  `json:"user_id"`
	NameAr       string `json:"name_ar"`
	NameEn       string `json:"name_en"`
	Role         string `json:"role"`
	Bio          string `json:"bio"`
	Link         string `json:"link"`
	ProfilePic   string `json:"profile_pic"`
	DisplayOrder int    `json:"display_order"`
}
