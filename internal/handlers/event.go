package handlers

import (
	"log/slog"
	"net/http"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/models/eventmodels"
	"sea-api/internal/response"
	"sea-api/internal/services/eventservice"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	service *eventservice.EventService
}

func NewEventHandler(service *eventservice.EventService) *EventHandler {
	return &EventHandler{service: service}
}

// ======== PUBLIC EVENTS ========

// ApplyForEvent godocs
//
//	@Summary		Apply for event
//	@Description	Submit an application to participate in an event
//	@Tags			Events:v2:Public
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int								true	"Event ID"
//	@Success		201		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/account/event/{id} [post]
//
//	@Security		ApiKeyAuth
func (h *EventHandler) ApplyForEvent(ctx *gin.Context) {
	idStr := ctx.Param("id")
	eventId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	value, exists := ctx.Get("user")
	claims, ok := value.(*models.ManagedClaims)
	if !exists || !ok {
		slog.Debug("NoN", "Claims", claims, "Value", value)
		ctx.Error(errs.New(errs.Unauthorized, "Unauthorized", nil))
		return
	}

	res, err := h.service.Apply(eventId, *claims)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, res)
}

// GetEventViewList godocs
//
//	@Summary		Get public events list
//	@Description	Get a list of all published events for public viewing
//	@Tags			Events:v2:Public
//	@Produce		json
//	@Param			limit	query		int	false	"Content count limit"
//	@Param			page	query		int	false	"Page number"
//	@Param			search-name	query		string	false	"Search query"
//	@Param			belonging	query		models.Secretariat	false	"Secretariat hosting this event"
//	@Success		200		{object}	eventmodels.EventViewListResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/event [get]
func (h *EventHandler) GetEventViewList(ctx *gin.Context) {
	var req eventmodels.EventListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	events, err := h.service.GetEventViewList(ctx.Request.Context(), &req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.PureJSON(http.StatusOK, events)
}

// GetEventView godocs
//
//	@Summary		Get public event details
//	@Description	Get details of a specific event for public viewing
//	@Tags			Events:v2:Public
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Success		200	{object}	eventmodels.EventViewResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		404	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/event/{id} [get]
func (h *EventHandler) GetEventView(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	event, err := h.service.GetEventView(ctx.Request.Context(), id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.PureJSON(http.StatusOK, event)
}

// ======== ADMIN EVENTS ========

// CreateEvent godocs
//
//	@Summary		Create event
//	@Description	Create a new event
//	@Tags			Events:v2
//	@Accept			json
//	@Produce		json
//	@Param			body	body		eventmodels.EventRequest	true	"Event creation data"
//	@Success		201		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/event [post]
//
//	@Security		ApiKeyAuth
func (h *EventHandler) CreateEvent(ctx *gin.Context) {
	var req eventmodels.EventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	id, err := h.service.Create(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(201, "Event created successfully", id, ctx)
}

// UpdateEvent godocs
//
//	@Summary		Update event
//	@Description	Update an existing event's details
//	@Tags			Events:v2
//	@Accept			json
//	@Produce		json
//	@Param			body	body		eventmodels.EventUpdateRequest	true	"Event update data"
//	@Success		200		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		404		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/event [put]
//
//	@Security		ApiKeyAuth
func (h *EventHandler) UpdateEvent(ctx *gin.Context) {
	var req eventmodels.EventUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err := h.service.Update(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Event updated successfully", req.ID, ctx)
}

// GetEvent godocs
//
//	@Summary		Get event
//	@Description	Get complete event details for administration
//	@Tags			Events:v2
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Success		200	{object}	eventmodels.EventUpdateRequest
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		404	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/event/{id} [get]
//
//	@Security		ApiKeyAuth
func (h *EventHandler) GetEvent(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	event, err := h.service.Get(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.PureJSON(http.StatusOK, event)
}

// GetEventList godocs
//
//	@Summary		Get all events
//	@Description	Get a list of all events for administration
//	@Tags			Events:v2
//	@Produce		json
//	@Param			limit	query		int	false	"Content count limit"
//	@Param			page	query		int	false	"Page number"
//	@Param			search-name	query		string	false	"Search query"
//	@Param			belonging	query		models.Secretariat	false	"Secretariat hosting this event"
//	@Success		200		{object}	eventmodels.EventListResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/event [get]
//
//	@Security		ApiKeyAuth
func (h *EventHandler) GetEventList(ctx *gin.Context) {
	var req eventmodels.EventListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	events, err := h.service.GetList(ctx.Request.Context(), &req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.PureJSON(http.StatusOK, events)
}

// DeleteEvent godocs
//
//	@Summary		Delete event
//	@Description	Delete an event and its related data
//	@Tags			Events:v2
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
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Event deleted successfully", id, ctx)
}

// ======== COORDS ========

// GetCoordList godocs
//
//	@Summary		Get Coordinators
//	@Description	Get all coordinators for event
//	@Tags			Events:v2:coordinator
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Success		200		{array}	eventmodels.EventCoord
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/event/{id}/coord [get]
//
//	@Security		ApiKeyAuth
func (h *EventHandler) GetCoordList(ctx *gin.Context) {
	idStr := ctx.Param("id")
	eventId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	coords, err := h.service.GetCoordinators(eventId)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.PureJSON(http.StatusOK, coords)
}

// CreateCoords godocs
//
//	@Summary		Create coordinators
//	@Description	Assign coordinators to an event
//	@Tags			Events:v2:coordinator
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Param			body	body		eventmodels.AddCoordsRequest	true	"Coordinator data"
//	@Success		201		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/event/{id}/coord [post]
//
//	@Security		ApiKeyAuth
func (h *EventHandler) CreateCoords(ctx *gin.Context) {
	var req eventmodels.AddCoordsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	idStr := ctx.Param("id")
	eventId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err = h.service.AddCoordinators(eventId, &req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(201, "Coordinator assigned successfully", eventId, ctx)
}

// UpdateCoord godocs
//
//	@Summary		Update coordinator
//	@Description	Update coordinator details/roles
//	@Tags			Events:v2:coordinator
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Param			body	body		eventmodels.CoordRequest	true	"Coordinator update data"
//	@Success		200		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/event/{id}/coord/{coord_id} [put]
//
//	@Security		ApiKeyAuth
func (h *EventHandler) UpdateCoord(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	coordIdStr := ctx.Param("coord_id")
	coordId, err := strconv.ParseInt(coordIdStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	var req eventmodels.CoordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err = h.service.UpdateCoordinator(id, coordId, &req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Coordinator updated successfully", coordId, ctx)
}

// DeleteCoord godocs
//
//	@Summary		Remove coordinator
//	@Description	Remove a coordinator from an event
//	@Tags			Events:v2:coordinator
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Param			coord_id	path		int	true	"Coordinator assignment ID"
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/event/{id}/coord/{coord_id} [delete]
//
//	@Security		ApiKeyAuth
func (h *EventHandler) DeleteCoord(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	coordIdStr := ctx.Param("coord_id")
	coordId, err := strconv.ParseInt(coordIdStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err = h.service.DeleteCoordinator(coordId, id)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Coordinator removed successfully", id, ctx)
}

// DeleteCoords godocs
//
//	@Summary		Remove coordinators
//	@Description	Remove all coordinator from an event
//	@Tags			Events:v2:coordinator
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/event/{id}/coord [delete]
//
//	@Security		ApiKeyAuth
func (h *EventHandler) DeleteCoords(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err = h.service.DeleteCoordinatorsByEventID(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "All coordinators for event removed successfully", id, ctx)
}

// ======== PARTICIPATION ========

// GetApplicationList godocs
//
//	@Summary		Get event applications
//	@Description	Get a list of applications for a specific event
//	@Tags			Events:v2:Participation
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Param			limit		query		int	false	"Content count limit"
//	@Param			page		query		int	false	"Page number"
//	@Success		200		{object}	eventmodels.EventApplicationListResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/event/{id}/application [get]
//
//	@Security		ApiKeyAuth
func (h *EventHandler) GetApplicationList(ctx *gin.Context) {
	idStr := ctx.Param("id")
	eventId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	var req eventmodels.EventApplicationListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	res, err := h.service.GetApplications(eventId, &req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.PureJSON(http.StatusOK, res)
}

// AcceptApplication godocs
//
//	@Summary		Accept application
//	@Description	Accept application to an event
//	@Tags			Events:v2:Participation
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Param			application_id	path		int	true	"application ID"
//	@Success		200		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/event/{id}/application/{application_id} [post]
//
//	@Security		ApiKeyAuth
func (h *EventHandler) AcceptApplication(ctx *gin.Context) {
	idStr := ctx.Param("id")
	eventId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	appIdStr := ctx.Param("application_id")
	appId, err := strconv.ParseInt(appIdStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err = h.service.AcceptApplication(eventId, appId)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Application Accepted successfully", appId, ctx)
}

// RejectApplication godocs
//
//	@Summary		Reject application
//	@Description	reject an event application record
//	@Tags			Events:v2:Participation
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Param			application_id	path		int	true	"application ID"
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		404	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/event/{id}/application/{application_id} [delete]
//
//	@Security		ApiKeyAuth
func (h *EventHandler) RejectApplication(ctx *gin.Context) {
	idStr := ctx.Param("id")
	eventId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	appIdStr := ctx.Param("application_id")
	appId, err := strconv.ParseInt(appIdStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err = h.service.RejectApplication(eventId, appId)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Application Rejected successfully", appId, ctx)
}

// GetParticipantList godocs
//
//	@Summary		Get event participants
//	@Description	Get a list of confirmed participants for a specific event
//	@Tags			Events:v2:Participation
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Param			limit		query		int	false	"Content count limit"
//	@Param			page		query		int	false	"Page number"
//	@Success		200		{object}	eventmodels.EventParticipantListResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/event/{id}/participant [get]
//
//	@Security		ApiKeyAuth
func (h *EventHandler) GetParticipantList(ctx *gin.Context) {
	idStr := ctx.Param("id")
	eventId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	var req eventmodels.EventParticipantListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	participants, err := h.service.GetParticipants(eventId, &req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.PureJSON(http.StatusOK, participants)
}

// RemoveParticipant godocs
//
//	@Summary		Remove participant
//	@Description	Remove a participant from an event
//	@Tags			Events:v2:Participation
//	@Produce		json
//	@Param			id	path		int	true	"Event ID"
//	@Param			participation_id	path		int	true	"Participation ID"
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		404	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/event/{id}/participant/{participation_id} [delete]
//
//	@Security		ApiKeyAuth
func (h *EventHandler) RemoveParticipant(ctx *gin.Context) {
	idStr := ctx.Param("id")
	eventId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	userIdStr := ctx.Param("participation_id")
	partID, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err = h.service.RemoveParticipant(eventId, partID)
	if err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Participant removed successfully", partID, ctx)
}
