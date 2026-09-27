package analytics

import (
	"sea-api/internal/models/analyticsmodel"
	"sea-api/internal/repositories/analyticsrepo"
)

type Analytics struct {
	repo *analyticsrepo.AnalyticsRepo
}

func NewAnalytics(repo *analyticsrepo.AnalyticsRepo) *Analytics {
	return &Analytics{
		repo: repo,
	}
}

func (s *Analytics) General() (*analyticsmodel.General, error) {
	return s.repo.General()
}
