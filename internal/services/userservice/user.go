package userservice

import (
	"context"
	"fmt"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/models/usermodels"
	"sea-api/internal/repositories"
	"sea-api/internal/repositories/userrepo"
	"sea-api/internal/services/storage"
	"sea-api/internal/utils"
	"sea-api/internal/utils/valid"
	"slices"
	"time"
)

type UserService struct {
	repo            *userrepo.UserRepository
	suspensionsRepo *repositories.SuspensionsRepo
	S3              *storage.S3
}

func NewUserService(repo *userrepo.UserRepository, suspensionsRepo *repositories.SuspensionsRepo, S3 *storage.S3) *UserService {
	return &UserService{repo: repo, suspensionsRepo: suspensionsRepo, S3: S3}
}

// ======== GET ALL ========

func (s *UserService) GetAll(req *models.ListRequest) (*usermodels.UserListResponse, error) {
	total, err := s.repo.GetTotal()
	if err != nil {
		return nil, err
	}
	pages := valid.Limit(req, total)

	users, err := s.repo.GetAll(req.Limit, req.Page)
	if err != nil {
		return nil, errs.New(errs.InternalServerError, "Error getting users: "+err.Error(), nil)
	}
	ids := utils.ExtractField(users, func(u usermodels.UserModel) int64 { return u.ID })

	roles, err := s.repo.GetAllRolesByUserIDs(ids)
	if err != nil {
		return nil, errs.New(errs.InternalServerError, "Error getting riles: "+err.Error(), nil)
	}
	rolesMap := extractRoles(roles)

	var userResponses []usermodels.UserListItemResponse
	for _, u := range users {
		user := parseUserListResponse(&u, rolesMap[u.ID])
		userResponses = append(userResponses, *user)
	}

	return &usermodels.UserListResponse{
		Users:   userResponses,
		Current: req.Page,
		Pages:   pages,
	}, nil
}

func (s *UserService) GetAllUserDetailsByIndices(indices []int64) ([]usermodels.UserDetails, error) {
	users, err := s.repo.GetAllByIndices(indices)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return []usermodels.UserDetails{}, nil
	}

	var userResponses []usermodels.UserDetails
	for _, user := range users {
		url := ""
		if user.ProfileImageID.Valid {
			link, err := s.S3.GenerateDownloadUrlByID(context.Background(), user.ProfileImageID.Int64)
			if err == nil {
				url = link
			}
		}
		userResponses = append(userResponses, usermodels.UserDetails{
			UserProfileResponse: usermodels.UserProfileResponse{
				ID:         user.ID,
				UniID:      *user.UniID,
				Username:   *user.Username,
				NameAr:     *user.NameAr,
				NameEn:     *user.NameEn,
				Email:      *user.Email,
				Phone:      *user.Phone,
				Gender:     *user.Gender,
				Department: *user.Department,
				ProfilePic: url,
			},
		})
	}

	return userResponses, nil
}

func (s *UserService) GetAllByIndices(indices []int64) ([]usermodels.UserListItemResponse, error) {
	users, err := s.repo.GetAllByIndices(indices)
	if err != nil {
		return nil, err
	}
	roles, err := s.repo.GetAllRolesByUserIDs(indices)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return []usermodels.UserListItemResponse{}, nil
	}

	rolesMap := extractRoles(roles)

	var userResponse []usermodels.UserListItemResponse
	for _, u := range users {
		user := parseUserListResponse(&u, rolesMap[u.ID])
		userResponse = append(userResponse, *user)
	}

	return userResponse, nil
}

// ======== GET ========

func (s *UserService) GetUserDetails(id int64) (*usermodels.UserDetails, error) {
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

	return &usermodels.UserDetails{
		UserProfileResponse: usermodels.UserProfileResponse{
			ID:         user.ID,
			UniID:      *user.UniID,
			Username:   *user.Username,
			NameAr:     *user.NameAr,
			NameEn:     *user.NameEn,
			Email:      *user.Email,
			Phone:      *user.Phone,
			Gender:     *user.Gender,
			Department: *user.Department,
			ProfilePic: url,
		},
	}, nil
}

func (s *UserService) GetByUserID(ctx context.Context, userID int64) (*usermodels.UserResponse, error) {
	user, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	rolesModels, err := s.repo.GetRolesByUserID(userID)
	if err != nil {
		return nil, err
	}
	roles := utils.ExtractField(rolesModels, func(r usermodels.UserRole) usermodels.Role { return r.Role })

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

func (s *UserService) GetByUsername(ctx context.Context, username string) (*usermodels.UserResponse, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	rolesModels, err := s.repo.GetRolesByUserID(user.ID)
	if err != nil {
		return nil, err
	}
	roles := utils.ExtractField(rolesModels, func(r usermodels.UserRole) usermodels.Role { return r.Role })

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

func (s *UserService) GetRolesByUserID(userID int64) ([]usermodels.Role, error) {
	rolesModels, err := s.repo.GetRolesByUserID(userID)
	if err != nil {
		return nil, err
	}
	if len(rolesModels) == 0 {
		return []usermodels.Role{}, nil
	}
	var roles []usermodels.Role
	for _, r := range rolesModels {
		roles = append(roles, r.Role)
	}
	return roles, nil
}

// ======== UPDATE ========

func (s *UserService) Update(req *usermodels.UpdateProfileRequest) error {
	user, err := s.repo.GetByUserID(req.ID)
	if err != nil {
		return err
	}
	roles, err := s.GetRolesByUserID(req.ID)
	if err != nil {
		return err
	}
	// Check is user is Super Admin to abort changes
	if slices.Contains(roles, usermodels.RoleSystemSuperAdmin) {
		return errs.New(errs.Forbidden, "Cannot update a Super Admin profile", nil)
	}
	user.UniID = &[]string{req.UniID}[0]
	user.NameAr = &[]string{string(req.NameAr)}[0]
	user.NameEn = &[]string{string(req.NameEn)}[0]
	user.Phone = &[]string{string(req.Phone)}[0]
	user.Department = &[]models.Department{req.Department}[0]
	user.Gender = &[]usermodels.Gender{req.Gender}[0]
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
		roles := utils.ExtractField(rolesModels, func(r usermodels.UserRole) usermodels.Role { return r.Role })
		if slices.Contains(roles, usermodels.RoleSystemSuperAdmin) {
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

func extractRoles(roles []usermodels.UserRole) map[int64][]usermodels.Role {
	rolesMap := make(map[int64][]usermodels.Role)
	for _, r := range roles {
		if _, ok := rolesMap[r.UserID]; !ok {
			rolesMap[r.UserID] = []usermodels.Role{}
		}
		rolesMap[r.UserID] = append(rolesMap[r.UserID], r.Role)
	}
	return rolesMap
}

func parseUserListResponse(user *usermodels.UserModel, roles []usermodels.Role) *usermodels.UserListItemResponse {
	return &usermodels.UserListItemResponse{
		ID:         user.ID,
		UniID:      *user.UniID,
		Username:   *user.Username,
		Email:      *user.Email,
		Verified:   user.Verified,
		Status:     user.Status,
		Gender:     *user.Gender,
		Department: *user.Department,
		Roles:      roles,
	}
}

func parseUserResponse(user *usermodels.UserModel, roles []usermodels.Role, url string) *usermodels.UserResponse {
	return &usermodels.UserResponse{
		ID:         user.ID,
		UniID:      *user.UniID,
		Username:   *user.Username,
		ProfilePic: url,
		NameAr:     *user.NameAr,
		NameEn:     *user.NameEn,
		Email:      *user.Email,
		Phone:      *user.Phone,
		Department: *user.Department,
		Gender:     *user.Gender,
		Verified:   user.Verified,
		Status:     user.Status,
		Roles:      roles,
	}
}
