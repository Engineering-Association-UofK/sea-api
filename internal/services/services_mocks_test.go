package services

import (
	"context"
	"sea-api/internal/models"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByUserID(id int64) (*models.UserModel, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserModel), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(username string) (*models.UserModel, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserModel), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*models.UserModel, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserModel), args.Error(1)
}

func (m *MockUserRepository) GetByUniID(uniID int64) (*models.UserModel, error) {
	args := m.Called(uniID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserModel), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.UserModel, tx *sqlx.Tx) error {
	args := m.Called(user, tx)
	return args.Error(0)
}

func (m *MockUserRepository) GetRolesByUserID(id int64) ([]models.UserRole, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.UserRole), args.Error(1)
}

func (m *MockUserRepository) GetTempUser(id int64) (*models.TempUserModel, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TempUserModel), args.Error(1)
}

func (m *MockUserRepository) Create(user *models.UserModel) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteTempUser(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) Verify(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockStorageService struct {
	mock.Mock
}

func (m *MockStorageService) Upload(ctx context.Context, key string, data []byte, contentType string) (int64, error) {
	args := m.Called(ctx, key, data, contentType)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockStorageService) Download(ctx context.Context, id int64) ([]byte, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockStorageService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockStorageService) GenerateDownloadUrlByID(ctx context.Context, id int64) (string, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}
	return args.Get(0).(string), args.Error(1)
}

func (m *MockStorageService) GenerateDownloadUrlByKey(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}
	return args.Get(0).(string), args.Error(1)
}

func (m *MockStorageService) MigrateFidToS3(ctx context.Context, s3Service *S3StorageService) error {
	args := m.Called(ctx, s3Service)
	return args.Error(0)
}

type MockMailService struct {
	mock.Mock
}

func (m *MockMailService) SendVerificationCode(to string, data models.VerifyEmail) error {
	args := m.Called(to, data)
	return args.Error(0)
}

type MockVerificationRepo struct {
	mock.Mock
}

func (m *MockVerificationRepo) GetByCode(code string) (*models.VerificationCodeModel, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.VerificationCodeModel), args.Error(1)
}

func (m *MockVerificationRepo) GetByUserID(userID int64) (*models.VerificationCodeModel, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.VerificationCodeModel), args.Error(1)
}

func (m *MockVerificationRepo) Create(code *models.VerificationCodeModel) error {
	args := m.Called(code)
	return args.Error(0)
}

func (m *MockVerificationRepo) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
