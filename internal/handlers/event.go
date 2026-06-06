package handlers

import (
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/response"
	"sea-api/internal/services/event"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	EventService *event.EventService
}

func NewEventHandler(eventService *event.EventService) *EventHandler {
	return &EventHandler{
		EventService: eventService,
	}
}

// GetEventsList godocs
//
//	@Summary		Get events list
//	@Description	Get a paginated list of events with filters for public view
//	@Tags			Public
//	@Produce		json
//	@Param			limit	query		int		true	"Content count limit"
//	@Param			page	query		int		true	"Page number"
//	@Param			type	query		string	false	"Event type filter"
//	@Param			status	query		string	false	"Event status filter"
//	@Success		200		{object}	models.EventViewListResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/event [get]
func (h *EventHandler) GetEventsList(ctx *gin.Context) {
	var query models.QueryEventPublicRequest
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Invalid query parameters", nil))
		return
	}

	resp, err := h.EventService.GetViewList(query)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, resp)
}

// GetEventByID godocs
//
//	@Summary		Get event by ID
//	@Description	Get event details by ID for administration
//	@Tags			Public
//	@Produce		json
//	@Param			id	path	int	true	"Event ID"
//	@Success		200	{object}	models.EventViewDetailsResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		404	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/event/{id} [get]
func (h *EventHandler) GetEventByID(ctx *gin.Context) {
	id := ctx.Param("id")
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	event, err := h.EventService.GetEventViewDetails(ctx.Request.Context(), intId)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(200, event)
}

// GetEventDetailsAdmin godocs
//
//	@Summary		Get event details for admin
//	@Description	Get full event details including grading scheme and participants count for administration
//	@Tags			Events
//	@Produce		json
//	@Param			id	path	int	true	"Event ID"
//	@Success		200	{object}	models.EventDetailsAdminResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		404	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/event/{id} [get]
//
//	@Security		ApiKeyAuth
func (h *EventHandler) GetEventDetailsAdmin(ctx *gin.Context) {
	id := ctx.Param("id")
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Invalid event ID", nil))
		return
	}

	event, err := h.EventService.GetEventDetailsAdmin(ctx.Request.Context(), intId)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(200, event)
}

// GetEventParticipants godocs
//
//	@Summary		Get event participants
//	@Description	Get a paginated list of participants for a specific event
//	@Tags			Events
//	@Produce		json
//	@Param			id		path		int	true	"Event ID"
//	@Param			limit	query		int	true	"Content count limit"
//	@Param			page	query		int	true	"Page number"
//	@Success		200		{array}		models.EventParticipantsResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/event/{id}/participants [get]
//
//	@Security		ApiKeyAuth
func (h *EventHandler) GetEventParticipants(ctx *gin.Context) {
	id := ctx.Param("id")
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Invalid event ID", nil))
		return
	}

	var req models.ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Invalid pagination parameters", nil))
		return
	}

	resp, err := h.EventService.GetEventParticipants(intId, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, resp)
}

// UpdateEventParticipants godocs
//
//	@Summary		Update event participants
//	@Description	Batch update status, completion, and grades for event participants
//	@Tags			Events
//	@Accept			json
//	@Produce		json
//	@Param			id				path		int	true	"Event ID"
//	@Param			body	body		[]models.ParticipantUpdateRequest	true	"Participants update data"
//	@Success		200		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/event/{id}/participants [put]
//
//	@Security		ApiKeyAuth
func (h *EventHandler) UpdateEventParticipants(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Invalid event ID", nil))
		return
	}

	var req []models.ParticipantUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Invalid request body", nil))
		return
	}

	if err := h.EventService.BatchUpdateParticipant(id, req); err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Participants updated successfully", 0, ctx)
}

// CreateEvent godocs
//
//	@Summary		Create event
//	@Description	Create a new event
//	@Tags			Events
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.EventCreateRequest	true	"Event data"
//	@Success		201		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/event [post]
//
//	@Security		ApiKeyAuth
func (h *EventHandler) CreateEvent(ctx *gin.Context) {
	var event models.EventCreateRequest
	if err := ctx.ShouldBindJSON(&event); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	id, err := h.EventService.CreateEvent(&event)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(201, "Event created successfully", id, ctx)
}

// UpdateEvent godocs
//
//	@Summary		Update event
//	@Description	Update an existing event
//	@Tags			Events
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.EventUpdateRequest	true	"Event update data"
//	@Success		200		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		404		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/event [put]
//
//	@Security		ApiKeyAuth
func (h *EventHandler) UpdateEvent(ctx *gin.Context) {
	var event models.EventUpdateRequest
	if err := ctx.ShouldBindJSON(&event); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	if err := h.EventService.UpdateEvent(&event); err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Event updated successfully", event.ID, ctx)
}

// DeleteEvent godocs
//
//	@Summary		Delete event
//	@Description	Delete an event by its ID
//	@Tags			Events
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		404	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/event/{id} [delete]
//
//	@Security		ApiKeyAuth
func (h *EventHandler) DeleteEvent(ctx *gin.Context) {
	id := ctx.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}
	if err := h.EventService.DeleteEvent(int64(intId)); err != nil {
		ctx.Error(err)
		return
	}
	response.NewTransactionResponse(200, "Event deleted successfully", int64(intId), ctx)
}

///////////////////////
///   Application   ///
///////////////////////

// / GetApplicationStatus godocs
//
//	@Summary		Get application status
//	@Description	Get the registration status and details for all events for the current user
//	@Tags			Applications
//	@Produce		json
//	@Param			limit	query		int	true	"Content count limit"
//	@Param			page	query		int	true	"Page number"
//	@Param			eventID	query		int	true	"Event ID (Send 0 for all events)"
//	@Success		200	{object}	models.ApplicationStatusList
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/account/event/status [get]
//
//	@Security		ApiKeyAuth
func (h *EventHandler) GetApplicationStatus(ctx *gin.Context) {
	var req models.ListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request, need limit number", nil))
		return
	}

	eventId := ctx.Query("eventID")
	intId, err := strconv.ParseInt(eventId, 10, 64)
	if eventId != "0" {
		intId = 0
	} else if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	value, exists := ctx.Get("user")
	claims, ok := value.(*models.ManagedClaims)
	if !exists || !ok {
		ctx.Error(errs.New(errs.Unauthorized, "Unauthorized", nil))
		return
	}

	resp, err := h.EventService.Status(claims.UserID, intId, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, resp)
}

// ApplyForEvent godocs
//
//	@Summary		Apply for event
//	@Description	Submit an application for a specific event for the current user
//	@Tags			Applications
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Success		200	{object}	models.ApplyResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/account/event/apply/{id} [post]
//
//	@Security		ApiKeyAuth
func (h *EventHandler) ApplyForEvent(ctx *gin.Context) {
	eventId := ctx.Param("id")
	intId, err := strconv.ParseInt(eventId, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	value, exists := ctx.Get("user")
	claims, ok := value.(*models.ManagedClaims)
	if !exists || !ok {
		ctx.Error(errs.New(errs.Unauthorized, "Unauthorized", nil))
		return
	}

	resp, err := h.EventService.Apply(claims.UserID, intId)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, resp)
}

// CancelApplicationForEvent godocs
//
//	@Summary		Cancel event application
//	@Description	Cancel an existing application for a specific event for the current user
//	@Tags			Applications
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/account/event/cancel/{id} [post]
//
//	@Security		ApiKeyAuth
func (h *EventHandler) CancelApplicationForEvent(ctx *gin.Context) {
	eventId := ctx.Param("id")
	intId, err := strconv.ParseInt(eventId, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	value, exists := ctx.Get("user")
	claims, ok := value.(*models.ManagedClaims)
	if !exists || !ok {
		ctx.Error(errs.New(errs.Unauthorized, "Unauthorized", nil))
		return
	}

	if err := h.EventService.Cancel(claims.UserID, intId); err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Event Application Canceled successfully", int64(intId), ctx)
}
