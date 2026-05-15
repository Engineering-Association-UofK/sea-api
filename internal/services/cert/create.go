package cert

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"sea-api/internal/config"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/utils"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

func (c *CertificateService) SignPDF(ctx context.Context, req models.SignPdfRequest) ([]byte, error) {
	event, err := c.eventService.GetEventByID(req.EventID)
	if err != nil {
		return nil, err
	}

	metadataMap := make(map[string]string)
	if req.Metadata != "" {
		err = json.Unmarshal([]byte(req.Metadata), &metadataMap)
		if err != nil {
			return nil, errs.New(errs.BadRequest, "invalid json in metadata: "+err.Error(), nil)
		}
	}
	if len(metadataMap) == 0 {
		return nil, errs.New(errs.BadRequest, "no metadata provided", nil)
	}

	stringToHash := req.File.Filename + "|" + fmt.Sprint(req.File.Size) + "|" + event.Name + "|" + event.StartDate.Format("02-01-2006") + "|" + event.EndDate.Format("02-01-2006") + "|" + config.App.SecretSalt
	hash := sha256.Sum256([]byte(stringToHash))
	hashString := hex.EncodeToString(hash[:])
	url := DOC_VERIFICATION_PATH + hashString

	if _, err := c.documentRepository.GetByHash(hashString); err == nil {
		return nil, errs.New(errs.BadRequest, "Document already exists", nil)
	}

	qr, err := utils.GenerateGearQR(url, 512, 512)
	if err != nil {
		return nil, err
	}

	file, err := req.File.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var output bytes.Buffer
	rs := bytes.NewReader(data)

	QrS := req.QrS
	QrX := req.QrX
	QrY := -req.QrY

	desc := fmt.Sprintf("pos:tl, off: %.2f %.2f, scale:%.2f rel, op: 1.0, rot: 0.0", QrX, QrY, QrS/100)

	wm, err := api.ImageWatermarkForReader(
		bytes.NewReader(qr),
		desc,
		true,
		false,
		types.POINTS,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create watermark: %w", err)
	}

	err = api.AddWatermarks(rs, &output, []string{"1-"}, wm, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to stamp PDF: %w", err)
	}

	storeId, err := c.S3StoreService.Upload(ctx, c.storePath+"/direct/"+hashString+".pdf", output.Bytes(), "application/pdf")
	if err != nil {
		slog.Error("error uploading file", "error", err, "s3 stored file", storeId)
		return nil, err
	}

	tx, err := c.documentRepository.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	id, err := c.documentRepository.Create(&models.DocumentModel{
		DocHash:  hashString,
		FileID:   storeId,
		Type:     req.Type,
		CreateAt: time.Now(),
	}, tx)
	if err != nil {
		c.S3StoreService.Delete(ctx, storeId)
		slog.Error("error creating certificate", "error", err, "stored file", storeId)
		return nil, err
	}

	_, err = c.documentRepository.CreateRelation(&models.DocumentRelationModel{
		DocumentID:  id,
		Description: "Certificate of gratitude for event",
		ObjectType:  models.ObjEvent,
		ObjectID:    req.EventID,
	}, tx)
	if err != nil {
		c.S3StoreService.Delete(ctx, storeId)
		return nil, err
	}

	if len(req.Metadata) > 0 {
		for key, value := range metadataMap {
			_, err := c.documentRepository.CreateMetadata(&models.DocumentMetadataModel{
				DocumentID: id,
				Key:        key,
				Value:      value,
			}, tx)
			if err != nil {
				return nil, err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		c.S3StoreService.Delete(ctx, storeId)
		return nil, err
	}

	return output.Bytes(), nil
}

func (c *CertificateService) MakeCertificatesForEvent(ctx context.Context, req *models.MakeCertificatesForEventRequest, progressChan chan string) error {
	defer close(progressChan)
	progressChan <- "started"

	eventId := req.EventID
	certTemplateVersion := req.CertificateVersion

	slog.Debug("making certificates for event", "event_id", eventId)

	event, err := c.eventService.GetEventByID(eventId)
	if err != nil {
		slog.Debug("error getting event", "error", err, "event_id", eventId)
		return err
	}

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
	slog.Debug("got participants", "count", len(ids))
	if len(ids) == 0 {
		progressChan <- "done"
		return nil
	}

	users, err := c.userRepo.GetAllByIndices(ids)
	if err != nil {
		slog.Error("error getting users", "error", err, "event_id", eventId)
		return err
	}
	slog.Debug("got users", "count", len(users))

	for i, user := range users {
		hash, _, err := c.CreateWorkshopCertificate(ctx, user.ID, eventId, certTemplateVersion, req.CertificateType)
		if err != nil {
			slog.Error("error creating certificate", "error", err, "user_id", user.ID, "event_id", eventId)
			utils.ParseProgressStruct(len(ids), i+1, user.ID, false, user.NameAr, progressChan)
			continue
		}
		utils.ParseProgressStruct(len(ids), i+1, user.ID, true, user.NameAr, progressChan)
		c.NotificationService.CreateNotification(&models.NotificationRequest{
			UserID:  user.ID,
			Title:   "Your certificate is ready",
			Message: "Your certificate for the event " + event.Name + " is ready.",
			Type:    models.NotifyCertificate,
			Data: models.NotifyCertificateData{
				EventID:         eventId,
				CertificateHash: hash,
			},
		})
	}
	progressChan <- "done"
	return nil
}

func (c *CertificateService) CreateWorkshopCertificate(ctx context.Context, userUserID, eventId int64, version models.CertVersion, certType models.CertType) (string, int64, error) {
	cert, err := c.certificateRepository.GetByUserIDAndEventID(userUserID, eventId)
	if err == nil {
		slog.Debug("certificate already exists", "user_id", userUserID, "event_id", eventId)
		return cert.Hash, cert.ID, nil
	}
	participant, err := c.eventService.EventRepo.GetParticipantByEventAndUserIDs(eventId, userUserID)
	if err != nil {
		slog.Error("error getting participant", "error", err, "user_id", userUserID, "event_id", eventId)
		return "", 0, err
	}
	if !participant.Completed {
		slog.Debug("participant did not complete event", "user_id", userUserID, "event_id", eventId)
		return "", 0, errs.New(errs.NotFound, fmt.Sprintf("Participant %d did not complete event %d yet", userUserID, eventId), nil)
	}
	event, err := c.eventService.GetEventByID(eventId)
	if err != nil {
		slog.Error("error getting event", "error", err, "event_id", eventId)
		return "", 0, err
	}
	user, err := c.userRepo.GetByUserID(userUserID)
	if err != nil {
		slog.Error("error getting user", "error", err, "user_id", userUserID)
		return "", 0, err
	}

	stringToHash := user.NameEn + "|" + event.Name + "|" + event.StartDate.Format("02-01-2006") + "|" + event.EndDate.Format("02-01-2006") + "|" + config.App.SecretSalt
	hash := sha256.Sum256([]byte(stringToHash))
	hashString := hex.EncodeToString(hash[:])
	url := CERT_VERIFICATION_PATH + hashString
	slog.Debug("generating qr", "url", url)

	qr, err := utils.GenerateGearQR(url, 512, 512)
	if err != nil {
		slog.Error("error generating qr", "error", err)
		return "", 0, err
	}

	// Get the right template dynamically
	pdfAr, pdfEn, err := CertTypeMap[certType][version](c, ctx, event, participant, user, qr)
	if err != nil {
		slog.Error("error generating en pdf", "error", err)
		return "", 0, err
	}

	storeIdAr, err := c.S3StoreService.Upload(ctx, c.storePath+"/ar/"+hashString+".pdf", pdfAr, "application/pdf")
	if err != nil {
		slog.Error("error uploading ar file", "error", err, "s3 stored file", storeIdAr)
		return "", 0, err
	}
	storeIdEn, err := c.S3StoreService.Upload(ctx, c.storePath+"/en/"+hashString+".pdf", pdfEn, "application/pdf")
	if err != nil {
		c.S3StoreService.Delete(ctx, storeIdAr)
		slog.Error("error uploading en file", "error", err, "s3 stored file", storeIdEn)
		return "", 0, err
	}

	id, err := c.certificateRepository.Create(models.CertificateModel{
		Hash:        hashString,
		UserID:      userUserID,
		EventID:     eventId,
		Type:        certType,
		CertVersion: version,
		Grade:       participant.Grade,
		IssueDate:   time.Now(),
		Status:      models.CertActive,
	})
	if err != nil {
		c.S3StoreService.Delete(ctx, storeIdAr)
		c.S3StoreService.Delete(ctx, storeIdEn)
		slog.Error("error creating certificate", "error", err, "stored file", storeIdAr, "stored file", storeIdEn)
		return "", 0, err
	}

	_, err = c.certificateRepository.CreateFile(models.CertificateFileModel{
		CertificateID: id,
		StoreID:       storeIdAr,
		Lang:          "ar",
	})
	if err != nil {
		c.S3StoreService.Delete(ctx, storeIdAr)
		c.S3StoreService.Delete(ctx, storeIdEn)
		slog.Error("error creating ar certificate file", "error", err, "stored file", storeIdAr, "stored file", storeIdEn)
		return "", 0, err
	}

	_, err = c.certificateRepository.CreateFile(models.CertificateFileModel{
		CertificateID: id,
		StoreID:       storeIdEn,
		Lang:          "en",
	})
	if err != nil {
		slog.Error("error creating en certificate file", "error", err, "stored file", storeIdEn)
		c.S3StoreService.Delete(ctx, storeIdEn)
		return "", 0, err
	}

	return hashString, id, nil
}
