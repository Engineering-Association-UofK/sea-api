package storage

import (
	"bytes"
	"context"
	"io"
	"sea-api/internal/config"
	"sea-api/internal/models"
	"sea-api/internal/repositories"
	"time"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
	return &S3{
		FilesRepo:     repo,
		Client:        internalClient,
		PresignClient: s3.NewPresignClient(externalClient),
		Bucket:        bucket,
	}
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

	return s.GenerateDownloadUrlByKey(ctx, file.Key)
}

func (s *S3) GenerateDownloadUrlByKey(ctx context.Context, key string) (string, error) {
	request, err := s.PresignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.Bucket,
		Key:    &key,
	}, func(opts *s3.PresignOptions) {
		opts.Expires = 15 * time.Minute
	})

	if err != nil {
		return "", err
	}

	return request.URL, nil
}
