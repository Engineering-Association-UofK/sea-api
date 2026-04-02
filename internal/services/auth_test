package services

import (
	"sea-api/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockMailService := new(MockMailService)
	mockVerificationRepo := new(MockVerificationRepo)
	authService := NewAuthService(mockUserRepo, mockMailService, mockVerificationRepo)

	// Test case 1: Successful login with username
	req := models.LoginRequest{Username: "testuser", Password: "password123"}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &models.UserModel{
		ID:       1,
		Username: "testuser",
		Password: string(hashedPassword),
		Verified: true,
		Status:   models.STATUS_ACTIVE,
	}

	mockUserRepo.On("GetByUsername", req.Username).Return(user, nil).Once()
	mockUserRepo.On("GetRolesByUserID", user.ID).Return([]models.UserRole{{Role: "user"}}, nil).Once()

	res, err := authService.Login(req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.Token)
	assert.Equal(t, user.ID, res.UserID)
	assert.True(t, res.IsVerified)

	mockUserRepo.AssertExpectations(t)
}
