package notificationsrepo

import (
	"fmt"
	"sea-api/internal/models"
	"sea-api/internal/models/notificationsmodels"
	"sea-api/internal/models/tables"
	"time"

	"github.com/jmoiron/sqlx"
)

type NotificationRepository struct {
	db *sqlx.DB
}

func NewNotificationRepository(db *sqlx.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(notification *notificationsmodels.Notification) (int64, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (user_id, title, message, type, data, created_at, is_read)
	VALUES (:user_id, :title, :message, :type, :data, :created_at, :is_read)
	`, tables.Notifications)
	res, err := r.db.NamedExec(query, notification)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *NotificationRepository) BulkCreateForUsers(req *notificationsmodels.NotificationInsertRow) error {
	rawQuery := fmt.Sprintf(`
	INSERT INTO %s (user_id, title, message, type, data, created_at, is_read)
	SELECT id, ?, ?, ?, ?, ?, false
	FROM %s
	`, tables.Notifications, tables.Users)

	var err error
	if req.IDs != nil {
		rawQuery += ` WHERE id IN (?)`

		query, args, err := sqlx.In(rawQuery, req.Title, req.Message, req.Type, req.Data, time.Now(), req.IDs)
		if err != nil {
			return err
		}
		query = r.db.Rebind(query)
		_, err = r.db.Exec(query, args...)
	} else {
		_, err = r.db.Exec(rawQuery, req.Title, req.Message, req.Type, req.Data, time.Now())
	}

	return err
}

func (r *NotificationRepository) GetByUserIDWithLimit(userID int64, limit models.ListRequest) ([]notificationsmodels.Notification, error) {
	offset := (limit.Page - 1) * limit.Limit
	query := fmt.Sprintf(`
	SELECT * FROM %s 
	WHERE user_id = ? 
	ORDER BY created_at DESC
	LIMIT ? OFFSET ? 
	`, tables.Notifications)
	var notifications = []notificationsmodels.Notification{}
	err := r.db.Select(&notifications, query, userID, limit.Limit, offset)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *NotificationRepository) GetTotalWithUserID(userID int64) (int64, error) {
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE user_id = ?`, tables.Notifications)
	var count int64
	err := r.db.Get(&count, query, userID)
	return count, err
}

func (r *NotificationRepository) MarkAsRead(userId, id int64) (int64, error) {
	query := fmt.Sprintf(`UPDATE %s SET is_read = true WHERE id = ? AND user_id = ?`, tables.Notifications)
	res, err := r.db.Exec(query, id, userId)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (r *NotificationRepository) MarkAllAsRead(userID int64) error {
	query := fmt.Sprintf(`UPDATE %s SET is_read = true WHERE user_id = ?`, tables.Notifications)
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *NotificationRepository) Delete(userId, id int64) (int64, error) {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = ? AND user_id = ?`, tables.Notifications)
	res, err := r.db.Exec(query, id, userId)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (r *NotificationRepository) GetUnreadCount(userID int64) (int, error) {
	var count int
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE user_id = ? AND is_read = false`, tables.Notifications)
	err := r.db.Get(&count, query, userID)
	return count, err
}
