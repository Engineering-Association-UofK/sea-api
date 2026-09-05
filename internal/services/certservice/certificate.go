package certservice

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/models/certmodels"
	"sea-api/internal/repositories/certrepo"
	"sea-api/internal/services/event"
	"sea-api/internal/services/storage"
	"sea-api/internal/services/user"
	"sea-api/internal/utils/valid"
	"slices"
	"time"
)

type CertService struct {
	repo *certrepo.CertRepository
	s3   *storage.S3

	eventService *event.EventService
	userService  *user.UserService

	storePath string
}

func NewCertService(repo *certrepo.CertRepository, s3 *storage.S3, eventService *event.EventService, userService *user.UserService) *CertService {
	return &CertService{
		repo: repo,
		s3:   s3,

		eventService: eventService,
		userService:  userService,

		storePath: "certificate",
	}
}

func (s *CertService) GetCertificateList(req *models.ListRequest) (*certmodels.CertListResponse, error) {
	count := s.repo.GetCount()
	if count == 0 {
		return &certmodels.CertListResponse{
			ListResponse: models.ListResponse{
				TotalPages:  1,
				CurrentPage: 1,
				Count:       0,
			},
			List: []certmodels.CertResponse{},
		}, nil
	}

	totalPages := valid.Limit(req, count)

	certs, err := s.repo.GetCertsList(*req)
	if err != nil {
		return nil, err
	}

	return &certmodels.CertListResponse{
		ListResponse: models.ListResponse{
			TotalPages:  totalPages,
			CurrentPage: req.Page,
			Count:       count,
		},
		List: certs,
	}, nil
}

func (s *CertService) Update(req *certmodels.UpdateRequest) error {
	cert, err := s.repo.GetCertWithID(req.ID)
	if err != nil {
		return err
	}

	if req.RecipientEmail != cert.RecipientEmail {
		slog.Debug("Updating Email", "Email", req.RecipientEmail)
		cert.RecipientEmail = req.RecipientEmail
	}

	if (cert.EventID.Valid && req.EventID != cert.EventID.Int64) ||
		(!cert.EventID.Valid) {
		slog.Debug("Updating Event ID", "Event ID", req.EventID)
		if req.EventID != 0 {
			_, err = s.eventService.GetEventByID(req.EventID)
			if err != nil {
				return errs.New(errs.NotFound, "Event with ID not found", nil)
			}
		}
		cert.EventID.Valid = req.EventID != 0
		cert.EventID.Int64 = req.EventID
	}

	if (cert.RecipientUserID.Valid && req.RecipientUserID != cert.RecipientUserID.Int64) ||
		(!cert.RecipientUserID.Valid) {
		slog.Debug("Updating User ID", "User ID", req.RecipientUserID)
		if req.RecipientUserID != 0 {
			_, err = s.userService.GetUserDetails(req.RecipientUserID)
			if err != nil {
				return errs.New(errs.NotFound, "User with ID not found", nil)
			}
		}
		cert.RecipientUserID.Valid = req.RecipientUserID != 0
		cert.RecipientUserID.Int64 = req.RecipientUserID
	}

	return s.repo.Update(cert)
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

func (s *CertService) Issue(ctx context.Context, req *certmodels.IssueRequest, issuerID int64) (int64, error) {
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
		FileKey:    fileKey,
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

func (s *CertService) Verify(hash string) (*certmodels.VerifyResponse, error) {
	cert, err := s.repo.GetCertWithHash(hash)
	if err != nil {
		return nil, err
	}

	var eventID int64
	if cert.EventID.Valid {
		eventID = cert.EventID.Int64
	}

	return &certmodels.VerifyResponse{
		Valid:     true,
		ID:        cert.ID,
		Name:      cert.RecipientName,
		IssueDate: cert.IssuedDate,
		EventID:   eventID,
	}, nil
}

func (s *CertService) Download(zw *zip.Writer, ctx context.Context, id int64, claims *models.ManagedClaims) error {
	cert, err := s.repo.GetCertWithID(id)
	if err != nil {
		return err
	}

	if (cert.RecipientUserID.Valid && cert.RecipientUserID.Int64 == claims.UserID) ||
		(!slices.Contains(claims.Roles, models.RoleCertMgr) && !slices.Contains(claims.Roles, models.RoleSystemSuperAdmin)) {
		return errs.New(errs.Forbidden, "Forbidden", nil)
	}

	file, err := s.s3.DownloadWithKey(ctx, cert.FileKey)
	if err != nil {
		return err
	}

	w, err := zw.Create("certificate-" + cert.RecipientName + ".pdf")
	if err != nil {
		return err
	}

	_, err = io.Copy(w, bytes.NewReader(file))
	if err != nil {
		return err
	}

	return nil
}
