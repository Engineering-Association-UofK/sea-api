package repositories

import (
	"sea-api/internal/models"

	"github.com/jmoiron/sqlx"
)

type CmsRepository struct {
	db *sqlx.DB
}

func NewCmsRepository(db *sqlx.DB) *CmsRepository {
	return &CmsRepository{db: db}
}

// ======== BLOG POSTS ========

func (r *CmsRepository) CreateBlogPost(post *models.BlogPostModel) (int64, error) {
	query := `
	INSERT INTO blog_posts (cover_image_id, title, slug, content, author_id, is_published, created_at, updated_at)
	VALUES (:cover_image_id, :title, :slug, :content, :author_id, :is_published, :created_at, :updated_at)
	`
	res, err := r.db.NamedExec(query, post)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *CmsRepository) GetBlogPostByID(id int64) (*models.BlogPostModel, error) {
	var post models.BlogPostModel
	err := r.db.Get(&post, `SELECT * FROM blog_posts WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *CmsRepository) GetBlogPostBySlug(slug string) (*models.BlogPostModel, error) {
	var post models.BlogPostModel
	err := r.db.Get(&post, `SELECT * FROM blog_posts WHERE slug = ?`, slug)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *CmsRepository) GetAllBlogPosts(publishedOnly bool) ([]models.BlogPostModel, error) {
	var posts []models.BlogPostModel
	query := `SELECT * FROM blog_posts`
	if publishedOnly {
		query += ` WHERE is_published = TRUE`
	}
	query += ` ORDER BY created_at DESC`
	err := r.db.Select(&posts, query)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *CmsRepository) UpdateBlogPost(post *models.BlogPostModel) error {
	query := `
	UPDATE blog_posts
	SET cover_image_id = :cover_image_id, title = :title, slug = :slug, content = :content, 
	    author_id = :author_id, is_published = :is_published, updated_at = :updated_at
	WHERE id = :id
	`
	_, err := r.db.NamedExec(query, post)
	return err
}

func (r *CmsRepository) DeleteBlogPost(id int64) error {
	_, err := r.db.Exec(`DELETE FROM blog_posts WHERE id = ?`, id)
	return err
}

// ======== TEAM MEMBERS ========

func (r *CmsRepository) CreateTeamMember(member *models.TeamMemberModel) (int64, error) {
	query := `
	INSERT INTO team_members (user_id, role, bio, display_order, is_active, created_at)
	VALUES (:user_id, :role, :bio, :display_order, :is_active, :created_at)
	`
	res, err := r.db.NamedExec(query, member)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *CmsRepository) GetTeamMemberByID(id int64) (*models.TeamMemberModel, error) {
	var member models.TeamMemberModel
	err := r.db.Get(&member, `SELECT * FROM team_members WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *CmsRepository) GetTeamMemberByUserID(userID int64) (*models.TeamMemberModel, error) {
	var member models.TeamMemberModel
	err := r.db.Get(&member, `SELECT * FROM team_members WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *CmsRepository) GetAllTeamMembers(activeOnly bool) ([]models.TeamMemberModel, error) {
	var members []models.TeamMemberModel
	query := `SELECT * FROM team_members`
	if activeOnly {
		query += ` WHERE is_active = TRUE`
	}
	query += ` ORDER BY display_order ASC`
	err := r.db.Select(&members, query)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (r *CmsRepository) UpdateTeamMember(member *models.TeamMemberModel) error {
	query := `
	UPDATE team_members
	SET user_id = :user_id, role = :role, bio = :bio, display_order = :display_order, is_active = :is_active
	WHERE id = :id
	`
	_, err := r.db.NamedExec(query, member)
	return err
}

func (r *CmsRepository) DeleteTeamMember(id int64) error {
	_, err := r.db.Exec(`DELETE FROM team_members WHERE id = ?`, id)
	return err
}
