package usermodels

import "sea-api/internal/models"

type UserResponse struct {
	ID       int64  `json:"id"`
	UniID    string `json:"uni_id"`
	Username string `json:"username"`

	ProfilePic string `json:"profile_pic"`

	NameAr string `json:"name_ar"`
	NameEn string `json:"name_en"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`

	Department models.Department `json:"department"`
	Gender     Gender            `json:"gender"`

	Verified bool   `json:"verified"`
	Status   Status `json:"status"`
	Roles    []Role `json:"roles"`
}

type UserListItemResponse struct {
	ID       int64  `json:"id"`
	UniID    string `json:"uni_id"`
	Username string `json:"username"`

	Email      string            `json:"email"`
	Department models.Department `json:"department"`
	Gender     Gender            `json:"gender"`

	Verified bool   `json:"verified"`
	Status   Status `json:"status"`
	Roles    []Role `json:"roles"`
}

type UserListResponse struct {
	Users   []UserListItemResponse `json:"users"`
	Current int64                  `json:"current"`
	Pages   int64                  `json:"pages"`
}
