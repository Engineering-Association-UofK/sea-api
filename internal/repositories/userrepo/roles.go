package userrepo

import (
	"fmt"
	"sea-api/internal/models"
	"sea-api/internal/models/tables"

	"github.com/jmoiron/sqlx"
)

func (r *UserRepository) GetAllRolesByUserIDs(ids []int64) ([]models.UserRole, error) {
	var roles []models.UserRole
	if len(ids) == 0 {
		return roles, nil
	}

	query, args, err := sqlx.In(fmt.Sprintf(`SELECT * FROM %s WHERE user_id IN (?)`, tables.UserRoles), ids)
	if err != nil {
		return nil, err
	}

	query = r.DB.Rebind(query)
	err = r.DB.Select(&roles, query, args...)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *UserRepository) GetRolesByUserID(id int64) ([]models.UserRole, error) {
	var roles []models.UserRole
	err := r.DB.Select(&roles, fmt.Sprintf(`SELECT * FROM %s WHERE user_id = ?`, tables.UserRoles), id)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *UserRepository) CreateRole(role *models.UserRole) error {
	_, err := r.DB.NamedExec(fmt.Sprintf(`INSERT INTO %s (user_id, role) VALUES (:user_id, :role)`, tables.UserRoles), role)
	return err
}

func (r *UserRepository) UpdateRole(role *models.UserRole, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`UPDATE %s SET role = :role WHERE user_id = :id`, tables.UserRoles)
	if tx != nil {
		_, err := tx.NamedExec(query, role)
		return err
	}
	_, err := r.DB.NamedExec(query, role)
	return err
}

func (r *UserRepository) RemoveRole(id int64, role models.Role, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE user_id = ? AND role = ?`, tables.UserRoles)
	if tx != nil {
		_, err := tx.Exec(query, id, role)
		return err
	}
	_, err := r.DB.Exec(query, id, role)
	return err
}

func (r *UserRepository) ReplaceRoles(id int64, roles []models.Role, tx *sqlx.Tx) error {
	deleteQuery := fmt.Sprintf(`DELETE FROM %s WHERE user_id = ?`, tables.UserRoles)
	insertQuery := fmt.Sprintf(`INSERT INTO %s (user_id, role) VALUES (?, ?)`, tables.UserRoles)

	if tx != nil {
		if _, err := tx.Exec(deleteQuery, id); err != nil {
			return err
		}
		for _, role := range roles {
			if _, err := tx.Exec(insertQuery, id, role); err != nil {
				return err
			}
		}
		return nil
	}

	newTx, err := r.DB.Beginx()
	if err != nil {
		return err
	}
	defer newTx.Rollback()

	if _, err := newTx.Exec(deleteQuery, id); err != nil {
		return err
	}
	for _, role := range roles {
		if _, err := newTx.Exec(insertQuery, id, role); err != nil {
			return err
		}
	}
	return newTx.Commit()
}

func (r *UserRepository) DeleteRole(id int64, role models.Role, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE user_id = ? AND role = ?`, tables.UserRoles)
	if tx != nil {
		_, err := tx.Exec(query, id, role)
		return err
	}
	_, err := r.DB.Exec(query, id, role)
	return err
}

func (r *UserRepository) DeleteRolesByUserID(id int64, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE user_id = ?`, tables.UserRoles)
	if tx != nil {
		_, err := tx.Exec(query, id)
		return err
	}
	_, err := r.DB.Exec(query, id)
	return err
}
