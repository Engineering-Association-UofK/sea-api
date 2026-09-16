package certmodels

var AllowedCertVersions = map[string]bool{
	"v0.1": true,
}

type V0_1 struct {
	Title     string
	Subtitle  string
	Statement string

	Name string

	CollabNameOne string
	CollabNameTwo string
	CollabRoleOne string
	CollabRoleTwo string

	QRBase64      string
	SignOneBase64 string
	SignTwoBase64 string
	StampBase64   string
}
