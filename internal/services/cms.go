package services

import (
	"log/slog"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/repositories"
	"sea-api/internal/utils"
	"strings"
	"time"
)

type CmsService struct {
	CmsRepo        *repositories.CmsRepository
	UserService    *UserService
	GalleryService *GalleryService
}

func NewCmsService(CmsRepo *repositories.CmsRepository, userService *UserService, galleryService *GalleryService) *CmsService {
	return &CmsService{
		CmsRepo:        CmsRepo,
		UserService:    userService,
		GalleryService: galleryService,
	}
}

// ======== BLOG POSTS ========

func (s *CmsService) CreateBlogPost(userId int64, post *models.BlogPostRequest) (int64, error) {
	if post.Slug == "" {
		post.Slug = strings.ToLower(strings.ReplaceAll(post.Title, " ", "-"))
	}
	if _, err := s.CmsRepo.GetBlogPostBySlug(post.Slug); err == nil {
		return 0, errs.New(errs.Conflict, "Slug already exists", nil)
	}
	if _, err := s.GalleryService.GetAssetByID(post.CoverImageID); err != nil {
		return 0, errs.New(errs.BadRequest, "invalid image ID provided", nil)
	}

	model := &models.BlogPostModel{
		CoverImageID: post.CoverImageID,
		Title:        post.Title,
		Slug:         post.Slug,
		Content:      post.Content,
		AuthorID:     userId,
		IsPublished:  post.IsPublished,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	id, err := s.CmsRepo.CreateBlogPost(model)
	if err != nil {
		return 0, err
	}
	err = s.GalleryService.AttachAssetToObject(post.CoverImageID, models.BlogPost, id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *CmsService) GetBlogPostById(id int64) (*models.BlogPostResponse, error) {
	post, err := s.CmsRepo.GetBlogPostByID(id)
	if err != nil {
		slog.Info("Post not found using ID")
		return nil, err
	}
	return s.getBlogPost(post)
}

func (s *CmsService) GetBlogPostBySlug(slug string) (*models.BlogPostResponse, error) {
	post, err := s.CmsRepo.GetBlogPostBySlug(slug)
	if err != nil {
		slog.Info("Post not found using slug")
		return nil, err
	}
	return s.getBlogPost(post)
}

func (s *CmsService) getBlogPost(post *models.BlogPostModel) (*models.BlogPostResponse, error) {
	url, err := s.GalleryService.GetLinkByAssetID(post.CoverImageID)
	if err != nil {
		slog.Info("Failed to generate url", "store id", post.CoverImageID)
		return nil, err
	}

	user, err := s.UserService.GetUserDetails(post.AuthorID)
	if err != nil {
		slog.Info("User not found for getting post author", "user id", post.AuthorID)
		return nil, err
	}

	return &models.BlogPostResponse{
		ID:           post.ID,
		CoverImageID: post.CoverImageID,
		ImageUrl:     url,
		Title:        post.Title,
		Slug:         post.Slug,
		Content:      post.Content,
		AuthorID:     post.AuthorID,
		AuthorName:   user.NameAr,
		IsPublished:  post.IsPublished,
		CreatedAt:    post.CreatedAt,
		UpdatedAt:    post.UpdatedAt,
	}, nil
}

func (s *CmsService) GetAllBlogPosts(publishedOnly bool) ([]models.BlogPostResponse, error) {
	posts, err := s.CmsRepo.GetAllBlogPosts(publishedOnly)
	if err != nil {
		return []models.BlogPostResponse{}, err
	}
	if len(posts) == 0 {
		return []models.BlogPostResponse{}, nil
	}
	usersIds := utils.FromSlice(posts, func(p models.BlogPostModel) int64 { return p.AuthorID }).Keys()

	users, err := s.UserService.GetAllUserDetailsByIndices(usersIds)
	if err != nil {
		return []models.BlogPostResponse{}, err
	}
	usersMap := utils.FromSlice(users, func(u models.UserDetails) int64 { return u.ID })

	var responses []models.BlogPostResponse
	for _, post := range posts {
		url, _ := s.GalleryService.GetLinkByAssetID(post.CoverImageID)

		responses = append(responses, models.BlogPostResponse{
			ID:           post.ID,
			CoverImageID: post.CoverImageID,
			ImageUrl:     url,
			Title:        post.Title,
			Slug:         post.Slug,
			Content:      post.Content,
			AuthorID:     post.AuthorID,
			AuthorName:   usersMap[post.AuthorID].NameAr,
			IsPublished:  post.IsPublished,
			CreatedAt:    post.CreatedAt,
			UpdatedAt:    post.UpdatedAt,
		})
	}

	return responses, nil

}

func (s *CmsService) UpdateBlogPost(req *models.BlogPostUpdateRequest) error {
	post, err := s.CmsRepo.GetBlogPostByID(req.ID)
	if err != nil {
		return err
	}
	if req.CoverImageID != post.CoverImageID && req.CoverImageID != 0 {
		s.GalleryService.RemoveReference(models.BlogPost, post.ID)
		s.GalleryService.AttachAssetToObject(req.CoverImageID, models.BlogPost, req.ID)
	}

	post.CoverImageID = req.CoverImageID
	post.Title = req.Title
	if req.Slug != "" {
		_, err := s.CmsRepo.GetBlogPostBySlug(req.Slug)
		if err == nil {
			return errs.New(errs.Conflict, "Slug already exists", nil)
		}
		post.Slug = req.Slug
	}
	post.Content = req.Content
	post.IsPublished = req.IsPublished
	post.UpdatedAt = time.Now()

	return s.CmsRepo.UpdateBlogPost(post)
}

func (s *CmsService) DeleteBlogPost(id int64) error {
	if _, err := s.CmsRepo.GetBlogPostByID(id); err != nil {
		return err
	}
	s.GalleryService.RemoveReference(models.BlogPost, id)
	return s.CmsRepo.DeleteBlogPost(id)
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

	return &models.TeamMemberResponse{
		ID:           member.ID,
		UserID:       member.UserID,
		NameAr:       user.NameAr,
		NameEn:       user.NameEn,
		Role:         member.Role,
		Bio:          member.Bio,
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

	userIds := make([]int64, len(members))
	for i, m := range members {
		userIds[i] = m.UserID
	}

	users, err := s.UserService.GetAllUserDetailsByIndices(userIds)
	if err != nil {
		return nil, err
	}
	usersMap := utils.FromSlice(users, func(u models.UserDetails) int64 { return u.ID })

	var dtos []models.TeamMemberResponse
	for _, m := range members {
		user := usersMap[m.UserID]
		dtos = append(dtos, models.TeamMemberResponse{
			ID:           m.ID,
			UserID:       m.UserID,
			NameAr:       user.NameAr,
			NameEn:       user.NameEn,
			Role:         m.Role,
			Bio:          m.Bio,
			DisplayOrder: m.DisplayOrder,
			CreatedAt:    m.CreatedAt,
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
	member.DisplayOrder = req.DisplayOrder
	member.IsActive = req.IsActive

	return s.CmsRepo.UpdateTeamMember(member)
}

func (s *CmsService) DeleteTeamMember(id int64) error {
	return s.CmsRepo.DeleteTeamMember(id)
}
