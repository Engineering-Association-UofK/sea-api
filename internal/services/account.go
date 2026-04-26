package services

import (
	"context"
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
	"regexp"
	"sea-api/internal/config"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/repositories"
	"sea-api/internal/services/storage"
	"sea-api/internal/utils"
	"sea-api/internal/utils/valid"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AccountService struct {
	UserRepo              *repositories.UserRepository
	store                 *storage.S3
	certificateRepository *repositories.CertificateRepository

	profilePath string
}

func NewAccountService(UserRepo *repositories.UserRepository, store *storage.S3, certificateRepository *repositories.CertificateRepository) *AccountService {
	return &AccountService{
		UserRepo:              UserRepo,
		store:                 store,
		certificateRepository: certificateRepository,
		profilePath:           "public/profiles",
	}
}

func (s *AccountService) GetProfileSummary(ctx context.Context, claims *models.ManagedClaims) (*models.UserProfileSummaryResponse, error) {
	user, err := s.UserRepo.GetUserRow(claims.UserID)
	if err != nil {
		return nil, err
	}

	url := ""
	if user.ProfilePicKey.Valid {
		url, err = s.store.GenerateDownloadUrlByKey(ctx, user.ProfilePicKey.String)
		if err != nil {
			return nil, err
		}
	}

	return &models.UserProfileSummaryResponse{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		ProfilePic: url,
	}, nil
}

func (s *AccountService) GetProfile(ctx context.Context, claims *models.ManagedClaims) (*models.UserProfileResponse, error) {
	user, err := s.UserRepo.GetByUserID(claims.UserID)
	if err != nil {
		return nil, err
	}
	url := ""
	if user.ProfileImageID.Valid {
		url, err = s.store.GenerateDownloadUrlByID(ctx, user.ProfileImageID.Int64)
		if err != nil {
			return nil, err
		}
	}

	return &models.UserProfileResponse{
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
	}, nil
}

func (s *AccountService) GetCertificates(claims *models.ManagedClaims, req *models.ListRequest) ([]models.CertificateListResponse, error) {
	total, err := s.certificateRepository.GetTotalByUserID(claims.UserID)
	if err != nil {
		return nil, err
	}
	valid.Limit(req, total)

	certs, err := s.certificateRepository.GetByUserID(claims.UserID, req)
	if err != nil {
		return nil, err
	}

	var responses = []models.CertificateListResponse{}
	for _, cert := range certs {
		responses = append(responses, models.CertificateListResponse{
			Hash:      cert.Hash,
			UserID:    cert.UserID,
			EventID:   cert.EventID,
			Grade:     cert.Grade,
			IssueDate: cert.IssueDate,
			Status:    cert.Status,
		})
	}

	return responses, nil
}

func (s *AccountService) UpdateProfile(claims *models.ManagedClaims, req models.UpdateProfileRequest) error {
	user, err := s.UserRepo.GetByUserID(req.ID)
	if err != nil {
		return err
	}
	errsMap := utils.Mpp[string, string]{}

	if len(strings.Split(string(req.NameAr), " ")) < 2 {
		errsMap.Add("name_ar", "Name in Arabic is not valid")
	}
	if len(strings.Split(string(req.NameEn), " ")) < 2 {
		errsMap.Add("name_en", "Name in English is not valid")
	}
	if _, err := s.UserRepo.GetByUniID(req.UniID); err != nil {
		errsMap.Add("uni_id", "University ID is not valid")
	}
	if len(errsMap) != 0 {
		return errs.New(errs.MultiBadRequest, "Invalid fields", errsMap)
	}

	user.UniID = req.UniID
	user.NameAr = string(req.NameAr)
	user.NameEn = string(req.NameEn)
	user.Gender = req.Gender
	user.Department = req.Department
	user.Phone = string(req.Phone)
	return s.UserRepo.Update(user, nil)
}

func (s *AccountService) UpdateProfilePicture(ctx context.Context, claims *models.ManagedClaims, file io.Reader) error {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	user, err := s.UserRepo.GetByUserID(claims.UserID)
	if err != nil {
		return err
	}

	contentType := http.DetectContentType(fileBytes)
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}
	if !allowedTypes[contentType] {
		return errs.New(errs.BadRequest, "unsupported file type:"+contentType, nil)
	}

	hash := fnv.New64()
	hash.Write([]byte(fmt.Sprint(claims.UserID) + config.App.SecretSalt))
	fileKey := fmt.Sprintf("%s/%d/%d-%d.%s", s.profilePath, time.Now().Year(), hash.Sum64(), claims.UserID, contentType[6:])

	if user.ProfileImageID.Valid {
		s.store.Delete(ctx, user.ProfileImageID.Int64)
	}
	id, err := s.store.Upload(ctx, fileKey, fileBytes, contentType)
	if err != nil {
		return err
	}
	user.ProfileImageID.Valid = true
	user.ProfileImageID.Int64 = id

	return s.UserRepo.Update(user, nil)
}

func (s *AccountService) UpdatePassword(claims *models.ManagedClaims, req models.UpdatePasswordRequest) error {
	user, err := s.UserRepo.GetByUserID(claims.UserID)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword))
	if err != nil {
		return errs.New(errs.Forbidden, "Invalid old password", nil)
	}
	if req.NewPassword != req.ConfirmPassword {
		return errs.New(errs.BadRequest, "passwords do not match", nil)
	}
	if msgs := isPasswordStrongEnough(req.NewPassword); msgs != nil {
		return errs.New(errs.MultiBadRequest, "password is not strong enough", msgs)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.UserRepo.Update(user, nil)
}

func (s *AccountService) UpdateEmail(claims *models.ManagedClaims, req models.UpdateEmailRequest) error {
	user, err := s.UserRepo.GetByUserID(claims.UserID)
	if err != nil {
		return err
	}
	if !isEmailFormatCorrect(req.Email) {
		return errs.New(errs.BadRequest, "Email is not valid", nil)
	}
	if _, err := s.UserRepo.GetByEmail(req.Email); err == nil {
		return errs.New(errs.Conflict, "Email is already in use", nil)
	}
	user.Email = req.Email
	user.Verified = false
	return s.UserRepo.Update(user, nil)
}

func (s *AccountService) UpdateUsername(claims *models.ManagedClaims, req models.UpdateUsernameRequest) error {
	if err := ValidateUsername(string(req.Username)); err != nil {
		return err
	}
	if b, err := s.IsUsernameAvailable(req); err != nil {
		return err
	} else if !b {
		return errs.New(errs.Conflict, "Username is already in use", nil)
	}

	user, err := s.UserRepo.GetByUserID(claims.UserID)
	if err != nil {
		return err
	}
	user.Username = string(req.Username)
	return s.UserRepo.Update(user, nil)
}

func (s *AccountService) IsUsernameAvailable(req models.UpdateUsernameRequest) (bool, error) {
	available, err := s.UserRepo.IsUsernameAvailable(string(req.Username))
	if err != nil {
		return false, err
	}
	return available, nil
}

// ====== HELPERS ======

func isEmailFormatCorrect(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func isPasswordStrongEnough(password string) map[string]string {
	msgs := make(map[string]string)
	if len(password) < 8 {
		msgs["length"] = "Password must be at least 8 characters long"
	}
	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case 'a' <= char && char <= 'z':
			hasLower = true
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case strings.ContainsAny(string(char), "!@#$%^&*()-_=+[]{}|;:,.<>?"):
			hasSpecial = true
		}
	}
	if !hasUpper {
		msgs["upper"] = "Password must contain at least one uppercase letter"
	}
	if !hasLower {
		msgs["lower"] = "Password must contain at least one lowercase letter"
	}
	if !hasNumber {
		msgs["number"] = "Password must contain at least one number"
	}
	if !hasSpecial {
		msgs["special"] = "Password must contain at least one special character"
	}
	if len(msgs) == 0 {
		return nil
	}
	return msgs
}

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9._]+$`)

func ValidateUsername(username string) error {
	// Length check
	if len(username) < 3 || len(username) > 20 {
		return errs.New(errs.BadRequest, "username must be between 3 and 20 characters", nil)
	}

	// Allowed characters
	if !usernameRegex.MatchString(username) {
		return errs.New(errs.BadRequest, "username can only contain letters, numbers, dots, and underscores", nil)
	}

	// Cannot start or end with dot or underscore
	if username[0] == '.' || username[0] == '_' ||
		username[len(username)-1] == '.' || username[len(username)-1] == '_' {
		return errs.New(errs.BadRequest, "username cannot start or end with '.' or '_'", nil)
	}

	// No consecutive dots or underscores
	for i := 0; i < len(username)-1; i++ {
		if (username[i] == '.' && username[i+1] == '.') ||
			(username[i] == '_' && username[i+1] == '_') {
			return errs.New(errs.BadRequest, "username cannot contain consecutive dots or underscores", nil)
		}
	}

	return nil
}
