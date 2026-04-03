package repositories

import (
	"sea-api/internal/models"

	"github.com/jmoiron/sqlx"
)

type RateLimitRepository struct {
	db *sqlx.DB
}

func NewRateLimitRepository(db *sqlx.DB) *RateLimitRepository {
	return &RateLimitRepository{db: db}
}

func (r *RateLimitRepository) GetRateLimit(ip string, endpoint models.RateLimitEndpoints) (*models.RateLimitModel, error) {
	var limit models.RateLimitModel
	query := `SELECT * FROM rate_limits WHERE ip_address = ? AND endpoint = ?`
	err := r.db.Get(&limit, query, ip, endpoint)
	if err != nil {
		return nil, err
	}
	return &limit, nil
}

func (r *RateLimitRepository) UpsertRateLimit(limit *models.RateLimitModel) error {
	query := `
	INSERT INTO rate_limits (ip_address, endpoint, request_count, last_request)
	VALUES (:ip_address, :endpoint, :request_count, :last_request)
	ON DUPLICATE KEY UPDATE 
		request_count = VALUES(request_count),
		last_request = VALUES(last_request)
	`
	_, err := r.db.NamedExec(query, limit)
	return err
}

func (r *RateLimitRepository) ResetRateLimit(ip string, endpoint models.RateLimitEndpoints) error {
	query := `DELETE FROM rate_limits WHERE ip_address = ? AND endpoint = ?`
	_, err := r.db.Exec(query, ip, endpoint)
	return err
}

func (r *RateLimitRepository) CleanExpired(endpoint models.RateLimitEndpoints, multiplier int) error {
	query := `
	DELETE FROM rate_limits 
	WHERE endpoint = ? 
	  AND last_request < NOW() - (POW(2, request_count - 1) * ? * INTERVAL '1 minute');
	`
	_, err := r.db.Exec(query, endpoint, multiplier)
	return err
}
