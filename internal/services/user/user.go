package user

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/repositories"
	"sea-api/internal/services/storage"
	"sea-api/internal/utils"
	"sea-api/internal/utils/valid"
	"slices"
	"time"
)

type UserService struct {
	repo            *repositories.UserRepository
	suspensionsRepo *repositories.SuspensionsRepo
	S3              *storage.S3
}

func NewUserService(repo *repositories.UserRepository, suspensionsRepo *repositories.SuspensionsRepo, S3 *storage.S3) *UserService {
	return &UserService{repo: repo, suspensionsRepo: suspensionsRepo, S3: S3}
}

// ======== GET ALL ========

func (s *UserService) GetAll(req *models.ListRequest) (*models.UserListResponse, error) {
	total, err := s.repo.GetTotal(req.Limit, false)
	if err != nil {
		return nil, err
	}
	valid.Limit(req, total)

	users, err := s.repo.GetAll(req.Limit, req.Page)
	if err != nil {
		return nil, errs.New(errs.InternalServerError, "Error getting users: "+err.Error(), nil)
	}
	ids := utils.ExtractField(users, func(u models.UserModel) int64 { return u.ID })

	roles, err := s.repo.GetAllRolesByUserIDs(ids)
	if err != nil {
		return nil, errs.New(errs.InternalServerError, "Error getting riles: "+err.Error(), nil)
	}
	rolesMap := extractRoles(roles)

	var userResponses []models.UserListItemResponse
	for _, u := range users {
		user := parseUserListResponse(&u, rolesMap[u.ID])
		userResponses = append(userResponses, *user)
	}

	return &models.UserListResponse{
		Users:   userResponses,
		Current: req.Page,
		Pages:   total / req.Limit,
	}, nil
}

func (s *UserService) GetAllUserDetailsByIndices(indices []int64) ([]models.UserDetails, error) {
	users, err := s.repo.GetAllByIndices(indices)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return []models.UserDetails{}, nil
	}

	var userResponses []models.UserDetails
	for _, user := range users {
		url := ""
		if user.ProfileImageID.Valid {
			link, err := s.S3.GenerateDownloadUrlByID(context.Background(), user.ProfileImageID.Int64)
			if err == nil {
				url = link
			}
		}
		userResponses = append(userResponses, models.UserDetails{
			UserProfileResponse: models.UserProfileResponse{
				ID:         user.ID,
				UniID:      user.UniID,
				Username:   user.Username,
				NameAr:     user.NameAr,
				NameEn:     user.NameEn,
				Email:      user.Email,
				Phone:      user.Phone,
				Gender:     user.Gender,
				Department: user.Department,
				ProfilePic: url,
			},
		})
	}

	return userResponses, nil
}

func (s *UserService) GetAllByIndices(indices []int64) ([]models.UserListItemResponse, error) {
	users, err := s.repo.GetAllByIndices(indices)
	if err != nil {
		return nil, err
	}
	roles, err := s.repo.GetAllRolesByUserIDs(indices)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return []models.UserListItemResponse{}, nil
	}

	rolesMap := extractRoles(roles)

	var userResponse []models.UserListItemResponse
	for _, u := range users {
		user := parseUserListResponse(&u, rolesMap[u.ID])
		userResponse = append(userResponse, *user)
	}

	return userResponse, nil
}

// ======== GET ========

func (s *UserService) GetUserDetails(id int64) (*models.UserDetails, error) {
	user, err := s.repo.GetByUserID(id)
	if err != nil {
		return nil, err
	}

	url := ""
	if user.ProfileImageID.Valid {
		link, err := s.S3.GenerateDownloadUrlByID(context.Background(), user.ProfileImageID.Int64)
		if err != nil {
			return nil, err
		}
		url = link
	}

	return &models.UserDetails{
		UserProfileResponse: models.UserProfileResponse{
			ID:         user.ID,
			UniID:      user.UniID,
			Username:   user.Username,
			NameAr:     user.NameAr,
			NameEn:     user.NameEn,
			Email:      user.Email,
			Phone:      user.Phone,
			Gender:     user.Gender,
			Department: user.Department,
			ProfilePic: url,
		},
	}, nil
}

func (s *UserService) GetByUserID(ctx context.Context, userID int64) (*models.UserResponse, error) {
	user, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	rolesModels, err := s.repo.GetRolesByUserID(userID)
	if err != nil {
		return nil, err
	}
	roles := utils.ExtractField(rolesModels, func(r models.UserRole) models.Role { return r.Role })

	url := ""
	if user.ProfileImageID.Valid {
		link, err := s.S3.GenerateDownloadUrlByID(ctx, user.ProfileImageID.Int64)
		if err != nil {
			return nil, err
		}
		url = link
	}

	return parseUserResponse(user, roles, url), nil
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*models.UserResponse, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	rolesModels, err := s.repo.GetRolesByUserID(user.ID)
	if err != nil {
		return nil, err
	}
	roles := utils.ExtractField(rolesModels, func(r models.UserRole) models.Role { return r.Role })

	url := ""
	if user.ProfileImageID.Valid {
		link, err := s.S3.GenerateDownloadUrlByID(ctx, user.ProfileImageID.Int64)
		if err != nil {
			return nil, err
		}
		url = link
	}

	return parseUserResponse(user, roles, url), nil
}

func (s *UserService) GetRolesByUserID(userID int64) ([]models.Role, error) {
	rolesModels, err := s.repo.GetRolesByUserID(userID)
	if err != nil {
		return nil, err
	}
	if len(rolesModels) == 0 {
		return []models.Role{}, nil
	}
	var roles []models.Role
	for _, r := range rolesModels {
		roles = append(roles, r.Role)
	}
	return roles, nil
}

// ======== UPDATE ========

func (s *UserService) Update(req *models.UpdateProfileRequest) error {
	user, err := s.repo.GetByUserID(req.ID)
	if err != nil {
		return err
	}
	roles, err := s.GetRolesByUserID(req.ID)
	if err != nil {
		return err
	}
	// Check is user is Super Admin to abort changes
	if slices.Contains(roles, models.RoleSystemSuperAdmin) {
		return errs.New(errs.Forbidden, "Cannot update a Super Admin profile", nil)
	}
	user.UniID = req.UniID
	user.NameAr = string(req.NameAr)
	user.NameEn = string(req.NameEn)
	user.Phone = string(req.Phone)
	user.Department = req.Department
	user.Gender = req.Gender
	if err := s.repo.Update(user, nil); err != nil {
		return err
	}
	return nil
}

// ======== SPECIAL ========

func (s *UserService) Suspend(req *models.SuspensionRequest, adminId int64) error {
	if u, err := s.suspensionsRepo.GetByUserID(req.UserID); err == nil {
		return errs.New(
			errs.Conflict,
			fmt.Sprintf("User is already suspended until %s", u.EndedAt.Format("2006-01-02 15:04:05")),
			nil,
		)
	}
	if rolesModels, err := s.repo.GetRolesByUserID(req.UserID); err == nil {
		roles := utils.ExtractField(rolesModels, func(r models.UserRole) models.Role { return r.Role })
		if slices.Contains(roles, models.RoleSystemSuperAdmin) {
			return errs.New(errs.Forbidden, "Cannot suspend a Super Admin", nil)
		}
	}
	tx, err := s.repo.BeginTransaction()
	if err != nil {
		return err
	}
	err = s.repo.Suspend(req.UserID, tx)
	if err != nil {
		return err
	}
	suspension := &models.SuspensionModel{
		UserID:    req.UserID,
		AdminID:   adminId,
		Reason:    req.Reason,
		StartedAt: time.Now(),
		EndedAt:   time.Now().Add(time.Millisecond * time.Duration(req.Duration)),
	}
	if _, err := s.suspensionsRepo.Create(suspension, tx, false); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := s.suspensionsRepo.Create(suspension, tx, true); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// ======== Helpers ========

func extractRoles(roles []models.UserRole) map[int64][]models.Role {
	rolesMap := make(map[int64][]models.Role)
	for _, r := range roles {
		if _, ok := rolesMap[r.UserID]; !ok {
			rolesMap[r.UserID] = []models.Role{}
		}
		rolesMap[r.UserID] = append(rolesMap[r.UserID], r.Role)
	}
	return rolesMap
}

func parseUserListResponse(user *models.UserModel, roles []models.Role) *models.UserListItemResponse {
	return &models.UserListItemResponse{
		ID:         user.ID,
		UniID:      user.UniID,
		Username:   user.Username,
		Email:      user.Email,
		Verified:   user.Verified,
		Status:     user.Status,
		Gender:     user.Gender,
		Department: user.Department,
		Roles:      roles,
	}
}

func parseUserResponse(user *models.UserModel, roles []models.Role, url string) *models.UserResponse {
	return &models.UserResponse{
		ID:         user.ID,
		UniID:      user.UniID,
		Username:   user.Username,
		ProfilePic: url,
		NameAr:     user.NameAr,
		NameEn:     user.NameEn,
		Email:      user.Email,
		Phone:      user.Phone,
		Department: user.Department,
		Gender:     user.Gender,
		Verified:   user.Verified,
		Status:     user.Status,
		Roles:      roles,
	}
}

func generatePasscode(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
