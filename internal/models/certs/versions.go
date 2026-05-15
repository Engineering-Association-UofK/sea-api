package certs

type V1_0 struct {
	Name        string
	EventName   string
	Grade       float64
	TaskColumns [][]string
	QRCode      string

	CoordinatorName      string
	CoordinatorTitle     string
	CoordinatorSignature string

	PresenterName      string
	PresenterTitle     string
	PresenterSignature string

	Stamp string

	StartDate string
	EndDate   string
	NowDate   string
}
