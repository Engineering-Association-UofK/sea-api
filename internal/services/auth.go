package services

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"log/slog"
	"math/big"
	"sea-api/internal/config"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/utils"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo         IUserRepository
	MailService      IMailService
	VerificationRepo IVerificationRepo

	SecretKey  []byte
	Issuer     string
	ExpiryTime time.Duration
}

func NewAuthService(userRepo IUserRepository, mailService IMailService, verificationRepo IVerificationRepo) *AuthService {
	return &AuthService{
		UserRepo:         userRepo,
		MailService:      mailService,
		VerificationRepo: verificationRepo,
		SecretKey:        []byte(config.App.JwtSecret),
		Issuer:           "SEA - UofK API Server Authorization",
		ExpiryTime:       time.Hour * 168,
	}
}

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

	if !user.Verified {
		return &models.LoginResponse{
			Token:      "",
			UserID:     user.ID,
			IsVerified: false,
		}, nil
	}

	if user.Status != "active" {
		return nil, errs.New(errs.Forbidden, "User is not active. This can happen for a lot of reasons, please contact the administration to resolve this issue.", nil)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errs.New(errs.Forbidden, "Invalid credentials", nil)
	}

	rolesModels, err := s.UserRepo.GetRolesByUserID(user.ID)
	if err != nil {
		return nil, err
	}
	roles := utils.ExtractField(rolesModels, func(r models.UserRole) models.Role { return r.Role })

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

func (s *AuthService) Register(req models.RegisterRequest) error {

	tempUser, err := s.UserRepo.GetTempUser(req.UserID)
	if err != nil {
		return errs.New(errs.NotFound, "Student UserID was not found in out database, please contact administration", nil)
	}
	_, err = s.UserRepo.GetByUserID(req.UserID)
	if err == nil {
		return errs.New(errs.Conflict, "User with userID already exists", nil)
	}
	if tempUser.Password.Valid && tempUser.Password.String != req.Passcode {
		return errs.New(errs.BadRequest, "Passcode is not valid", nil)
	}
	_, err = s.UserRepo.GetByEmail(string(req.Email))
	if err == nil {
		return errs.New(errs.Conflict, "Email already in use", nil)
	}
	_, err = s.UserRepo.GetByUsername(string(req.Username))
	if err == nil {
		return errs.New(errs.Conflict, "Username already in use", nil)
	}
	if req.Password != req.ConfirmPassword || len(req.Password) < 8 || len(req.Password) > 32 {
		return errs.New(errs.BadRequest, "There was an error with your password, please try again", nil)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.UserRepo.Create(&models.UserModel{
		ID:         req.UserID,
		UniID:      req.UniID,
		Username:   string(req.Username),
		NameAr:     req.NameAr,
		NameEn:     req.NameEn,
		Email:      string(req.Email),
		Phone:      req.Phone,
		Password:   string(hashedPassword),
		Gender:     req.Gender,
		Department: req.Department,
		Verified:   false,
		Status:     models.STATUS_ACTIVE,
	}, nil)
	if err != nil {
		return err
	}

	err = s.UserRepo.DeleteTempUser(req.UserID, nil)
	if err != nil {
		slog.Error("error deleting temp user", "error", err, "user_id", req.UserID)
	}
	return nil
}

func (s *AuthService) Verify(req models.VerifyRequest) error {
	code, err := s.VerificationRepo.GetByUserID(req.UserID)
	if err != nil {
		return err
	}
	user, err := s.UserRepo.GetByUserID(req.UserID)
	if err != nil {
		return err
	}
	if user.Verified {
		return errs.New(errs.Conflict, "User is already verified", nil)
	}
	if bcrypt.CompareHashAndPassword([]byte(code.Code), []byte(req.Code)) != nil {
		return errs.New(errs.BadRequest, "Invalid verification code", nil)
	}
	if time.Now().After(code.CreatedAt.Add(time.Minute * 15)) {
		return errs.New(errs.BadRequest, "Verification code has expired", nil)
	}
	err = s.UserRepo.Verify(req.UserID)
	if err != nil {
		return err
	}
	err = s.VerificationRepo.Delete(code.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) SendVerificationCode(userID int64) error {
	user, err := s.UserRepo.GetByUserID(userID)
	if err != nil {
		return err
	}
	if user.Verified {
		return errs.New(errs.Conflict, "User is already verified", nil)
	}

	for {
		oldCode, err := s.VerificationRepo.GetByUserID(userID)
		if err != nil {
			if err == sql.ErrNoRows {
				break
			}
			return err
		}
		if err := s.VerificationRepo.Delete(oldCode.ID); err != nil {
			return err
		}
	}

	code, err := generateVerifyCode()
	if err != nil {
		return err
	}
	hashedCode, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	codeModel := &models.VerificationCodeModel{
		Code:      string(hashedCode),
		UserID:    userID,
		CreatedAt: time.Now(),
	}
	err = s.VerificationRepo.Create(codeModel)
	if err != nil {
		return err
	}
	err = s.MailService.SendVerificationCode(user.Email, models.VerifyEmail{
		Input: code,
		Year:  time.Now().Year(),
	})
	if err != nil {
		return err
	}
	return nil
}

// ====== HELPERS ======

func generateVerifyCode() (string, error) {
	max := big.NewInt(900000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}
