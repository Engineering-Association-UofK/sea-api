package handlers

import (
	"log/slog"
	"sea-api/internal/exception"
	"sea-api/internal/models"
	"sea-api/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type EventHandler struct {
	EventService *services.EventService
}

func NewEventHandler(db *sqlx.DB) *EventHandler {
	return &EventHandler{
		EventService: services.NewEventService(db),
	}
}

func (h *EventHandler) GetAllEvents(ctx *gin.Context) {
	events, err := h.EventService.GetAllEvents()
	if err != nil {
		slog.Error("error getting events", "error", err)
		exception.InternalServerError(ctx)
		return
	}
	ctx.JSON(200, events)
}

func (h *EventHandler) GetEventByID(ctx *gin.Context) {
	id := ctx.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		exception.BadRequest(ctx)
		return
	}

	event, err := h.EventService.GetEventByID(int64(intId))
	if err != nil {
		if err.Error() == "event not found" {
			exception.NotFound(ctx)
			return
		}
		slog.Error("error getting event", "error", err)
		exception.InternalServerError(ctx)
		return
	}
	ctx.JSON(200, event)
}

func (h *EventHandler) CreateEvent(ctx *gin.Context) {
	var event models.EventDTO
	if err := ctx.ShouldBindJSON(&event); err != nil {
		exception.BadRequest(ctx)
		return
	}

	id, err := h.EventService.CreateEvent(&event)
	if err != nil {
		if err.Error() == "event not found" {
			exception.NotFound(ctx)
			return
		}
		slog.Error("error creating event", "error", err)
		exception.InternalServerError(ctx)
		return
	}
	event.ID = id
	ctx.JSON(200, gin.H{"message": "Event created successfully", "event": event})
}

func (h *EventHandler) UpdateEvent(ctx *gin.Context) {
	var event models.EventDTO
	if err := ctx.ShouldBindJSON(&event); err != nil {
		exception.BadRequest(ctx)
		return
	}

	if err := h.EventService.UpdateEvent(&event); err != nil {
		slog.Error("error updating event", "error", err)
		exception.InternalServerError(ctx)
		return
	}

	ctx.JSON(200, gin.H{"message": "Event updated successfully", "event": event})
}

func (h *EventHandler) DeleteEvent(ctx *gin.Context) {
	id := ctx.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		exception.BadRequest(ctx)
		return
	}
	if err := h.EventService.DeleteEvent(int64(intId)); err != nil {
		slog.Error("error deleting event", "error", err)
		exception.InternalServerError(ctx)
		return
	}
	ctx.JSON(200, gin.H{"message": "Event deleted successfully"})
}
