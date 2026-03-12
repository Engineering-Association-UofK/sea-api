package services

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"sea-api/internal/config"
	"sea-api/internal/models"
	"sea-api/internal/repositories"
	"sea-api/internal/utils"
	"strings"
	"time"
)

type CertificateService struct {
	userRepo     *repositories.UserRepository
	eventService *EventService
	storeService *StorageService
	pdfService   *PDFService
	mailService  *MailService

	certificateRepository *repositories.CertificateRepository
}

func NewCertificateService(userRepo *repositories.UserRepository, eventService *EventService, storeService *StorageService, pdfService *PDFService, mailService *MailService, CertificateRepository *repositories.CertificateRepository) *CertificateService {
	return &CertificateService{
		userRepo:              userRepo,
		pdfService:            pdfService,
		eventService:          eventService,
		storeService:          storeService,
		mailService:           mailService,
		certificateRepository: CertificateRepository,
	}
}

func (c *CertificateService) MakeCertificatesForEvent(eventId int64, progressChan chan string) error {
	defer close(progressChan)
	participants, err := c.eventService.EventRepo.GetParticipantByEventID(eventId)
	if err != nil {
		slog.Error("error getting participants", "error", err, "event_id", eventId)
		return err
	}
	slog.Info("participants", "participants", len(participants))
	type progress struct {
		Total   int    `json:"total"`
		Current int    `json:"current"`
		ID      int64  `json:"id"`
		Success bool   `json:"success"`
		Name    string `json:"name"`
	}

	var ids []int64
	for _, p := range participants {
		if p.Status == models.COMPLETED && p.Completed {
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
	slog.Info("users", "users", len(users))

	for i, user := range users {
		_, err := c.CreateWorkshopCertificate(user.Index, eventId)
		if err != nil {
			slog.Error("error creating certificate", "error", err, "user_id", user.Index, "event_id", eventId)
			s, err := parseToJsonString(progress{
				Total:   len(ids),
				Current: i + 1,
				ID:      user.Index,
				Success: false,
				Name:    user.NameAr,
			})
			if err != nil {
				slog.Error("Error parsing progress to JSON string", "error", err, "user_id", user.Index, "event_id", eventId)
			}
			progressChan <- s
			continue
		}
		s, err := parseToJsonString(progress{
			Total:   len(ids),
			Current: i + 1,
			ID:      user.Index,
			Success: true,
			Name:    user.NameAr,
		})
		if err != nil {
			slog.Error("Error parsing progress to JSON string", "error", err, "user_id", user.Index, "event_id", eventId)
		}
		progressChan <- s
	}
	progressChan <- "done"
	return nil
}

func (c *CertificateService) CreateWorkshopCertificate(userIndex, eventId int64) (int64, error) {
	cert, err := c.certificateRepository.GetByUserIDAndEventID(userIndex, eventId)
	if err == nil {
		return cert.ID, nil
	}

	participant, err := c.eventService.EventRepo.GetParticipantByEventAndUserIDs(eventId, userIndex)
	if err != nil {
		return 0, err
	}
	if participant.Status != models.COMPLETED || !participant.Completed {
		return 0, fmt.Errorf("Participant with ID %d did not complete the event yet", userIndex)
	}
	slog.Info("Participant Grade", "grade", participant.Grade)
	event, err := c.eventService.GetEventByID(eventId)
	if err != nil {
		return 0, err
	}
	user, err := c.userRepo.GetByIndex(userIndex)
	if err != nil {
		return 0, err
	}

	stringToHash := user.NameEn + "|" + event.Name + "|" + event.StartDate.Format("02-01-2006") + "|" + event.EndDate.Format("02-01-2006") + "|" + config.App.SecretSalt
	hash := sha256.Sum256([]byte(stringToHash))
	hashString := hex.EncodeToString(hash[:])
	url := "https://sea.uofk.edu/cert/verify/" + hashString

	qr, err := utils.GenerateGearQR(url, 512, 512)
	if err != nil {
		return 0, err
	}

	pdfAr, err := c.getFile(
		user.NameAr,
		event.Name,
		base64.StdEncoding.EncodeToString(qr),
		toArabicDate(event.StartDate, "02 January, 2006"),
		toArabicDate(event.EndDate, "02 January, 2006"),
		toArabicDate(time.Now(), "Monday الموافق January 02, 2006"),
		event.Outcomes,
		participant.Grade,
		utils.GetArabicCertificateTemplate,
	)
	if err != nil {
		return 0, err
	}

	pdfEn, err := c.getFile(
		user.NameEn,
		event.Name,
		base64.StdEncoding.EncodeToString(qr),
		event.StartDate.Format("January 01, 2006"),
		event.EndDate.Format("January 01, 2006"),
		time.Now().Format("Monday, Jan 02, 2006"),
		event.Outcomes,
		participant.Grade,
		utils.GetEnglishCertificateTemplate,
	)
	if err != nil {
		return 0, err
	}

	storeIdAr, err := c.storeService.UploadFileMaster("cert.pdf", pdfAr)
	if err != nil {
		return 0, err
	}
	storeIdEn, err := c.storeService.UploadFileMaster("cert.pdf", pdfEn)
	if err != nil {
		c.storeService.DeleteFile(storeIdAr)
		slog.Error("error uploading file", "error", err, "stored file", storeIdAr)
		return 0, err
	}

	id, err := c.certificateRepository.Create(models.CertificateModel{
		Hash:      hashString,
		UserID:    userIndex,
		EventID:   eventId,
		Grade:     participant.Grade,
		IssueDate: time.Now(),
		Status:    models.ACTIVE,
	})
	if err != nil {
		c.storeService.DeleteFile(storeIdAr)
		c.storeService.DeleteFile(storeIdEn)
		slog.Error("error creating certificate", "error", err, "stored file", storeIdAr, "stored file", storeIdEn)
		return 0, err
	}

	_, err = c.certificateRepository.CreateFile(models.CertificateFileModel{
		CertificateID: id,
		StoreID:       storeIdAr,
		Lang:          "ar",
	})
	if err != nil {
		c.storeService.DeleteFile(storeIdAr)
		c.storeService.DeleteFile(storeIdEn)
		return 0, err
	}

	c.certificateRepository.CreateFile(models.CertificateFileModel{
		CertificateID: id,
		StoreID:       storeIdEn,
		Lang:          "en",
	})
	if err != nil {
		c.storeService.DeleteFile(storeIdEn)
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
	user, err := c.userRepo.GetByIndex(cert.UserID)
	if err != nil {
		return nil, err
	}

	return &models.CertificateVerify{
		Valid:     true,
		ID:        cert.ID,
		NameAr:    user.NameAr,
		NameEn:    user.NameEn,
		EventName: event.Name,
		Status:    cert.Status,
		Grade:     cert.Grade,
		Outcomes:  event.Outcomes,
		EndDate:   event.EndDate,
		IssueDate: cert.IssueDate,
	}, nil
}

func (c *CertificateService) GetCertificates(zw *zip.Writer, id int64) error {
	cert, err := c.certificateRepository.GetByID(id)
	if err != nil {
		slog.Error("error getting certificate", "error", err, "certificate_id", id)
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
		data, err := c.storeService.DownloadFile(file.StoreID)
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

func (c *CertificateService) getFile(name, event, qr, startDate, endDate, nowDate string, tasks []string, grade float64, f func(data any) (string, error)) ([]byte, error) {
	data := models.DefaultCertificateData{
		Name:        name,
		EventName:   event,
		Grade:       grade,
		TaskColumns: make3x3Grid(tasks),
		QRCode:      fmt.Sprintf("data:image/png;base64,%s", qr),

		StartDate: startDate,
		EndDate:   endDate,
		NowDate:   nowDate,
	}

	html, err := f(data)
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

func parseToJsonString(data any) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
