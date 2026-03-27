package services

import (
	"context"
	"sea-api/internal/models"

	"github.com/jmoiron/sqlx"
)

type Auth interface {
	Login(req models.LoginRequest) (*models.LoginResponse, error)
	Register(req models.RegisterRequest) error
	Verify(req models.VerifyRequest) error
	SendVerificationCode(userID int64) error
}

type IUserRepository interface {
	GetByUserID(id int64) (*models.UserModel, error)
	GetByUsername(username string) (*models.UserModel, error)
	GetByEmail(email string) (*models.UserModel, error)
	GetByUniID(uniID int64) (*models.UserModel, error)
	Update(user *models.UserModel, tx *sqlx.Tx) error
	GetRolesByUserID(id int64) ([]models.UserRole, error)
	GetTempUser(id int64) (*models.TempUserModel, error)
	Create(user *models.UserModel) error
	DeleteTempUser(id int64) error
	Verify(id int64) error
}

type IStorageService interface {
	UploadFileFiler(path string, filename string, data []byte, contentType string) error
	DownloadFileFiler(path, fileName string) ([]byte, error)
}

type IS3StorageService interface {
	Upload(ctx context.Context, key string, data []byte, contentType string) (int64, error)
	Download(ctx context.Context, id int64) ([]byte, error)
	Delete(ctx context.Context, id int64) error
	GenerateDownloadUrlByID(ctx context.Context, id int64) (string, error)
	GenerateDownloadUrlByKey(ctx context.Context, key string) (string, error)
}

type IMailService interface {
	SendVerificationCode(to string, data models.VerifyEmail) error
}

type IVerificationRepo interface {
	GetByCode(code string) (*models.VerificationCodeModel, error)
	GetByUserID(userID int64) (*models.VerificationCodeModel, error)
	Create(code *models.VerificationCodeModel) error
	Delete(id int64) error
}
