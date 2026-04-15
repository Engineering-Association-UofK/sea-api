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

// CreatePost godocs
//
//	@Summary		Create blog post
//	@Description	Create a new blog post
//	@Tags			CMS:blog
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.PostRequest	true	"Post data"
//	@Success		201		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/blog [post]
//
//	@Security		ApiKeyAuth
func (h *CmsHandler) CreatePost(ctx *gin.Context) {
	var req models.PostRequest
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

	id, err := h.CmsService.CreatePost(claims.UserID, &req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(201, "Blog post created successfully", id, ctx)
}

// GetPostById godocs
//
//	@Summary		Get blog post by ID
//	@Description	Get a blog post by its ID
//	@Tags			CMS:blog
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		200	{object}	models.PostResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		404	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/blog/{id} [get]
//
//	@Security		ApiKeyAuth
func (h *CmsHandler) GetPostById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	var id int64
	if idStr != "" {
		id, _ = strconv.ParseInt(idStr, 10, 64)
	}

	post, err := h.CmsService.GetPostById(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.PureJSON(200, post)
}

// GetViewPostBySlug godocs
//
//	@Summary		Get blog post by slug
//	@Description	Get a blog post by its slug for public view
//	@Tags			CMS:blog
//	@Produce		json
//	@Param			slug	path		string	true	"Post slug"
//	@Success		200		{object}	models.PostResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		404		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/cms/blogs/{slug} [get]
func (h *CmsHandler) GetViewPostBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	post, err := h.CmsService.GetViewPostBySlug(slug)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.PureJSON(200, post)
}

// GetViewPostsList godocs
//
//	@Summary		Get blog posts list
//	@Description	Get a list of blog posts for public view
//	@Tags			CMS:blog
//	@Produce		json
//	@Param			limit	query		int	true	"Content count limit"
//	@Param			page	query		int	true	"Page number"
//	@Success		200		{object}	models.PostListViewResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/cms/blogs [get]
func (h *CmsHandler) GetViewPostsList(ctx *gin.Context) {
	var req models.ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request, need limit number", nil))
		return
	}

	posts, err := h.CmsService.GetViewPostList(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.PureJSON(200, posts)
}

// GetAllPosts godocs
//
//	@Summary		Get all blog posts
//	@Description	Get a list of all blog posts for administration
//	@Tags			CMS:blog
//	@Produce		json
//	@Param			limit	query		int	true	"Content count limit"
//	@Param			page	query		int	true	"Page number"
//	@Success		200		{object}	models.PostListResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/blog [get]
//
//	@Security		ApiKeyAuth
func (h *CmsHandler) GetAllPosts(ctx *gin.Context) {
	var req models.ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request, need limit number", nil))
		return
	}

	posts, err := h.CmsService.GetAllPosts(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.PureJSON(200, posts)
}

// UpdatePost godocs
//
//	@Summary		Update blog post
//	@Description	Update an existing blog post
//	@Tags			CMS:blog
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.PostUpdateRequest	true	"Post update data"
//	@Success		200		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		404		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/blog [put]
//
//	@Security		ApiKeyAuth
func (h *CmsHandler) UpdatePost(ctx *gin.Context) {
	var req models.PostUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err := h.CmsService.UpdatePost(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Blog post updated successfully", req.ID, ctx)
}

// DeletePost godocs
//
//	@Summary		Delete blog post
//	@Description	Delete a blog post by its ID
//	@Tags			CMS:blog
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		404	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/blog/{id} [delete]
//
//	@Security		ApiKeyAuth
func (h *CmsHandler) DeletePost(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err = h.CmsService.DeletePost(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Blog post deleted successfully", id, ctx)
}

// CreateTeamMember godocs
//
//	@Summary		Create team member
//	@Description	Create a new member of the Thirtieth Council
//	@Tags			CMS:team
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.TeamMemberRequest	true	"Team member data"
//	@Success		201		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/team [post]
//
//	@Security		ApiKeyAuth
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

// GetTeamMemberByID godocs
//
//	@Summary		Get team member by ID
//	@Description	Get a team member by their ID
//	@Tags			CMS:team
//	@Produce		json
//	@Param			id	path		int	true	"Member ID"
//	@Success		200	{object}	models.TeamMemberResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		404	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/team/{id} [get]
//
//	@Security		ApiKeyAuth
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

// GetAllTeamMembers godocs
//
//	@Summary		Get all team members
//	@Description	Get a list of all team members for administration
//	@Tags			CMS:team
//	@Produce		json
//	@Param			active_only	query		bool	false	"Filter by active status"
//	@Success		200			{array}		models.TeamMemberResponse
//	@Failure		401			{object}	response.BaseError
//	@Failure		500			{object}	response.BaseError
//	@Router			/admin/team [get]
//
//	@Security		ApiKeyAuth
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

// GetViewTeamMembers godocs
//
//	@Summary		Get team members for view
//	@Description	Get a list of active team members for public view
//	@Tags			CMS:team
//	@Produce		json
//	@Success		200	{array}		models.TeamMemberViewResponse
//	@Failure		500	{object}	response.BaseError
//	@Router			/cms/team [get]
func (h *CmsHandler) GetViewTeamMembers(ctx *gin.Context) {
	members, err := h.CmsService.GetViewTeamMembers()
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.PureJSON(200, members)
}

// UpdateTeamMember godocs
//
//	@Summary		Update team member
//	@Description	Update an existing team member
//	@Tags			CMS:team
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.TeamMemberUpdateRequest	true	"Team member update data"
//	@Success		200		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		404		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/team [put]
//
//	@Security		ApiKeyAuth
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

// DeleteTeamMember godocs
//
//	@Summary		Delete team member
//	@Description	Delete a team member by their ID
//	@Tags			CMS:team
//	@Produce		json
//	@Param			id	path		int	true	"Member ID"
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		404	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/team/{id} [delete]
//
//	@Security		ApiKeyAuth
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
