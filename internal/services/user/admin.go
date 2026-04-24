package user

import (
	"context"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/utils/valid"
	"slices"
)

func (s *UserService) GetAdmins(ctx context.Context, req *models.ListRequest) (*models.AdminResponseList, error) {
	total, err := s.repo.GetAdminsCount()
	if err != nil {
		return nil, err
	}
	valid.Limit(req, total)

	admins, err := s.repo.GetAdmins(req)
	if err != nil {
		return nil, err
	}

	adminMap := map[int64]models.AdminRow{}
	adminRoles := map[int64][]models.Role{}
	for _, a := range admins {
		role := a.Role
		if _, ok := adminRoles[a.ID]; !ok {
			adminRoles[a.ID] = []models.Role{}
		}
		adminRoles[a.ID] = append(adminRoles[a.ID], role)
		if _, ok := adminMap[a.ID]; !ok {
			adminMap[a.ID] = a
		}
	}

	var adminResponses []models.AdminResponse
	for _, a := range admins {
		url := ""
		if a.ProfilePic.Valid {
			url, err = s.S3.GenerateDownloadUrlByKey(ctx, a.ProfilePic.String)
			if err != nil {
				return nil, err
			}
		}
		adminResponses = append(adminResponses, models.AdminResponse{
			ID:         a.ID,
			Email:      a.Email,
			ProfilePic: url,
			NameAr:     a.NameAr,
			Username:   a.Username,
			Gender:     a.Gender,
			Roles:      adminRoles[a.ID],
		})
	}

	return &models.AdminResponseList{
		Admins:  adminResponses,
		Total:   total,
		Current: req.Page,
		Pages:   total / req.Limit,
	}, nil
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
