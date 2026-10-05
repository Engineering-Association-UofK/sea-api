package userrepo

import (
	"fmt"
	"sea-api/internal/models"
	"sea-api/internal/models/tables"
)

func (r *UserRepository) GetAdmins(req *models.ListRequest) ([]models.AdminRow, error) {
	var admins []models.AdminRow
	offset := (req.Page - 1) * req.Limit
	query := fmt.Sprintf(`
		SELECT 
			u.id, u.email, u.name_ar, u.username, u.gender, ur.role,
			f.file_key AS profile_pic
		FROM %s u
		JOIN %s ur ON u.id = ur.user_id
		LEFT JOIN %s f ON u.profile_image_id = f.id
		WHERE u.id IN (
			SELECT user_id FROM %s WHERE role = ?
		)
		ORDER BY u.id DESC
		LIMIT ? OFFSET ?
	`, tables.Users, tables.UserRoles, tables.Files, tables.UserRoles)

	err := r.DB.Select(&admins, query, models.RoleSystemAdmin, req.Limit, offset)
	if err != nil {
		return nil, err
	}
	return admins, nil
}

func (r *UserRepository) GetAdminsCount() (int64, error) {
	var count int64
	err := r.DB.Get(&count, fmt.Sprintf(`
		SELECT COUNT(*) FROM %s u 
		JOIN %s ur ON u.id = ur.user_id 
		WHERE ur.role = ?
	`, tables.Users, tables.UserRoles), models.RoleSystemAdmin)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *UserRepository) AddAdmin(id int64) error {
	_, err := r.DB.Exec(fmt.Sprintf(`INSERT INTO %s (user_id, role) VALUES (?, ?)`, tables.UserRoles), id, models.RoleSystemAdmin)
	return err
}

func (r *UserRepository) RemoveAdmin(id int64) error {
	_, err := r.DB.Exec(fmt.Sprintf(`DELETE FROM %s WHERE user_id = ? AND role = ?`, tables.UserRoles), id, models.RoleSystemAdmin)
	return err
}
