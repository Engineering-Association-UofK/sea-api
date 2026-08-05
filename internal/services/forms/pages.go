package forms

import (
	"sea-api/internal/errs"
	"sea-api/internal/models"
)

func (s *FormService) CreatePage(req *models.CreatePageRequest) (int64, error) {
	if _, err := s.formRepo.GetFormByID(req.FormID); err != nil {
		return 0, errs.New(errs.NotFound, "form not found", nil)
	}
	if req.PageNumber <= 0 {
		return 0, errs.New(errs.BadRequest, "invalid page number", nil)
	}
	if _, err := s.formRepo.GetPageByFormIdAndPageNumber(req.FormID, req.PageNumber); err == nil {
		return 0, errs.New(errs.Conflict, "page number already in use", nil)
	}
	page := &models.FormPageModel{
		FormID:     req.FormID,
		PageNumber: req.PageNumber,
	}
	return s.formRepo.CreatePage(page)
}

func (s *FormService) UpdatePage(req *models.UpdatePageRequest) error {
	page, err := s.formRepo.GetPageByID(req.ID)
	if err != nil {
		return errs.New(errs.NotFound, "page not found", nil)
	}

	if req.PageNumber != page.PageNumber {
		if _, err := s.formRepo.GetPageByFormIdAndPageNumber(page.FormID, req.PageNumber); err == nil {
			return errs.New(errs.Conflict, "page number already in use", nil)
		}
	}

	page.PageNumber = req.PageNumber
	return s.formRepo.UpdatePage(page)
}

func (s *FormService) DeletePage(id int64) error {
	if _, err := s.formRepo.GetPageByID(id); err != nil {
		return errs.New(errs.NotFound, "page not found", nil)
	}
	return s.formRepo.DeletePage(id)
}
