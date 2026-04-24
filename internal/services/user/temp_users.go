package user

import (
	"sea-api/internal/models"
	"sea-api/internal/utils"
	"sea-api/internal/utils/valid"
)

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
			Users:   []models.TempUserResponse{},
			Current: req.Page,
			Pages:   total / req.Limit,
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
		Users:   userResponses,
		Current: req.Page,
		Pages:   total / req.Limit,
	}, nil
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
