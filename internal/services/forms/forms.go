package forms

import (
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/repositories"
	"sea-api/internal/services"
	"strings"
	"time"
)

type FormService struct {
	formRepo       *repositories.FormRepository
	galleryService *services.GalleryService
}

func NewFormService(formRepo *repositories.FormRepository, galleryService *services.GalleryService) *FormService {
	return &FormService{formRepo: formRepo, galleryService: galleryService}
}

// ======== CREATE ========

func (s *FormService) CreateForm(userID int64, req *models.CreateFormRequest) (int64, error) {
	if err := s.isValidForm(req); err != nil {
		return 0, err
	}

	form := &models.FormModel{
		Title:                req.Title,
		Description:          req.Description,
		AllowMultipleEntries: req.AllowMultipleEntries,
		StartDate:            req.StartDate,
		EndDate:              req.EndDate,
		Type:                 req.Type,
		CreatedBy:            userID,
		CreatedAt:            time.Now(),
	}
	id, err := s.formRepo.CreateForm(form)
	if err != nil {
		return 0, err
	}

	return s.CreatePage(&models.CreatePageRequest{
		FormID:     id,
		PageNumber: 1,
	})
}

// ======== UPDATE ========

func (s *FormService) UpdateForm(req *models.UpdateFormRequest) error {
	form, err := s.formRepo.GetFormByID(req.ID)
	if err != nil {
		return errs.New(errs.NotFound, "form not found", nil)
	}
	if err := s.isValidForm(&req.CreateFormRequest); err != nil {
		return err
	}

	form.Title = req.Title
	form.Description = req.Description
	form.AllowMultipleEntries = req.AllowMultipleEntries
	form.StartDate = req.StartDate
	form.EndDate = req.EndDate
	form.Type = req.Type

	return s.formRepo.UpdateForm(form)
}

// ======== GET ONE ========

func (s *FormService) GetFormByID(id int64) (*models.FormModel, error) {
	return s.formRepo.GetFormByID(id)
}

// ======== GET MANY ========

func (s *FormService) GetAllForms() ([]models.FormSummaryResponse, error) {
	forms, err := s.formRepo.GetAllForms()
	if err != nil {
		return nil, err
	}

	var responses = []models.FormSummaryResponse{}
	for _, f := range forms {
		responses = append(responses, models.FormSummaryResponse{
			ID:                   f.ID,
			Title:                f.Title,
			Description:          f.Description,
			StartDate:            f.StartDate,
			EndDate:              f.EndDate,
			AllowMultipleEntries: f.AllowMultipleEntries,
			Type:                 f.Type,
			CreatedAt:            f.CreatedAt,
		})
	}
	return responses, nil
}

// ======== DELETE ========

func (s *FormService) DeleteForm(id int64) error {
	if _, err := s.formRepo.GetFormByID(id); err != nil {
		return errs.New(errs.NotFound, "form not found", nil)
	}
	s.galleryService.RemoveReference(models.ObjForm, id)
	return s.formRepo.DeleteForm(id)
}

// ====== CHECKS ======

func (s *FormService) isValidForm(form *models.CreateFormRequest) error {
	if strings.TrimSpace(form.Title) == "" {
		return errs.New(errs.BadRequest, "Title is not provided", nil)
	}
	if form.Description == "" {
		return errs.New(errs.BadRequest, "Description is not provided", nil)
	}
	// if _, err := s.galleryService.GetAssetByID(form.HeaderImageID); err != nil {
	// 	return errs.New(errs.BadRequest, "invalid image ID provided", nil)
	// }
	return nil
}
