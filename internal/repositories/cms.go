package repositories

import (
	"fmt"
	"log/slog"
	"sea-api/internal/models"
	"strings"

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
	INSERT INTO %s (cover_image_id, title, slug, summary, content, author_id, is_published, created_at, updated_at)
	VALUES (:cover_image_id, :title, :slug, :summary, :content, :author_id, :is_published, :created_at, :updated_at)
	`, models.TablePosts)
	res, err := r.db.NamedExec(query, post)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

/////////////////
//// GET ONE ////
/////////////////

func (r *CmsRepository) GetPostModelByID(id int64) (*models.PostModel, error) {
	var post models.PostModel
	err := r.db.Get(&post, fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, models.TablePosts), id)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *CmsRepository) GetPostByID(id int64) (*models.PostAdminViewRow, error) {
	var post models.PostAdminViewRow
	query := fmt.Sprintf(`
	SELECT 
		p.id,
		p.cover_image_id,
		f.file_key AS image_file_key,
		p.title,
		p.slug,
		p.summary,
		p.content,
		p.post_type,
		p.author_id,
		u.name_en AS author_name,
		p.is_published,
		p.created_at,
		p.updated_at
	FROM %s p
	LEFT JOIN %s g ON p.cover_image_id = g.id
	LEFT JOIN %s f ON g.file_id = f.id
	LEFT JOIN %s u ON p.author_id = u.id
	WHERE p.id = ?
	`, models.TablePosts, models.TableGalleryAssets, models.TableFiles, models.TableUsers)

	err := r.db.Get(&post, query, id)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *CmsRepository) GetPostDetailsBySlug(slug string) (*models.PostViewRow, error) {
	var post models.PostViewRow
	query := fmt.Sprintf(`
	SELECT 
		f.file_key AS image_file_key,
		p.title,
		p.slug,
		p.summary,
		p.content,
		u.name_en AS author_name,
		p.updated_at
	FROM %s p
	LEFT JOIN %s g ON p.cover_image_id = g.id
	LEFT JOIN %s f ON g.file_id = f.id
	LEFT JOIN %s u ON p.author_id = u.id
	WHERE p.slug = ? AND p.is_published = TRUE
	`, models.TablePosts, models.TableGalleryAssets, models.TableFiles, models.TableUsers)

	err := r.db.Get(&post, query, slug)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

////////////////
/// GET MANY ///
////////////////

func (r *CmsRepository) GetAllPostModels(req *models.ListRequest, publishedOnly bool) ([]models.PostModel, error) {
	var posts []models.PostModel
	query := fmt.Sprintf(`SELECT * FROM %s`, models.TablePosts)
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

func (r *CmsRepository) GetPostsAdminListByType(req *models.ListRequest, postType models.PostType) ([]models.PostAdminViewRow, error) {
	offset := (req.Page - 1) * req.Limit
	query := fmt.Sprintf(`
	SELECT 
		p.id,
		p.cover_image_id,
		f.file_key AS image_file_key,
		p.title,
		p.slug,
		p.summary,
		p.content,
		p.post_type,
		p.author_id,
		u.name_en AS author_name,
		p.is_published,
		p.created_at,
		p.updated_at
	FROM %s p
	LEFT JOIN %s g ON p.cover_image_id = g.id
	LEFT JOIN %s f ON g.file_id = f.id
	LEFT JOIN %s u ON p.author_id = u.id
	`, models.TablePosts, models.TableGalleryAssets, models.TableFiles, models.TableUsers)

	var args []interface{}
	if postType != "" {
		query += ` WHERE p.post_type = ?`
		args = append(args, postType)
	}

	query += ` ORDER BY p.created_at DESC LIMIT ? OFFSET ?`
	args = append(args, req.Limit, offset)

	var posts []models.PostAdminViewRow
	err := r.db.Select(&posts, query, args...)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *CmsRepository) GetPostsViewListByType(req *models.ListRequest, postType models.PostType) ([]models.PostListViewRaw, error) {
	var args []interface{}
	offset := (req.Page - 1) * req.Limit
	query := fmt.Sprintf(`
	SELECT 
		f.file_key AS image_file_key,
		p.title,
		p.slug,
		p.summary,
		u.name_en AS author_name,
		p.updated_at
	FROM %s p
	LEFT JOIN %s g ON p.cover_image_id = g.id
	LEFT JOIN %s f ON g.file_id = f.id
	LEFT JOIN %s u ON p.author_id = u.id
	WHERE p.is_published = TRUE
	`, models.TablePosts, models.TableGalleryAssets, models.TableFiles, models.TableUsers)
	if postType != "" {
		query += ` AND p.post_type = ?`
		args = append(args, postType)
	}
	query += ` ORDER BY p.created_at DESC LIMIT ? OFFSET ?`
	args = append(args, req.Limit, offset)

	var posts []models.PostListViewRaw
	err := r.db.Select(&posts, query, args...)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *CmsRepository) GetTotalPosts(postType models.PostType, published bool) (int64, error) {
	var conditions []string
	var args []interface{}

	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s p`, models.TablePosts)

	if postType != "" {
		conditions = append(conditions, "p.post_type = ?")
		args = append(args, postType)
	}

	if published {
		conditions = append(conditions, "p.is_published = TRUE")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var count int64
	err := r.db.Get(&count, query, args...)
	if err != nil {
		slog.Error("Database error in GetTotalPosts", "query", query, "error", err)
		return 0, err
	}

	return count, nil
}

///////////////////
/// Transaction ///
///////////////////

func (r *CmsRepository) UpdatePost(post *models.PostModel) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET cover_image_id = :cover_image_id, title = :title, slug = :slug, summary = :summary, content = :content, 
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
	INSERT INTO %s (user_id, role, bio, link, display_order, is_active, created_at)
	VALUES (:user_id, :role, :bio, :link, :display_order, :is_active, :created_at)
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

func (r *CmsRepository) GetAllTeamMembers(activeOnly bool) ([]models.TeamMemberRow, error) {
	var members []models.TeamMemberRow
	query := fmt.Sprintf(`
	SELECT 
		tm.id,
		tm.user_id,
		u.name_ar,
		u.name_en,
		tm.role,
		tm.bio,
		tm.link,
		tm.display_order,
		tm.is_active,
		tm.created_at,
		f.file_key AS profile_pic_key
	FROM %s tm
	JOIN %s u ON tm.user_id = u.id
	LEFT JOIN %s f ON u.profile_image_id = f.id
	`, models.TableTeamMembers, models.TableUsers, models.TableFiles)

	if activeOnly {
		query += ` WHERE tm.is_active = TRUE`
	}
	query += ` ORDER BY tm.display_order ASC`

	err := r.db.Select(&members, query)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (r *CmsRepository) UpdateTeamMember(member *models.TeamMemberModel) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET user_id = :user_id, role = :role, bio = :bio, link = :link, display_order = :display_order, is_active = :is_active
	WHERE id = :id
	`, models.TableTeamMembers)
	_, err := r.db.NamedExec(query, member)
	return err
}

func (r *CmsRepository) DeleteTeamMember(id int64) error {
	_, err := r.db.Exec(fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, models.TableTeamMembers), id)
	return err
}
