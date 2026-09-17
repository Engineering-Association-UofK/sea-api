package analyticsmodel

type General struct {
	Users        int64 `json:"users"`
	Events       int64 `json:"events"`
	Certificates int64 `json:"certificates"`
	Forms        int64 `json:"forms"`
	Posts        int64 `json:"posts"`
}
