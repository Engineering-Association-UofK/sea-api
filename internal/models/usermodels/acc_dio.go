package usermodels

import (
	"sea-api/internal/models"
	"sea-api/internal/utils"
)

type UserProfileSummaryResponse struct {
	// The unique identifier university Index for the user
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`

	// Link to the user profile picture
	ProfilePic string            `json:"profile_pic"`
	NameAr     string            `json:"name_ar"`
	NameEn     string            `json:"name_en"`
	Department models.Department `json:"department"`
	Gender     Gender            `json:"gender"`
}

type UserProfileResponse struct {
	ID       int64  `json:"id"`
	UniID    string `json:"uni_id"`
	Username string `json:"username"`

	NameAr string `json:"name_ar"`
	NameEn string `json:"name_en"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`

	Department models.Department `json:"department"`
	Gender     Gender            `json:"gender"`
	ProfilePic string            `json:"profile_pic"`
}

type UpdateProfileRequest struct {
	ID         int64               `json:"id" binding:"required"`
	UniID      string              `json:"uni_id" binding:"required"`
	NameAr     utils.TrimmedString `json:"name_ar" binding:"required"`
	NameEn     utils.TrimmedString `json:"name_en" binding:"required"`
	Phone      utils.TrimmedString `json:"phone" binding:"required"`
	Department models.Department   `json:"department" binding:"required"`
	Gender     Gender              `json:"gender" binding:"required"`
}

type UserDetails struct {
	UserProfileResponse
}

type UpdateUsernameRequest struct {
	Username utils.TrimmedString `json:"username" binding:"required"`
}

type UpdateEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type UpdatePasswordRequest struct {
	OldPassword     string `json:"old_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type CheckUsername struct {
	Available bool `json:"available"`
}
