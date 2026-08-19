package handlers

import (
	"archive/zip"
	"fmt"
	"io"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	_ "sea-api/internal/response"
	"sea-api/internal/services/cert"

	"github.com/gin-gonic/gin"
)

type CertificateHandler struct {
	service *cert.CertificateService
}

func NewCertificateHandler(service *cert.CertificateService) *CertificateHandler {
	return &CertificateHandler{service: service}
}

// VerifyCertificate godocs
//
//	@Summary		Verify Certificate
//	@Description	Validate certificate
//	@Tags			Certificate
//	@Produce		json
//	@Param			hash	path	string	true	"Certificate hash"
//	@Success		200	{object}	models.CertificateVerify
//	@Failure		500	{object}	response.BaseError
//	@Router			/cert/verify/{hash} [get]
func (h *CertificateHandler) VerifyCertificate(ctx *gin.Context) {
	hash := ctx.Param("hash")
	cert, err := h.service.VerifyCertificate(hash)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, cert)
}

// VerifyDocument godocs
//
//	@Summary		Verify Document
//	@Description	Validate document
//	@Tags			Certificate
//	@Produce		json
//	@Param			hash	path	string	true	"Document hash"
//	@Success		200	{object}	models.DocumentVerifyResponse
//	@Failure		500	{object}	response.BaseError
//	@Router			/cert/verify-document/{hash} [get]
func (h *CertificateHandler) VerifyDocument(ctx *gin.Context) {

	hash := ctx.Param("hash")
	doc, err := h.service.VerifyDocument(hash)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, doc)
}

// MakeCertificatesForEvent godocs
//
//	@Summary		Make Certificates For Event
//	@Description	Generate certificates for all eligible participants of an event, the stream always starts with "started" and ends with "done" if no errors occurred.
//	@Tags			Certificate
//	@Produce		text/event-stream
//	@Param			body	body	models.MakeCertificatesForEventRequest	true	"Request body"
//	@Success		200		{string}	string	"SSE stream"
//	@Success		200		{object}	models.Progress
//	@Failure		200		{object}	models.ProgressError
//	@Failure		400		{object}	response.BaseError
//	@Router			/admin/event/generate-certs [get]
//
//	@Security		ApiKeyAuth
func (h *CertificateHandler) MakeCertificatesForEvent(ctx *gin.Context) {
	var req models.MakeCertificatesForEventRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")

	progressChan := make(chan string)

	go h.service.MakeCertificatesForEvent(ctx.Request.Context(), &req, progressChan)

	ctx.Stream(func(w io.Writer) bool {
		msg, ok := <-progressChan
		if !ok {
			return false
		}

		ctx.SSEvent("message", msg)
		return true
	})
}

// SignPDF godocs
//
//	@Summary		Sign PDF
//	@Description	Sign a PDF certificate
//	@Tags			Certificate
//	@Accept			multipart/form-data
//	@Produce		application/pdf
//	@Param			event_id	formData	integer	true	"Event ID"
//	@Param			type		formData	string	true	"Document Type"
//	@Param			metadata	formData	string	true	"Metadata JSON or string"
//	@Param			qr_x		formData	number	false	"QR X coordinate"
//	@Param			qr_y		formData	number	false	"QR Y coordinate"
//	@Param			qr_s		formData	number	true	"QR Scale/Size"
//	@Param			file		formData	file	true	"PDF file to sign"
//	@Success		200			{file}		binary
//	@Failure		400			{object}	response.BaseError
//	@Failure		500			{object}	response.BaseError
//	@Router			/admin/certificate/sign [post]
//
//	@Security		ApiKeyAuth
func (h *CertificateHandler) SignPDF(ctx *gin.Context) {
	var req models.SignPdfRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}
	signedPdf, err := h.service.SignPDF(ctx.Request.Context(), req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Header("Content-Disposition", "attachment; filename=signed_certificate.pdf")
	ctx.Header("Content-Type", "application/pdf")
	ctx.Data(200, "application/pdf", signedPdf)
}

// SendCertificatesEmailsForEvent godocs
//
//	@Summary		Send Certificates Emails For Event
//	@Description	Send certificate emails to all eligible participants of an event, the stream always starts with "started" and ends with "done" if no errors occurred.
//	@Tags			Certificate
//	@Produce		text/event-stream
//	@Param			body	body	models.CertificateSendEmailData	true	"Request body"
//	@Success		200		{string}	string	"SSE stream"
//	@Success		200		{object}	models.Progress
//	@Failure		200		{object}	models.ProgressError
//	@Failure		400		{object}	response.BaseError
//	@Router			/admin/event/send-all-emails [post]
//
//	@Security		ApiKeyAuth
func (h *CertificateHandler) SendCertificatesEmailsForEvent(ctx *gin.Context) {
	var req models.CertificateSendEmailData

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")

	progressChan := make(chan string)

	go h.service.SendCertificatesEmailsForEvent(req, progressChan)

	ctx.Stream(func(w io.Writer) bool {
		msg, ok := <-progressChan
		if !ok {
			return false
		}

		ctx.SSEvent("message", msg)
		return true
	})
}

// GenerateAndDownloadDebugCert godocs
//
//	@Summary		Generate and Download Debug Certificate Pack
//	@Description	Directly generates a certificate and streams back the zip containing both languages
//	@Tags			Certificate
//	@Produce		application/zip
//	@Param			user_id		query		int		true	"User ID"
//	@Param			event_id	query		int		true	"Event ID"
//	@Success		200			{file}		binary
//	@Failure		400			{object}	response.BaseError
//	@Failure		500			{object}	response.BaseError
//	@Router			/api/v1/cert/debug/generate [get]
func (h *CertificateHandler) GenerateAndDownloadDebugCert(ctx *gin.Context) {
	var query struct {
		UserID  int64 `form:"user_id" binding:"required"`
		EventID int64 `form:"event_id" binding:"required"`
	}

	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Missing or invalid user_id or event_id", nil))
		return
	}

	value, exists := ctx.Get("user")
	claims, ok := value.(*models.ManagedClaims)
	if !exists || !ok {
		ctx.Error(errs.New(errs.Unauthorized, "Unauthorized", nil))
		return
	}

	// 1. Fetch dependencies via User and Event repositories/services
	// Note: CertificateService already holds pointers to these internally
	userModel, err := h.service.UserRepo.GetByUserID(query.UserID)
	if err != nil {
		ctx.Error(errs.New(errs.NotFound, fmt.Sprintf("User not found: %v", err), nil))
		return
	}

	eventModel, err := h.service.EventService.GetEventByID(query.EventID)
	if err != nil {
		ctx.Error(errs.New(errs.NotFound, fmt.Sprintf("Event not found: %v", err), nil))
		return
	}

	participants, err := h.service.EventService.EventRepo.GetParticipantByEventID(query.EventID)
	if err != nil {
		ctx.Error(err)
		return
	}

	var targetParticipant *models.EventParticipantModel
	for _, p := range participants {
		if p.UserID == query.UserID {
			targetParticipant = &p
			break
		}
	}

	// Fallback mock participant for testing if user isn't assigned to event
	if targetParticipant == nil {
		targetParticipant = &models.EventParticipantModel{
			UserID:    query.UserID,
			EventID:   query.EventID,
			Grade:     100.0,
			Completed: true,
		}
	}

	// 2. Generate the Certificate pair (triggers PDF service via Docker/HTTP)
	hashString, _, err := h.service.CreateWorkshopCertificate(
		ctx.Request.Context(),
		userModel,
		targetParticipant,
		eventModel,
		models.V0_1,
		models.CertParticipation,
	)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 3. Stream the ZIP back using io.Pipe (matching your GetCertificates pattern)
	pr, pw := io.Pipe()
	go func() {
		zipWriter := zip.NewWriter(pw)

		err := h.service.GetCertificates(zipWriter, hashString, claims)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		zipWriter.Close()
		pw.Close()
	}()

	filename := fmt.Sprintf("debug-cert-event%d-user%d.zip", query.EventID, query.UserID)
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Header("Content-Type", "application/zip")
	ctx.DataFromReader(200, -1, "application/zip", pr, nil)
}
