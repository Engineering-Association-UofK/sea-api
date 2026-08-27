package auth

import (
	"fmt"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/utils"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	var user *models.UserModel = nil
	var userID int64 = 0
	var err error = nil

	if userID, err = strconv.ParseInt(string(req.Username), 10, 64); err == nil {
		user, err = s.UserRepo.GetByUserID(userID)
	} else if strings.Contains(string(req.Username), "@") {
		user, err = s.UserRepo.GetByEmail(string(req.Username))
	} else {
		user, err = s.UserRepo.GetByUsername(string(req.Username))
	}
	if err != nil {
		return nil, err
	}

	State, err := s.AuthRepository.GetStateWithID(user.ID)
	if err != nil {
		return nil, err
	}

	// The number of the registration step is 4 in code
	//   4 means all registration steps are done
	if State.Step != 4 && user.Password == nil {
		return &models.LoginResponse{
			Token:       "",
			UserID:      user.ID,
			IsVerified:  false,
			RedirectURL: fmt.Sprintf(`https://sea.uofk.edu/registration/%s`, State.RegCode),
		}, nil
	}

	if user.Status != "active" {
		return nil, errs.New(errs.Forbidden, "User is not active. This can happen for a lot of reasons, please contact the administration to resolve this issue.", nil)
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(req.Password))
	if err != nil {
		return nil, errs.New(errs.Forbidden, "Invalid credentials", nil)
	}

	rolesModels, err := s.UserRepo.GetRolesByUserID(user.ID)
	if err != nil {
		return nil, err
	}
	roles := []models.Role{}
	roles = utils.ExtractField(rolesModels, func(r models.UserRole) models.Role { return r.Role })

	claims := &models.ManagedClaims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.ExpiryTime)),
			Issuer:    s.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.SecretKey)
	if err != nil {
		return nil, err
	}
	return &models.LoginResponse{
		Token:      tokenString,
		UserID:     user.ID,
		Roles:      roles,
		IsVerified: user.Verified,
	}, nil
}

// Checks the user details and sends them a link to complete the registration process if
// it was not done yet, or take them to the extra step for resetting password.
func (s *AuthService) ForgotPassword(req *models.ForgotPasswordRequest) error {
	var user *models.UserModel = nil
	var err error = nil

	if req.UserID != nil {
		user, err = s.UserRepo.GetByUserID(*req.UserID)
	} else if req.Email != nil {
		user, err = s.UserRepo.GetByEmail(*req.Email)
	} else if req.Username != nil {
		user, err = s.UserRepo.GetByUsername(*req.Username)
	} else {
		return errs.New(errs.BadRequest, "Fields empty", nil)
	}
	if err != nil {
		return errs.New(errs.BadRequest, "account not found, try another details or contact support", nil)
	}

	var stringToReturn string

	state, err := s.AuthRepository.GetStateWithID(user.ID)
	if err != nil {
		return err
	}

	if user.Verified && state.Step == 4 {
		err = s.AuthRepository.IncrementStep(state.RegCode)
	}

	stringToReturn = fmt.Sprintf("https://sea.uofk.edu/registration/%s", state.RegCode)

	return s.MailService.SendForgotPassMail(user.Email, stringToReturn, req.Lang)

}
