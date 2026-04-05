package services

import (
	"context"
	"database/sql"
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
	"sea-api/internal/config"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/repositories"
	"strings"
)

type CollaboratorService struct {
	repo             *repositories.CollaboratorRepo
	s3StorageService *S3StorageService

	signaturePath string
}

func NewCollaboratorService(repo *repositories.CollaboratorRepo, s3StorageService *S3StorageService) *CollaboratorService {
	return &CollaboratorService{
		repo:             repo,
		s3StorageService: s3StorageService,
		signaturePath:    "internal/collaborators",
	}
}

func (s *CollaboratorService) GetAll(ctx context.Context) ([]models.CollaboratorResponse, error) {
	collaborators, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	if len(collaborators) == 0 {
		return []models.CollaboratorResponse{}, nil
	}

	var collaboratorResponses []models.CollaboratorResponse
	for _, collaborator := range collaborators {
		url := ""
		if collaborator.SignatureID.Valid {
			link, err := s.s3StorageService.GenerateDownloadUrlByID(ctx, collaborator.SignatureID.Int64)
			if err != nil {
				return nil, err
			}
			url = link
		}
		collaboratorResponses = append(collaboratorResponses, models.CollaboratorResponse{
			ID:           collaborator.ID,
			NameAr:       collaborator.NameAr,
			NameEn:       collaborator.NameEn,
			Email:        collaborator.Email,
			SignatureUrl: url,
		})
	}
	return collaboratorResponses, nil
}

func (s *CollaboratorService) GetByID(ctx context.Context, id int64) (*models.CollaboratorResponse, error) {
	collaborator, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	url := ""
	if collaborator.SignatureID.Valid {
		link, err := s.s3StorageService.GenerateDownloadUrlByID(ctx, collaborator.SignatureID.Int64)
		if err != nil {
			return nil, err
		}
		url = link
	}

	return &models.CollaboratorResponse{
		ID:           collaborator.ID,
		NameAr:       collaborator.NameAr,
		NameEn:       collaborator.NameEn,
		Email:        collaborator.Email,
		SignatureUrl: url,
	}, nil
}

func (s *CollaboratorService) Create(ctx context.Context, req *models.CollaboratorCreateRequest, file io.Reader) (int64, error) {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return 0, err
	}

	contentType := http.DetectContentType(fileBytes)
	if "image/png" != contentType {
		return 0, errs.New(errs.BadRequest, "unsupported file type:"+contentType, nil)
	}

	hash := fnv.New64()
	hash.Write([]byte(string(req.NameAr) + string(req.NameEn) + string(req.Email) + config.App.SecretSalt))
	fileKey := fmt.Sprintf("%s/%d-%s.%s", s.signaturePath, hash.Sum64(), strings.Split(req.NameEn, " ")[0], contentType[6:])

	id, err := s.s3StorageService.Upload(ctx, fileKey, fileBytes, contentType)
	if err != nil {
		return 0, err
	}

	return s.repo.Create(&models.CollaboratorModel{
		NameAr: req.NameAr,
		NameEn: req.NameEn,
		Email:  string(req.Email),
		SignatureID: sql.NullInt64{
			Int64: id,
			Valid: true,
		},
	})
}

func (s *CollaboratorService) Update(ctx context.Context, req *models.CollaboratorUpdateRequest, file io.Reader) error {
	collaborator, err := s.repo.GetByID(req.ID)
	if err != nil {
		return err
	}

	if file != nil {
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		contentType := http.DetectContentType(fileBytes)
		if contentType != "image/png" {
			return errs.New(errs.BadRequest, "unsupported file type:"+contentType, nil)
		}

		hash := fnv.New64()
		hash.Write([]byte(string(req.NameAr) + string(req.NameEn) + string(req.Email) + config.App.SecretSalt))
		fileKey := fmt.Sprintf("%s/%d-%s.%s", s.signaturePath, hash.Sum64(), strings.Split(req.NameEn, " ")[0], contentType[6:])

		newFileID, err := s.s3StorageService.Upload(ctx, fileKey, fileBytes, contentType)
		if err != nil {
			return err
		}

		if collaborator.SignatureID.Valid {
			s.s3StorageService.Delete(ctx, collaborator.SignatureID.Int64)
		}
		collaborator.SignatureID = sql.NullInt64{Int64: newFileID, Valid: true}
	}

	collaborator.NameAr = req.NameAr
	collaborator.NameEn = req.NameEn
	collaborator.Email = string(req.Email)

	return s.repo.Update(collaborator)
}

func (s *CollaboratorService) Delete(ctx context.Context, id int64) error {
	collaborator, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if collaborator.SignatureID.Valid {
		err = s.s3StorageService.Delete(ctx, collaborator.SignatureID.Int64)
		if err != nil {
			return err
		}
	}

	return s.repo.Delete(id)
}
