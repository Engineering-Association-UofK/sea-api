package certservice

import (
	"encoding/json"
	"sea-api/internal/models"
	"sea-api/internal/models/certmodels"
	"sea-api/internal/utils/valid"
	"time"
)

func (s *CertService) CreateTemplate(req *certmodels.CreateTemplateRequest) (int64, error) {
	data, err := json.Marshal(req.Layout)
	if err != nil {
		return 0, err
	}

	return s.repo.CreateTemplate(&certmodels.CertificateTemplate{
		Name:         req.Name,
		Language:     req.Language,
		Version:      req.Version,
		LayoutConfig: data,
		CreatedAt:    time.Now(),
	})
}

func (s *CertService) UpdateTemplate(req *certmodels.UpdateTemplateRequest) error {
	data, err := json.Marshal(req.Layout)
	if err != nil {
		return err
	}

	return s.repo.UpdateTemplate(&certmodels.CertificateTemplate{
		ID:           req.ID,
		Name:         req.Name,
		Language:     req.Language,
		Version:      req.Version,
		LayoutConfig: data,
		CreatedAt:    time.Now(),
	})
}

func (s *CertService) GetTemplate(id int64) (*certmodels.TemplateResponse, error) {
	template, err := s.repo.GetTemplateByID(id)
	if err != nil {
		return nil, err
	}

	var layout certmodels.Layout
	if err := json.Unmarshal(template.LayoutConfig, &layout); err != nil {
		return nil, err
	}

	return &certmodels.TemplateResponse{
		Name:      template.Name,
		Version:   template.Version,
		Layout:    layout,
		CreatedAt: template.CreatedAt,
	}, nil
}

func (s *CertService) GetTemplateList(req *models.ListRequest) (*certmodels.TemplateListResponse, error) {
	count := s.repo.GetCount()

	pages := valid.Limit(req, count)

	templates, err := s.repo.ListTemplates(req)
	if err != nil {
		return nil, err
	}

	var list = []certmodels.TemplateResponse{}
	for _, temp := range templates {
		var layout certmodels.Layout
		err = json.Unmarshal(temp.LayoutConfig, &layout)
		if err != nil {
			return nil, err
		}

		list = append(list, certmodels.TemplateResponse{
			ID:        temp.ID,
			Name:      temp.Name,
			Version:   temp.Version,
			Layout:    layout,
			CreatedAt: temp.CreatedAt,
		})
	}

	return &certmodels.TemplateListResponse{
		Pages:   pages,
		Current: req.Page,
		Total:   count,

		List: list,
	}, nil
}
