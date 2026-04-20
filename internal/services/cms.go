package services

import (
	"log/slog"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/repositories"
	"sea-api/internal/services/user"
	"sea-api/internal/utils"
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
	if _, err := s.CmsRepo.GetPostBySlug(post.Slug); err == nil {
		return 0, errs.New(errs.Conflict, "Slug already exists", nil)
	}
	if _, err := s.GalleryService.GetAssetByID(post.CoverImageID); err != nil {
		return 0, errs.New(errs.BadRequest, "invalid image ID provided", nil)
	}

	model := &models.PostModel{
		CoverImageID: post.CoverImageID,
		Title:        post.Title,
		Slug:         post.Slug,
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

func (s *CmsService) GetPostById(id int64) (*models.PostResponse, error) {
	post, err := s.CmsRepo.GetPostByID(id)
	if err != nil {
		slog.Info("Post not found using ID")
		return nil, err
	}
	return s.getPost(post)
}

func (s *CmsService) GetViewPostBySlug(slug string) (*models.PostViewResponse, error) {
	post, err := s.CmsRepo.GetPostBySlug(slug)
	if err != nil {
		return nil, err
	}
	if !post.IsPublished {
		return nil, errs.New(errs.NotFound, "post not found", nil)
	}
	model, err := s.getPost(post)
	if err != nil {
		return nil, err
	}

	return &models.PostViewResponse{
		ImageUrl:   model.ImageUrl,
		Title:      model.Title,
		Content:    model.Content,
		UpdatedAt:  model.UpdatedAt,
		AuthorName: model.AuthorName,
	}, nil
}

func (s *CmsService) GetViewPostList(req *models.ListRequest) (*models.PostListViewResponse, error) {
	total, err := s.CmsRepo.GetPublishedTotalByType(models.PostBlog)
	if err != nil {
		return nil, err
	}
	valid.ValidateListRequest(req, total)

	posts, err := s.CmsRepo.GetPostsListByType(req, models.PostBlog)
	if err != nil {
		return nil, err
	}
	response := models.PostListViewResponse{}
	if len(posts) == 0 {
		return &response, nil
	}

	var usersIds []int64
	for _, post := range posts {
		usersIds = append(usersIds, post.AuthorID)
	}

	authors, err := s.UserService.GetAllUserDetailsByIndices(usersIds)
	if err != nil {
		return nil, err
	}
	authorsMap := utils.FromSlice(authors, func(u models.UserDetails) int64 { return u.ID })

	responses := []models.PostViewListResponse{}
	for _, post := range posts {
		url, err := s.GalleryService.GetLinkByAssetID(post.CoverImageID)
		if err != nil {
			slog.Info("Failed to generate url", "store id", post.CoverImageID)
			return nil, err
		}
		responses = append(responses, models.PostViewListResponse{
			ImageUrl:   url,
			Title:      post.Title,
			AuthorName: authorsMap[post.AuthorID].NameAr,
			UpdatedAt:  post.UpdatedAt,
			Slug:       post.Slug,
		})
	}

	response = models.PostListViewResponse{
		Posts:   responses,
		Current: req.Page,
		Pages:   total / req.Limit,
	}

	return &response, nil
}

func (s *CmsService) getPost(post *models.PostModel) (*models.PostResponse, error) {
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

	return &models.PostResponse{
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

func (s *CmsService) GetAllPosts(req *models.ListRequest) (*models.PostListResponse, error) {
	total, err := s.CmsRepo.GetTotalPosts()
	if err != nil {
		return &models.PostListResponse{}, err
	}
	valid.ValidateListRequest(req, total)

	posts, err := s.CmsRepo.GetAllPosts(req, false)
	if err != nil {
		return &models.PostListResponse{}, err
	}
	if len(posts) == 0 {
		return &models.PostListResponse{}, nil
	}

	usersIds := utils.FromSlice(posts, func(p models.PostModel) int64 { return p.AuthorID }).Keys()

	users, err := s.UserService.GetAllUserDetailsByIndices(usersIds)
	if err != nil {
		return &models.PostListResponse{}, err
	}
	usersMap := utils.FromSlice(users, func(u models.UserDetails) int64 { return u.ID })

	var responses = []models.PostResponse{}
	for _, post := range posts {
		url, _ := s.GalleryService.GetLinkByAssetID(post.CoverImageID)

		responses = append(responses, models.PostResponse{
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

	return &models.PostListResponse{
		Posts:   responses,
		Current: req.Page,
		Pages:   total / req.Limit,
	}, nil
}

func (s *CmsService) UpdatePost(req *models.PostUpdateRequest) error {
	post, err := s.CmsRepo.GetPostByID(req.ID)
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
		_, err := s.CmsRepo.GetPostBySlug(req.Slug)
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

func (s *CmsService) GetViewTeamMembers() ([]models.TeamMemberViewResponse, error) {
	members, err := s.CmsRepo.GetAllTeamMembers(true)
	if err != nil {
		return nil, err
	}

	if len(members) == 0 {
		return []models.TeamMemberViewResponse{}, nil
	}

	var usersIds []int64
	for _, member := range members {
		usersIds = append(usersIds, member.UserID)
	}

	users, err := s.UserService.GetAllUserDetailsByIndices(usersIds)
	if err != nil {
		return nil, err
	}
	usersMap := utils.FromSlice(users, func(u models.UserDetails) int64 { return u.ID })

	var dtos []models.TeamMemberViewResponse
	for _, m := range members {
		dtos = append(dtos, models.TeamMemberViewResponse{
			UserID:       m.UserID,
			NameAr:       usersMap[m.UserID].NameAr,
			NameEn:       usersMap[m.UserID].NameEn,
			Role:         m.Role,
			Bio:          m.Bio,
			DisplayOrder: m.DisplayOrder,
			ProfilePic:   usersMap[m.UserID].ProfilePic,
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
