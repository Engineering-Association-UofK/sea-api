package models

import (
	"encoding/json"
	"time"

	_ "github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type RegistrationStepModel struct {
	RegCode string `db:"reg_code"`
	UserID  int64  `db:"user_id"`
	Step    int64  `db:"step"`
}

type RegInitCreate struct {
	ID    int64         `db:"id"`
	Email TrimmedString `db:"email"`
}

type RegDetailsUpdate struct {
	UserID int64  `db:"user_id"`
	NameAr string `db:"name_ar"`
	NameEn string `db:"name_en"`
	Gender Gender `db:"gender"`

	UniID      string     `db:"uni_id"`
	Department Department `db:"department"`

	Phone string `db:"phone"`
}

type LoginRequest struct {
	Username TrimmedString `json:"username" binding:"required"`
	Password string        `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token       string `json:"token"`
	UserID      int64  `json:"user_id"`
	Roles       []Role `json:"roles"`
	IsVerified  bool   `json:"is_verified"`
	RedirectURL string `json:"redirect_url"`
}

type ForgotPasswordRequest struct {
	UserID   *int64   `json:"user_id"`
	Email    *string  `json:"email"`
	Lang     Language `json:"lang" binding:"required"`
	Username *string  `json:"username"`
}

type ManagedClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Roles    []Role `json:"roles"`

	jwt.RegisteredClaims
}

// Registration structs

type RegistrationRequest struct {
	Step int64           `json:"step"`
	Data json.RawMessage `json:"data"`
}

type InitialRegistrationRequest struct {
	UserID   int64         `json:"user_id" validate:"required"`
	Passcode string        `json:"passcode" validate:"required"`
	Email    TrimmedString `json:"email" validate:"required,email"`
	Lang     Language      `json:"lang" validate:"required"`
}

type PasswordRegistrationRequest struct {
	RegCode         string `json:"reg_code" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

type DetailsRegistrationRequest struct {
	RegCode string `json:"reg_code" validate:"required"`
	NameAr  string `json:"name_ar" validate:"required"`
	NameEn  string `json:"name_en" validate:"required"`
	Gender  Gender `json:"gender" validate:"required"`

	UniID      string     `json:"uni_id" validate:"required"`
	Department Department `json:"department" validate:"required"`

	Phone string `json:"phone"`
}

type UsernameRegistrationRequest struct {
	RegCode  string        `json:"reg_code" validate:"required"`
	Username TrimmedString `json:"username" validate:"required,min=3,max=20"`
}

type CheckRegistrationRequest struct {
	RegCode string `json:"reg_code" binding:"required"`
}

type CheckRegistrationResponse struct {
	RegStep int64 `json:"reg_step"`
}

// End

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
