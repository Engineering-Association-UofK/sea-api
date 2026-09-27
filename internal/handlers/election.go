package handlers

import (
	"net/http"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/models/electionsmodels"
	"sea-api/internal/response"
	"sea-api/internal/services/electionservice"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ElectionHandler struct {
	service *electionservice.ElectionService
}

func NewElectionHandler(service *electionservice.ElectionService) *ElectionHandler {
	return &ElectionHandler{service: service}
}

// GetElectionStatistics godocs
//
//	@Summary		Get election statistics
//	@Description	Retrieve real-time voter metrics and election cycle timelines
//	@Tags			Election:Public
//	@Produce		json
//	@Success		200	{object}	electionsmodels.VoteStatistics
//	@Failure		500	{object}	response.BaseError
//	@Router			/election/statistics [get]
func (h *ElectionHandler) GetElectionStatistics(c *gin.Context) {
	stats, err := h.service.ElectionStatistics()
	if err != nil {
		c.Error(err)
		return
	}

	c.PureJSON(http.StatusOK, stats)
}

// ======== CANDIDATE MANAGEMENT (SETUP PHASE) ========

// CreateCandidate godocs
//
//	@Summary		Create candidate
//	@Description	Add a new candidate for the active election cycle
//	@Tags			Election:Admin
//	@Accept			json
//	@Produce		json
//	@Param			body	body		electionsmodels.CandidateCreateRequest	true	"Candidate data"
//	@Success		201		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/election/candidate [post]
//
//	@Security		ApiKeyAuth
func (h *ElectionHandler) CreateCandidate(c *gin.Context) {
	var req electionsmodels.CandidateCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	id, err := h.service.CreateCandidate(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(201, "Candidate created successfully", id, c)
}

// UpdateCandidate godocs
//
//	@Summary		Update candidate
//	@Description	Update an existing candidate for the active election cycle
//	@Tags			Election:Admin
//	@Accept			json
//	@Produce		json
//	@Param			body	body		electionsmodels.CandidateUpdateRequest	true	"Candidate data"
//	@Success		200		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/election/candidate/{id} [put]
//
//	@Security		ApiKeyAuth
func (h *ElectionHandler) UpdateCandidate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	var req electionsmodels.CandidateUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err = h.service.UpdateCandidate(id, req.Belonging)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Candidate updated successfully", id, c)
}

// GetCandidateList godocs
//
//	@Summary		Get candidate list
//	@Description	Get the list of all candidates for the active election cycle
//	@Tags			Election:Public
//	@Produce		json
//	@Success		200	{array}		electionsmodels.CandidateResponse
//	@Failure		500	{object}	response.BaseError
//	@Router			/election/candidates [get]
func (h *ElectionHandler) GetCandidateList(c *gin.Context) {
	candidates, err := h.service.GetCandidateList(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.PureJSON(http.StatusOK, candidates)
}

// RemoveCandidate godocs
//
//	@Summary		Remove candidate
//	@Description	Remove a candidate by their ID
//	@Tags			Election:Admin
//	@Produce		json
//	@Param			id	path		int	true	"Candidate ID"
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		400	{object}	response.BaseError
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/election/candidate/{id} [delete]
//
//	@Security		ApiKeyAuth
func (h *ElectionHandler) RemoveCandidate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err = h.service.RemoveCandidate(id)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Candidate removed successfully", id, c)
}

// ======== TICKETING (PRE-VOTING PHASE) ========

// GetTicket godocs
//
//	@Summary		Get election ticket
//	@Description	Generate and retrieve a unique voting ticket for the requesting user
//	@Tags			Election:Voting
//	@Produce		json
//	@Success		200	{object}	electionsmodels.TicketResponse
//	@Failure		401	{object}	response.BaseError
//	@Failure		403	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/account/election/ticket [get]
//
//	@Security		ApiKeyAuth
func (h *ElectionHandler) GetTicket(c *gin.Context) {
	value, exists := c.Get("user")
	claims, ok := value.(*models.ManagedClaims)
	if !exists || !ok {
		c.Error(errs.New(errs.Unauthorized, "Unauthorized", nil))
		return
	}

	ticket, err := h.service.GetTicket(c.Request.Context(), claims.UserID)
	if err != nil {
		c.Error(err)
		return
	}

	c.PureJSON(http.StatusOK, ticket)
}

// ======== VOTING PHASE ========

// Vote godocs
//
//	@Summary		Submit votes
//	@Description	Submit a batch of candidate votes using a valid ticket
//	@Tags			Election:Voting
//	@Accept			json
//	@Produce		json
//	@Param			body	body		electionsmodels.VoteRequest	true	"Vote submission data"
//	@Success		200		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		403		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/election/vote [post]
func (h *ElectionHandler) Vote(c *gin.Context) {
	var req electionsmodels.VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	err := h.service.Vote(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Votes submitted successfully", 0, c)
}

// ======== RESULTS RESOLUTION (POST-VOTING PHASE) ========

// ResetElection godocs
//
//	@Summary		Reset and resolve election
//	@Description	Resolve results (applying Department Quotas and Dense Rank), save to history, and clear active votes/tickets
//	@Tags			Election:Admin
//	@Produce		json
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/election/reset [post]
//
//	@Security		ApiKeyAuth
func (h *ElectionHandler) ResolveElection(c *gin.Context) {
	err := h.service.ResolveAndReset(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Election resolved, saved, and reset successfully", 0, c)
}

// GetResults godocs
//
//	@Summary		Get election results by cycle
//	@Description	Retrieve resolved election results for a specific cycle number
//	@Tags			Election:Public
//	@Produce		json
//	@Param			cycle	query		int	true	"Election Cycle Number"
//	@Success		200		{array}		electionsmodels.Result
//	@Failure		400		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/election/results [get]
func (h *ElectionHandler) GetResults(c *gin.Context) {
	cycleStr := c.Query("cycle")
	cycle, err := strconv.Atoi(cycleStr)
	if err != nil || cycle <= 0 {
		c.Error(errs.New(errs.BadRequest, "Invalid or missing cycle parameter", nil))
		return
	}

	results, err := h.service.GetResults(cycle)
	if err != nil {
		c.Error(err)
		return
	}

	c.PureJSON(http.StatusOK, results)
}
