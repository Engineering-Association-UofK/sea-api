package userrepo

import (
	"fmt"
	"sea-api/internal/models"
	"sea-api/internal/models/authmodels"
	"sea-api/internal/models/tables"
	"sea-api/internal/models/usermodels"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(DB *sqlx.DB) *UserRepository {
	return &UserRepository{DB: DB}
}

// ======== GET ALL ========

func (r *UserRepository) GetUserRows(req *models.ListRequest) ([]usermodels.UserRow, error) {
	var users []usermodels.UserRow
	offset := (req.Page - 1) * req.Limit
	query := fmt.Sprintf(`
		SELECT 
			u.id, u.uni_id, u.username, u.profile_image_id, u.name_ar, u.name_en, 
			u.email, u.phone, u.department, u.gender, u.verified, u.status,
			f.file_key AS profile_pic_key
		FROM %s u
		LEFT JOIN %s f ON u.profile_image_id = f.id
		WHERE u.is_anonymous = false
		ORDER BY u.id DESC
		LIMIT ? OFFSET ?
	`, tables.Users, tables.Files)

	err := r.DB.Select(&users, query, req.Limit, offset)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetAll(limit int64, page int64) ([]usermodels.UserModel, error) {
	var users []usermodels.UserModel
	offset := (page - 1) * limit
	err := r.DB.Select(&users, fmt.Sprintf(`
		SELECT * FROM %s 
		WHERE is_anonymous = false
		AND username IS NOT NULL
		LIMIT ? OFFSET ? 
	`, tables.Users), limit, offset)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetTotal() (int64, error) {
	var count int64
	err := r.DB.Get(&count, fmt.Sprintf(`SELECT COUNT(*) FROM %s`, tables.Users))
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *UserRepository) GetAllByIndices(indices []int64) ([]usermodels.UserModel, error) {
	var users []usermodels.UserModel

	if len(indices) == 0 {
		return users, nil
	}

	query, args, err := sqlx.In(fmt.Sprintf(`SELECT * FROM %s WHERE id IN (?)`, tables.Users), indices)
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

// ======== GET ========

func (r *UserRepository) GetUserRow(id int64) (*usermodels.UserRow, error) {
	var user usermodels.UserRow
	query := fmt.Sprintf(`
		SELECT 
			u.id, u.uni_id, u.username, u.profile_image_id, u.name_ar, u.name_en, 
			u.email, u.phone, u.department, u.gender, u.verified, u.status,
			f.file_key AS profile_pic_key
		FROM %s u
		LEFT JOIN %s f ON u.profile_image_id = f.id
		WHERE u.id = ? AND u.is_anonymous = false
	`, tables.Users, tables.Files)
	err := r.DB.Get(&user, query, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUserID(id int64) (*usermodels.UserModel, error) {
	var user usermodels.UserModel
	err := r.DB.Get(&user, fmt.Sprintf(`SELECT * FROM %s WHERE id = ? AND is_anonymous = false`, tables.Users), id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*usermodels.UserModel, error) {
	var user usermodels.UserModel
	err := r.DB.Get(&user, fmt.Sprintf(`SELECT * FROM %s WHERE username = ? AND is_anonymous = false`, tables.Users), username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) IsUsernameAvailable(username string) (bool, error) {
	var count int
	err := r.DB.Get(&count, fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE username = ?`, tables.Users), username)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func (r *UserRepository) GetByEmail(email string) (*usermodels.UserModel, error) {
	var user usermodels.UserModel
	err := r.DB.Get(&user, fmt.Sprintf(`SELECT * FROM %s WHERE email = ? AND is_anonymous = false`, tables.Users), email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUniID(uniID string) (*usermodels.UserModel, error) {
	var user usermodels.UserModel
	err := r.DB.Get(&user, fmt.Sprintf(`SELECT * FROM %s WHERE uni_id = ? AND is_anonymous = false`, tables.Users), uniID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ======= CREATE ========

func (r *UserRepository) DetailedCreate(user *usermodels.UserModel, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`
	INSERT INTO %s (
		id, uni_id, username, profile_image_id, name_ar, name_en,
		email, phone, password, verified, status,
		department, gender, is_editable, is_loggable, is_anonymous
	) VALUES (
		:id, :uni_id, :username, :profile_image_id, :name_ar, :name_en,
		:email, :phone, :password, :verified, :status,
		:department, :gender, :is_editable, :is_loggable, :is_anonymous
	)`, tables.Users)
	if tx != nil {
		_, err := tx.NamedExec(query, user)
		return err
	}
	_, err := r.DB.NamedExec(query, user)
	return err
}

func (r *UserRepository) Create(user *usermodels.UserModel, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`
	INSERT INTO %s (
		id, uni_id, username, profile_image_id, name_ar, name_en,
		email, phone, password, verified, status,
		department, gender
	) VALUES (
		:id, :uni_id, :username, :profile_image_id, :name_ar, :name_en,
		:email, :phone, :password, :verified, :status,
		:department, :gender
	)`, tables.Users)
	if tx != nil {
		_, err := tx.NamedExec(query, user)
		return err
	}
	_, err := r.DB.NamedExec(query, user)
	return err
}

// ======== UPDATE ========

func (r *UserRepository) Update(user *usermodels.UserModel, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET uni_id = :uni_id, username = :username, profile_image_id = :profile_image_id, name_ar = :name_ar, name_en = :name_en,
	email = :email, phone = :phone, password = :password, verified = :verified,
	status = :status, department = :department, gender = :gender
	WHERE id = :id
	AND is_editable = true
	`, tables.Users)

	if tx != nil {
		_, err := tx.NamedExec(query, user)
		return err
	}
	_, err := r.DB.NamedExec(query, user)
	return err
}

// Update using email as key, updating the id as well
func (r *UserRepository) UpdateWithID(user *usermodels.UserModel, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`
	UPDATE %s
	SET id = :id, uni_id = :uni_id, username = :username, profile_image_id = :profile_image_id, name_ar = :name_ar, name_en = :name_en,
	phone = :phone, password = :password, verified = :verified,
	status = :status, department = :department, gender = :gender
	WHERE email = :email
	AND is_editable = true
	`, tables.Users)

	if tx != nil {
		_, err := tx.NamedExec(query, user)
		return err
	}
	_, err := r.DB.NamedExec(query, user)
	return err
}

func (r *UserRepository) VerifyUser(userID int64) error {
	query := fmt.Sprintf(`UPDATE %s SET verified = 1 WHERE id = ?`, tables.Users)
	_, err := r.DB.Exec(query, userID)
	return err
}

func (r *UserRepository) UnVerifyUser(userID int64) error {
	query := fmt.Sprintf(`UPDATE %s SET verified = 0 WHERE id = ?`, tables.Users)
	_, err := r.DB.Exec(query, userID)
	return err
}

func (r *UserRepository) UpdatePassword(userID int64, passwordHash string) error {
	query := fmt.Sprintf(`UPDATE %s SET password = ? WHERE id = ?`, tables.Users)
	_, err := r.DB.Exec(query, passwordHash, userID)
	return err
}

func (r *UserRepository) UpdateDetails(data *authmodels.RegDetailsUpdate) error {
	query := fmt.Sprintf(`
	UPDATE %s SET uni_id = :uni_id, name_ar = :name_ar, name_en = :name_en,
	phone = :phone, department = :department, gender = :gender
	WHERE id = :user_id
	`, tables.Users)
	_, err := r.DB.NamedExec(query, data)
	return err
}

func (r *UserRepository) UpdateUsername(userID int64, username string) error {
	query := fmt.Sprintf(`UPDATE %s SET username = ? WHERE id = ?`, tables.Users)
	_, err := r.DB.Exec(query, username, userID)
	return err
}

// ======== DELETE ========

func (r *UserRepository) Delete(id int64, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = ? AND is_editable = true`, tables.Users)
	if tx != nil {
		_, err := tx.Exec(query, id)
		return err
	}
	_, err := r.DB.Exec(query, id)
	return err
}

// ======== SPECIAL ========

func (r *UserRepository) Verify(id int64) error {
	_, err := r.DB.Exec(fmt.Sprintf(`UPDATE %s SET verified = ? WHERE id = ?`, tables.Users), true, id)
	return err
}

func (r *UserRepository) Suspend(id int64, tx *sqlx.Tx) error {
	query := fmt.Sprintf(`UPDATE %s SET status = ? WHERE id = ?`, tables.Users)
	if tx != nil {
		_, err := tx.Exec(query, usermodels.STATUS_SUSPENDED, id)
		return err
	}
	_, err := r.DB.Exec(query, usermodels.STATUS_SUSPENDED, id)
	return err
}

func (r *UserRepository) Activate(id int64) error {
	_, err := r.DB.Exec(fmt.Sprintf(`UPDATE %s SET status = ? WHERE id = ?`, tables.Users), usermodels.STATUS_ACTIVE, id)
	return err
}

func (r *UserRepository) RemoveSuspensionState(id int64) error {
	_, err := r.DB.Exec(fmt.Sprintf(`UPDATE %s SET status = ? WHERE id = ? AND status = ?`, tables.Users), usermodels.STATUS_ACTIVE, id, usermodels.STATUS_SUSPENDED)
	return err
}

func (r *UserRepository) BeginTransaction() (*sqlx.Tx, error) {
	return r.DB.Beginx()
}
