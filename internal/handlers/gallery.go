package handlers

import (
	"fmt"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/response"
	"sea-api/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GalleryHandler struct {
	GalleryService *services.GalleryService
}

func NewGalleryHandler(galleryService *services.GalleryService) *GalleryHandler {
	return &GalleryHandler{
		GalleryService: galleryService,
	}
}

// Upload godocs
//
//	@Summary		Upload gallery asset
//	@Description	Upload a new asset to the gallery
//	@Tags			Gallery
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file_name	formData	string	true	"Name of the file"
//	@Param			alt_text	formData	string	true	"Alternative text for the image"
//	@Param			file		formData	file	true	"The actual asset file"
//	@Success		201			{object}	response.TransactionResponse
//	@Failure		400			{object}	response.BaseError
//	@Failure		401			{object}	response.BaseError
//	@Failure		500			{object}	response.BaseError
//	@Router			/admin/gallery [post]
//
//	@Security		ApiKeyAuth
func (h *GalleryHandler) Upload(ctx *gin.Context) {
	var req models.NewGalleryAssetRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}
	req.File = file

	value, exists := ctx.Get("user")
	claims, ok := value.(*models.ManagedClaims)
	if !exists || !ok {
		ctx.Error(errs.New(errs.Unauthorized, "Unauthorized", nil))
		return
	}

	id, err := h.GalleryService.UploadToGallery(ctx.Request.Context(), claims.UserID, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(201, "Asset uploaded successfully", id, ctx)
}

// GetAll godocs
//
//	@Summary		Get all gallery assets
//	@Description	Get a list of all assets in the gallery
//	@Tags			Gallery
//	@Produce		json
//	@Param			limit	query		int	false	"Number of items per page"
//	@Param			page	query		int	false	"Page number"
//	@Success		200	{object}		models.GalleryListRequest
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/gallery [get]
//
//	@Security		ApiKeyAuth
func (h *GalleryHandler) GetAll(ctx *gin.Context) {
	var req models.ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	assets, err := h.GalleryService.GetAllAssets(&req)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.PureJSON(200, assets)
}

// GetByID godocs
//
//	@Summary		Get gallery asset by ID
//	@Description	Get a gallery asset by its ID
//	@Tags			Gallery
//	@Produce		json
//	@Param			id	path		int	true	"Asset ID"
//	@Success		200	{object}	models.GalleryAssetResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		404	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/gallery/{id} [get]
//
//	@Security		ApiKeyAuth
func (h *GalleryHandler) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	asset, err := h.GalleryService.GetAssetByID(id)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.PureJSON(200, asset)
}

// CleanGallery godocs
//
//	@Summary		Clean Gallery
//	@Description	Delete all assets from the gallery
//	@Tags			Gallery
//	@Produce		json
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/gallery [delete]
//
//	@Security		ApiKeyAuth
func (h *GalleryHandler) CleanGallery(ctx *gin.Context) {
	num, err := h.GalleryService.CleanGallery()
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, fmt.Sprintf("%d assets deleted", num), 0, ctx)
}
