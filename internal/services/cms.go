package services

import (
	"database/sql"
	"log/slog"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/repositories"
	"sea-api/internal/services/user"
	"sea-api/internal/utils/valid"
	"strings"
	"time"
)

type CmsService struct {
	CmsRepo        *repositories.CmsRepository
	UserService    *user.UserService
	GalleryService *GalleryService
}

func NewCmsService(CmsRepo *repositories.CmsRepository, userService *user.UserService, galleryService *GalleryService) *CmsService {
	return &CmsService{
		CmsRepo:        CmsRepo,
		UserService:    userService,
		GalleryService: galleryService,
	}
}

// ======== BLOG POSTS ========

func (s *CmsService) CreatePost(userId int64, post *models.PostRequest) (int64, error) {
	if post.Slug == "" {
		post.Slug = strings.ToLower(strings.ReplaceAll(post.Title, " ", "-"))
	}
	if _, err := s.CmsRepo.GetPostDetailsBySlug(post.Slug); err == nil {
		return 0, errs.New(errs.Conflict, "Slug already exists", nil)
	}
	if _, err := s.GalleryService.GetAssetByID(post.CoverImageID); err != nil {
		return 0, errs.New(errs.BadRequest, "invalid image ID provided", nil)
	}

	model := &models.PostModel{
		CoverImageID: post.CoverImageID,
		Title:        post.Title,
		Slug:         post.Slug,
		Summary:      sql.NullString{String: post.Summary, Valid: post.Summary != ""},
		Content:      post.Content,
		AuthorID:     userId,
		IsPublished:  post.IsPublished,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	id, err := s.CmsRepo.CreatePost(model)
	if err != nil {
		return 0, err
	}
	err = s.GalleryService.AttachAssetToObject(post.CoverImageID, models.ObjPost, id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *CmsService) GetPostById(id int64) (*models.PostAdminViewResponse, error) {
	post, err := s.CmsRepo.GetPostByID(id)
	if err != nil {
		slog.Info("Post not found using ID")
		return nil, err
	}

	url, err := s.GalleryService.GetLinkByAssetKey(post.ImageFileKey)
	if err != nil {
		return nil, err
	}

	return &models.PostAdminViewResponse{
		ID:           post.ID,
		CoverImageID: post.CoverImageID,
		ImageUrl:     url,
		Title:        post.Title,
		Slug:         post.Slug,
		Summary:      post.Summary.String,
		Content:      post.Content,
		PostType:     post.PostType,
		AuthorID:     post.AuthorID,
		AuthorName:   post.AuthorName,
		IsPublished:  post.IsPublished,
		CreatedAt:    post.CreatedAt,
		UpdatedAt:    post.UpdatedAt,
	}, nil
}

func (s *CmsService) GetViewPostBySlug(slug string) (*models.PostViewResponse, error) {
	post, err := s.CmsRepo.GetPostDetailsBySlug(slug)
	if err != nil {
		return nil, err
	}

	url, err := s.GalleryService.GetLinkByAssetKey(post.ImageFileKey)
	if err != nil {
		return nil, err
	}

	return &models.PostViewResponse{
		ImageUrl:   url,
		Title:      post.Title,
		Slug:       post.Slug,
		Summary:    post.Summary.String,
		Content:    post.Content,
		AuthorName: post.AuthorName,
		UpdatedAt:  post.UpdatedAt,
	}, nil
}

func (s *CmsService) GetViewPostList(req *models.ListRequest) (*models.BatchPostListViewResponse, error) {
	if !models.AllowedPostTypes[models.PostType(req.Type)] {
		return nil, errs.New(errs.BadRequest, "invalid post type", nil)
	}
	total, err := s.CmsRepo.GetTotalPosts(models.PostType(req.Type), true)
	if err != nil {
		return nil, err
	}
	slog.Info("Post Details",
		"Limit", req.Limit,
		"Page", req.Page,
		"Total", total,
		"Type", req.Type,
	)

	valid.Limit(req, total)

	postsRows, err := s.CmsRepo.GetPostsViewListByType(req, models.PostType(req.Type))
	if err != nil {
		return nil, err
	}

	responses := []models.PostViewListResponse{}
	for _, post := range postsRows {
		url, err := s.GalleryService.GetLinkByAssetKey(post.ImageFileKey)
		if err != nil {
			slog.Info("Failed to generate url", "store key", post.ImageFileKey)
			return nil, err
		}
		summary := ""
		if post.Summary.Valid {
			summary = post.Summary.String
		}
		responses = append(responses, models.PostViewListResponse{
			ImageUrl:   url,
			Title:      post.Title,
			Slug:       post.Slug,
			Summary:    summary,
			AuthorName: post.AuthorName,
			UpdatedAt:  post.UpdatedAt,
		})
	}

	response := models.BatchPostListViewResponse{
		Posts:   responses,
		Current: req.Page,
		Pages:   total / req.Limit,
	}

	return &response, nil
}

func (s *CmsService) GetAllPosts(req *models.ListRequest) (*models.BatchPostAdminListViewResponse, error) {
	if !models.AllowedPostTypes[models.PostType(req.Type)] {
		return nil, errs.New(errs.BadRequest, "invalid post type", nil)
	}
	total, err := s.CmsRepo.GetTotalPosts("", false)
	if err != nil {
		return nil, err
	}
	valid.Limit(req, total)

	posts, err := s.CmsRepo.GetPostsAdminListByType(req, "")
	if err != nil {
		return nil, err
	}
	if len(posts) == 0 {
		return &models.BatchPostAdminListViewResponse{
			Posts:   []models.PostAdminListViewResponse{},
			Current: req.Page,
			Pages:   total / req.Limit,
		}, nil
	}

	var responses = []models.PostAdminListViewResponse{}

	for _, post := range posts {
		url, _ := s.GalleryService.GetLinkByAssetID(post.CoverImageID)

		responses = append(responses, models.PostAdminListViewResponse{
			ID:          post.ID,
			ImageUrl:    url,
			Title:       post.Title,
			Summary:     post.Summary.String,
			PostType:    post.PostType,
			AuthorName:  post.AuthorName,
			IsPublished: post.IsPublished,
			CreatedAt:   post.CreatedAt,
			UpdatedAt:   post.UpdatedAt,
		})
	}

	return &models.BatchPostAdminListViewResponse{
		Posts:   responses,
		Current: req.Page,
		Pages:   total / req.Limit,
	}, nil
}

func (s *CmsService) UpdatePost(req *models.PostUpdateRequest) error {
	post, err := s.CmsRepo.GetPostModelByID(req.ID)
	if err != nil {
		return err
	}
	if req.CoverImageID != post.CoverImageID {
		s.GalleryService.RemoveReference(models.ObjPost, post.ID)
		if _, err := s.GalleryService.GetAssetByID(req.CoverImageID); err != nil {
			return errs.New(errs.BadRequest, "invalid image ID provided", nil)
		}
		s.GalleryService.AttachAssetToObject(req.CoverImageID, models.ObjPost, req.ID)
	}

	post.CoverImageID = req.CoverImageID
	post.Title = req.Title
	if req.Slug != "" {
		_, err := s.CmsRepo.GetPostDetailsBySlug(req.Slug)
		if err == nil {
			return errs.New(errs.Conflict, "Slug already exists", nil)
		}
		post.Slug = req.Slug
	}
	post.Content = req.Content
	post.IsPublished = req.IsPublished
	post.UpdatedAt = time.Now()

	return s.CmsRepo.UpdatePost(post)
}

func (s *CmsService) DeletePost(id int64) error {
	if _, err := s.CmsRepo.GetPostByID(id); err != nil {
		return err
	}
	s.GalleryService.RemoveReference(models.ObjPost, id)
	return s.CmsRepo.DeletePost(id)
}

// ======== TEAM MEMBERS ========

func (s *CmsService) CreateTeamMember(member *models.TeamMemberRequest) (int64, error) {
	if _, err := s.UserService.GetUserDetails(member.UserID); err != nil {
		return 0, err
	}
	if _, err := s.CmsRepo.GetTeamMemberByUserID(member.UserID); err == nil {
		return 0, errs.New(errs.Conflict, "User already has a team member profile", nil)
	}
	return s.CmsRepo.CreateTeamMember(&models.TeamMemberModel{
		UserID:       member.UserID,
		Role:         member.Role,
		Bio:          member.Bio,
		Link:         sql.NullString{String: member.Link, Valid: member.Link != ""},
		DisplayOrder: member.DisplayOrder,
		IsActive:     member.IsActive,
		CreatedAt:    time.Now(),
	})
}

func (s *CmsService) GetTeamMemberByID(id int64) (*models.TeamMemberResponse, error) {
	member, err := s.CmsRepo.GetTeamMemberByID(id)
	if err != nil {
		return nil, err
	}

	user, err := s.UserService.GetUserDetails(member.UserID)
	if err != nil {
		return nil, err
	}

	link := ""
	if member.Link.Valid {
		link = member.Link.String
	}

	return &models.TeamMemberResponse{
		ID:           member.ID,
		UserID:       member.UserID,
		NameAr:       user.NameAr,
		NameEn:       user.NameEn,
		Role:         member.Role,
		Bio:          member.Bio,
		Link:         link,
		ProfilePic:   user.ProfilePic,
		DisplayOrder: member.DisplayOrder,
		CreatedAt:    member.CreatedAt,
	}, nil
}

func (s *CmsService) GetAllTeamMembers(activeOnly bool) ([]models.TeamMemberResponse, error) {
	members, err := s.CmsRepo.GetAllTeamMembers(activeOnly)
	if err != nil {
		return nil, err
	}
	if len(members) == 0 {
		return []models.TeamMemberResponse{}, nil
	}

	var dtos []models.TeamMemberResponse
	for _, m := range members {
		url := ""
		if m.ProfilePicKey.Valid {
			url, _ = s.GalleryService.GetLinkByAssetKey(m.ProfilePicKey.String)
		}

		link := ""
		if m.Link.Valid {
			link = m.Link.String
		}

		dtos = append(dtos, models.TeamMemberResponse{
			ID:           m.ID,
			UserID:       m.UserID,
			NameAr:       m.NameAr,
			NameEn:       m.NameEn,
			Role:         m.Role,
			Bio:          m.Bio,
			Link:         link,
			ProfilePic:   url,
			DisplayOrder: m.DisplayOrder,
			CreatedAt:    m.CreatedAt,
		})
	}

	return dtos, nil
}

func (s *CmsService) GetViewTeamMembers() ([]models.TeamMemberViewResponse, error) {
	members, err := s.CmsRepo.GetAllTeamMembers(true)
	if err != nil {
		return nil, err
	}

	if len(members) == 0 {
		return []models.TeamMemberViewResponse{}, nil
	}

	var dtos []models.TeamMemberViewResponse
	for _, m := range members {
		url := ""
		if m.ProfilePicKey.Valid {
			url, _ = s.GalleryService.GetLinkByAssetKey(m.ProfilePicKey.String)
		}

		link := ""
		if m.Link.Valid {
			link = m.Link.String
		}

		dtos = append(dtos, models.TeamMemberViewResponse{
			UserID:       m.UserID,
			NameAr:       m.NameAr,
			NameEn:       m.NameEn,
			Role:         m.Role,
			Bio:          m.Bio,
			Link:         link,
			ProfilePic:   url,
			DisplayOrder: m.DisplayOrder,
		})
	}

	return dtos, nil
}

func (s *CmsService) UpdateTeamMember(req *models.TeamMemberUpdateRequest) error {
	member, err := s.CmsRepo.GetTeamMemberByID(req.ID)
	if err != nil {
		return err
	}

	member.Role = req.Role
	member.Bio = req.Bio
	member.Link = sql.NullString{String: req.Link, Valid: req.Link != ""}
	member.DisplayOrder = req.DisplayOrder
	member.IsActive = req.IsActive

	return s.CmsRepo.UpdateTeamMember(member)
}

func (s *CmsService) DeleteTeamMember(id int64) error {
	return s.CmsRepo.DeleteTeamMember(id)
}
