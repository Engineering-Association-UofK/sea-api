package handlers

import (
	"archive/zip"
	"io"
	"sea-api/internal/models"
	"sea-api/internal/response"
	"sea-api/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CertificateHandler struct {
	service *services.CertificateService
}

func NewCertificateHandler(service *services.CertificateService) *CertificateHandler {
	return &CertificateHandler{service: service}
}

func (h *CertificateHandler) CreateWorkshopCertificate(ctx *gin.Context) {
	var req struct {
		UserID  int64 `json:"user_id" binding:"required"`
		EventID int64 `json:"event_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx)
		return
	}

	id, err := h.service.CreateWorkshopCertificate(req.UserID, req.EventID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, gin.H{
		"message":        "Certificate created successfully",
		"certificate_id": id,
	})
}

func (h *CertificateHandler) VerifyCertificate(ctx *gin.Context) {

	hash := ctx.Param("hash")
	cert, err := h.service.VerifyCertificate(hash)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, cert)
}

func (h *CertificateHandler) MakeCertificatesForEvent(ctx *gin.Context) {
	var req struct {
		EventID int64 `json:"event_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx)
		return
	}

	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")

	progressChan := make(chan string)

	go h.service.MakeCertificatesForEvent(req.EventID, progressChan)

	ctx.Stream(func(w io.Writer) bool {
		msg, ok := <-progressChan
		if !ok {
			return false
		}

		ctx.SSEvent("message", msg)
		return true
	})
}

func (h *CertificateHandler) SendCertificatesEmailsForEvent(ctx *gin.Context) {
	var req models.CertificateSendEmailData

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx)
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

func (h *CertificateHandler) GetCertificates(ctx *gin.Context) {
	id := ctx.Param("id")
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		response.BadRequest(ctx)
		return
	}

	pr, pw := io.Pipe()
	go func() {
		zipWriter := zip.NewWriter(pw)

		err = h.service.GetCertificates(zipWriter, intId)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		zipWriter.Close()
		pw.Close()
	}()

	ctx.Header("Content-Disposition", "attachment; filename=certificates.zip")
	ctx.Header("Content-Type", "application/zip")
	ctx.DataFromReader(200, -1, "application/zip", pr, nil)
}
