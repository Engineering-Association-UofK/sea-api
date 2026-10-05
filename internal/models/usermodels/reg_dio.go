package usermodels

type TempUserResponse struct {
	ID       int64  `json:"id"`
	NameAr   string `json:"name_ar"`
	Passcode string `json:"passcode"`
}

type TempUserListResponse struct {
	Users   []TempUserResponse `json:"users"`
	Current int64              `json:"current"`
	Pages   int64              `json:"pages"`
}
