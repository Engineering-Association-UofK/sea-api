package handlers

import (
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/response"
	"sea-api/internal/services/forms"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FormHandler struct {
	service *forms.FormService
}

func NewFormHandler(service *forms.FormService) *FormHandler {
	return &FormHandler{service: service}
}

// ======== ANALYSIS ========

// GetFormAnalysis godocs
//
//	@Summary		Get form analysis
//	@Description	Get statistical analysis of form responses
//	@Tags			Forms
//	@Produce		json
//	@Param			id	path		int	true	"Form ID"
//	@Success		200	{object}	models.FormAnalysisResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		404	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/form/analysis/{id} [get]				// <------- add admin endpoint
//
//	@Security		ApiKeyAuth
func (h *FormHandler) GetFormAnalysis(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	analysis, err := h.service.GetFormAnalysis(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, analysis)
}

// GetFormDetailedResponses godocs
//
//	@Summary		Get detailed form responses
//	@Description	Get all detailed responses and structure for a specific form
//	@Tags			Forms
//	@Produce		json
//	@Param			id	path		int	true	"Form ID"
//	@Success		200	{object}	models.FormDerailedResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		404	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/form/detailed-responses/{id} [get]			// <------- add admin endpoint
//
//	@Security		ApiKeyAuth
func (h *FormHandler) GetFormDetailedResponses(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	details, err := h.service.GetEntireFormDetails(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, details)
}

// ======== SPECIAL ========

// GetEntireForEditForm godocs
//
//	@Summary		Get form for editing
//	@Description	Get the complete structure of a form for administrative editing
//	@Tags			Forms
//	@Produce		json
//	@Param			id	path		int	true	"Form ID"
//	@Success		200	{object}	models.FormForEditDTO
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		404	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/form/{id} [get]			// <------- add admin endpoint
//
//	@Security		ApiKeyAuth
func (h *FormHandler) GetEntireForEditForm(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	form, err := h.service.GetFormForEdit(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, form)
}

// GetEntireForUserForm godocs
//
//	@Summary		Get form for user
//	@Description	Get the complete structure of a form for public user submission
//	@Tags			Forms
//	@Produce		json
//	@Param			id	path		int	true	"Form ID"
//	@Success		200	{object}	models.FormForUserDTO
//	@Failure		400	{object}	response.BaseError
//	@Failure		404	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/event/form/{id} [get]
func (h *FormHandler) GetEntireForUserForm(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	form, err := h.service.GetFormForUser(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, form)
}

// SubmitForm godocs
//
//	@Summary		Submit form
//	@Description	Submit a response to a form
//	@Tags			Forms
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.SubmitFormRequest	true	"Form submission data"
//	@Success		201		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/form/submit [post]			// <------- add event endpoint ----------
func (h *FormHandler) SubmitForm(ctx *gin.Context) {
	var req models.SubmitFormRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	// value, exists := ctx.Get("user")
	// claims, ok := value.(*models.ManagedClaims)
	// if !exists || !ok {
	// 	ctx.Error(errs.New(errs.Unauthorized, "Unauthorized", nil))
	// 	return
	// }

	id, err := h.service.SubmitForm(0, &req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(201, "Form submitted successfully", id, ctx)
}

// ======== CREATE ========

// CreateForm godocs
//
//	@Summary		Create form
//	@Description	Create a new form structure
//	@Tags			Forms
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.CreateFormRequest	true	"Form creation data"
//	@Success		201		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/form [post]			// <------- add admin endpoint
//
//	@Security		ApiKeyAuth
func (h *FormHandler) CreateForm(ctx *gin.Context) {
	var req models.CreateFormRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	// value, exists := ctx.Get("user")
	// claims, ok := value.(*models.ManagedClaims)
	// if !exists || !ok {
	// 	ctx.Error(errs.New(errs.Unauthorized, "Unauthorized", nil))
	// 	return
	// }

	id, err := h.service.CreateForm(0, &req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(201, "Form created successfully", id, ctx)
}

// CreatePage godocs
//
//	@Summary		Create form page
//	@Description	Create a new page within a form
//	@Tags			Forms
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.CreatePageRequest	true	"Page creation data"
//	@Success		201		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/form/page [post]			// <------- add admin endpoint
//
//	@Security		ApiKeyAuth
func (h *FormHandler) CreatePage(ctx *gin.Context) {
	var req models.CreatePageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	id, err := h.service.CreatePage(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(201, "Page created successfully", id, ctx)
}

// CreateQuestion godocs
//
//	@Summary		Create form question
//	@Description	Create a new question within a form page
//	@Tags			Forms
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.CreateQuestionRequest	true	"Question creation data"
//	@Success		201		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/form/question [post]			// <------- add admin endpoint
//
//	@Security		ApiKeyAuth
func (h *FormHandler) CreateQuestion(ctx *gin.Context) {
	var req models.CreateQuestionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	id, err := h.service.CreateQuestion(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(201, "Question created successfully", id, ctx)
}

// ======== UPDATE ========

// UpdateForm godocs
//
//	@Summary		Update form
//	@Description	Update an existing form's basic information
//	@Tags			Forms
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.UpdateFormRequest	true	"Form update data"
//	@Success		200		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		404		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/form [put]			// <------- add admin endpoint
//
//	@Security		ApiKeyAuth
func (h *FormHandler) UpdateForm(ctx *gin.Context) {
	var req models.UpdateFormRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err := h.service.UpdateForm(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Form updated successfully", req.ID, ctx)
}

// UpdatePage godocs
//
//	@Summary		Update form page
//	@Description	Update an existing page's details within a form
//	@Tags			Forms
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.UpdatePageRequest	true	"Page update data"
//	@Success		200		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		404		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/form/page [put]			// <------- add admin endpoint
//
//	@Security		ApiKeyAuth
func (h *FormHandler) UpdatePage(ctx *gin.Context) {
	var req models.UpdatePageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err := h.service.UpdatePage(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Page updated successfully", req.ID, ctx)
}

// UpdateQuestion godocs
//
//	@Summary		Update form question
//	@Description	Update an existing question's details within a form page
//	@Tags			Forms
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.UpdateQuestionRequest	true	"Question update data"
//	@Success		200		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		404		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/form/question [put]			// <------- add admin endpoint
//
//	@Security		ApiKeyAuth
func (h *FormHandler) UpdateQuestion(ctx *gin.Context) {
	var req models.UpdateQuestionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err := h.service.UpdateQuestion(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Question updated successfully", req.ID, ctx)
}

// UpdateResponseStatus godocs

// 	@Summary		Update response status
// 	@Description	Update the status of a form response
// 	@Tags			Forms
// 	@Accept			json
// 	@Produce		json
// 	@Param			body	body		models.UpdateResponseStatusRequest	true	"Status update data"
// 	@Success		200		{object}	response.TransactionResponse
// 	@Failure		400		{object}	response.BaseError
// 	@Failure		401		{object}	response.BaseError
// 	@Failure		404		{object}	response.BaseError
// 	@Failure		500		{object}	response.BaseError
// 	@Router			/form/response/status [put]			// <------- add admin endpoint

// @Security		ApiKeyAuth

func (h *FormHandler) UpdateResponseStatus(ctx *gin.Context) {
	var req models.UpdateResponseStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err := h.service.UpdateResponseStatus(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Response status updated successfully", req.ID, ctx)
}

// ======== GET ONE ========

// GetResponseByID godocs
//
//	@Summary		Get response by ID
//	@Description	Get a specific form response by its ID
//	@Tags			Forms
//	@Produce		json
//	@Param			id	path		int	true	"Response ID"
//	@Success		200	{object}	models.FormResponseDTO
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		404	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/form/response/{id} [get]			// <------- add admin endpoint
//
//	@Security		ApiKeyAuth

func (h *FormHandler) GetResponseByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	resp, err := h.service.GetResponseByID(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, resp)
}

// ======== GET MANY ========

// GetAllForms godocs
//
//	@Summary		Get all forms
//	@Description	Get a list of all forms for administration
//	@Tags			Forms
//	@Produce		json
//	@Success		200	{array}		models.FormSummaryResponse
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/form [get]			// <------- add admin endpoint
//
//	@Security		ApiKeyAuth
func (h *FormHandler) GetAllForms(ctx *gin.Context) {
	forms, err := h.service.GetAllForms()
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, forms)
}

// GetResponsesByFormID godocs
//
//	@Summary		Get responses by form ID
//	@Description	Get all responses for a specific form
//	@Tags			Forms
//	@Produce		json
//	@Param			id	path		int	true	"Form ID"
//	@Success		200	{array}		models.FormResponseDTO
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/form/responses/{id} [get]			// <------- add admin endpoint
//
//	@Security		ApiKeyAuth

func (h *FormHandler) GetResponsesByFormID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	responses, err := h.service.GetResponsesByFormID(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, responses)
}

// GetUserResponsesForForm godocs
//
//	@Summary		Get user responses for form
//	@Description	Get all responses submitted by the current user for a specific form
//	@Tags			Forms
//	@Produce		json
//	@Param			id	path		int	true	"Form ID"
//	@Success		200	{array}		models.FormResponseDTO
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/form/responses/{id} [get]
//
//	@Security		ApiKeyAuth

func (h *FormHandler) GetUserResponsesForForm(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	// value, exists := ctx.Get("user")
	// claims, ok := value.(*models.ManagedClaims)
	// if !exists || !ok {
	// 	ctx.Error(errs.New(errs.Unauthorized, "Unauthorized", nil))
	// 	return
	// }

	responses, err := h.service.GetUserResponsesForForm(0, id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, responses)
}

// ======== DELETE ========

// DeleteForm godocs
//
//	@Summary		Delete form
//	@Description	Delete a form and all its associated pages and questions
//	@Tags			Forms
//	@Produce		json
//	@Param			id	path		int	true	"Form ID"
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		404	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/form/{id} [delete]			// <------- add admin endpoint
//
//	@Security		ApiKeyAuth
func (h *FormHandler) DeleteForm(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err = h.service.DeleteForm(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Form deleted successfully", id, ctx)
}

// DeletePage godocs
//
//	@Summary		Delete form page
//	@Description	Delete a page and all its associated questions
//	@Tags			Forms
//	@Produce		json
//	@Param			id	path		int	true	"Page ID"
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		404	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/form/page/{id} [delete]			// <------- add admin endpoint
//
//	@Security		ApiKeyAuth
func (h *FormHandler) DeletePage(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err = h.service.DeletePage(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Page deleted successfully", id, ctx)
}

// DeleteQuestion godocs
//
//	@Summary		Delete form question
//	@Description	Delete a specific question from a form page
//	@Tags			Forms
//	@Produce		json
//	@Param			id	path		int	true	"Question ID"
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		404	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/form/question/{id} [delete]			// <------- add admin endpoint
//
//	@Security		ApiKeyAuth
func (h *FormHandler) DeleteQuestion(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err = h.service.DeleteQuestion(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Question deleted successfully", id, ctx)
}

// DeleteResponse godocs
//
//	@Summary		Delete form response
//	@Description	Delete a specific form response by its ID
//	@Tags			Forms
//	@Produce		json
//	@Param			id	path		int	true	"Response ID"
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		404	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/form/response/{id} [delete]			// <------- add admin endpoint
//
//	@Security		ApiKeyAuth

func (h *FormHandler) DeleteResponse(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err = h.service.DeleteResponse(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Response deleted successfully", id, ctx)
}
