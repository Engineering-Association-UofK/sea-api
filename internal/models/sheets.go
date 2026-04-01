package models

type EventUsersImport struct {
	Index  string `excel:"index"`
	NameAr string `excel:"name_ar"`
	NameEn string `excel:"name_en"`
	Email  string `excel:"email"`
	Grade  string `excel:"grade"`
}

type ImportUserUpdate struct {
	Index string `excel:"index"`
	Email string `excel:"email"`
	Phone string `excel:"phone"`
}
