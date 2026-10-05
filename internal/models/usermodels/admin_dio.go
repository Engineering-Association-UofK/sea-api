package usermodels

type AdminResponse struct {
	ID         int64  `json:"id"`
	Email      string `json:"email"`
	ProfilePic string `json:"profile_pic"`
	NameAr     string `json:"name_ar"`
	Username   string `json:"username"`
	Gender     Gender `json:"gender"`
	Roles      []Role `json:"roles"`
}

type AdminResponseList struct {
	Admins  []AdminResponse `json:"admins"`
	Total   int64           `json:"total"`
	Pages   int64           `json:"pages"`
	Current int64           `json:"current"`
}

type AdminRequest struct {
	ID    int64  `json:"id"`
	Roles []Role `json:"roles"`
}
