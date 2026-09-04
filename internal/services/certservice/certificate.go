package certservice

import (
	"context"
	"fmt"
	"sea-api/internal/models/certmodels"
	"sea-api/internal/repositories/certrepo"
	"sea-api/internal/services/storage"
	"time"
)

type CertService struct {
	repo *certrepo.CertRepository
	s3   *storage.S3

	storePath string
}

func NewCertService(repo *certrepo.CertRepository, s3 *storage.S3) *CertService {
	return &CertService{
		repo: repo,
		s3:   s3,

		storePath: "certificate",
	}
}

func (s *CertService) TestGenerateCert(ctx context.Context, req *certmodels.IssueRequest) (*certmodels.TestImageResponse, error) {
	url, err := s.GenerateTestImage(req, ctx)
	if err != nil {
		return nil, err
	}

	return &certmodels.TestImageResponse{
		Url: url,
	}, nil
}

func (s *CertService) IssueCertificate(ctx context.Context, req *certmodels.IssueRequest, issuerID int64) (int64, error) {
	cert, hash, err := s.GeneratePdf(req)
	if err != nil {
		return 0, fmt.Errorf("Failed to generate certificate PDF: %v", err)
	}

	fileKey := fmt.Sprintf("%s/%s.%s", s.storePath, hash, ".pdf")
	_, err = s.s3.Upload(ctx, fileKey, cert, "application/pdf")
	if err != nil {
		return 0, fmt.Errorf("Failed to upload certificate PDF: %v", err)
	}

	id, err := s.repo.CreateCert(&certmodels.Certificate{
		Hash:       hash,
		TemplateID: req.TemplateID,
		IssuerID:   issuerID,

		RecipientName:  req.RecipientName,
		RecipientEmail: req.RecipientEmail,
		IssuedDate:     time.Now(),
	})
	if err != nil {
		err2 := s.s3.DeleteWithKey(ctx, fileKey)
		if err2 != nil {
			return 0, fmt.Errorf("Failed to Create certificate row and delete uploaded PDF: %v --- Failed to delete uploaded certificate PDF: %v", err, err2)
		}
		return 0, fmt.Errorf("Failed to Create certificate row: %v", err)
	}

	return id, nil
}
