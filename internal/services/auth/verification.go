package auth

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Deprecated: No longer needed with new registration system as of v1.0.1
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
	if time.Now().After(code.CreatedAt.Add(time.Minute * 60)) {
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

// Deprecated: No longer needed with new registration system as of v1.0.1
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
	err = s.MailService.SendVerificationCode(*user.Email, models.VerifyEmail{
		Input: code,
		Year:  time.Now().Year(),
	})
	if err != nil {
		return err
	}
	return nil
}

// ====== HELPERS ======

// Deprecated: No longer needed with new registration system as of v1.0.1
func generateVerifyCode() (string, error) {
	max := big.NewInt(900000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}
