package services

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"sea-api/internal/config"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/repositories"
	"sea-api/internal/utils"
	"strings"
	"time"
)

const VERIFICATION_PATH = `https://sea.uofk.edu/cert/verify/`

type CertificateService struct {
	userRepo            repositories.IUserRepository
	eventService        *EventService
	S3StoreService      IS3StorageService
	pdfService          IPDFService
	mailService         IMailService
	CollaboratorService *CollaboratorService

	certificateRepository repositories.ICertificateRepository

	storePath string
}

func NewCertificateService(userRepo repositories.IUserRepository, eventService *EventService, S3StoreService IS3StorageService, pdfService IPDFService, mailService IMailService, CollaboratorService *CollaboratorService, CertificateRepository repositories.ICertificateRepository) *CertificateService {
	return &CertificateService{
		userRepo:              userRepo,
		pdfService:            pdfService,
		eventService:          eventService,
		S3StoreService:        S3StoreService,
		mailService:           mailService,
		CollaboratorService:   CollaboratorService,
		certificateRepository: CertificateRepository,
		storePath:             "public/certificates",
	}
}

func (c *CertificateService) MakeCertificatesForEvent(ctx context.Context, eventId int64, progressChan chan string) error {
	defer close(progressChan)
	progressChan <- "started"

	participants, err := c.eventService.EventRepo.GetParticipantByEventID(eventId)
	if err != nil {
		slog.Error("error getting participants", "error", err, "event_id", eventId)
		return err
	}

	var ids []int64
	for _, p := range participants {
		if p.Completed {
			if p.Grade >= 40 || p.Grade == 0 {
				ids = append(ids, p.UserID)
			}
		}
	}

	users, err := c.userRepo.GetAllByIndices(ids)
	if err != nil {
		slog.Error("error getting users", "error", err, "event_id", eventId)
		return err
	}

	for i, user := range users {
		_, err := c.CreateWorkshopCertificate(ctx, user.ID, eventId)
		if err != nil {
			slog.Error("error creating certificate", "error", err, "user_id", user.ID, "event_id", eventId)
			utils.ParseProgressStruct(len(ids), i+1, user.ID, false, user.NameAr, progressChan)
			continue
		}
		utils.ParseProgressStruct(len(ids), i+1, user.ID, true, user.NameAr, progressChan)
	}
	progressChan <- "done"
	return nil
}

func (c *CertificateService) SendCertificatesEmailsForEvent(request models.CertificateSendEmailData, progressChan chan string) error {
	eventId := request.EventID
	defer close(progressChan)
	progressChan <- "started"

	participants, err := c.eventService.EventRepo.GetParticipantByEventID(eventId)
	if err != nil {
		slog.Error("error getting participants", "error", err, "event_id", eventId)
		return err
	}

	event, err := c.eventService.GetEventByID(eventId)
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

	users, err := c.userRepo.GetAllByIndices(ids)
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
			Username:  name[0] + " " + name[1],
			EventName: event.Name,
			EventType: string(event.EventType),
			CertURL:   VERIFICATION_PATH + certificate.Hash,
			Year:      time.Now().Year(),
		}
		temp, err := utils.GetTemplate(string(utils.EmailCertificateAr), data)
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

func (c *CertificateService) CreateWorkshopCertificate(ctx context.Context, userUserID, eventId int64) (int64, error) {
	cert, err := c.certificateRepository.GetByUserIDAndEventID(userUserID, eventId)
	if err == nil {
		return cert.ID, nil
	}

	participant, err := c.eventService.EventRepo.GetParticipantByEventAndUserIDs(eventId, userUserID)
	if err != nil {
		return 0, err
	}
	if !participant.Completed {
		return 0, errs.New(errs.NotFound, fmt.Sprintf("Participant %d did not complete event %d yet", userUserID, eventId), nil)
	}
	event, err := c.eventService.GetEventByID(eventId)
	if err != nil {
		return 0, err
	}
	user, err := c.userRepo.GetByUserID(userUserID)
	if err != nil {
		return 0, err
	}

	collab, err := c.CollaboratorService.repo.GetByID(event.PresenterID)
	if err != nil {
		return 0, err
	}

	signature := ""
	if collab.SignatureID.Valid {
		signatureImage, err := c.S3StoreService.Download(ctx, collab.SignatureID.Int64)
		if err != nil {
			return 0, err
		}
		signature = base64.StdEncoding.EncodeToString(signatureImage)
	}

	stringToHash := user.NameEn + "|" + event.Name + "|" + event.StartDate.Format("02-01-2006") + "|" + event.EndDate.Format("02-01-2006") + "|" + config.App.SecretSalt
	hash := sha256.Sum256([]byte(stringToHash))
	hashString := hex.EncodeToString(hash[:])
	url := VERIFICATION_PATH + hashString

	qr, err := utils.GenerateGearQR(url, 512, 512)
	if err != nil {
		return 0, err
	}

	pdfAr, err := c.getFile(
		user.NameAr,
		event.Name,
		base64.StdEncoding.EncodeToString(qr),
		collab.NameAr,
		signature,
		toArabicDate(event.StartDate, "02 January, 2006"),
		toArabicDate(event.EndDate, "02 January, 2006"),
		toArabicDate(time.Now(), "Monday الموافق January 02, 2006"),
		event.Outcomes,
		participant.Grade,
		string(utils.EventCertificateAr),
		utils.GetTemplate,
	)
	if err != nil {
		return 0, err
	}

	pdfEn, err := c.getFile(
		user.NameEn,
		event.Name,
		base64.StdEncoding.EncodeToString(qr),
		collab.NameEn,
		signature,
		event.StartDate.Format("January 01, 2006"),
		event.EndDate.Format("January 01, 2006"),
		time.Now().Format("Monday, Jan 02, 2006"),
		event.Outcomes,
		participant.Grade,
		string(utils.EventCertificateEn),
		utils.GetTemplate,
	)
	if err != nil {
		return 0, err
	}

	storeIdAr, err := c.S3StoreService.Upload(ctx, c.storePath+"/ar/"+hashString+".pdf", pdfAr, "application/pdf")
	if err != nil {
		return 0, err
	}
	storeIdEn, err := c.S3StoreService.Upload(ctx, c.storePath+"/en/"+hashString+".pdf", pdfEn, "application/pdf")
	if err != nil {
		c.S3StoreService.Delete(ctx, storeIdAr)
		slog.Error("error uploading file", "error", err, "s3 stored file", storeIdEn)
		return 0, err
	}

	id, err := c.certificateRepository.Create(models.CertificateModel{
		Hash:      hashString,
		UserID:    userUserID,
		EventID:   eventId,
		Grade:     participant.Grade,
		IssueDate: time.Now(),
		Status:    models.ACTIVE,
	})
	if err != nil {
		c.S3StoreService.Delete(ctx, storeIdAr)
		c.S3StoreService.Delete(ctx, storeIdEn)
		slog.Error("error creating certificate", "error", err, "stored file", storeIdAr, "stored file", storeIdEn)
		return 0, err
	}

	_, err = c.certificateRepository.CreateFile(models.CertificateFileModel{
		CertificateID: id,
		StoreID:       storeIdAr,
		Lang:          "ar",
	})
	if err != nil {
		c.S3StoreService.Delete(ctx, storeIdAr)
		c.S3StoreService.Delete(ctx, storeIdEn)
		return 0, err
	}

	_, err = c.certificateRepository.CreateFile(models.CertificateFileModel{
		CertificateID: id,
		StoreID:       storeIdEn,
		Lang:          "en",
	})
	if err != nil {
		c.S3StoreService.Delete(ctx, storeIdEn)
		return 0, err
	}

	return id, nil
}

func (c *CertificateService) VerifyCertificate(hash string) (*models.CertificateVerify, error) {
	cert, err := c.certificateRepository.GetByHash(hash)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return &models.CertificateVerify{
				Valid: false,
			}, nil
		}
		return nil, err
	}
	event, err := c.eventService.GetEventByID(cert.EventID)
	if err != nil {
		return nil, err
	}
	user, err := c.userRepo.GetByUserID(cert.UserID)
	if err != nil {
		return nil, err
	}
	valid := true
	status := cert.Status
	if status == models.REVOKED {
		valid = false
	}

	return &models.CertificateVerify{
		Valid:     valid,
		ID:        fmt.Sprint(cert.ID),
		NameAr:    user.NameAr,
		NameEn:    user.NameEn,
		EventName: event.Name,
		Status:    status,
		Grade:     fmt.Sprintf("%.2f", cert.Grade),
		Outcomes:  event.Outcomes,
		EndDate:   event.EndDate,
		IssueDate: cert.IssueDate,
	}, nil
}

func (c *CertificateService) GetCertificates(zw *zip.Writer, hash string) error {
	cert, err := c.certificateRepository.GetByHash(hash)
	if err != nil {
		return err
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

func (c *CertificateService) getFile(name, event, qr, collabName, signature, startDate, endDate, nowDate string, tasks []string, grade float64, filename string, f func(name string, data any) (string, error)) ([]byte, error) {
	data := models.DefaultCertificateData{
		Name:        name,
		EventName:   event,
		Grade:       grade,
		TaskColumns: make3x3Grid(tasks),
		QRCode:      fmt.Sprintf("data:image/png;base64,%s", qr),

		CollabName: collabName,
		Signature:  fmt.Sprintf("data:image/png;base64,%s", signature),

		StartDate: startDate,
		EndDate:   endDate,
		NowDate:   nowDate,
	}

	html, err := f(filename, data)
	if err != nil {
		return nil, err
	}

	pdf, err := c.pdfService.GeneratePDFFromHTML(context.Background(), html)
	if err != nil {
		return nil, err
	}
	return pdf, nil
}

// ======== HELPERS ========

func make3x3Grid(input []string) [][]string {
	limit := len(input)
	if limit == 0 || input[0] == "" {
		return nil
	}
	if limit > 9 {
		limit = 9
	}

	grid := [][]string{}

	for i := 0; i < limit; i += 3 {
		end := i + 3
		if end > limit {
			end = limit
		}
		grid = append(grid, input[i:end])
	}

	return grid
}

func toArabicDate(t time.Time, layout string) string {
	// Arabic translation maps
	months := map[string]string{
		"January": "يناير", "February": "فبراير", "March": "مارس",
		"April": "أبريل", "May": "مايو", "June": "يونيو",
		"July": "يوليو", "August": "أغسطس", "September": "سبتمبر",
		"October": "أكتوبر", "November": "نوفمبر", "December": "ديسمبر",
	}
	days := map[string]string{
		"Monday": "الاثنين", "Tuesday": "الثلاثاء", "Wednesday": "الأربعاء",
		"Thursday": "الخميس", "Friday": "الجمعة", "Saturday": "السبت", "Sunday": "الأحد",
	}

	numbers := map[string]string{
		"0": "٠", "1": "١", "2": "٢", "3": "٣", "4": "٤", "5": "٥", "6": "٦", "7": "٧", "8": "٨", "9": "٩",
	}

	// Get the English formatted string
	formatted := t.Format(layout)

	// Replace English names with Arabic
	for en, ar := range months {
		formatted = strings.ReplaceAll(formatted, en, ar)
	}
	for en, ar := range days {
		formatted = strings.ReplaceAll(formatted, en, ar)
	}
	for en, ar := range numbers {
		formatted = strings.ReplaceAll(formatted, en, ar)
	}

	return formatted
}
