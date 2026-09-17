package analytics

import (
	"sea-api/internal/models/analyticsmodel"
	"sea-api/internal/repositories/analyticsrepo"

	"github.com/jmoiron/sqlx"
)

type Analytics struct {
	repo analyticsrepo.AnalyticsRepo
}

func NewAnalytics(db sqlx.DB) Analytics {
	return Analytics{
		repo: analyticsrepo.NewAnalyticsRepo(db),
	}
}

func (s *Analytics) General() (*analyticsmodel.General, error) {
	return s.repo.General()
}
