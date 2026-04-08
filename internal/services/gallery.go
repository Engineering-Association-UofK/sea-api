package services

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"sea-api/internal/models"
	"sea-api/internal/repositories"
	"time"
)

type GalleryService struct {
	Repo      *repositories.GalleryRepository
	S3Service *S3StorageService

	path string
}

func NewGalleryService(galleryRepo *repositories.GalleryRepository, s3Service *S3StorageService) *GalleryService {
	return &GalleryService{
		Repo:      galleryRepo,
		S3Service: s3Service,
		path:      "gallery",
	}
}

func (s *GalleryService) UploadToGallery(ctx context.Context, userID int64, req models.NewGalleryAssetRequest) (int64, error) {
	file, err := req.File.Open()
	if err != nil {
		return 0, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return 0, err
	}

	ext := filepath.Ext(req.File.Filename)
	key := fmt.Sprintf("%s/%d/%d_%s%s", s.path, time.Now().Year(), time.Now().Unix(), req.FileName, ext)
	contentType := req.File.Header.Get("Content-Type")

	fileID, err := s.S3Service.Upload(ctx, key, data, contentType)
	if err != nil {
		return 0, err
	}

	asset := &models.GalleryAssetModel{
		FileID:     fileID,
		FileName:   req.FileName,
		AltText:    req.AltText,
		UploadedBy: userID,
		Showcase:   false,
		CreatedAt:  time.Now(),
	}

	assetID, err := s.Repo.CreateAsset(asset)
	if err != nil {
		s.S3Service.Delete(ctx, fileID)
		return 0, err
	}
	return assetID, nil
}

func (s *GalleryService) AttachAssetToObject(assetID int64, objType models.ObjectType, objID int64) error {
	ref := &models.GalleryReferenceModel{
		AssetID:    assetID,
		ObjectType: objType,
		ObjectID:   objID,
	}
	_, err := s.Repo.CreateReference(ref)
	return err
}

func (s *GalleryService) GetAssetByID(id int64) (*models.GalleryAssetResponse, error) {
	asset, err := s.Repo.GetAssetByID(id)
	if err != nil {
		return nil, err
	}

	url, err := s.S3Service.GenerateDownloadUrlByID(context.Background(), asset.FileID)
	if err != nil {
		return nil, err
	}

	return &models.GalleryAssetResponse{
		ID:         asset.ID,
		URL:        url,
		FileName:   asset.FileName,
		AltText:    asset.AltText,
		UploadedBy: asset.UploadedBy,
		CreatedAt:  asset.CreatedAt,
	}, nil
}

func (s *GalleryService) GetAllAssets() ([]models.GalleryAssetResponse, error) {
	assets, err := s.Repo.GetAllGallery()
	if err != nil {
		return nil, err
	}

	var responses []models.GalleryAssetResponse
	for _, asset := range assets {
		url, _ := s.S3Service.GenerateDownloadUrlByID(context.Background(), asset.FileID)
		responses = append(responses, models.GalleryAssetResponse{
			ID:             asset.ID,
			URL:            url,
			ReferenceTimes: asset.ReferenceTimes,
			FileName:       asset.FileName,
			AltText:        asset.AltText,
			UploadedBy:     asset.UploadedBy,
			CreatedAt:      asset.CreatedAt,
		})
	}

	return responses, nil
}

func (s *GalleryService) CleanGallery() (int, error) {
	assets, err := s.Repo.GetUnreferencedAssetIDs()
	if err != nil {
		return 0, err
	}

	if len(assets) == 0 {
		return 0, nil
	}

	var idsToDelete []int64
	for _, asset := range assets {
		err := s.S3Service.Delete(context.Background(), asset.FileID)
		if err != nil {
			slog.Error("Failed to delete asset file", "Asset ID", asset.ID, "File ID", asset.FileID, "error", err)
		} else {
			idsToDelete = append(idsToDelete, asset.ID)
		}
	}

	for _, id := range idsToDelete {
		err = s.Repo.DeleteAsset(id)
		if err != nil {
			slog.Error("Failed to delete asset record", "Asset ID", id, "error", err)
		}
	}

	return len(idsToDelete), nil
}

func (s *GalleryService) GetLinkByAssetID(assetID int64) (string, error) {
	asset, err := s.Repo.GetAssetByID(assetID)
	if err != nil {
		slog.Info("Failed to get gallery asset", "id", assetID)
		return "", err
	}

	return s.S3Service.GenerateDownloadUrlByID(context.Background(), asset.FileID)
}

func (s *GalleryService) RemoveReference(objType models.ObjectType, objID int64) error {
	return s.Repo.DeleteReferencesByObject(objType, objID)
}
