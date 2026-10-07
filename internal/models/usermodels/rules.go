package usermodels

type Status string
type Gender string
type Role string

const (
	STATUS_ACTIVE    Status = "active"
	STATUS_INACTIVE  Status = "inactive"
	STATUS_SUSPENDED Status = "suspended"
	STATUS_GRADUATED Status = "graduated"
	STATUS_DROPPED   Status = "dropped"
	STATUS_WITHDRAWN Status = "withdrawn"
	STATUS_DELETED   Status = "deleted"
)

const (
	MALE   Gender = "male"
	FEMALE Gender = "female"
)

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
