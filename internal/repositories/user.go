package repositories

import (
	"fmt"
	"sea-api/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// ======== GET ALL ========

func (r *UserRepository) GetAll(limit int, page int) ([]models.UserModel, error) {
	var users []models.UserModel
	offset := (page - 1) * limit
	err := r.db.Select(&users, `SELECT * FROM users LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetAllTempUsers(limit int, page int) ([]models.TempUserModel, error) {
	var users []models.TempUserModel
	offset := (page - 1) * limit
	err := r.db.Select(&users, `SELECT * FROM users_temp LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetTempUsersWithNullPasswords() ([]models.TempUserModel, error) {
	var users []models.TempUserModel
	err := r.db.Select(&users, `SELECT * FROM users_temp WHERE password IS NULL OR password = ''`)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetPagesCount(limit int, isTempUser bool) (int, error) {
	table := "users"
	if isTempUser {
		table = "users_temp"
	}
	var count int
	err := r.db.Get(&count, fmt.Sprintf(`SELECT COUNT(*) FROM %s`, table))
	if err != nil {
		return 0, err
	}
	pages := (count + limit - 1) / limit
	return pages, nil
}

func (r *UserRepository) GetAllByIndices(indices []int64) ([]models.UserModel, error) {
	var users []models.UserModel

	if len(indices) == 0 {
		return users, nil
	}

	query, args, err := sqlx.In(`SELECT * FROM users WHERE id IN (?)`, indices)
	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)

	err = r.db.Select(&users, query, args...)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) GetAllRolesByUserIDs(ids []int64) ([]models.UserRole, error) {
	var roles []models.UserRole
	if len(ids) == 0 {
		return roles, nil
	}

	query, args, err := sqlx.In(`SELECT * FROM user_roles WHERE user_id IN (?)`, ids)
	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)
	err = r.db.Select(&roles, query, args...)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *UserRepository) GetRolesByUserID(id int64) ([]models.UserRole, error) {
	var roles []models.UserRole
	err := r.db.Select(&roles, `SELECT role FROM user_roles WHERE user_id = ?`, id)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *UserRepository) GetAdmins() ([]models.UserModel, error) {
	var users []models.UserModel
	err := r.db.Select(&users, `
		SELECT u.* FROM users u
		JOIN user_roles ur ON u.id = ur.user_id
		WHERE ur.role = ?
	`, models.ROLE_ADMIN)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// ======== GET ========

func (r *UserRepository) GetByUserID(id int64) (*models.UserModel, error) {
	var user models.UserModel
	err := r.db.Get(&user, `SELECT * FROM users WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*models.UserModel, error) {
	var user models.UserModel
	err := r.db.Get(&user, `SELECT * FROM users WHERE username = ?`, username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.UserModel, error) {
	var user models.UserModel
	err := r.db.Get(&user, `SELECT * FROM users WHERE email = ?`, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUniID(uniID int64) (*models.UserModel, error) {
	var user models.UserModel
	err := r.db.Get(&user, `SELECT * FROM users WHERE uni_id = ?`, uniID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetTempUser(id int64) (*models.TempUserModel, error) {
	var user models.TempUserModel
	err := r.db.Get(&user, `SELECT * FROM users_temp WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ======= CREATE ========

func (r *UserRepository) Create(user *models.UserModel) error {
	query := `
	INSERT INTO users (
		id, uni_id, username, profile_image_id, name_ar, name_en,
		email, phone, password, verified, status,
		department, gender
	) VALUES (
		:id, :uni_id, :username, :profile_image_id, :name_ar, :name_en,
		:email, :phone, :password, :verified, :status,
		:department, :gender
	)`
	_, err := r.db.NamedExec(query, user)
	return err
}

func (r *UserRepository) CreateRole(role *models.UserRole) error {
	_, err := r.db.NamedExec(`INSERT INTO user_roles (id, role) VALUES (:id, :role)`, role)
	return err
}

// ======== UPDATE ========

func (r *UserRepository) Update(user *models.UserModel, tx *sqlx.Tx) error {
	query := `
	UPDATE users
	SET uni_id = :uni_id, username = :username, profile_image_id = :profile_image_id, name_ar = :name_ar, name_en = :name_en,
	email = :email, phone = :phone, password = :password, verified = :verified,
	status = :status, department = :department, gender = :gender
	WHERE id = :id`

	if tx != nil {
		_, err := tx.NamedExec(query, user)
		return err
	}
	_, err := r.db.NamedExec(query, user)
	return err
}

func (r *UserRepository) UpdateTempPasscode(id int64, passcode string) error {
	_, err := r.db.Exec(`UPDATE users_temp SET password = ? WHERE id = ?`, passcode, id)
	return err
}

func (r *UserRepository) UpdateRole(role *models.UserRole, tx *sqlx.Tx) error {
	query := `UPDATE user_roles SET role = :role WHERE id = :id`
	if tx != nil {
		_, err := tx.NamedExec(query, role)
		return err
	}
	_, err := r.db.NamedExec(query, role)
	return err
}

func (r *UserRepository) RemoveRole(id int64, role models.Role, tx *sqlx.Tx) error {
	query := `DELETE FROM user_roles WHERE id = ? AND role = ?`
	if tx != nil {
		_, err := tx.Exec(query, id, role)
		return err
	}
	_, err := r.db.Exec(query, id, role)
	return err
}

func (r *UserRepository) ReplaceRoles(id int64, roles []models.Role, tx *sqlx.Tx) error {
	deleteQuery := `DELETE FROM user_roles WHERE user_id = ?`
	insertQuery := `INSERT INTO user_roles (user_id, role) VALUES (?, ?)`

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

	newTx, err := r.db.Beginx()
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

func (r *UserRepository) AddAdmin(id int64) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`INSERT INTO user_roles (user_id, role) VALUES (?, ?)`, id, models.ROLE_ADMIN)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *UserRepository) RemoveAdmin(id int64) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM user_roles WHERE user_id = ? AND role IN (?)`, id, models.AdminRoles)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// ======== DELETE ========

func (r *UserRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id = ?`, id)
	return err
}

func (r *UserRepository) DeleteTempUser(id int64) error {
	_, err := r.db.Exec(`DELETE FROM users_temp WHERE id = ?`, id)
	return err
}

func (r *UserRepository) DeleteRole(id int64, role models.Role) error {
	_, err := r.db.Exec(`DELETE FROM user_roles WHERE user_id = ? AND role = ?`, id, role)
	return err
}

func (r *UserRepository) DeleteRolesByUserID(id int64) error {
	_, err := r.db.Exec(`DELETE FROM user_roles WHERE user_id = ?`, id)
	return err
}

// ======== SPECIAL ========

func (r *UserRepository) Verify(id int64) error {
	_, err := r.db.Exec(`UPDATE users SET verified = ? WHERE id = ?`, true, id)
	return err
}

func (r *UserRepository) Suspend(id int64, tx *sqlx.Tx) error {
	query := `UPDATE users SET status = ? WHERE id = ?`
	if tx != nil {
		_, err := tx.Exec(query, models.STATUS_SUSPENDED, id)
		return err
	}
	_, err := r.db.Exec(query, models.STATUS_SUSPENDED, id)
	return err
}

func (r *UserRepository) Activate(id int64) error {
	_, err := r.db.Exec(`UPDATE users SET status = ? WHERE id = ?`, models.STATUS_ACTIVE, id)
	return err
}

func (r *UserRepository) RemoveSuspensionState(id int64) error {
	_, err := r.db.Exec(`UPDATE users SET status = ? WHERE id = ? AND status = ?`, models.STATUS_ACTIVE, id, models.STATUS_SUSPENDED)
	return err
}

func (r *UserRepository) BeginTransaction() (*sqlx.Tx, error) {
	return r.db.Beginx()
}
