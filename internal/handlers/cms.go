package handlers

import (
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/response"
	"sea-api/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CmsHandler struct {
	CmsService *services.CmsService
}

func NewCmsHandler(cmsService *services.CmsService) *CmsHandler {
	return &CmsHandler{
		CmsService: cmsService,
	}
}

func (h *CmsHandler) CreateBlogPost(ctx *gin.Context) {
	var req models.BlogPostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	value, exists := ctx.Get("user")
	claims, ok := value.(*models.ManagedClaims)
	if !exists || !ok {
		ctx.Error(errs.New(errs.Unauthorized, "Unauthorized", nil))
		return
	}

	id, err := h.CmsService.CreateBlogPost(claims.UserID, &req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(201, "Blog post created successfully", id, ctx)
}

func (h *CmsHandler) GetBlogPostById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	var id int64
	if idStr != "" {
		id, _ = strconv.ParseInt(idStr, 10, 64)
	}

	post, err := h.CmsService.GetBlogPostById(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.PureJSON(200, post)
}

func (h *CmsHandler) GetViewBlogPostBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	post, err := h.CmsService.GetViewBlogPostBySlug(slug)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.PureJSON(200, post)
}

func (h *CmsHandler) GetBlogPostsList(ctx *gin.Context) {
	var req models.ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	posts, err := h.CmsService.GetViewBlogPostList(req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.PureJSON(200, posts)
}

func (h *CmsHandler) GetAllBlogPosts(ctx *gin.Context) {
	posts, err := h.CmsService.GetAllBlogPosts(false)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.PureJSON(200, posts)
}

func (h *CmsHandler) UpdateBlogPost(ctx *gin.Context) {
	var req models.BlogPostUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err := h.CmsService.UpdateBlogPost(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Blog post updated successfully", req.ID, ctx)
}

func (h *CmsHandler) DeleteBlogPost(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err = h.CmsService.DeleteBlogPost(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Blog post deleted successfully", id, ctx)
}

func (h *CmsHandler) CreateTeamMember(ctx *gin.Context) {
	var req models.TeamMemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	id, err := h.CmsService.CreateTeamMember(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(201, "Team member created successfully", id, ctx)
}

func (h *CmsHandler) GetTeamMemberByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	member, err := h.CmsService.GetTeamMemberByID(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.PureJSON(200, member)
}

func (h *CmsHandler) GetAllTeamMembers(ctx *gin.Context) {
	activeOnlyStr := ctx.Query("active_only")
	activeOnly, _ := strconv.ParseBool(activeOnlyStr)

	members, err := h.CmsService.GetAllTeamMembers(activeOnly)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.PureJSON(200, members)
}

func (h *CmsHandler) GetViewTeamMembers(ctx *gin.Context) {
	members, err := h.CmsService.GetViewTeamMembers()
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.PureJSON(200, members)
}

func (h *CmsHandler) UpdateTeamMember(ctx *gin.Context) {
	var req models.TeamMemberUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err := h.CmsService.UpdateTeamMember(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Team member updated successfully", req.ID, ctx)
}

func (h *CmsHandler) DeleteTeamMember(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err = h.CmsService.DeleteTeamMember(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Team member deleted successfully", id, ctx)
}
