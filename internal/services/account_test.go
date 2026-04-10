package services

import (
	"bytes"
	"context"
	"sea-api/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestGetProfile(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockStorageService := new(MockStorageService)
	accountService := NewAccountService(mockUserRepo, mockStorageService)

	claims := &models.ManagedClaims{UserID: 1}
	user := &models.UserModel{
		ID:    1,
		UniID: 12345,
	}

	mockUserRepo.On("GetByUserID", claims.UserID).Return(user, nil)

	profile, err := accountService.GetProfile(context.Background(), claims)

	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, user.ID, profile.ID)
	assert.Equal(t, user.UniID, profile.UniID)

	mockUserRepo.AssertExpectations(t)
}

func TestUpdateProfile(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	accountService := NewAccountService(mockUserRepo, nil)
	claims := &models.ManagedClaims{UserID: 1}

	req := models.UpdateProfileRequest{
		ID:         1,
		UniID:      54321,
		NameAr:     "اسم جديد",
		NameEn:     "New Name",
		Gender:     "Male",
		Department: "IT",
		Phone:      "123456789",
	}

	user := &models.UserModel{ID: 1}

	mockUserRepo.On("GetByUserID", req.ID).Return(user, nil)
	mockUserRepo.On("GetByUniID", req.UniID).Return(nil, nil)
	mockUserRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	err := accountService.UpdateProfile(claims, req)
	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)

}

func TestUpdateProfilePicture(t *testing.T) {
	mockStorageService := new(MockStorageService)
	accountService := NewAccountService(nil, mockStorageService)
	claims := &models.ManagedClaims{UserID: 1}

	// 1x1 transparent PNG

	file := bytes.NewReader([]byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4, 0x89, 0x00, 0x00, 0x00, 0x0a, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9c, 0x63, 0x00, 0x01, 0x00, 0x00, 0x05, 0x00, 0x01, 0x0d, 0x0a, 0x2d, 0xb4, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82})

	mockStorageService.On("UploadFileFiler", mock.Anything, mock.Anything, mock.Anything, "image/png").Return(nil)
	err := accountService.UpdateProfilePicture(context.Background(), claims, file)
	assert.NoError(t, err)
	mockStorageService.AssertExpectations(t)
}

func TestUpdatePassword(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	accountService := NewAccountService(mockUserRepo, nil)

	claims := &models.ManagedClaims{UserID: 1}
	req := models.UpdatePasswordRequest{
		OldPassword:     "OldPassword1!",
		NewPassword:     "NewPassword1!",
		ConfirmPassword: "NewPassword1!",
	}

	hashedOldPassword, _ := bcrypt.GenerateFromPassword([]byte("OldPassword1!"), bcrypt.DefaultCost)
	user := &models.UserModel{ID: 1, Password: string(hashedOldPassword)}

	mockUserRepo.On("GetByUserID", claims.UserID).Return(user, nil)
	mockUserRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	err := accountService.UpdatePassword(claims, req)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

func TestUpdateEmail(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	accountService := NewAccountService(mockUserRepo, nil)

	claims := &models.ManagedClaims{UserID: 1}
	req := models.UpdateEmailRequest{Email: "new@example.com"}
	user := &models.UserModel{ID: 1}

	mockUserRepo.On("GetByUserID", claims.UserID).Return(user, nil)
	mockUserRepo.On("GetByEmail", req.Email).Return(nil, assert.AnError)
	mockUserRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	err := accountService.UpdateEmail(claims, req)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

func TestUpdateUsername(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	accountService := NewAccountService(mockUserRepo, nil)

	claims := &models.ManagedClaims{UserID: 1}
	req := models.UpdateUsernameRequest{Username: "new_username"}
	user := &models.UserModel{ID: 1}

	mockUserRepo.On("GetByUsername", req.Username).Return(nil, assert.AnError)
	mockUserRepo.On("GetByUserID", claims.UserID).Return(user, nil)
	mockUserRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	err := accountService.UpdateUsername(claims, req)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

func TestIsUsernameAvailable(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	accountService := NewAccountService(mockUserRepo, nil)

	req := models.UpdateUsernameRequest{Username: "new_username"}

	// Case 1: Username is available
	mockUserRepo.On("GetByUsername", req.Username).Return(nil, assert.AnError).Once()
	available, err := accountService.IsUsernameAvailable(req)
	assert.NoError(t, err)
	assert.True(t, available)

	// Case 2: Username is not available
	mockUserRepo.On("GetByUsername", req.Username).Return(&models.UserModel{}, nil).Once()
	available, err = accountService.IsUsernameAvailable(req)
	assert.NoError(t, err)
	assert.False(t, available)

	mockUserRepo.AssertExpectations(t)
}

func TestIsPasswordStrongEnough(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{"Valid password", "Password1!", true},
		{"Invalid password", "password", false},
		{"Too short password", "short", false},
		{"No uppercase letter", "password1!", false},
		{"No lowercase letter", "PASSWORD1!", false},
		{"No number", "Password!", false},
		{"No special character", "Password1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPasswordStrongEnough(tt.password)
			if result != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}

func TestIsEmailFormatCorrect(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{"Valid email", "test@example.com", true},
		{"Missing @", "testexample.com", false},
		{"Missing .", "test@examplecom", false},
		{"Empty email", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isEmailFormatCorrect(tt.email)
			if result != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name        string
		username    string
		expectedErr bool
	}{
		{"Valid username", "valid_user.123", false},
		{"Too short", "us", true},
		{"Too long", "a_very_long_username_that_is_not_valid", true},
		{"Invalid characters", "user name", true},
		{"Starts with underscore", "_user", true},
		{"Starts with dot", ".user", true},
		{"Ends with underscore", "user_", true},
		{"Ends with dot", "user.", true},
		{"Consecutive dots", "user..name", true},
		{"Consecutive underscores", "user__name", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsername(tt.username)
			if (err != nil) != tt.expectedErr {
				t.Errorf("Expected error: %v, but got: %v", tt.expectedErr, err)
			}
		})
	}
}
