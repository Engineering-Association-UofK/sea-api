package repositories

import (
	"fmt"
	"sea-api/internal/models"

	"github.com/jmoiron/sqlx"
)

type NotificationRepository struct {
	db *sqlx.DB
}

func NewNotificationRepository(db *sqlx.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(notification *models.Notification) (int64, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (user_id, title, message, type, data, created_at, is_read)
	VALUES (:user_id, :title, :message, :type, :data, :created_at, :is_read)
	`, models.TableNotifications)
	res, err := r.db.NamedExec(query, notification)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *NotificationRepository) GetByUserIDWithLimit(userID int64, limit models.ListRequest) ([]models.Notification, error) {
	offset := (limit.Page - 1) * limit.Limit
	query := fmt.Sprintf(`
	SELECT * FROM %s 
	WHERE user_id = ? 
	ORDER BY created_at DESC
	LIMIT ? OFFSET ? 
	`, models.TableNotifications)
	var notifications = []models.Notification{}
	err := r.db.Select(&notifications, query, userID, limit.Limit, offset)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *NotificationRepository) GetTotalWithUserID(userID int64) (int, error) {
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE user_id = ?`, models.TableNotifications)
	var count int
	err := r.db.Get(&count, query, userID)
	return count, err
}

func (r *NotificationRepository) MarkAsRead(userId, id int64) (int64, error) {
	query := fmt.Sprintf(`UPDATE %s SET is_read = true WHERE id = ? AND user_id = ?`, models.TableNotifications)
	res, err := r.db.Exec(query, id, userId)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (r *NotificationRepository) MarkAllAsRead(userID int64) error {
	query := fmt.Sprintf(`UPDATE %s SET is_read = true WHERE user_id = ?`, models.TableNotifications)
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *NotificationRepository) Delete(userId, id int64) (int64, error) {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = ? AND user_id = ?`, models.TableNotifications)
	res, err := r.db.Exec(query, id, userId)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (r *NotificationRepository) GetUnreadCount(userID int64) (int, error) {
	var count int
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE user_id = ? AND is_read = false`, models.TableNotifications)
	err := r.db.Get(&count, query, userID)
	return count, err
}
