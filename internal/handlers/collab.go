package handlers

import (
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
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
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
	var req models.CollaboratorCreateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	if req.SignatureFile.Size > 2<<20 {
		errs.New(errs.BadRequest, "Signature size too big, should be less than 2MB", nil)
		return
	}

	id, err := h.service.Create(ctx.Request.Context(), &req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(201, "Collaborator created successfully", id, ctx)
}

func (h *CollaboratorHandler) Update(ctx *gin.Context) {
	var req models.CollaboratorUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	if req.SignatureFile.Size > 2<<20 {
		errs.New(errs.BadRequest, "Signature size too big, should be less than 2MB", nil)
		return
	}

	err := h.service.Update(ctx.Request.Context(), &req)
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
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err = h.service.Delete(ctx.Request.Context(), id)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Collaborator deleted successfully", id, ctx)
}
