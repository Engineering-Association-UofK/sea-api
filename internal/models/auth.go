package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token      string `json:"token"`
	UserID     int64  `json:"user_id"`
	IsVerified bool   `json:"is_verified"`
}

type ManagedClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Roles    []Role `json:"roles"`

	jwt.RegisteredClaims
}

type RegisterRequest struct {
	UserID   int64         `json:"user_id" binding:"required"`
	UniID    int64         `json:"uni_id" binding:"required"`
	Username TrimmedString `json:"username" binding:"required,min=3,max=20"`

	NameAr string `json:"name_ar" binding:"required"`
	NameEn string `json:"name_en" binding:"required"`
	Email  string `json:"email" binding:"required,email"`
	Phone  string `json:"phone"`

	Passcode        string     `json:"passcode" binding:"required"`
	Password        string     `json:"password" binding:"required"`
	ConfirmPassword string     `json:"confirm_password" binding:"required"`
	Department      Department `json:"department" binding:"required"`
	Gender          Gender     `json:"gender" binding:"required"`
}

type VerificationCodeModel struct {
	ID        int64     `db:"id"`
	Code      string    `db:"code"`
	UserID    int64     `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}

type VerifyEmailRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
}

type VerifyRequest struct {
	UserID int64  `json:"user_id" binding:"required"`
	Code   string `json:"code" binding:"required"`
}

type VerifyEmail struct {
	Input string
	Year  int
}
