package repositories

import (
	"sea-api/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAll() ([]models.UserModel, error) {
	var users []models.UserModel
	err := r.db.Select(&users, `SELECT * FROM users`)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetAllByIndices(indices []int64) ([]models.UserModel, error) {
	var users []models.UserModel

	if len(indices) == 0 {
		return users, nil
	}

	query, args, err := sqlx.In(
		`SELECT * FROM users WHERE idx IN (?)`,
		indices,
	)
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

func (r *UserRepository) GetAllRoles() ([]models.UserRole, error) {
	var roles []models.UserRole
	err := r.db.Select(&roles, `SELECT * FROM user_roles`)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *UserRepository) GetRolesByUserID(userID int64) ([]string, error) {
	var roles []string
	err := r.db.Select(&roles, `SELECT role FROM user_roles WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *UserRepository) GetByIndex(index int64) (*models.UserModel, error) {
	var user models.UserModel
	err := r.db.Get(&user, `SELECT * FROM users WHERE id = ?`, index)
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

func (r *UserRepository) GetByUniID(uniID string) (*models.UserModel, error) {
	var user models.UserModel
	err := r.db.Get(&user, `SELECT * FROM users WHERE uni_id = ?`, uniID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *models.UserModel) error {
	_, err := r.db.Exec(`INSERT INTO users (idx, uni_id, username, name_ar, name_en, email, phone, password, verified, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.Index, user.UniID, user.Username, user.NameAr, user.NameEn, user.Email, user.Phone, user.Password, user.Verified, user.Status)
	return err
}

func (r *UserRepository) Update(user *models.UserModel) error {
	_, err := r.db.Exec(`UPDATE users SET uni_id = ?, username = ?, name_ar = ?, name_en = ?, email = ?, phone = ?, password = ?, verified = ?, status = ? WHERE id = ?`,
		user.UniID, user.Username, user.NameAr, user.NameEn, user.Email, user.Phone, user.Password, user.Verified, user.Status, user.Index)
	return err
}

func (r *UserRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id = ?`, id)
	return err
}
