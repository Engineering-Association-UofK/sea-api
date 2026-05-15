package cert

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"
	"sea-api/internal/config"
	"sea-api/internal/models"
	"sea-api/internal/models/certs"
	"sea-api/internal/utils"
	"strings"
	"time"
)

var CertTypeMap = map[models.CertType]map[models.CertVersion]func(
	*CertificateService,
	context.Context,
	*models.EventDTO,
	*models.EventParticipantModel,
	*models.UserModel,
	[]byte,
) ([]byte, []byte, error){
	models.CertParticipation: certParticipationVersionMap,
}

var certParticipationVersionMap = map[models.CertVersion]func(
	*CertificateService,
	context.Context,
	*models.EventDTO,
	*models.EventParticipantModel,
	*models.UserModel,
	[]byte,
) ([]byte, []byte, error){
	models.V0_1: generateParticipationV0_1,
}

func generateParticipationV0_1(
	c *CertificateService,
	ctx context.Context,
	event *models.EventDTO,
	participant *models.EventParticipantModel,
	user *models.UserModel,
	qrCode []byte,
) ([]byte, []byte, error) {
	presenter, err := c.CollaboratorService.GetModelByID(event.PresenterID)
	if err != nil {
		slog.Error("error getting collaborator", "error", err, "collaborator_id", event.PresenterID)
		return nil, nil, err
	}

	coordinator, err := c.CollaboratorService.GetModelByID(event.CoordinatorID)
	if err != nil {
		slog.Error("error getting collaborator", "error", err, "collaborator_id", event.PresenterID)
		return nil, nil, err
	}

	preSignature := ""
	if presenter.SignatureID.Valid {
		slog.Debug("getting signature", "signature_id", presenter.SignatureID.Int64)
		signatureImage, err := c.S3StoreService.Download(ctx, presenter.SignatureID.Int64)
		if err != nil {
			slog.Error("error getting signature", "error", err, "signature_id", presenter.SignatureID.Int64)
			return nil, nil, err
		}
		preSignature = base64.StdEncoding.EncodeToString(signatureImage)
	}

	coordSignature := ""
	if coordinator.SignatureID.Valid {
		slog.Debug("getting signature", "signature_id", coordinator.SignatureID.Int64)
		signatureImage, err := c.S3StoreService.Download(ctx, coordinator.SignatureID.Int64)
		if err != nil {
			slog.Error("error getting signature", "error", err, "signature_id", coordinator.SignatureID.Int64)
			return nil, nil, err
		}
		coordSignature = base64.StdEncoding.EncodeToString(signatureImage)
	}

	// Get stamp png from resources folder and then convert it to base64
	stampPath := fmt.Sprintf("%s/stamp.png", config.App.ResourcesDir)
	stampBytes, err := os.ReadFile(stampPath)
	if err != nil {
		slog.Error("error reading stamp file", "error", err, "path", stampPath)
		return nil, nil, err
	}
	stampBase64 := fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(stampBytes))

	dataEN := certs.V1_0{
		Name:        user.NameEn,
		EventName:   event.Name,
		Grade:       participant.Grade,
		TaskColumns: make3x3Grid(event.Outcomes),
		QRCode:      fmt.Sprintf("data:image/png;base64,%s", qrCode),

		CoordinatorName:      coordinator.NameEn,
		CoordinatorTitle:     coordinator.TitleEn,
		CoordinatorSignature: fmt.Sprintf("data:image/png;base64,%s", coordSignature),

		PresenterName:      presenter.NameEn,
		PresenterTitle:     presenter.TitleEn,
		PresenterSignature: fmt.Sprintf("data:image/png;base64,%s", preSignature),

		Stamp: stampBase64,

		StartDate: event.StartDate.Format("January 02, 2006"),
		EndDate:   event.EndDate.Format("January 02, 2006"),
		NowDate:   time.Now().Format("Monday, Jan 02, 2006"),
	}

	dataAR := certs.V1_0{
		Name:      user.NameAr,
		EventName: event.Name,
		Grade:     participant.Grade, TaskColumns: make3x3Grid(event.Outcomes),
		QRCode: fmt.Sprintf("data:image/png;base64,%s", qrCode),

		CoordinatorName:      coordinator.NameAr,
		CoordinatorTitle:     coordinator.TitleAr,
		CoordinatorSignature: fmt.Sprintf("data:image/png;base64,%s", coordSignature),

		PresenterName:      presenter.NameAr,
		PresenterTitle:     presenter.TitleAr,
		PresenterSignature: fmt.Sprintf("data:image/png;base64,%s", preSignature),

		Stamp: stampBase64,

		StartDate: toArabicDate(event.StartDate, "02 January, 2006"),
		EndDate:   toArabicDate(event.EndDate, "02 January, 2006"),
		NowDate:   toArabicDate(time.Now(), "Monday الموافق January 02, 2006"),
	}

	pdfAR, err := c.getFile(models.CertParticipation, "v0.1", models.Arabic, dataAR)
	if err != nil {
		return nil, nil, err
	}

	pdfEN, err := c.getFile(models.CertParticipation, "v0.1", models.English, dataEN)
	if err != nil {
		return nil, nil, err
	}

	return pdfAR, pdfEN, nil
}

func (c *CertificateService) getFile(certType models.CertType, certVersion string, lang models.Language, data any) ([]byte, error) {
	html, err := utils.GetCertTemplate(certType, certVersion, lang, data)
	if err != nil {
		return nil, err
	}

	pdf, err := c.pdfService.GeneratePDFFromHTML(context.Background(), html)
	if err != nil {
		return nil, err
	}
	return pdf, nil
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
