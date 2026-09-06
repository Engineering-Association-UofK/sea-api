package storage

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log/slog"
	"sea-api/internal/config"
	"sea-api/internal/models"
	"sea-api/internal/repositories"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3 struct {
	FilesRepo     *repositories.FileRepository
	Client        *s3.Client
	PresignClient *s3.PresignClient
	Bucket        string
}

func NewS3Service(repo *repositories.FileRepository) *S3 {
	bucket := "sea-api"
	internalClient := s3.New(s3.Options{
		Region:       "us-east-1",
		BaseEndpoint: &config.App.StoreS3ApiUrl,
		UsePathStyle: true,
		Credentials:  credentials.NewStaticCredentialsProvider(config.App.S3AccessKey, config.App.S3SecretKey, ""),
	})
	externalClient := s3.New(s3.Options{
		Region:       "us-east-1",
		BaseEndpoint: &config.App.StoreUrl,
		UsePathStyle: true,
		Credentials:  credentials.NewStaticCredentialsProvider(config.App.S3AccessKey, config.App.S3SecretKey, ""),
	})
	s := &S3{
		FilesRepo:     repo,
		Client:        internalClient,
		PresignClient: s3.NewPresignClient(externalClient),
		Bucket:        bucket,
	}
	err := s.ConfigureTempLifecycle(context.Background())
	if err != nil {
		panic(err)
	}
	return s
}

func (s *S3) ConfigureTempLifecycle(ctx context.Context) error {
	_, err := s.Client.PutBucketLifecycleConfiguration(ctx, &s3.PutBucketLifecycleConfigurationInput{
		Bucket: &s.Bucket,
		LifecycleConfiguration: &types.BucketLifecycleConfiguration{
			Rules: []types.LifecycleRule{
				{
					ID:     aws.String("DeleteTempFilesAfter1Day"),
					Status: types.ExpirationStatusEnabled,
					Filter: &types.LifecycleRuleFilter{
						Prefix: &[]string{"temp/"}[0],
					},
					Expiration: &types.LifecycleExpiration{
						Days: aws.Int32(1),
					},
				},
			},
		},
	})
	return err
}

// generateRandomKey generates a unique path inside temp/
func generateRandomKey(extension string) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	if extension != "" && extension[0] != '.' {
		extension = "." + extension
	}
	return "temp/" + hex.EncodeToString(b) + extension, nil
}

// UploadTemp uploads a file directly to temp/ without DB persistence and returns its presigned URL
func (s *S3) UploadTemp(ctx context.Context, data []byte, contentType string, extension string) (string, error) {
	key, err := generateRandomKey(extension)
	if err != nil {
		return "", err
	}

	_, err = s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &s.Bucket,
		Key:         &key,
		Body:        bytes.NewReader(data),
		ContentType: &contentType,
	})
	if err != nil {
		return "", err
	}

	url, err := s.GenerateDownloadUrlByKey(ctx, key)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (s *S3) Upload(ctx context.Context, key string, data []byte, contentType string) (int64, error) {
	_, err := s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &s.Bucket,
		Key:         &key,
		Body:        bytes.NewReader(data),
		ContentType: &contentType,
	})
	if err != nil {
		return 0, err
	}

	id, err := s.FilesRepo.CreateFile(models.FileModel{
		Key:      key,
		FileSize: int64(len(data)),
		MimeType: contentType,
	})
	if err != nil {
		s.DeleteWithKey(context.Background(), key)
		return 0, err
	}

	return id, nil
}

func (s *S3) DownloadWithKey(ctx context.Context, key string) ([]byte, error) {
	result, err := s.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.Bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	return io.ReadAll(result.Body)
}

func (s *S3) Download(ctx context.Context, id int64) ([]byte, error) {
	file, err := s.FilesRepo.GetFileById(id)
	if err != nil {
		return nil, err
	}
	slog.Debug("Download file found", "file ID", file.ID, "File key", file.Key)

	return s.DownloadWithKey(ctx, file.Key)
}

func (s *S3) Delete(ctx context.Context, id int64) error {
	file, err := s.FilesRepo.GetFileById(id)
	if err != nil {
		return err
	}
	err = s.DeleteWithKey(ctx, file.Key)
	if err != nil {
		return err
	}

	return s.FilesRepo.DeleteFile(id)
}

func (s *S3) DeleteWithKey(ctx context.Context, key string) error {
	_, err := s.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &s.Bucket,
		Key:    &key,
	})
	return err
}

func (s *S3) GenerateDownloadUrlByID(ctx context.Context, id int64) (string, error) {
	file, err := s.FilesRepo.GetFileById(id)
	if err != nil {
		return "", err
	}

	slog.Debug("Link file found", "file ID", file.ID, "File key", file.Key)

	return s.GenerateDownloadUrlByKey(ctx, file.Key)
}

func (s *S3) GenerateDownloadUrlByKey(ctx context.Context, key string) (string, error) {
	request, err := s.PresignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.Bucket,
		Key:    &key,
	}, func(opts *s3.PresignOptions) {
		opts.Expires = 10 * time.Hour
	})

	if err != nil {
		return "", err
	}

	return request.URL, nil
}
