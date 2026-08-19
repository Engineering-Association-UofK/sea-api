package cert

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/repositories"
	"sea-api/internal/services"
	"sea-api/internal/services/event"
	"sea-api/internal/services/storage"
	"sea-api/internal/utils"
	"slices"
	"strings"
	"time"
)

const CERT_VERIFICATION_PATH = `https://sea.uofk.edu/cert/verify/`
const DOC_VERIFICATION_PATH = `https://sea.uofk.edu/doc/verify/`

type CertificateService struct {
	UserRepo            *repositories.UserRepository
	pdfService          *services.PDFService
	EventService        *event.EventService
	S3StoreService      *storage.S3
	mailService         *services.MailService
	NotificationService *services.NotificationService
	CollaboratorService *services.CollaboratorService

	certificateRepository *repositories.CertificateRepository
	documentRepository    *repositories.DocumentRepository

	storePath string
}

func NewCertificateService(
	userRepo *repositories.UserRepository,
	eventService *event.EventService,
	S3StoreService *storage.S3,
	pdfService *services.PDFService,
	mailService *services.MailService,
	CollaboratorService *services.CollaboratorService,
	NotificationService *services.NotificationService,
	CertificateRepository *repositories.CertificateRepository,
	DocumentRepository *repositories.DocumentRepository,
) *CertificateService {
	return &CertificateService{
		UserRepo:              userRepo,
		pdfService:            pdfService,
		EventService:          eventService,
		S3StoreService:        S3StoreService,
		mailService:           mailService,
		NotificationService:   NotificationService,
		CollaboratorService:   CollaboratorService,
		certificateRepository: CertificateRepository,
		documentRepository:    DocumentRepository,

		storePath: "public/certificates",
	}
}

func (c *CertificateService) SendCertificatesEmailsForEvent(request models.CertificateSendEmailData, progressChan chan string) error {
	eventId := request.EventID
	defer close(progressChan)
	progressChan <- "started"

	participants, err := c.EventService.EventRepo.GetParticipantByEventID(eventId)
	if err != nil {
		slog.Error("error getting participants", "error", err, "event_id", eventId)
		return err
	}

	event, err := c.EventService.GetEventByID(eventId)
	if err != nil {
		slog.Error("error getting event", "error", err, "event_id", eventId)
		return err
	}

	var ids []int64
	for _, p := range participants {
		if p.Status == models.COMPLETED && p.Completed {
			if p.Grade >= 40 || p.Grade == 0 {
				ids = append(ids, p.UserID)
			}
		}
	}

	certificates, err := c.certificateRepository.GetByEventIDAndUserIDs(eventId, ids)
	if err != nil {
		slog.Error("error getting certificates", "error", err, "event_id", eventId)
		return err
	}

	users, err := c.UserRepo.GetAllByIndices(ids)
	if err != nil {
		slog.Error("error getting users", "error", err, "event_id", eventId)
		return err
	}

	usersMap := utils.FromSlice(users, func(c models.UserModel) int64 { return c.ID })

	for i, certificate := range certificates {
		user, err := usersMap.Value(certificate.UserID)
		if err != nil {
			slog.Error("error getting user", "error", err, "event_id", eventId)
			utils.ParseProgressStruct(len(ids), i+1, user.ID, false, user.NameAr, progressChan)
			continue
		}
		name := strings.Split(user.NameAr, " ")
		data := models.CertificateEmailData{
			Name:      name[0] + " " + name[1],
			EventName: event.Name,
			EventType: string(event.EventType),
			Year:      time.Now().Year(),
		}
		temp, err := utils.GetEmailTemplate(models.EmailEventCertificate, models.Arabic, data)
		if err != nil {
			slog.Error("error reading template", "error", err, "event_id", eventId)
			utils.ParseProgressStruct(len(ids), i+1, user.ID, false, user.NameAr, progressChan)
			continue
		}
		err = c.mailService.SendEmail(models.Email{
			To:      []string{user.Email},
			Cc:      request.Cc,
			Bcc:     request.Bcc,
			Subject: "Certificate Completion",
			HTML:    temp,
		})
		if err != nil {
			slog.Error("error sending email", "error", err, "event_id", eventId)
			utils.ParseProgressStruct(len(ids), i+1, user.ID, false, user.NameAr, progressChan)
			continue
		}
		utils.ParseProgressStruct(len(ids), i+1, user.ID, true, user.NameAr, progressChan)
	}
	progressChan <- "done"
	return nil
}

func (c *CertificateService) GetCertificates(zw *zip.Writer, hash string, claims *models.ManagedClaims) error {
	cert, err := c.certificateRepository.GetByHash(hash)
	if err != nil {
		return err
	}

	// Check if the requester have the permissions needed to view this file
	if cert.UserID != claims.UserID ||
		!slices.Contains(claims.Roles, models.RoleCertMgr) ||
		!slices.Contains(claims.Roles, models.RoleSystemSuperAdmin) {
		return errs.New(errs.Forbidden, "Action not allowed", nil)
	}

	files, err := c.certificateRepository.GetFilesByCertificateID(cert.ID)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no files found")
	}

	for i, file := range files {
		data, err := c.S3StoreService.Download(context.Background(), file.StoreID)
		if err != nil {
			return err
		}
		w, err := zw.Create("certificate-" + files[i].Lang + ".pdf")
		if err != nil {
			return err
		}
		_, err = io.Copy(w, bytes.NewReader(data))
		if err != nil {
			return err
		}
	}
	return nil
}
