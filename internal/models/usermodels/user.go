package usermodels

import (
	"database/sql"
	"sea-api/internal/models"
)

type UserModel struct {
	ID       int64   `db:"id"`
	UniID    *string `db:"uni_id"`
	Username *string `db:"username"`

	ProfileImageID sql.NullInt64 `db:"profile_image_id"`

	NameAr *string `db:"name_ar"`
	NameEn *string `db:"name_en"`
	Email  *string `db:"email"`
	Phone  *string `db:"phone"`

	Department *models.Department `db:"department"`
	Gender     *Gender            `db:"gender"`

	Password *string `db:"password"`
	Verified bool    `db:"verified"`
	Status   Status  `db:"status"`

	IsEditable  bool `db:"is_editable"`
	IsLoggable  bool `db:"is_loggable"`
	IsAnonymous bool `db:"is_anonymous"`
}

type UserRole struct {
	UserID int64 `db:"user_id"`
	Role   Role  `db:"role"`
}

type Passcode struct {
	ID       int64  `db:"id"`
	Passcode string `db:"passcode"`
}
