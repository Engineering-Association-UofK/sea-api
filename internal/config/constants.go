package config

const domain = "sea.uofk.edu"
const httpsDomain = "https://" + domain

var Links = struct {
	URL    string
	Domain string

	Register string

	CertVerify string
	DocVerify  string
}{
	URL:    httpsDomain,
	Domain: domain,

	Register: httpsDomain + "/register",

	CertVerify: httpsDomain + "/cert/verify",
	DocVerify:  httpsDomain + "/doc/verify",
}
