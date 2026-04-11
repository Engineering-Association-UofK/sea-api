package repositories

import (
	"fmt"
	"sea-api/internal/models"

	"github.com/jmoiron/sqlx"
)

type SuspensionsRepo struct {
	db *sqlx.DB
}

func NewSuspensionsRepo(db *sqlx.DB) *SuspensionsRepo {
	return &SuspensionsRepo{db: db}
}

func (r *SuspensionsRepo) Create(suspension *models.SuspensionModel, tx *sqlx.Tx, isHistory bool) (int64, error) {
	table := models.TableSuspensions
	if isHistory {
		table = models.TableSuspensionHistory
	}
	query := fmt.Sprintf(`
	INSERT INTO %s (user_id, admin_id, reason, started_at, ended_at)
	VALUES (:user_id, :admin_id, :reason, :started_at, :ended_at)
	`, table)
	if tx != nil {
		res, err := tx.NamedExec(query, suspension)
		if err != nil {
			return 0, err
		}
		return res.LastInsertId()
	}
	res, err := r.db.NamedExec(query, suspension)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *SuspensionsRepo) GetByUserID(user_id int64) (*models.SuspensionModel, error) {
	var suspension models.SuspensionModel
	err := r.db.Get(&suspension, fmt.Sprintf(`SELECT * FROM %s WHERE user_id = ?`, models.TableSuspensions), user_id)
	if err != nil {
		return nil, err
	}
	return &suspension, nil
}

func (r *SuspensionsRepo) Delete(id int64) error {
	_, err := r.db.Exec(fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, models.TableSuspensions), id)
	return err
}

func (r *SuspensionsRepo) CleanExpired() ([]int64, error) {
	var ids []int64

	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	err = tx.Select(&ids, fmt.Sprintf(`SELECT id FROM %s WHERE ended_at < NOW() FOR UPDATE`, models.TableSuspensions))
	if err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		err = tx.Commit()
		return ids, err
	}

	query, args, err := sqlx.In(fmt.Sprintf(`DELETE FROM %s WHERE id IN (?)`, models.TableSuspensions), ids)
	if err != nil {
		return nil, err
	}
	query = tx.Rebind(query)

	_, err = tx.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return ids, nil
}
