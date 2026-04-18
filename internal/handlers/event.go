package handlers

import (
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/response"
	"sea-api/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	EventService *services.EventService
}

func NewEventHandler(eventService *services.EventService) *EventHandler {
	return &EventHandler{
		EventService: eventService,
	}
}

// GetAllEvents godocs
//
//	@Summary		Get all events
//	@Description	Get a list of all events for administration
//	@Tags			Events
//	@Produce		json
//	@Success		200	{array}		models.EventListResponse
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/event [get]
//
//	@Security		ApiKeyAuth
func (h *EventHandler) GetAllEvents(ctx *gin.Context) {
	events, err := h.EventService.GetAllEvents()
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(200, events)
}

// GetEventByID godocs
//
//	@Summary		Get event by ID
//	@Description	Get event details by ID for administration
//	@Tags			Events
//	@Produce		json
//	@Param			id	path	int	true	"Event ID"
//	@Success		200	{object}	models.EventDTO
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		404	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/event/{id} [get]
//
//	@Security		ApiKeyAuth
func (h *EventHandler) GetEventByID(ctx *gin.Context) {
	id := ctx.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	event, err := h.EventService.GetEventByID(int64(intId))
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(200, event)
}

// CreateEvent godocs
//
//	@Summary		Create event
//	@Description	Create a new event
//	@Tags			Events
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.EventDTO	true	"Event data"
//	@Success		201		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/event [post]
//
//	@Security		ApiKeyAuth
func (h *EventHandler) CreateEvent(ctx *gin.Context) {
	var event models.EventDTO
	if err := ctx.ShouldBindJSON(&event); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	id, err := h.EventService.CreateEvent(&event)
	if err != nil {
		ctx.Error(err)
		return
	}
	event.ID = id
	response.NewTransactionResponse(201, "Event created successfully", id, ctx)
}

// UpdateEvent godocs
//
//	@Summary		Update event
//	@Description	Update an existing event
//	@Tags			Events
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.EventDTO	true	"Event update data"
//	@Success		200		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		404		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/event [put]
//
//	@Security		ApiKeyAuth
func (h *EventHandler) UpdateEvent(ctx *gin.Context) {
	var event models.EventDTO
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
