package handlers

import (
	"log/slog"
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
		response.BadRequest(ctx)
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
		response.BadRequest(ctx)
		return
	}

	id, err := h.EventService.CreateEvent(&event)
	if err != nil {
		slog.Error("error creating event", "error", err)
		ctx.Error(err)
		return
	}
	event.ID = id
	ctx.JSON(200, gin.H{"message": "Event created successfully", "event": event})
}

func (h *EventHandler) UpdateEvent(ctx *gin.Context) {
	var event models.EventDTO
	if err := ctx.ShouldBindJSON(&event); err != nil {
		response.BadRequest(ctx)
		return
	}

	if err := h.EventService.UpdateEvent(&event); err != nil {
		slog.Error("error updating event", "error", err)
		response.InternalServerError(ctx)
		return
	}

	ctx.JSON(200, gin.H{"message": "Event updated successfully", "event": event})
}

func (h *EventHandler) DeleteEvent(ctx *gin.Context) {
	id := ctx.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		response.BadRequest(ctx)
		return
	}
	if err := h.EventService.DeleteEvent(int64(intId)); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(200, gin.H{"message": "Event deleted successfully"})
}
