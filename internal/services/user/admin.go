package user

import (
	"context"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/utils"
	"slices"
)

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
			link, err := s.S3.GenerateDownloadUrlByID(context.Background(), a.ProfileImageID.Int64)
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
