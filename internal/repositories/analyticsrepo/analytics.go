package analyticsrepo

import (
	"fmt"
	"sea-api/internal/models"
	"sea-api/internal/models/analyticsmodel"

	"github.com/jmoiron/sqlx"
)

type AnalyticsRepo struct {
	db sqlx.DB
}

func NewAnalyticsRepo(db sqlx.DB) AnalyticsRepo {
	return AnalyticsRepo{db: db}
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
		models.TableUsers,
		models.TableNewEvents,
		models.TableNewCertificates,
		models.TableForms,
		models.TablePosts,
	)

	err := r.db.Get(&stats, query)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
