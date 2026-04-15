package repositories

import (
	"fmt"
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

func (r *CmsRepository) CreatePost(post *models.PostModel) (int64, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (cover_image_id, title, slug, content, author_id, is_published, created_at, updated_at)
	VALUES (:cover_image_id, :title, :slug, :content, :author_id, :is_published, :created_at, :updated_at)
	`, models.TablePosts)
	res, err := r.db.NamedExec(query, post)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *CmsRepository) GetPostByID(id int64) (*models.PostModel, error) {
	var post models.PostModel
	err := r.db.Get(&post, fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, models.TablePosts), id)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *CmsRepository) GetPostBySlug(slug string) (*models.PostModel, error) {
	var post models.PostModel
	err := r.db.Get(&post, fmt.Sprintf(`SELECT * FROM %s WHERE slug = ?`, models.TablePosts), slug)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *CmsRepository) GetAllPosts(req *models.ListRequest, publishedOnly bool) ([]models.PostModel, error) {
	query := fmt.Sprintf(`SELECT * FROM %s`, models.TablePosts)
	var posts []models.PostModel
	if publishedOnly {
		query += ` WHERE is_published = TRUE`
	}
	query += `
	ORDER BY created_at DESC
	LIMIT ? OFFSET ?
	`

	offset := (req.Page - 1) * req.Limit
	err := r.db.Select(&posts, query, req.Limit, offset)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *CmsRepository) GetTotalPosts() (int64, error) {
	var total int64
	err := r.db.Get(&total, fmt.Sprintf(`SELECT COUNT(*) FROM %s`, models.TablePosts))
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *CmsRepository) GetPostsListByType(req *models.ListRequest, postType models.PostType) ([]models.PostModel, error) {
	query := fmt.Sprintf(`
	SELECT * FROM %s 
	WHERE post_type = ?
	AND is_published = TRUE
	ORDER BY created_at DESC
	LIMIT ? OFFSET ?
	`, models.TablePosts)

	offset := (req.Page - 1) * req.Limit
	var posts []models.PostModel
	err := r.db.Select(&posts, query, postType, req.Limit, offset)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *CmsRepository) GetPublishedTotalByType(postType models.PostType) (int64, error) {
	var total int64
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE post_type = ? AND is_published = TRUE`, models.TablePosts)
	err := r.db.Get(&total, countQuery, postType)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *CmsRepository) UpdatePost(post *models.PostModel) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET cover_image_id = :cover_image_id, title = :title, slug = :slug, content = :content, 
	    author_id = :author_id, is_published = :is_published, updated_at = :updated_at
	WHERE id = :id
	`, models.TablePosts)
	_, err := r.db.NamedExec(query, post)
	return err
}

func (r *CmsRepository) DeletePost(id int64) error {
	_, err := r.db.Exec(fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, models.TablePosts), id)
	return err
}

// ======== TEAM MEMBERS ========

func (r *CmsRepository) CreateTeamMember(member *models.TeamMemberModel) (int64, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (user_id, role, bio, display_order, is_active, created_at)
	VALUES (:user_id, :role, :bio, :display_order, :is_active, :created_at)
	`, models.TableTeamMembers)
	res, err := r.db.NamedExec(query, member)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *CmsRepository) GetTeamMemberByID(id int64) (*models.TeamMemberModel, error) {
	var member models.TeamMemberModel
	err := r.db.Get(&member, fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, models.TableTeamMembers), id)
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *CmsRepository) GetTeamMemberByUserID(userID int64) (*models.TeamMemberModel, error) {
	var member models.TeamMemberModel
	err := r.db.Get(&member, fmt.Sprintf(`SELECT * FROM %s WHERE user_id = ?`, models.TableTeamMembers), userID)
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *CmsRepository) GetAllTeamMembers(activeOnly bool) ([]models.TeamMemberModel, error) {
	var members []models.TeamMemberModel
	query := fmt.Sprintf(`SELECT * FROM %s`, models.TableTeamMembers)
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
	query := fmt.Sprintf(`
	UPDATE %s
	SET user_id = :user_id, role = :role, bio = :bio, display_order = :display_order, is_active = :is_active
	WHERE id = :id
	`, models.TableTeamMembers)
	_, err := r.db.NamedExec(query, member)
	return err
}

func (r *CmsRepository) DeleteTeamMember(id int64) error {
	_, err := r.db.Exec(fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, models.TableTeamMembers), id)
	return err
}
