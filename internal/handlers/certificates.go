package handlers

import (
	"archive/zip"
	"io"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/models/certmodels"
	"sea-api/internal/response"
	"sea-api/internal/services/certservice"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CertificatesHandler struct {
	service *certservice.CertService
}

func NewCertificatesHandler(service *certservice.CertService) *CertificatesHandler {
	return &CertificatesHandler{service: service}
}

// Templates

// CreateTemplate godocs
//
//	@Summary		Create a certificate template
//	@Description	Create a template for the specific certificate version
//	@Tags			cert
//	@Produce		json
//	@Param			body	body		certmodels.CreateTemplateRequest	true	"Template data"
//	@Success		201		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/certificate/template [post]
func (h *CertificatesHandler) CreateTemplate(ctx *gin.Context) {
	var req certmodels.CreateTemplateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	id, err := h.service.CreateTemplate(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Created successfully", id, ctx)
}

// UpdateTemplate godocs
//
//	@Summary		Update a certificate template
//	@Description	Update a template for the specific certificate version
//	@Tags			cert
//	@Produce		json
//	@Param			body	body		certmodels.UpdateTemplateRequest	true	"Template data"
//	@Success		200		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/certificate/template [put]
func (h *CertificatesHandler) UpdateTemplate(ctx *gin.Context) {
	var req certmodels.UpdateTemplateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err := h.service.UpdateTemplate(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Updated successfully", req.ID, ctx)
}

// GetTemplate godocs
//
//	@Summary		Get a certificate template
//	@Description	Get a template for Viewing
//	@Tags			cert
//	@Produce		json
//	@Param			id	path		int	true	"Template ID"
//	@Success		200		{object}	certmodels.TemplateResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/certificate/template/{id} [get]
func (h *CertificatesHandler) GetTemplate(ctx *gin.Context) {
	idStr := ctx.Param("id")
	var id int64
	if idStr != "" {
		id, _ = strconv.ParseInt(idStr, 10, 64)
	}

	post, err := h.service.GetTemplate(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.PureJSON(200, post)
}

// GetViewPostsList godocs
//
//	@Summary		Get a list of template
//	@Description	Get a list of all available templates
//	@Tags			cert
//	@Produce		json
//	@Param			limit	query		int				false	"Content count limit"
//	@Param			page	query		int				false	"Page number"
//	@Success		200		{object}	certmodels.TemplateListResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/certificate/template [get]
func (h *CertificatesHandler) GetTemplatesList(ctx *gin.Context) {
	var req = models.ListRequest{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	posts, err := h.service.GetTemplateList(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.PureJSON(200, posts)
}

// Certificates

// TestGeneration godocs
//
//	@Summary		Test Certificate Generation
//	@Description	Creates and image of a certificate without saving it to the database and sends it in the response
//	@Tags			cert
//	@Accepts		mbfd
//	@Produce		json
//	@Param			template_id		formData		int			true	"Template ID to use"
//	@Param			recipient_name	formData		string		true	"Name written on the certificate"
//	@Param			recipient_email	formData		string		true	"Email of the recipient for automated sending"
//	@Param			signer_name_one	formData		string		true	"Signer #1 name"
//	@Param			signer_role_one	formData		string		true	"Signer #1 role"
//	@Param			signer_name_two	formData		string		true	"Signer #2 name"
//	@Param			signer_role_two	formData		string		true	"Signer #2 role"
//	@Param			signer_signature_one	formData	file	true	"Signature of Signer #1"
//	@Param			signer_signature_two	formData	file	true	"Signature of signer #2"
//	@Param			event_id				formData	int		false	"Event ID if exists"
//	@Param			recipient_user_id		formData	int		false	"User ID if linked to a user"
//	@Success		200		{object}	certmodels.TestImageResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/certificate/test [post]
func (h *CertificatesHandler) TestGeneration(ctx *gin.Context) {
	var req certmodels.IssueRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	img, err := h.service.TestGenerateCert(ctx.Request.Context(), &req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, img)
}

// IssueCertificate godocs
//
//	@Summary		Certificate Generation
//	@Description	Issue a certificate and save it to the database
//	@Tags			cert
//	@Accepts		mbfd
//	@Produce		json
//	@Param			template_id		formData		int			true	"Template ID to use"
//	@Param			recipient_name	formData		string		true	"Name written on the certificate"
//	@Param			recipient_email	formData		string		true	"Email of the recipient for automated sending"
//	@Param			signer_name_one	formData		string		true	"Signer #1 name"
//	@Param			signer_role_one	formData		string		true	"Signer #1 role"
//	@Param			signer_name_two	formData		string		true	"Signer #2 name"
//	@Param			signer_role_two	formData		string		true	"Signer #2 role"
//	@Param			signer_signature_one	formData	file	true	"Signature of Signer #1"
//	@Param			signer_signature_two	formData	file	true	"Signature of signer #2"
//	@Param			event_id				formData	int		false	"Event ID if exists"
//	@Param			recipient_user_id		formData	int		false	"User ID if linked to a user"
//	@Success		200		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/certificate [post]
func (h *CertificatesHandler) IssueCertificate(ctx *gin.Context) {
	var req certmodels.IssueRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	value, exists := ctx.Get("user")
	claims, ok := value.(*models.ManagedClaims)
	if !exists || !ok {
		ctx.Error(errs.New(errs.Unauthorized, "Unauthorized", nil))
		return
	}

	id, err := h.service.Issue(ctx.Request.Context(), &req, claims.UserID)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Issued successfully", id, ctx)
}

// UpdateCertificate godocs
//
//	@Summary		Update Certificate
//	@Description	Update specific fields in a certificate
//	@Tags			cert
//	@Produce		json
//	@Param			body	body		certmodels.UpdateRequest	true	"Update data"
//	@Success		200		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/certificate [put]
func (h *CertificatesHandler) UpdateCertificate(ctx *gin.Context) {
	var req certmodels.UpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err := h.service.Update(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Updated successfully", req.ID, ctx)
}

// GetCertificateList godocs
//
//	@Summary		Get Certificate List
//	@Description	Get a filtered list of all certificates
//	@Tags			cert
//	@Produce		json
//	@Param			body	body		models.ListRequest	true	"Limit data"
//	@Success		200		{object}	certmodels.CertListResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/certificate [get]
func (h *CertificatesHandler) GetCertificateList(ctx *gin.Context) {
	var req models.ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	list, err := h.service.GetCertificateList(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, list)
}

// VerifyCertificate godocs
//
//	@Summary		Verify Certificate
//	@Description	Verify if certificate is valid
//	@Tags			cert
//	@Produce		json
//	@Param			hash	path		string	true	"Certificate Hash"
//	@Success		200		{object}	certmodels.CertListResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/certificate/verify [get]
func (h *CertificatesHandler) VerifyCertificate(ctx *gin.Context) {
	hash := ctx.Param("hash")

	list, err := h.service.Verify(hash)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, list)
}

// DownloadCertificate godocs
//
//	@Summary		Download Certificate
//	@Description	Download certificate as a ZIP using it's ID
//	@Tags			cert
//	@Produce		json
//	@Param			id	path		int	true	"Certificate Hash"
//	@Success		200		{object}	certmodels.CertListResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/certificate/{id} [get]
func (h *CertificatesHandler) DownloadCertificate(ctx *gin.Context) {
	idStr := ctx.Param("id")
	var id int64
	if idStr != "" {
		id, _ = strconv.ParseInt(idStr, 10, 64)
	}

	value, exists := ctx.Get("user")
	claims, ok := value.(*models.ManagedClaims)
	if !exists || !ok {
		ctx.Error(errs.New(errs.Unauthorized, "Unauthorized", nil))
		return
	}

	pr, pw := io.Pipe()
	go func() {
		zipWriter := zip.NewWriter(pw)

		err := h.service.Download(zipWriter, ctx.Request.Context(), id, claims)
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
