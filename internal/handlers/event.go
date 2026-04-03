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

func (h *EventHandler) GetAllEvents(ctx *gin.Context) {
	events, err := h.EventService.GetAllEvents()
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(200, events)
}

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

func (h *EventHandler) ImportUsers(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}
	defer file.Close()

	err = h.EventService.ImportUsers(id, file)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Users imported successfully", id, ctx)
}
