package models

import "time"

type PostType string

const (
	PostNews         PostType = "NEWS"
	PostIssues       PostType = "ISSUES"
	PostAnnouncement PostType = "ANNOUNCEMENT"
	PostEvent        PostType = "EVENT"
	PostBlog         PostType = "BLOG"
	PostPhilosophy   PostType = "PHILANTHROPY"
)

// Models

type PostModel struct {
	ID           int64     `db:"id" json:"id"`
	CoverImageID int64     `db:"cover_image_id" json:"cover_image_id"`
	Title        string    `db:"title" json:"title"`
	Slug         string    `db:"slug" json:"slug"`
	Content      string    `db:"content" json:"content"`
	PostType     PostType  `db:"post_type" json:"post_type"`
	AuthorID     int64     `db:"author_id" json:"author_id"`
	IsPublished  bool      `db:"is_published" json:"is_published"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type TeamMemberModel struct {
	ID           int64     `db:"id" json:"id"`
	UserID       int64     `db:"user_id" json:"user_id"`
	Role         string    `db:"role" json:"role"`
	Bio          string    `db:"bio" json:"bio"`
	DisplayOrder int       `db:"display_order" json:"display_order"`
	IsActive     bool      `db:"is_active" json:"is_active"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

// DTOs

type PostRequest struct {
	CoverImageID int64  `json:"cover_image_id"`
	Title        string `json:"title" binding:"required,min=3,max=255"`
	Slug         string `json:"slug"`
	Content      string `json:"content" binding:"required"`
	IsPublished  bool   `json:"is_published"`
}

type PostUpdateRequest struct {
	ID int64 `db:"id" json:"id"`
	PostRequest
}

type PostResponse struct {
	ID           int64     `json:"id"`
	CoverImageID int64     `json:"cover_image_id"`
	ImageUrl     string    `json:"image_url"`
	Title        string    `json:"title"`
	Slug         string    `json:"slug"`
	Content      string    `json:"content"`
	AuthorID     int64     `json:"author_id"`
	AuthorName   string    `json:"author_name"`
	IsPublished  bool      `json:"is_published"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PostListResponse struct {
	Posts   []PostResponse `json:"posts"`
	Current int64          `json:"current"`
	Pages   int64          `json:"pages"`
}

type PostViewResponse struct {
	ImageUrl   string    `json:"image_url"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	AuthorName string    `json:"author_name"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type PostViewListResponse struct {
	ImageUrl   string    `json:"image_url"`
	Title      string    `json:"title"`
	AuthorName string    `json:"author_name"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type PostListViewResponse struct {
	Posts   []PostViewListResponse `json:"posts"`
	Current int64                  `json:"current"`
	Pages   int64                  `json:"pages"`
}

type TeamMemberRequest struct {
	UserID       int64  `json:"user_id" binding:"required"`
	Role         string `json:"role" binding:"required"`
	Bio          string `json:"bio" binding:"required"`
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
	ProfilePic   string `json:"profile_pic"`
	DisplayOrder int    `json:"display_order"`
}
