package handlers

import (
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/response"
	_ "sea-api/internal/response"
	"sea-api/internal/services/bot"

	"github.com/gin-gonic/gin"
)

type BotHandler struct {
	BotService *bot.BotService
}

func NewBotHandler(botService *bot.BotService) *BotHandler {
	return &BotHandler{BotService: botService}
}

// GetBotGraph godocs
//
//	@Summary		Get bot graph
//	@Description	Get the entire bot conversation graph for the admin builder
//	@Tags			Bot
//	@Produce		json
//	@Success		200	{object}	models.BotGraphData
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/bot/graph [get]
//
//	@Security		ApiKeyAuth
func (h *BotHandler) GetBotGraph(ctx *gin.Context) {
	graph, err := h.BotService.GetBotGraph()
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.PureJSON(200, graph)
}

// UpdateBotGraph godocs
//
//	@Summary		Update bot graph
//	@Description	Update the entire bot conversation graph structure
//	@Tags			Bot
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.BotGraphData	true	"Bot graph data"
//	@Success		200		{object}	response.TransactionResponse
//	@Failure		400		{object}	response.ErrorResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		401		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/admin/bot/graph [put]
//
//	@Security		ApiKeyAuth
func (h *BotHandler) UpdateBotGraph(ctx *gin.Context) {
	var req models.BotGraphData
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	if err := h.BotService.UpdateBotGraph(&req); err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Bot graph updated successfully", 0, ctx)
}

// GetNodeView godocs
//
//	@Summary		Get bot node view
//	@Description	Get the current node content and options based on user message
//	@Tags			Bot
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.BotRequest	true	"Bot request data"
//	@Success		200		{object}	models.BotResponse
//	@Failure		400		{object}	response.BaseError
//	@Failure		500		{object}	response.BaseError
//	@Router			/open/bot [post]
func (h *BotHandler) GetNodeView(ctx *gin.Context) {
	var req models.BotRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(errs.New(errs.BadRequest, "Bad Request", nil))
		return
	}

	var claims *models.ManagedClaims
	if value, exists := ctx.Get("user"); exists {
		if c, ok := value.(*models.ManagedClaims); ok {
			claims = c
		}
	}

	resp, err := h.BotService.HandleSession(req, claims)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.PureJSON(200, resp)
}

// ResetDefault godocs
//
//	@Summary		Reset bot to default
//	@Description	Reset the bot configuration to its default state using the default SQL resource
//	@Tags			Bot
//	@Produce		json
//	@Success		200	{object}	response.TransactionResponse
//	@Failure		401	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/bot/reset [post]
//
//	@Security		ApiKeyAuth
func (h *BotHandler) ResetDefault(ctx *gin.Context) {
	if err := h.BotService.ResetDefault(); err != nil {
		ctx.Error(err)
		return
	}

	response.NewTransactionResponse(200, "Bot reset successfully", 0, ctx)
}
