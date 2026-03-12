package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"sea-api/internal/config"
	"sea-api/internal/models"
	"sea-api/internal/repositories"
	"time"
)

type StorageService struct {
	StoreRepo *repositories.StoreRepository

	MasterURL string
	PublicURL string
	FilerURL  string
	S3ApiURL  string
}

func NewSeaweedService(StoreRepo *repositories.StoreRepository) *StorageService {
	return &StorageService{
		StoreRepo: StoreRepo,
		MasterURL: config.App.StoreMasterURL,
		PublicURL: config.App.StorePublicURL,
		FilerURL:  config.App.StoreFilerUrl,
		S3ApiURL:  config.App.StoreS3ApiUrl,
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

func (s *StorageService) DownloadFile(id int64) ([]byte, error) {
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
