package models

import "github.com/golang-jwt/jwt/v5"

type UserModel struct {
	Index    int    `db:"idx"`
	UniID    int    `db:"uni_id"`
	Username string `db:"username"`

	NameAr string `db:"name_ar"`
	NameEn string `db:"name_en"`
	Email  string `db:"email"`
	Phone  string `db:"phone"`

	Password string `db:"password"`
	Verified bool   `db:"verified"`
	Status   string `db:"status"`
}

type UserRole struct {
	Index int    `db:"user_id"`
	Role  string `db:"role"`
}

type UserResponse struct {
	Index    int    `json:"index"`
	UniID    int    `json:"uni_id"`
	Username string `json:"username"`

	NameAr string `json:"name_ar"`
	NameEn string `json:"name_en"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`

	Verified bool     `db:"verified"`
	Status   string   `db:"status"`
	Roles    []string `db:"-"`
}

type UserClaims struct {
	Index    int      `json:"index"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`

	jwt.RegisteredClaims
}
