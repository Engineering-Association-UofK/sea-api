package userservice

import (
	"crypto/rand"
	"encoding/base64"
	"sea-api/internal/errs"
	"sea-api/internal/models/usermodels"
)

func (s *UserService) GetPasscode(userID int64) (*usermodels.GetPasscodeResponse, error) {
	passcode, err := s.repo.GetPasscode(userID)
	if err == nil {
		return &usermodels.GetPasscodeResponse{Passcode: passcode}, nil
	}
	if _, err = s.repo.GetUserRow(userID); err != nil {
		return nil, errs.New(errs.Conflict, "User already Registered", nil)
	}
	passcode, err = generateAndSavePasscode(8)
	if err != nil {
		return nil, err
	}

	err = s.repo.CreatePasscode(userID, passcode)
	if err != nil {
		return nil, err
	}

	return &usermodels.GetPasscodeResponse{Passcode: passcode}, nil
}

func generateAndSavePasscode(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
