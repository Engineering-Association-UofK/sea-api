package models

import (
	"database/sql"
	"encoding/json"
	"strings"
)

type TrimmedString string

func (ts *TrimmedString) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*ts = TrimmedString(strings.TrimSpace(s))
	return nil
}

var AllowedAdminRoles = map[Role]bool{
	RoleSystemUserMgr:   true,
	RoleSystemSupport:   true,
	RoleContentEditor:   true,
	RoleContentBlogMgr:  true,
	RoleContentEventMgr: true,
	RoleContentFormMgr:  true,
	RoleCertifier:       true,
	RoleCertMgr:         true,
	RolePaperViewer:     true,
}

var SpecialAdminRoles = map[Role]bool{
	RoleSystemAdmin:        true,
	RoleSystemAdminManager: true,
}

var AdminRoles = []Role{
	RoleSystemUserMgr,
	RoleSystemSupport,
	RoleContentEditor,
	RoleContentBlogMgr,
	RoleContentEventMgr,
	RoleContentFormMgr,
	RoleCertifier,
	RoleCertMgr,
	RolePaperViewer,
}

type Status string

const (
	STATUS_ACTIVE    Status = "active"
	STATUS_INACTIVE  Status = "inactive"
	STATUS_SUSPENDED Status = "suspended"
	STATUS_GRADUATED Status = "graduated"
	STATUS_DROPPED   Status = "dropped"
	STATUS_WITHDRAWN Status = "withdrawn"
	STATUS_DELETED   Status = "deleted"
)

type Gender string

const (
	MALE   Gender = "male"
	FEMALE Gender = "female"
)

type Department string

const (
	DEP_MECHANICAL   Department = "mechanical"
	DEP_CIVIL        Department = "civil"
	DEP_ELECTRICAL   Department = "electrical"
	DEP_CHEMICAL     Department = "chemical"
	DEP_PETROLEUM    Department = "petroleum"
	DEP_AGRICULTURAL Department = "agricultural"
	DEP_MINING       Department = "mining"
	DEP_SURVEYING    Department = "surveying"
)

type Role string

const (
	RoleSystemSuperAdmin   Role = "sys:super_admin"
	RoleSystemAdmin        Role = "sys:admin"
	RoleSystemAdminManager Role = "sys:admin_manager"
	RoleSystemUserMgr      Role = "sys:user_manager"
	RoleSystemSupport      Role = "sys:tech_support"

	RoleContentEditor   Role = "content:editor"
	RoleContentBlogMgr  Role = "content:blog_manager"
	RoleContentEventMgr Role = "content:event_manager"
	RoleContentFormMgr  Role = "content:form_manager"

	RoleCertifier   Role = "cert:certifier"
	RoleCertMgr     Role = "cert:manager"
	RolePaperViewer Role = "cert:viewer"

	RoleOrgOwner     Role = "org:owner"
	RoleOrgMember    Role = "org:member"
	RoleOrgModerator Role = "org:moderator"
)

type UserModel struct {
	ID       int64  `db:"id"`
	UniID    int64  `db:"uni_id"`
	Username string `db:"username"`

	ProfileImageID sql.NullInt64 `db:"profile_image_id"`

	NameAr string `db:"name_ar"`
	NameEn string `db:"name_en"`
	Email  string `db:"email"`
	Phone  string `db:"phone"`

	Department Department `db:"department"`
	Gender     Gender     `db:"gender"`

	Password string `db:"password"`
	Verified bool   `db:"verified"`
	Status   Status `db:"status"`

	IsEditable  bool `db:"is_editable"`
	IsLoggable  bool `db:"is_loggable"`
	IsAnonymous bool `db:"is_anonymous"`
}

type TempUserModel struct {
	ID       sql.NullInt64  `db:"id"`
	UniID    sql.NullInt64  `db:"uni_id"`
	Username sql.NullString `db:"username"`

	NameAr sql.NullString `db:"name_ar"`
	NameEn sql.NullString `db:"name_en"`
	Email  sql.NullString `db:"email"`
	Phone  sql.NullString `db:"phone"`

	Department sql.NullString `db:"department"`
	Gender     sql.NullString `db:"gender"`

	Password sql.NullString `db:"password"`
	Verified sql.NullBool   `db:"verified"`
	Status   sql.NullString `db:"status"`
}

type UserRole struct {
	UserID int64 `db:"user_id"`
	Role   Role  `db:"role"`
}

type UserResponse struct {
	ID       int64  `json:"id"`
	UniID    int64  `json:"uni_id"`
	Username string `json:"username"`

	ProfilePic string `json:"profile_pic"`

	NameAr string `json:"name_ar"`
	NameEn string `json:"name_en"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`

	Department Department `json:"department"`
	Gender     Gender     `json:"gender"`

	Verified bool   `json:"verified"`
	Status   Status `json:"status"`
	Roles    []Role `json:"roles"`
}

type UserListItemResponse struct {
	ID       int64  `json:"id"`
	UniID    int64  `json:"uni_id"`
	Username string `json:"username"`

	Email      string     `json:"email"`
	Department Department `json:"department"`
	Gender     Gender     `json:"gender"`

	Verified bool   `json:"verified"`
	Status   Status `json:"status"`
	Roles    []Role `json:"roles"`
}

type AdminResponse struct {
	ID         int64  `json:"id"`
	Email      string `json:"email"`
	ProfilePic string `json:"profile_pic"`
	NameAr     string `json:"name_ar"`
	Username   string `json:"username"`
	Gender     Gender `json:"gender"`
	Roles      []Role `json:"roles"`
}

type AdminRequest struct {
	ID    int64  `json:"id"`
	Roles []Role `json:"roles"`
}

type TempUserResponse struct {
	ID       int64  `json:"id"`
	NameAr   string `json:"name_ar"`
	Passcode string `json:"passcode"`
}

type UserListResponse struct {
	Users   []UserListItemResponse `json:"users"`
	Current int64                  `json:"current"`
	Pages   int64                  `json:"pages"`
}

type TempUserListResponse struct {
	Users   []TempUserResponse `json:"users"`
	Current int64              `json:"current"`
	Pages   int64              `json:"pages"`
}

type UserProfileResponse struct {
	ID       int64  `json:"id"`
	UniID    int64  `json:"uni_id"`
	Username string `json:"username"`

	NameAr string `json:"name_ar"`
	NameEn string `json:"name_en"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`

	Department Department `json:"department"`
	Gender     Gender     `json:"gender"`
	ProfilePic string     `json:"profile_pic"`
}

type UpdateProfileRequest struct {
	ID         int64         `json:"id" binding:"required"`
	UniID      int64         `json:"uni_id" binding:"required"`
	NameAr     TrimmedString `json:"name_ar" binding:"required"`
	NameEn     TrimmedString `json:"name_en" binding:"required"`
	Phone      TrimmedString `json:"phone" binding:"required"`
	Department Department    `json:"department" binding:"required"`
	Gender     Gender        `json:"gender" binding:"required"`
}

type UserDetails struct {
	UserProfileResponse
}

type UpdateUsernameRequest struct {
	Username TrimmedString `json:"username" binding:"required"`
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

type GetPasscodeResponse struct {
	Passcode string `json:"passcode"`
}
