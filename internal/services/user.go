package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/repositories"
	"sea-api/internal/utils"
	"slices"
	"time"
)

type UserService struct {
	repo             *repositories.UserRepository
	suspensionsRepo  *repositories.SuspensionsRepo
	S3StorageService *S3StorageService
}

func NewUserService(repo *repositories.UserRepository, suspensionsRepo *repositories.SuspensionsRepo, s3StorageService *S3StorageService) *UserService {
	return &UserService{repo: repo, suspensionsRepo: suspensionsRepo, S3StorageService: s3StorageService}
}

// ======== GET ALL ========

func (s *UserService) GetAllTempUsers(req *models.UserListRequest) (*models.TempUserListResponse, error) {
	pages, err := s.repo.GetPagesCount(int(req.Limit), true)
	if err != nil {
		return nil, err
	}
	if req.Page > pages {
		return &models.TempUserListResponse{
			Users: []models.TempUserResponse{},
			Pages: pages,
		}, nil
	}
	users, err := s.repo.GetAllTempUsers(int(req.Limit), req.Page)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return &models.TempUserListResponse{
			Users: []models.TempUserResponse{},
			Pages: pages,
		}, nil
	}

	var userResponses []models.TempUserResponse
	for _, u := range users {
		userResponses = append(userResponses, models.TempUserResponse{
			ID:       u.ID.Int64,
			NameAr:   u.NameAr.String,
			Passcode: u.Password.String,
		})
	}

	return &models.TempUserListResponse{
		Users: userResponses,
		Pages: pages,
	}, nil
}

func (s *UserService) GetAll(req *models.UserListRequest) (*models.UserListResponse, error) {
	if req.Page < 1 || req.Limit < 1 {
		return nil, errs.New(errs.BadRequest, "Given page and limit must be greater than 0", nil)
	}
	if !models.ListLimit[req.Limit] {
		return nil, errs.New(errs.BadRequest, "Given limit is not valid", nil)
	}

	pages, err := s.repo.GetPagesCount(int(req.Limit), false)
	if err != nil {
		return nil, errs.New(errs.InternalServerError, "Error getting pages count: "+err.Error(), nil)
	}
	if req.Page > pages {
		return &models.UserListResponse{
			Users: []models.UserListItemResponse{},
			Pages: pages,
		}, nil
	}
	users, err := s.repo.GetAll(int(req.Limit), req.Page)
	if err != nil {
		return nil, errs.New(errs.InternalServerError, "Error getting users: "+err.Error(), nil)
	}
	if len(users) == 0 {
		return &models.UserListResponse{
			Users: []models.UserListItemResponse{},
			Pages: pages,
		}, nil
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
		Users: userResponses,
		Pages: pages,
	}, nil
}

func (s *UserService) GetAllUserDetailsByIndices(indices []int64) ([]models.UserDetails, error) {
	users, err := s.repo.GetAllByIndices(indices)
	if err != nil {
		return nil, err
	}

	var userResponses []models.UserDetails
	for _, user := range users {
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

// Get UserDetails
func (s *UserService) GetUserDetails(id int64) (*models.UserDetails, error) {
	user, err := s.repo.GetByUserID(id)
	if err != nil {
		return nil, err
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
		link, err := s.S3StorageService.GenerateDownloadUrlByID(ctx, user.ProfileImageID.Int64)
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
		link, err := s.S3StorageService.GenerateDownloadUrlByID(ctx, user.ProfileImageID.Int64)
		if err != nil {
			return nil, err
		}
		url = link
	}

	return parseUserResponse(user, roles, url), nil
}

func (s *UserService) GetTempUserPasscode(userID int64) (*models.GetPasscodeResponse, error) {
	user, err := s.repo.GetTempUser(userID)
	if err != nil {
		return nil, err
	}
	return &models.GetPasscodeResponse{
		Passcode: user.Password.String,
	}, nil
}

func (s *UserService) GetRolesByUserID(userID int64) ([]models.Role, error) {
	rolesModels, err := s.repo.GetRolesByUserID(userID)
	if err != nil {
		return nil, err
	}
	return utils.ExtractField(rolesModels, func(u models.UserRole) models.Role { return u.Role }), nil
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
	if slices.Contains(roles, models.ROLE_SUPER_ADMIN) {
		return errs.New(errs.Forbidden, "Cannot update a Super Admin profile", nil)
	}
	user.UniID = req.UniID
	user.NameAr = req.NameAr
	user.NameEn = req.NameEn
	user.Phone = req.Phone
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
		if slices.Contains(roles, models.ROLE_SUPER_ADMIN) {
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

func (s *UserService) AssignPasscodes(progressChan chan string) error {
	defer close(progressChan)
	tempUsers, err := s.repo.GetTempUsersWithNullPasswords()
	if err != nil {
		return err
	}
	if len(tempUsers) == 0 {
		return nil
	}
	total := len(tempUsers)
	for i, u := range tempUsers {
		passcode, err := generatePasscode(6)
		if err != nil {
			utils.ParseProgressStruct(total, i+1, u.ID.Int64, false, u.NameAr.String, progressChan)
			continue
		}
		err = s.repo.UpdateTempPasscode(u.ID.Int64, passcode)
		if err != nil {
			utils.ParseProgressStruct(total, i+1, u.ID.Int64, false, u.NameAr.String, progressChan)
			continue
		}
		utils.ParseProgressStruct(total, i+1, u.ID.Int64, true, u.NameAr.String, progressChan)
	}
	progressChan <- "done"
	return nil
}

// ###### ADMIN MANAGEMENT ######

func (s *UserService) GetAdmins() ([]models.AdminResponse, error) {
	admins, err := s.repo.GetAdmins()
	if err != nil {
		return nil, err
	}
	if len(admins) == 0 {
		return []models.AdminResponse{}, nil
	}
	rolesModels, err := s.repo.GetAllRolesByUserIDs(utils.ExtractField(admins, func(u models.UserModel) int64 { return u.ID }))
	if err != nil {
		return nil, err
	}
	rolesMap := extractRoles(rolesModels)

	var adminResponses []models.AdminResponse
	for _, a := range admins {
		adminResponses = append(adminResponses, models.AdminResponse{
			ID:       a.ID,
			Email:    a.Email,
			NameAr:   a.NameAr,
			Username: a.Username,
			Gender:   a.Gender,
			Roles:    rolesMap[a.ID],
		})
	}

	return adminResponses, nil
}

func (s *UserService) AddAdmin(ID int64) error {
	roles, err := s.GetRolesByUserID(ID)
	if err != nil {
		return err
	}
	if slices.Contains(roles, models.ROLE_SUPER_ADMIN) {
		return errs.New(errs.Forbidden, "Cannot add a Super Admin as an admin", nil)
	}
	if slices.Contains(roles, models.ROLE_ADMIN) {
		return errs.New(errs.Conflict, "User is already an admin", nil)
	}
	err = s.repo.AddAdmin(ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) MakeAdminManager(ID int64) error {
	roles, err := s.GetRolesByUserID(ID)
	if err != nil {
		return err
	}
	if slices.Contains(roles, models.ROLE_SUPER_ADMIN) {
		return errs.New(errs.Forbidden, "Cannot add a Super Admin as an admin", nil)
	}
	if slices.Contains(roles, models.ROLE_ADMIN_MANAGER) {
		return errs.New(errs.Conflict, "User is already an admin manager", nil)
	}
	err = s.repo.CreateRole(&models.UserRole{UserID: ID, Role: models.ROLE_ADMIN_MANAGER})
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) UpdateAdminRoles(req *models.AdminRequest) error {
	roles, err := s.GetRolesByUserID(req.ID)
	if err != nil {
		return err
	}
	if slices.Contains(roles, models.ROLE_SUPER_ADMIN) {
		return errs.New(errs.Forbidden, "Cannot update a Super Admin's roles", nil)
	}
	for _, role := range req.Roles {
		if !models.AllowedAdminRoles[role] {
			return errs.New(errs.BadRequest, "Invalid role", nil)
		}
	}
	err = s.repo.ReplaceRoles(req.ID, req.Roles, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) RemoveAdmin(ID int64) error {
	roles, err := s.GetRolesByUserID(ID)
	if err != nil {
		return err
	}
	if slices.Contains(roles, models.ROLE_SUPER_ADMIN) {
		return errs.New(errs.Forbidden, "Cannot remove a Super Admin as an admin", nil)
	}
	if !slices.Contains(roles, models.ROLE_ADMIN) {
		return errs.New(errs.Conflict, "User is not an admin", nil)
	}
	err = s.repo.RemoveAdmin(ID)
	if err != nil {
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
