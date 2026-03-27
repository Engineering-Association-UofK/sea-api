package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"sea-api/internal/config"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/repositories"
	"time"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3StorageService struct {
	FilesRepo     *repositories.FileRepository
	Client        *s3.Client
	PresignClient *s3.PresignClient
	Bucket        string
}

func NewS3Service(repo *repositories.FileRepository) *S3StorageService {
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
	return &S3StorageService{
		FilesRepo:     repo,
		Client:        internalClient,
		PresignClient: s3.NewPresignClient(externalClient),
		Bucket:        "sea-api",
	}
}

func (s *S3StorageService) Upload(ctx context.Context, key string, data []byte, contentType string) (int64, error) {
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

func (s *S3StorageService) Download(ctx context.Context, id int64) ([]byte, error) {
	file, err := s.FilesRepo.GetFileById(id)
	if err != nil {
		return nil, err
	}

	result, err := s.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.Bucket,
		Key:    &file.Key,
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	return io.ReadAll(result.Body)
}

func (s *S3StorageService) Delete(ctx context.Context, id int64) error {
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

func (s *S3StorageService) DeleteWithKey(ctx context.Context, key string) error {
	_, err := s.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &s.Bucket,
		Key:    &key,
	})
	return err
}

func (s *S3StorageService) GenerateDownloadUrlByID(ctx context.Context, id int64) (string, error) {
	file, err := s.FilesRepo.GetFileById(id)
	if err != nil {
		return "", err
	}

	return s.GenerateDownloadUrlByKey(ctx, file.Key)
}

func (s *S3StorageService) GenerateDownloadUrlByKey(ctx context.Context, key string) (string, error) {
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

type StorageService struct {
	StoreRepo *repositories.StoreRepository

	MasterURL string
	PublicURL string
	FilerURL  string
}

func NewSeaweedService(StoreRepo *repositories.StoreRepository) *StorageService {

	return &StorageService{
		StoreRepo: StoreRepo,
		MasterURL: config.App.StoreMasterURL,
		PublicURL: config.App.StorePublicURL,
		FilerURL:  config.App.StoreFilerUrl,
	}
}

func (s *StorageService) UploadFileMaster(filename string, data []byte) (int64, error) {
	resp, err := http.Get(s.MasterURL + "/dir/assign")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var assign struct {
		Fid string `json:"fid"`
	}

	json.NewDecoder(resp.Body).Decode(&assign)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, _ := writer.CreateFormFile("file", filename)
	part.Write(data)
	writer.Close()

	uploadURL := fmt.Sprintf("%s/%s", s.PublicURL, assign.Fid)

	req, _ := http.NewRequest("POST", uploadURL, &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	var r struct {
		Size int64  `json:"size"`
		Mime string `json:"mime"`
	}
	json.NewDecoder(res.Body).Decode(&r)

	id, err := s.StoreRepo.Create(models.StoreModel{
		Fid:  assign.Fid,
		Size: r.Size,
		Mime: r.Mime,
	})
	if err != nil {
		return 0, err
	}

	return id, err
}

func (s *StorageService) UploadFileFiler(path string, filename string, data []byte, contentType string) error {
	uploadURL := fmt.Sprintf("%s/%s/%s", s.FilerURL, path, filename)

	req, err := http.NewRequest("PUT", uploadURL, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		return errs.New(errs.InternalServerError, "filer upload failed with status: "+res.Status, nil)
	}

	return nil
}

func (s *StorageService) DownloadFileFID(id int64) ([]byte, error) {
	store, err := s.StoreRepo.GetById(id)
	if err != nil {
		return nil, err
	}

	fid := store.Fid
	url := fmt.Sprintf("%s/%s", s.PublicURL, fid)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (s *StorageService) DownloadFileFiler(path, fileName string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/%s", s.FilerURL, path, fileName)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (s *StorageService) DeleteFile(id int64) error {
	store, err := s.StoreRepo.GetById(id)
	if err != nil {
		return err
	}

	fid := store.Fid
	url := fmt.Sprintf("%s/%s", s.PublicURL, fid)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	err = s.StoreRepo.DeleteStore(id)
	if err != nil {
		return err
	}

	return nil
}

// func (s *StorageService) MigrateFidToS3(ctx context.Context, s3Service *S3StorageService) error {
// 	// 1. Get all records that have a FID but no S3 Key yet
// 	records, err := s.StoreRepo.GetAllLegacyFidRecords()
// 	if err != nil {
// 		return err
// 	}

// 	for _, record := range records {
// 		// 2. Download from SeaweedFS Master using FID
// 		data, err := s.DownloadFileFID(record.Id)
// 		if err != nil {
// 			slog.Error("Failed to download", "fid", record.Fid, "error", err)
// 			continue
// 		}

// 		// 3. Generate a logical S3 Key
// 		// Suggestion: use their UniID or a UUID so it's unique
// 		newKey := fmt.Sprintf("migrated/%d_%s", record.Id, record.Fid)
// 		if record.Mime == "application/pdf" {
// 			newKey += ".pdf"
// 		}

// 		// 4. Upload via S3
// 		err = s3Service.Upload(ctx, newKey, data, record.Mime)
// 		if err != nil {
// 			slog.Error("Failed S3 upload", "key", newKey, "error", err)
// 			continue
// 		}

// 		// 5. Update DB: Save the Key and nullify/mark the Fid as migrated
// 		err = s.StoreRepo.UpdateToS3(record.Id, newKey)
// 		if err != nil {
// 			slog.Error("Failed DB update", "id", record.Id, "error", err)
// 			continue
// 		}

// 		slog.Info("Successfully migrated", "id", record.Id, "new_key", newKey)

// 		// 6. (Optional) Delete the old FID to save space
// 		// s.DeleteFile(record.Id)
// 	}

// 	return nil
// }
