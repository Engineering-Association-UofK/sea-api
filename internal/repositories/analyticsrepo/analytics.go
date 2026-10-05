package analyticsrepo

import (
	"fmt"
	"sea-api/internal/models/analyticsmodel"
	"sea-api/internal/models/tables"

	"github.com/jmoiron/sqlx"
)

type AnalyticsRepo struct {
	db *sqlx.DB
}

func NewAnalyticsRepo(db *sqlx.DB) *AnalyticsRepo {
	return &AnalyticsRepo{db: db}
}

func (r *AnalyticsRepo) General() (*analyticsmodel.General, error) {
	var stats = analyticsmodel.General{}

	query := fmt.Sprintf(`
		SELECT 
			(SELECT COUNT(*) FROM %s) AS users,
			(SELECT COUNT(*) FROM %s) AS events,
			(SELECT COUNT(*) FROM %s) AS certificates,
			(SELECT COUNT(*) FROM %s) AS forms,
			(SELECT COUNT(*) FROM %s) AS posts
	`,
		tables.Users,
		tables.NewEvents,
		tables.NewCertificates,
		tables.Forms,
		tables.Posts,
	)

	err := r.db.Get(&stats, query)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
