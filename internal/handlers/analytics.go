package handlers

import (
	"sea-api/internal/services/analytics"

	_ "sea-api/internal/models/analyticsmodel"
	_ "sea-api/internal/response"

	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	service *analytics.Analytics
}

func NewAnalyticsHandler(service *analytics.Analytics) *AnalyticsHandler {
	return &AnalyticsHandler{service: service}
}

// GetGeneralAnalytics godocs
//
//	@Summary		Get General Analytics
//	@Description	Get basic numbers of system components
//	@Tags			Analysis
//	@Produce		json
//	@Success		200	{object}	analyticsmodel.General
//	@Failure		400	{object}	response.BaseError
//	@Failure		500	{object}	response.BaseError
//	@Router			/admin/analysis [get]
//
//	@Security		ApiKeyAuth
func (h *AnalyticsHandler) GetGeneralAnalytics(ctx *gin.Context) {
	post, err := h.service.General()
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.PureJSON(200, post)
}
