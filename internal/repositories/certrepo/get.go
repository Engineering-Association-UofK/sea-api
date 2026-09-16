package certrepo

import "sea-api/internal/models/certmodels"

func (r *CertRepository) GetCountForUser(id int64) int64 {
	return 0
}

func (r *CertRepository) GetCertList(req *certmodels.CertListRequest) ([]certmodels.CertResponse, error) {
	return nil, nil
}
