package handlers

import (
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

func (h *GalleryHandler) Upload(ctx *gin.Context) {
	var req models.NewGalleryAssetRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.BadRequest(ctx)
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		response.BadRequest(ctx)
		return
	}
	req.File = file

	value, exists := ctx.Get("user")
	claims, ok := value.(*models.ManagedClaims)
	if !exists || !ok {
		response.Unauthorized(ctx)
		return
	}

	id, err := h.GalleryService.UploadToGallery(ctx.Request.Context(), claims.UserID, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(201, "Asset uploaded successfully", id, ctx)
}

func (h *GalleryHandler) GetAll(ctx *gin.Context) {
	assets, err := h.GalleryService.GetAllAssets()
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(200, assets)
}

func (h *GalleryHandler) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx)
		return
	}

	asset, err := h.GalleryService.GetAssetByID(id)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(200, asset)
}

func (h *GalleryHandler) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx)
		return
	}

	err = h.GalleryService.DeleteAsset(ctx.Request.Context(), id)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Asset deleted successfully", id, ctx)
}
