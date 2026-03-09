package services

import (
	"fmt"
	"sea-api/internal/models"
	"sea-api/internal/repositories"

	"github.com/jmoiron/sqlx"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(db *sqlx.DB) *UserService {
	return &UserService{repo: repositories.NewUserRepository(db)}
}

func (s *UserService) GetAll() ([]models.UserResponse, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("Error getting users: %s", err)
	}
	roles, err := s.repo.GetAllRoles()
	if err != nil {
		return nil, fmt.Errorf("Error getting riles: %s", err)
	}

	if len(users) == 0 {
		return []models.UserResponse{}, nil
	}

	rolesMap := extractRoles(roles)

	var userResponse []models.UserResponse
	for _, u := range users {
		user := models.UserResponse{
			Index:    u.Index,
			UniID:    u.UniID,
			Username: u.Username,
			NameAr:   u.NameAr,
			NameEn:   u.NameEn,
			Email:    u.Email,
			Phone:    u.Phone,
			Verified: u.Verified,
			Status:   u.Status,
			Roles:    rolesMap[u.Index],
		}
		userResponse = append(userResponse, user)
	}

	return userResponse, nil
}

func (s *UserService) GetByIndex(index int64) (*models.UserResponse, error) {
	user, err := s.repo.GetByIndex(index)
	if err != nil {
		return nil, err
	}
	roles, err := s.repo.GetRolesByUserID(index)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		Index:    user.Index,
		UniID:    user.UniID,
		Username: user.Username,
		NameAr:   user.NameAr,
		NameEn:   user.NameEn,
		Email:    user.Email,
		Phone:    user.Phone,
		Verified: user.Verified,
		Status:   user.Status,
		Roles:    roles,
	}, nil
}

func (s *UserService) Create(user *models.UserModel) error {
	return s.repo.Create(user)
}

func (s *UserService) Update(user *models.UserModel) error {
	return s.repo.Update(user)
}

func (s *UserService) Delete(index int64) error {
	return s.repo.Delete(index)
}

func extractRoles(roles []models.UserRole) map[int64][]string {
	rolesMap := make(map[int64][]string)
	for _, r := range roles {
		if _, ok := rolesMap[r.Index]; !ok {
			rolesMap[r.Index] = []string{}
		}
		rolesMap[r.Index] = append(rolesMap[r.Index], r.Role)
	}
	return rolesMap
}
