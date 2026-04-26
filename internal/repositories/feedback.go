package repositories

import (
	"fmt"
	"sea-api/internal/models"

	"github.com/jmoiron/sqlx"
)

type FeedbackRepository struct {
	db *sqlx.DB
}

func NewFeedbackRepository(db *sqlx.DB) *FeedbackRepository {
	return &FeedbackRepository{
		db: db,
	}
}

func (r *FeedbackRepository) Create(feedback *models.Feedback) (int64, error) {
	query := fmt.Sprintf(`
		INSERT INTO %s (message, user_id, type, created_at)
		VALUES (:message, :user_id, :type, :created_at)
	`, models.TableFeedback)
	res, err := r.db.NamedExec(query, feedback)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *FeedbackRepository) GetByID(id int64) (*models.Feedback, error) {
	var feedback models.Feedback
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, models.TableFeedback)
	err := r.db.Get(&feedback, query, id)
	if err != nil {
		return nil, err
	}
	return &feedback, nil
}

func (r *FeedbackRepository) GetAll(req *models.ListRequest) ([]models.Feedback, error) {
	var feedbacks []models.Feedback
	offset := (req.Page - 1) * req.Limit
	query := fmt.Sprintf(`SELECT * FROM %s ORDER BY created_at DESC LIMIT ? OFFSET ?`, models.TableFeedback)
	err := r.db.Select(&feedbacks, query, req.Limit, offset)
	if err != nil {
		return nil, err
	}
	return feedbacks, nil
}

func (r *FeedbackRepository) GetByType(fType models.FeedbackType, req *models.ListRequest) ([]models.Feedback, error) {
	var feedbacks []models.Feedback
	offset := (req.Page - 1) * req.Limit
	query := fmt.Sprintf(`SELECT * FROM %s WHERE type = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`, models.TableFeedback)
	err := r.db.Select(&feedbacks, query, fType, req.Limit, offset)
	if err != nil {
		return nil, err
	}
	return feedbacks, nil
}

func (r *FeedbackRepository) Delete(id int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, models.TableFeedback)
	_, err := r.db.Exec(query, id)
	return err
}
