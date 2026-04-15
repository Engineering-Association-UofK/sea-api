package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/repositories"
	"sea-api/internal/utils"
	"sea-api/internal/utils/sheets"
	"sea-api/internal/utils/valid"
	"slices"
	"strconv"
	"time"
)

type UserService struct {
	repo             *repositories.UserRepository
	suspensionsRepo  *repositories.SuspensionsRepo
	S3StorageService *S3StorageService
}

func NewUserService(repo *repositories.UserRepository, suspensionsRepo *repositories.SuspensionsRepo, s3StorageService *S3StorageService) *UserService {
	// if _, err := repo.GetByUserID(0); err != nil {
	// 	err = repo.DetailedCreate(&models.UserModel{
	// 		ID:             0,
	// 		UniID:          1000000000,
	// 		ProfileImageID: sql.NullInt64{Int64: 0, Valid: false},
	// 		Username:       "system",
	// 		NameAr:         "النظام",
	// 		NameEn:         "System",
	// 		Email:          "system@sea.uofk.edu",
	// 		Phone:          "0000000000",
	// 		Department:     "IT",
	// 		Gender:         models.MALE,
	// 		Verified:       true,
	// 		Password:       "",
	// 		Status:         models.STATUS_ACTIVE,
	// 		IsEditable:     false,
	// 		IsLoggable:     false,
	// 		IsAnonymous:    true,
	// 	}, nil)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	return &UserService{repo: repo, suspensionsRepo: suspensionsRepo, S3StorageService: s3StorageService}
}

// ======== GET ALL ========

func (s *UserService) GetAllTempUsers(req *models.ListRequest) (*models.TempUserListResponse, error) {
	total, err := s.repo.GetTotal(req.Limit, true)
	if err != nil {
		return nil, err
	}
	valid.ValidateListRequest(req, total)

	users, err := s.repo.GetAllTempUsers(req.Limit, req.Page)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return &models.TempUserListResponse{
			Users: []models.TempUserResponse{},
			Pages: total / req.Limit,
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
		Pages: total / req.Limit,
	}, nil
}

func (s *UserService) GetAll(req *models.ListRequest) (*models.UserListResponse, error) {
	total, err := s.repo.GetTotal(req.Limit, false)
	if err != nil {
		return nil, err
	}
	valid.ValidateListRequest(req, total)

	users, err := s.repo.GetAll(req.Limit, req.Page)
	if err != nil {
		return nil, errs.New(errs.InternalServerError, "Error getting users: "+err.Error(), nil)
	}
	if len(users) == 0 {
		return &models.UserListResponse{
			Users: []models.UserListItemResponse{},
			Pages: total / req.Limit,
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
		Pages: total / req.Limit,
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
			link, err := s.S3StorageService.GenerateDownloadUrlByID(context.Background(), user.ProfileImageID.Int64)
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

// Get UserDetails
func (s *UserService) GetUserDetails(id int64) (*models.UserDetails, error) {
	user, err := s.repo.GetByUserID(id)
	if err != nil {
		return nil, err
	}

	url := ""
	if user.ProfileImageID.Valid {
		link, err := s.S3StorageService.GenerateDownloadUrlByID(context.Background(), user.ProfileImageID.Int64)
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

func (s *UserService) UpdateUsersImport(file io.Reader) error {
	mods, err := sheets.ParseExcelToStructs[models.ImportUserUpdate](file)
	if err != nil {
		return err
	}

	users, err := s.repo.GetAll(100, 1)
	if err != nil {
		return err
	}

	modsMap := utils.FromSlice(mods, func(u models.ImportUserUpdate) string { return u.Email })

	for _, u := range users {
		if user, ok := modsMap[u.Email]; ok {
			index, err := strconv.ParseInt(user.Index, 10, 64)
			if err != nil {
				slog.Error("user "+user.Index+" failed to update", "error", err)
			}
			u.ID = index
			u.Phone = user.Phone
			u.Status = models.STATUS_INACTIVE
			s.repo.UpdateWithID(&u, nil)
		}
	}

	return nil
}

// ############################################################
// ##################    ADMIN MANAGEMENT    ##################
// ############################################################

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
		url := ""
		if a.ProfileImageID.Valid {
			link, err := s.S3StorageService.GenerateDownloadUrlByID(context.Background(), a.ProfileImageID.Int64)
			if err != nil {
				return nil, err
			}
			url = link
		}
		adminResponses = append(adminResponses, models.AdminResponse{
			ID:         a.ID,
			Email:      a.Email,
			ProfilePic: url,
			NameAr:     a.NameAr,
			Username:   a.Username,
			Gender:     a.Gender,
			Roles:      rolesMap[a.ID],
		})
	}

	return adminResponses, nil
}

func (s *UserService) AddAdmin(ID int64) error {
	roles, err := s.GetRolesByUserID(ID)
	if err != nil {
		return err
	}
	if slices.Contains(roles, models.RoleSystemSuperAdmin) {
		return errs.New(errs.Forbidden, "Cannot add a Super Admin as an admin", nil)
	}
	if slices.Contains(roles, models.RoleSystemAdmin) {
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
	if slices.Contains(roles, models.RoleSystemSuperAdmin) {
		return errs.New(errs.Forbidden, "Cannot add a Super Admin as an admin", nil)
	}
	if slices.Contains(roles, models.RoleSystemAdminManager) {
		return errs.New(errs.Conflict, "User is already an admin manager", nil)
	}
	err = s.repo.CreateRole(&models.UserRole{UserID: ID, Role: models.RoleSystemAdminManager})
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) RemoveAdminManager(ID int64) error {
	roles, err := s.GetRolesByUserID(ID)
	if err != nil {
		return err
	}
	if slices.Contains(roles, models.RoleSystemSuperAdmin) {
		return errs.New(errs.Forbidden, "Cannot remove a Super Admin as an admin manager", nil)
	}
	if !slices.Contains(roles, models.RoleSystemAdminManager) {
		return errs.New(errs.NotFound, "User is not an admin manager", nil)
	}
	err = s.repo.RemoveRole(ID, models.RoleSystemAdminManager, nil)
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
	specialRoles := []models.Role{}
	for _, role := range roles {
		if role == models.RoleSystemSuperAdmin {
			return errs.New(errs.Forbidden, "Cannot update a Super Admin's roles", nil)
		}
		if models.SpecialAdminRoles[role] {
			specialRoles = append(specialRoles, role)
		}
	}
	if !slices.Contains(roles, models.RoleSystemAdmin) {
		return errs.New(errs.NotFound, "User is not an admin", nil)
	}

	var rolesToAdd = []models.Role{}
	for _, role := range req.Roles {
		if models.AllowedAdminRoles[role] {
			rolesToAdd = append(rolesToAdd, role)
		}
	}

	if len(rolesToAdd) == 0 {
		return nil
	}

	for _, role := range specialRoles {
		rolesToAdd = append(rolesToAdd, role)
	}

	err = s.repo.ReplaceRoles(req.ID, rolesToAdd, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) RemoveAdmin(ID int64) error {
	err := s.UpdateAdminRoles(&models.AdminRequest{
		ID:    ID,
		Roles: []models.Role{},
	})
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
