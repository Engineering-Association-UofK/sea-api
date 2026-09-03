package certservice

import (
	"context"
	"sea-api/internal/config"
	"sea-api/internal/models/certmodels"
	"sea-api/internal/repositories/certrepo"
	"sea-api/internal/services/storage"
)

type CertService struct {
	repo *certrepo.CertRepository
	s3   *storage.S3

	arFontPath   string
	enFontPath   string
	templatePath string
}

func NewCertService(repo *certrepo.CertRepository, s3 *storage.S3) *CertService {
	return &CertService{
		repo: repo,
		s3:   s3,

		arFontPath:   config.App.ResourcesDir + "/fonts/ar/EmbeddedMohanad.ttf",
		enFontPath:   "", // TODO: add EN font
		templatePath: "internal/template",
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
