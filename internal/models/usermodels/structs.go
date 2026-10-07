package usermodels

import (
	"database/sql"
	"sea-api/internal/models"
)

type UserRow struct {
	ID             int64             `db:"id"`
	UniID          string            `db:"uni_id"`
	Username       string            `db:"username"`
	ProfileImageID sql.NullInt64     `db:"profile_image_id"`
	ProfilePicKey  sql.NullString    `db:"profile_pic_key"`
	NameAr         string            `db:"name_ar"`
	NameEn         string            `db:"name_en"`
	Email          string            `db:"email"`
	Phone          string            `db:"phone"`
	Department     models.Department `db:"department"`
	Gender         Gender            `db:"gender"`
	Verified       bool              `db:"verified"`
	Status         Status            `db:"status"`
}

type AdminRow struct {
	ID         int64          `db:"id"`
	Email      string         `db:"email"`
	ProfilePic sql.NullString `db:"profile_pic"`
	NameAr     string         `db:"name_ar"`
	Username   string         `db:"username"`
	Gender     Gender         `db:"gender"`
	Role       Role           `db:"role"`
}
