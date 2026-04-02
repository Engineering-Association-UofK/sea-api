package handlers

import (
	"net/http"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/response"
	"sea-api/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CollaboratorHandler struct {
	service *services.CollaboratorService
}

func NewCollaboratorHandler(service *services.CollaboratorService) *CollaboratorHandler {
	return &CollaboratorHandler{service: service}
}

func (h *CollaboratorHandler) GetAll(ctx *gin.Context) {
	collabs, err := h.service.GetAll(ctx.Request.Context())
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.PureJSON(200, collabs)
}

func (h *CollaboratorHandler) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx)
		return
	}

	collab, err := h.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.PureJSON(200, collab)
}

func (h *CollaboratorHandler) Create(ctx *gin.Context) {
	// file size limit
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 2<<20)
	if err := ctx.Request.ParseMultipartForm(2 << 20); err != nil {
		errs.New(errs.BadRequest, "Signature size too big, should be less than 2MB", nil)
		return
	}

	var req models.CollaboratorCreateRequest
	req.NameAr = ctx.PostForm("name_ar")
	req.NameEn = ctx.PostForm("name_en")
	req.Email = models.TrimmedString(ctx.PostForm("email"))

	file, _, err := ctx.Request.FormFile("file")
	if err != nil || req.NameAr == "" || req.NameEn == "" {
		response.BadRequest(ctx)
		return
	}
	defer file.Close()

	id, err := h.service.Create(ctx.Request.Context(), &req, file)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(201, "Collaborator created successfully", id, ctx)
}

func (h *CollaboratorHandler) Update(ctx *gin.Context) {
	// file size limit
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 2<<20)
	if err := ctx.Request.ParseMultipartForm(2 << 20); err != nil {
		errs.New(errs.BadRequest, "Signature size too big, should be less than 2MB", nil)
		return
	}

	var req models.CollaboratorUpdateRequest
	idStr := ctx.PostForm("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx)
		return
	}

	req.ID = id
	req.NameAr = ctx.PostForm("name_ar")
	req.NameEn = ctx.PostForm("name_en")
	req.Email = models.TrimmedString(ctx.PostForm("email"))
	if id == 0 || req.NameAr == "" || req.NameEn == "" {
		response.BadRequest(ctx)
		return
	}

	file, _, _ := ctx.Request.FormFile("file")

	err = h.service.Update(ctx.Request.Context(), &req, file)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Collaborator updated successfully", req.ID, ctx)
}

func (h *CollaboratorHandler) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx)
		return
	}

	err = h.service.Delete(ctx.Request.Context(), id)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Collaborator deleted successfully", id, ctx)
}
