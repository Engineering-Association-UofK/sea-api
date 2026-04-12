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

func (h *FormHandler) GetAllForms(ctx *gin.Context) {
	forms, err := h.service.GetAllForms()
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, forms)
}

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
