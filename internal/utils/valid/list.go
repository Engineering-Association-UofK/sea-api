package valid

import "sea-api/internal/models"

func Limit(req *models.ListRequest, total int64) int64 {
	if !models.AllowedListLimit[req.Limit] {
		req.Limit = 10
	}

	if req.Page < 1 {
		req.Page = 1
	}

	var numPages int64 = 1
	if total > 0 {
		numPages = total / req.Limit
		if total%req.Limit != 0 {
			numPages++
		}

		if req.Page > numPages {
			req.Page = numPages
		}
	}

	return numPages
}
