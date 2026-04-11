package repositories

import (
	"fmt"
	"sea-api/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(DB *sqlx.DB) *UserRepository {
	return &UserRepository{DB: DB}
}

// ======== GET ALL ========

func (r *UserRepository) GetAll(limit int, page int) ([]models.UserModel, error) {
	var users []models.UserModel
	offset := (page - 1) * limit
	err := r.DB.Select(&users, fmt.Sprintf(`
		SELECT * FROM %s 
		WHERE is_anonymous = false
		LIMIT ? OFFSET ? 
	`, models.TableUsers), limit, offset)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetAllTempUsers(limit int, page int) ([]models.TempUserModel, error) {
	var users []models.TempUserModel
	offset := (page - 1) * limit
	err := r.DB.Select(&users, fmt.Sprintf(`SELECT * FROM %s LIMIT ? OFFSET ?`, models.TableTempUsers), limit, offset)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetTempUsersWithNullPasswords() ([]models.TempUserModel, error) {
	var users []models.TempUserModel
	err := r.DB.Select(&users, fmt.Sprintf(`SELECT * FROM %s WHERE password IS NULL OR password = ''`, models.TableTempUsers))
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetPagesCount(limit int, isTempUser bool) (int, error) {
	table := models.TableUsers
	if isTempUser {
		table = models.TableTempUsers
	}
	var count int
	err := r.DB.Get(&count, fmt.Sprintf(`SELECT COUNT(*) FROM %s`, table))
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

	query, args, err := sqlx.In(fmt.Sprintf(`SELECT * FROM %s WHERE id IN (?)`, models.TableUsers), indices)
	if err != nil {
		return nil, err
	}

	query = r.DB.Rebind(query)

	err = r.DB.Select(&users, query, args...)
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

	query, args, err := sqlx.In(fmt.Sprintf(`SELECT * FROM %s WHERE user_id IN (?)`, models.TableUserRoles), ids)
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
	err := r.DB.Select(&roles, fmt.Sprintf(`SELECT * FROM %s WHERE user_id = ?`, models.TableUserRoles), id)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *UserRepository) GetAdmins() ([]models.UserModel, error) {
	var users []models.UserModel
	err := r.DB.Select(&users, fmt.Sprintf(`
		SELECT u.* FROM %s u
		JOIN %s ur ON u.id = ur.user_id
		WHERE ur.role = ?
	`, models.TableUserRoles, models.TableUserRoles), models.RoleSystemAdmin)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// ======== GET ========

func (r *UserRepository) GetByUserID(id int64) (*models.UserModel, error) {
	var user models.UserModel
	err := r.DB.Get(&user, fmt.Sprintf(`SELECT * FROM %s WHERE id = ? AND is_anonymous = false`, models.TableUsers), id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*models.UserModel, error) {
	var user models.UserModel
	err := r.DB.Get(&user, fmt.Sprintf(`SELECT * FROM %s WHERE username = ? AND is_anonymous = false`, models.TableUsers), username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) IsUsernameAvailable(username string) (bool, error) {
	var count int
	err := r.DB.Get(&count, fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE username = ?`, models.TableUsers), username)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.UserModel, error) {
	var user models.UserModel
	err := r.DB.Get(&user, fmt.Sprintf(`SELECT * FROM %s WHERE email = ? AND is_anonymous = false`, models.TableUsers), email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUniID(uniID int64) (*models.UserModel, error) {
	var user models.UserModel
	err := r.DB.Get(&user, fmt.Sprintf(`SELECT * FROM %s WHERE uni_id = ? AND is_anonymous = false`, models.TableUsers), uniID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetTempUser(id int64) (*models.TempUserModel, error) {
	var user models.TempUserModel
	err := r.DB.Get(&user, fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, models.TableTempUsers), id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ======= CREATE ========

func (r *UserRepository) DetailedCreate(user *models.UserModel, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`
	INSERT INTO %s (
		id, uni_id, username, profile_image_id, name_ar, name_en,
		email, phone, password, verified, status,
		department, gender, is_editable, is_loggable, is_anonymous
	) VALUES (
		:id, :uni_id, :username, :profile_image_id, :name_ar, :name_en,
		:email, :phone, :password, :verified, :status,
		:department, :gender, :is_editable, :is_loggable, :is_anonymous
	)`, models.TableUsers)
	if tx != nil {
		_, err := tx.NamedExec(query, user)
		return err
	}
	_, err := r.DB.NamedExec(query, user)
	return err
}

func (r *UserRepository) Create(user *models.UserModel, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`
	INSERT INTO %s (
		id, uni_id, username, profile_image_id, name_ar, name_en,
		email, phone, password, verified, status,
		department, gender
	) VALUES (
		:id, :uni_id, :username, :profile_image_id, :name_ar, :name_en,
		:email, :phone, :password, :verified, :status,
		:department, :gender
	)`, models.TableUsers)
	if tx != nil {
		_, err := tx.NamedExec(query, user)
		return err
	}
	_, err := r.DB.NamedExec(query, user)
	return err
}

func (r *UserRepository) CreateRole(role *models.UserRole) error {
	_, err := r.DB.NamedExec(fmt.Sprintf(`INSERT INTO %s (user_id, role) VALUES (:user_id, :role)`, models.TableUserRoles), role)
	return err
}

// ======== UPDATE ========

func (r *UserRepository) Update(user *models.UserModel, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET uni_id = :uni_id, username = :username, profile_image_id = :profile_image_id, name_ar = :name_ar, name_en = :name_en,
	email = :email, phone = :phone, password = :password, verified = :verified,
	status = :status, department = :department, gender = :gender
	WHERE id = :id
	AND is_editable = true
	`, models.TableUsers)

	if tx != nil {
		_, err := tx.NamedExec(query, user)
		return err
	}
	_, err := r.DB.NamedExec(query, user)
	return err
}

// Update using email as key, updating the id as well
func (r *UserRepository) UpdateWithID(user *models.UserModel, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET id = :id, uni_id = :uni_id, username = :username, profile_image_id = :profile_image_id, name_ar = :name_ar, name_en = :name_en,
	phone = :phone, password = :password, verified = :verified,
	status = :status, department = :department, gender = :gender
	WHERE email = :email
	AND is_editable = true
	`, models.TableUsers)

	if tx != nil {
		_, err := tx.NamedExec(query, user)
		return err
	}
	_, err := r.DB.NamedExec(query, user)
	return err
}

func (r *UserRepository) UpdateTempPasscode(id int64, passcode string) error {
	_, err := r.DB.Exec(fmt.Sprintf(`UPDATE %s SET password = ? WHERE id = ?`, models.TableTempUsers), passcode, id)
	return err
}

func (r *UserRepository) UpdateRole(role *models.UserRole, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`UPDATE %s SET role = :role WHERE user_id = :id`, models.TableUserRoles)
	if tx != nil {
		_, err := tx.NamedExec(query, role)
		return err
	}
	_, err := r.DB.NamedExec(query, role)
	return err
}

func (r *UserRepository) RemoveRole(id int64, role models.Role, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE user_id = ? AND role = ?`, models.TableUserRoles)
	if tx != nil {
		_, err := tx.Exec(query, id, role)
		return err
	}
	_, err := r.DB.Exec(query, id, role)
	return err
}

func (r *UserRepository) ReplaceRoles(id int64, roles []models.Role, tx *sqlx.Tx) error {
	deleteQuery := fmt.Sprintf(`DELETE FROM %s WHERE user_id = ?`, models.TableUserRoles)
	insertQuery := fmt.Sprintf(`INSERT INTO %s (user_id, role) VALUES (?, ?)`, models.TableUserRoles)

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

func (r *UserRepository) AddAdmin(id int64) error {
	_, err := r.DB.Exec(fmt.Sprintf(`INSERT INTO %s (user_id, role) VALUES (?, ?)`, models.TableUserRoles), id, models.RoleSystemAdmin)
	return err
}

func (r *UserRepository) RemoveAdmin(id int64) error {
	_, err := r.DB.Exec(fmt.Sprintf(`DELETE FROM %s WHERE user_id = ? AND role = ?`, models.TableUserRoles), id, models.RoleSystemAdmin)
	return err
}

// ======== DELETE ========

func (r *UserRepository) Delete(id int64, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = ? AND is_editable = true`, models.TableUsers)
	if tx != nil {
		_, err := tx.Exec(query, id)
		return err
	}
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *UserRepository) DeleteTempUser(id int64, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, models.TableTempUsers)
	if tx != nil {
		_, err := tx.Exec(query, id)
		return err
	}
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *UserRepository) DeleteRole(id int64, role models.Role, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE user_id = ? AND role = ?`, models.TableUserRoles)
	if tx != nil {
		_, err := tx.Exec(query, id, role)
		return err
	}
	_, err := r.DB.Exec(query, id, role)
	return err
}

func (r *UserRepository) DeleteRolesByUserID(id int64, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE user_id = ?`, models.TableUserRoles)
	if tx != nil {
		_, err := tx.Exec(query, id)
		return err
	}
	_, err := r.DB.Exec(query, id)
	return err
}

// ======== SPECIAL ========

func (r *UserRepository) Verify(id int64) error {
	_, err := r.DB.Exec(fmt.Sprintf(`UPDATE %s SET verified = ? WHERE id = ?`, models.TableUsers), true, id)
	return err
}

func (r *UserRepository) Suspend(id int64, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`UPDATE %s SET status = ? WHERE id = ?`, models.TableUsers)
	if tx != nil {
		_, err := tx.Exec(query, models.STATUS_SUSPENDED, id)
		return err
	}
	_, err := r.DB.Exec(query, models.STATUS_SUSPENDED, id)
	return err
}

func (r *UserRepository) Activate(id int64) error {
	_, err := r.DB.Exec(fmt.Sprintf(`UPDATE %s SET status = ? WHERE id = ?`, models.TableUsers), models.STATUS_ACTIVE, id)
	return err
}

func (r *UserRepository) RemoveSuspensionState(id int64) error {
	_, err := r.DB.Exec(fmt.Sprintf(`UPDATE %s SET status = ? WHERE id = ? AND status = ?`, models.TableUsers), models.STATUS_ACTIVE, id, models.STATUS_SUSPENDED)
	return err
}

func (r *UserRepository) BeginTransaction() (*sqlx.Tx, error) {
	return r.DB.Beginx()
}
