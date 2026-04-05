package handlers

import (
	"archive/zip"
	"io"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/services"

	"github.com/gin-gonic/gin"
)

type CertificateHandler struct {
	service *services.CertificateService
}

func NewCertificateHandler(service *services.CertificateService) *CertificateHandler {
	return &CertificateHandler{service: service}
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

func (h *CertificateHandler) VerifyDocument(ctx *gin.Context) {

	hash := ctx.Param("hash")
	doc, err := h.service.VerifyDocument(hash)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, doc)
}

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

	go h.service.MakeCertificatesForEvent(ctx.Request.Context(), req.EventID, progressChan)

	ctx.Stream(func(w io.Writer) bool {
		msg, ok := <-progressChan
		if !ok {
			return false
		}

		ctx.SSEvent("message", msg)
		return true
	})
}

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

func (h *CertificateHandler) GetCertificates(ctx *gin.Context) {
	hash := ctx.Param("hash")

	pr, pw := io.Pipe()
	go func() {
		zipWriter := zip.NewWriter(pw)

		err := h.service.GetCertificates(zipWriter, hash)
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
