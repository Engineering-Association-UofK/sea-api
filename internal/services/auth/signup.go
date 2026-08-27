package auth

import (
	"crypto/sha256"
	"fmt"
	"log/slog"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

/*
## There-steps Account Creation

  - 1. Initial registration:
    Where the user enters their ID, passcode, and email. If the passcode is correct,
    an email containing the next step will be sent to the email, effectively verifying
    the email while on it.

  - 2. Credentials registration:
    After following the link, the user is prompted to enter their password.

  - 3. Details registration:
    After the password, the user enters general student details.

  - 4. Username registration:
    Then the user is redirected to the Username page where hey choose a username.
    On this step they are set as verified on the database as well

And done. The rest of the details can be filled later on if the user wanted to enter an event.
*/
func (s *AuthService) CheckRegistration(req *models.CheckRegistrationRequest) (*models.CheckRegistrationResponse, int64, error) {
	state, err := s.AuthRepository.GetStateWithCode(req.RegCode)
	if err != nil {
		return nil, 0, err
	}
	return &models.CheckRegistrationResponse{
		RegStep: state.Step,
	}, state.UserID, nil
}

// # First Step
func (s *AuthService) InitialRegistration(req *models.InitialRegistrationRequest) error {
	slog.Debug("Initial Registration Started")

	// Get models and check them
	tempUser, err := s.UserRepo.GetTempUser(req.UserID)
	if err != nil {
		return errs.New(errs.NotFound, "Student UserID was not found in out database, please contact administration", nil)
	}
	_, err = s.UserRepo.GetByUserID(req.UserID)
	if err == nil {
		return errs.New(errs.Conflict, "User with userID already exists", nil)
	}
	slog.Debug("User found and not already registered")

	// Check passed values
	if tempUser.Password.Valid && tempUser.Password.String != req.Passcode {
		return errs.New(errs.BadRequest, "Passcode is not valid", nil)
	}
	_, err = s.UserRepo.GetByEmail(string(req.Email))
	if err == nil {
		return errs.New(errs.Conflict, "Email already in use", nil)
	}
	slog.Debug("User passcode and email are clear")

	// Create user model
	err = s.UserRepo.StartUserRegistration(&models.RegInitCreate{
		ID:    req.UserID,
		Email: req.Email,
	})
	if err != nil {
		return err
	}

	// Delete temp user model
	err = s.UserRepo.DeleteTempUser(req.UserID, nil)
	if err != nil {
		slog.Error("error deleting temp user", "error", err, "user_id", req.UserID)
	}
	slog.Debug("User temp profile deleted")

	// Start registration counter
	data := []byte(fmt.Sprintf("%s|%d|%s", req.Email, req.UserID, time.Now()))
	hash := sha256.Sum256(data)
	err = s.AuthRepository.StartRegistration(&models.RegistrationStepModel{
		RegCode: fmt.Sprintf("%x", hash),
		UserID:  req.UserID,
		Step:    1,
	})
	slog.Debug("Registration process started")

	// Send email
	link := fmt.Sprintf("https://sea.uofk.edu/registration/%x", hash)
	slog.Debug("Sending email")
	return s.MailService.SendRegistrationMail(string(req.Email), link, req.Lang)
}

// # Second Step
// Also acts as a password reset function
func (s *AuthService) CredentialsRegistration(req *models.PasswordRegistrationRequest) error {
	state, userID, err := s.CheckRegistration(&models.CheckRegistrationRequest{RegCode: req.RegCode})
	if err != nil {
		return err
	}
	if state.RegStep != 1 && state.RegStep != 5 {
		return errs.New(errs.Forbidden, "Forbidden", nil)
	}

	if req.Password != req.ConfirmPassword || len(req.Password) < 8 || len(req.Password) > 32 {
		return errs.New(errs.BadRequest, "Password does not follow security requirements", nil)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = s.UserRepo.UpdatePassword(userID, string(hashedPassword))
	if err != nil {
		return err
	}

	if state.RegStep == 1 {
		return s.AuthRepository.IncrementStep(req.RegCode)
	} else {
		return s.AuthRepository.DecrementStep(req.RegCode)
	}
}

// # Third step
func (s *AuthService) DetailsRegistration(req *models.DetailsRegistrationRequest) error {
	state, userID, err := s.CheckRegistration(&models.CheckRegistrationRequest{RegCode: req.RegCode})
	if err != nil {
		return err
	}
	if state.RegStep != 2 {
		return errs.New(errs.Forbidden, "Forbidden", nil)
	}

	_, err = strconv.ParseInt(req.UniID, 10, 64)
	if err != nil {
		return err
	}

	err = s.UserRepo.UpdateDetails(&models.RegDetailsUpdate{
		UserID:     userID,
		NameAr:     req.NameAr,
		NameEn:     req.NameEn,
		Gender:     req.Gender,
		UniID:      req.UniID,
		Department: req.Department,
		Phone:      req.Phone,
	})
	if err != nil {
		return err
	}

	return s.AuthRepository.IncrementStep(req.RegCode)
}

// # Fourth and final step
func (s *AuthService) UsernameRegistration(req *models.UsernameRegistrationRequest) error {
	state, userID, err := s.CheckRegistration(&models.CheckRegistrationRequest{RegCode: req.RegCode})
	if err != nil {
		return err
	}
	if state.RegStep != 3 {
		return errs.New(errs.Forbidden, "Forbidden", nil)
	}

	_, err = s.UserRepo.GetByUsername(string(req.Username))
	if err == nil {
		return errs.New(errs.Conflict, "Username already in use", nil)
	}
	err = s.UserRepo.UpdateUsername(userID, string(req.Username))
	if err != nil {
		return err
	}

	err = s.UserRepo.VerifyUser(userID)
	if err != nil {
		return err
	}

	return s.AuthRepository.IncrementStep(req.RegCode)
}
