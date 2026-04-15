package valid

import "sea-api/internal/models"

func ValidateListRequest(req *models.ListRequest, total int) {
	if !models.AllowedListLimit[req.Limit] {
		req.Limit = 10
	}
	totalPages := total / req.Limit
	if total%req.Limit != 0 {
		totalPages++
	}
	if req.Page > totalPages {
		req.Page = totalPages
	}
}
