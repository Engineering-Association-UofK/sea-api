package handlers

import (
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

// TestGeneration godocs
//
//	@Summary		Test Certificate Generation
//	@Description	Creates and image of a certificate without saving it to the database and sends it in the response
//	@Tags			cert
//	@Produce		json
//	@Param			body	body		certmodels.IssueRequest	true	"Template data"
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
