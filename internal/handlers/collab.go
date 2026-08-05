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

// GetAll godocs
//
//	@Summary		Get all collaborators
//	@Description	Get a list of all collaborators/presenters
//	@Tags			Collaborators
//	@Produce		json
//	@Success		200	{object}		models.CollaboratorListResponse
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/collabs [get]
//
//	@Security		ApiKeyAuth
func (h *CollaboratorHandler) GetAll(ctx *gin.Context) {
	var req models.ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Invalid pagination parameters", nil))
		return
	}

	collabs, err := h.service.GetAll(req, ctx.Request.Context())
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.PureJSON(200, collabs)
}

// GetByID godocs
//
//	@Summary		Get collaborator by ID
//	@Description	Get details of a specific collaborator/presenter
//	@Tags			Collaborators
//	@Produce		json
//	@Param			id	path		int	true	"Collaborator ID"
//	@Success		200	{object}	models.CollaboratorResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		404	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/collabs/{id} [get]
//
//	@Security		ApiKeyAuth
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

// Create godocs
//
//	@Summary		Create collaborator
//	@Description	Create a new collaborator/presenter with a signature file
//	@Tags			Collaborators
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			name_ar	formData	string	true	"Name in Arabic"
//	@Param			name_en	formData	string	true	"Name in English"
//	@Param			email	formData	string	false	"Email address"
//	@Param			file	formData	file	true	"Signature image file"
//	@Success		201		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/collabs [post]
//
//	@Security		ApiKeyAuth
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

// Update godocs
//
//	@Summary		Update collaborator
//	@Description	Update an existing collaborator/presenter details and optionally their signature
//	@Tags			Collaborators
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id		formData	int		true	"Collaborator ID"
//	@Param			name_ar	formData	string	true	"Name in Arabic"
//	@Param			name_en	formData	string	true	"Name in English"
//	@Param			email	formData	string	false	"Email address"
//	@Param			file	formData	file	false	"New signature image file (optional)"
//	@Success		200		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		404		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/collabs [put]
//
//	@Security		ApiKeyAuth
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

// Delete godocs
//
//	@Summary		Delete collaborator
//	@Description	Delete a collaborator/presenter and their signature file
//	@Tags			Collaborators
//	@Produce		json
//	@Param			id	path		int	true	"Collaborator ID"
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		404	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/collabs/{id} [delete]
//
//	@Security		ApiKeyAuth

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
