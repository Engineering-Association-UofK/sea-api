package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/repositories"
	"sea-api/internal/utils"
	"strconv"
	"strings"
	"time"
)

type FormService struct {
	formRepo       *repositories.FormRepository
	galleryService *GalleryService
}

func NewFormService(formRepo *repositories.FormRepository, galleryService *GalleryService) *FormService {
	return &FormService{formRepo: formRepo, galleryService: galleryService}
}

// ======== ANALYSIS ========

func (s *FormService) GetFormAnalysis(formID int64) (*models.FormAnalysisResponse, error) {
	form, err := s.formRepo.GetFormByID(formID)
	if err != nil {
		return nil, err
	}

	analysis, err := s.getPerQuestionAnalysisData(formID)
	if err != nil {
		return nil, err
	}

	responses, err := s.formRepo.GetNumberOfResponsesByFormID(formID)
	if err != nil {
		return nil, err
	}

	return &models.FormAnalysisResponse{
		FormID:         form.ID,
		Title:          form.Title,
		TotalResponses: responses,
		Questions:      analysis,
	}, nil
}

// ======== SPECIAL ========

func (s *FormService) GetFormForEdit(formID int64) (*models.FormForEditDTO, error) {
	rows, err := s.formRepo.GetFormWithQuestions(formID)
	if err != nil {
		return nil, err
	}

	url := ""
	if rows[0].HeaderImageID.Valid {
		url, err = s.galleryService.GetLinkByAssetID(rows[0].HeaderImageID.Int64)
		if err != nil {
			slog.Error("Unable to get form image url", "error", err, "Form ID", formID)
		}
	}

	form := models.UpdateFormRequest{
		ID: rows[0].FormID,
		CreateFormRequest: models.CreateFormRequest{
			Title:                rows[0].Title,
			Description:          rows[0].Description,
			AllowMultipleEntries: rows[0].AllowMultiple,
			IsActive:             rows[0].IsActive,
			HeaderImageID:        rows[0].HeaderImageID.Int64,
		},
	}

	pageMap := make(map[int64]models.UpdatePageRequest)
	var questions []models.UpdateQuestionRequest

	for _, row := range rows {
		if row.PageID == nil {
			continue
		}

		if _, exists := pageMap[*row.PageID]; !exists {
			pageMap[*row.PageID] = models.UpdatePageRequest{
				ID: *row.PageID,
				CreatePageRequest: models.CreatePageRequest{
					FormID:     row.FormID,
					PageNumber: *row.PageNum,
				},
			}
		}

		if row.QuestionID != nil {
			questions = append(questions, models.UpdateQuestionRequest{
				ID: *row.QuestionID,
				CreateQuestionRequest: models.CreateQuestionRequest{
					FormPageID:   *row.PageID,
					QuestionText: *row.QuestionText,
					Type:         *row.Type,
					Options:      *row.Options,
					IsRequired:   *row.IsRequired,
					DisplayOrder: *row.DisplayOrder,
				},
			})
		}
	}

	var pages []models.UpdatePageRequest
	for _, p := range pageMap {
		pages = append(pages, p)
	}

	return &models.FormForEditDTO{
		Url:       url,
		Form:      form,
		Pages:     pages,
		Questions: questions,
	}, nil
}

func (s *FormService) GetFormForUser(formID int64) (*models.FormForUserDTO, error) {
	rows, err := s.formRepo.GetFormWithQuestions(formID)
	if err != nil {
		return nil, err
	}

	url := ""
	if rows[0].HeaderImageID.Valid {
		url, err = s.galleryService.GetLinkByAssetID(rows[0].HeaderImageID.Int64)
		if err != nil {
			slog.Error("Unable to get form image url", "error", err, "Form ID", formID)
		}
	}
	dto := &models.FormForUserDTO{
		Form: models.FormDTO{
			ID:             rows[0].FormID,
			Title:          rows[0].Title,
			Description:    rows[0].Description,
			HeaderImageUrl: url,
		},
		Pages:     []models.FormPageDTO{},
		Questions: []models.FormQuestionDTO{},
	}

	pageMap := make(map[int64]bool)
	questionMap := make(map[int64]bool)

	for _, row := range rows {
		if row.PageID != nil && !pageMap[*row.PageID] {
			dto.Pages = append(dto.Pages, models.FormPageDTO{
				ID:         *row.PageID,
				FormID:     row.FormID,
				PageNumber: *row.PageNum,
			})
			pageMap[*row.PageID] = true
		}

		if row.QuestionID != nil && !questionMap[*row.QuestionID] {
			dto.Questions = append(dto.Questions, models.FormQuestionDTO{
				ID:           *row.QuestionID,
				FormPageID:   *row.PageID,
				QuestionText: *row.QuestionText,
				Type:         *row.Type,
				Options:      *row.Options,
				IsRequired:   *row.IsRequired,
				DisplayOrder: *row.DisplayOrder,
			})
			questionMap[*row.QuestionID] = true
		}
	}

	return dto, nil
}

func (s *FormService) SubmitForm(userID int64, req *models.SubmitFormRequest) (int64, error) {
	answers, err := s.isValidSubmitFormRequest(userID, req)
	if err != nil {
		return 0, err
	}

	response := &models.FormResponseModel{
		FormID:      req.FormID,
		UserID:      userID,
		Status:      models.FORM_SUBMITTED,
		SubmittedAt: time.Now(),
	}

	responseID, err := s.formRepo.CreateResponse(response)
	if err != nil {
		return 0, err
	}

	var answerModels []models.FormAnswerModel
	for _, a := range answers {
		answerModels = append(answerModels, models.FormAnswerModel{
			ResponseID:  responseID,
			QuestionID:  a.QuestionID,
			AnswerValue: a.AnswerValue,
		})
	}

	err = s.formRepo.CreateAnswersBatch(answerModels)
	if err != nil {
		return 0, err
	}

	return responseID, nil
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
		IsActive:             req.IsActive,
		HeaderImageID:        sql.NullInt64{Int64: req.HeaderImageID, Valid: req.HeaderImageID != 0},
		CreatedBy:            userID,
		CreatedAt:            time.Now(),
	}
	id, err := s.formRepo.CreateForm(form)
	if err != nil {
		return 0, err
	}
	if form.HeaderImageID.Valid {
		s.galleryService.AttachAssetToObject(req.HeaderImageID, models.ObjForm, id)
	}
	return s.CreatePage(&models.CreatePageRequest{
		FormID:     id,
		PageNumber: 1,
	})
}

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

func (s *FormService) CreateQuestion(req *models.CreateQuestionRequest) (int64, error) {
	q, err := s.isQuestionValid(req)
	if err != nil {
		return 0, err
	}

	question := &models.FormQuestionModel{
		FormPageID:   q.FormPageID,
		QuestionText: q.QuestionText,
		Type:         q.Type,
		Options:      &q.Options,
		IsRequired:   q.IsRequired,
		DisplayOrder: q.DisplayOrder,
	}
	return s.formRepo.CreateQuestion(question)
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

	// Options for if there is an image attached or not
	if req.HeaderImageID == 0 {
		if form.HeaderImageID.Valid {
			s.galleryService.RemoveReference(models.ObjForm, form.HeaderImageID.Int64)
		}
	} else if form.HeaderImageID.Valid {
		s.galleryService.RemoveReference(models.ObjForm, form.HeaderImageID.Int64)
		s.galleryService.AttachAssetToObject(req.HeaderImageID, models.ObjForm, req.ID)
	} else {
		s.galleryService.AttachAssetToObject(req.HeaderImageID, models.ObjForm, req.ID)
	}

	form.Title = req.Title
	form.Description = req.Description
	form.AllowMultipleEntries = req.AllowMultipleEntries
	form.IsActive = req.IsActive
	form.HeaderImageID = sql.NullInt64{Int64: req.HeaderImageID, Valid: req.HeaderImageID != 0}

	return s.formRepo.UpdateForm(form)
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

func (s *FormService) UpdateQuestion(req *models.UpdateQuestionRequest) error {
	question, err := s.formRepo.GetQuestionByID(req.ID)
	if err != nil {
		return errs.New(errs.NotFound, "question not found", nil)
	}

	q, err := s.isQuestionValid(&req.CreateQuestionRequest)
	if err != nil {
		return err
	}

	question.QuestionText = q.QuestionText
	question.Type = q.Type
	question.Options = &q.Options
	question.IsRequired = q.IsRequired
	question.DisplayOrder = q.DisplayOrder

	return s.formRepo.UpdateQuestion(question)
}

func (s *FormService) UpdateResponseStatus(req *models.UpdateResponseStatusRequest) error {
	if _, err := s.formRepo.GetResponseByID(req.ID); err != nil {
		return errs.New(errs.NotFound, "response not found", nil)
	}
	return s.formRepo.UpdateResponseStatus(req.ID, req.Status)
}

// ======== GET ONE ========

func (s *FormService) GetFormByID(id int64) (*models.FormModel, error) {
	return s.formRepo.GetFormByID(id)
}

func (s *FormService) GetResponseByID(id int64) (*models.FormResponseDTO, error) {
	resp, err := s.formRepo.GetResponseByID(id)
	if err != nil {
		return nil, err
	}

	answers, err := s.formRepo.GetAnswersByResponseID(id)
	if err != nil {
		return nil, err
	}

	answerDTOs := make([]models.FormAnswerDTO, len(answers))
	for i, a := range answers {
		answerDTOs[i] = models.FormAnswerDTO{
			ID:          a.ID,
			ResponseID:  a.ResponseID,
			QuestionID:  a.QuestionID,
			AnswerValue: a.AnswerValue,
		}
	}

	return &models.FormResponseDTO{
		ID:          resp.ID,
		FormID:      resp.FormID,
		UserID:      resp.UserID,
		Status:      resp.Status,
		SubmittedAt: resp.SubmittedAt,
		Answers:     answerDTOs,
	}, nil
}

// ======== GET MANY ========

func (s *FormService) GetAllForms() ([]models.FormSummaryResponse, error) {
	forms, err := s.formRepo.GetAllForms()
	if err != nil {
		return nil, err
	}

	if len(forms) == 0 {
		return []models.FormSummaryResponse{}, nil
	}

	var responses []models.FormSummaryResponse
	for _, f := range forms {
		responses = append(responses, models.FormSummaryResponse{
			ID:                   f.ID,
			Title:                f.Title,
			Description:          f.Description,
			IsActive:             f.IsActive,
			AllowMultipleEntries: f.AllowMultipleEntries,
			CreatedAt:            f.CreatedAt,
		})
	}
	return responses, nil
}

func (s *FormService) GetResponsesByFormID(formID int64) ([]models.FormResponseDTO, error) {
	responses, err := s.formRepo.GetResponsesByFormID(formID)
	if err != nil {
		return nil, err
	}

	if len(responses) == 0 {
		return []models.FormResponseDTO{}, nil
	}

	ids := utils.ExtractField(responses, func(r models.FormResponseModel) int64 { return r.ID })
	allAnswers, err := s.formRepo.GetAnswersByResponseIDs(ids)
	if err != nil {
		return nil, err
	}
	answerMap := make(map[int64][]models.FormAnswerModel)
	for _, a := range allAnswers {
		answerMap[a.ResponseID] = append(answerMap[a.ResponseID], a)
	}

	var dtos []models.FormResponseDTO
	for _, r := range responses {
		answers := answerMap[r.ID]
		answerDTOs := make([]models.FormAnswerDTO, len(answers))
		for i, a := range answers {
			answerDTOs[i] = models.FormAnswerDTO{
				ID:          a.ID,
				ResponseID:  a.ResponseID,
				QuestionID:  a.QuestionID,
				AnswerValue: a.AnswerValue,
			}
		}

		dtos = append(dtos, models.FormResponseDTO{
			ID:          r.ID,
			FormID:      r.FormID,
			UserID:      r.UserID,
			Status:      r.Status,
			SubmittedAt: r.SubmittedAt,
			Answers:     answerDTOs,
		})
	}
	return dtos, nil
}

func (s *FormService) GetUserResponsesForForm(userID, formID int64) ([]models.FormResponseDTO, error) {
	responses, err := s.formRepo.GetUserResponsesForForm(userID, formID)
	if err != nil {
		return nil, err
	}

	if len(responses) == 0 {
		return []models.FormResponseDTO{}, nil
	}

	ids := utils.ExtractField(responses, func(r models.FormResponseModel) int64 { return r.ID })
	allAnswers, err := s.formRepo.GetAnswersByResponseIDs(ids)
	if err != nil {
		return nil, err
	}
	answerMap := make(map[int64][]models.FormAnswerModel)
	for _, a := range allAnswers {
		answerMap[a.ResponseID] = append(answerMap[a.ResponseID], a)
	}

	var dtos []models.FormResponseDTO
	for _, r := range responses {
		answers := answerMap[r.ID]
		answerDTOs := make([]models.FormAnswerDTO, len(answers))
		for i, a := range answers {
			answerDTOs[i] = models.FormAnswerDTO{
				ID:          a.ID,
				ResponseID:  a.ResponseID,
				QuestionID:  a.QuestionID,
				AnswerValue: a.AnswerValue,
			}
		}

		dtos = append(dtos, models.FormResponseDTO{
			ID:          r.ID,
			FormID:      r.FormID,
			UserID:      r.UserID,
			Status:      r.Status,
			SubmittedAt: r.SubmittedAt,
			Answers:     answerDTOs,
		})
	}
	return dtos, nil
}

// ======== DELETE ========

func (s *FormService) DeleteForm(id int64) error {
	if _, err := s.formRepo.GetFormByID(id); err != nil {
		return errs.New(errs.NotFound, "form not found", nil)
	}
	s.galleryService.RemoveReference(models.ObjForm, id)
	return s.formRepo.DeleteForm(id)
}

func (s *FormService) DeletePage(id int64) error {
	if _, err := s.formRepo.GetPageByID(id); err != nil {
		return errs.New(errs.NotFound, "page not found", nil)
	}
	return s.formRepo.DeletePage(id)
}

func (s *FormService) DeleteQuestion(id int64) error {
	if _, err := s.formRepo.GetQuestionByID(id); err != nil {
		return errs.New(errs.NotFound, "question not found", nil)
	}
	return s.formRepo.DeleteQuestion(id)
}

func (s *FormService) DeleteResponse(id int64) error {
	if _, err := s.formRepo.GetResponseByID(id); err != nil {
		return errs.New(errs.NotFound, "response not found", nil)
	}
	return s.formRepo.DeleteResponse(id)
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

func (s *FormService) isQuestionValid(question *models.CreateQuestionRequest) (*models.CreateQuestionRequest, error) {
	if _, err := s.formRepo.GetPageByID(question.FormPageID); err != nil {
		return nil, errs.New(errs.NotFound, "Page with ID not found", nil)
	}
	switch question.Type {
	case models.FORM_CHECKBOX, models.FORM_DROPDOWN, models.FORM_RADIO:
		var options models.Options
		if err := json.Unmarshal(question.Options, &options); err != nil {
			return nil, errs.New(errs.BadRequest, "invalid options", nil)
		}
	case models.FORM_PARAGRAPH, models.FORM_TEXT:
		question.Options = json.RawMessage{byte('['), byte(']')}
	case models.FORM_NUMBER:
		var number models.NumberOption
		if err := json.Unmarshal(question.Options, &number); err != nil {
			return nil, errs.New(errs.BadRequest, "invalid options", nil)
		}
	}
	if strings.TrimSpace(question.QuestionText) == "" {
		return nil, errs.New(errs.BadRequest, "Empty question body", nil)
	}
	if !models.AllowedQuestionTypes[question.Type] {
		return nil, errs.New(errs.BadRequest, "Question type provided does not exist", nil)
	}
	return question, nil
}

func (s *FormService) arePagesValid(pages []models.CreatePageRequest) error {
	if len(pages) == 0 {
		return nil
	}
	FormID := pages[0].FormID
	for _, p := range pages {
		if p.FormID != FormID {
			return errs.New(errs.BadRequest, "All pages must belong to the same form", nil)
		}
		if p.PageNumber <= 0 {
			return errs.New(errs.BadRequest, "Invalid page number", nil)
		}
	}
	return nil
}

func (s *FormService) isValidSubmitFormRequest(userID int64, req *models.SubmitFormRequest) ([]models.AnswerRequest, error) {
	// Form should exist
	form, err := s.formRepo.GetFormByID(req.FormID)
	if err != nil {
		return nil, errs.New(errs.NotFound, "form not found", nil)
	}

	// Fetch questions for checks
	formQs, err := s.formRepo.GetQuestionsByFormID(req.FormID)
	if err != nil {
		return nil, err
	}
	formQsMap := utils.FromSlice(formQs, func(q models.FormQuestionModel) int64 { return q.ID })

	// Check for duplicate entries if not allowed
	if !form.AllowMultipleEntries {
		if res, err := s.formRepo.GetUserResponsesForForm(userID, req.FormID); err == nil {
			if len(res) > 0 {
				return nil, errs.New(errs.Conflict, "form already submitted", nil)
			}
		} else if err != sql.ErrNoRows {
			return nil, err
		}
	}

	// Check questions count
	if len(req.Answers) == 0 {
		return nil, errs.New(errs.BadRequest, "Cannot submit empty form", nil)
	}

	// Checks all required questions are answered
	answerMap := utils.FromSlice(req.Answers, func(a models.AnswerRequest) int64 {
		if a.AnswerValue != "" {
			return a.QuestionID
		} else {
			return 0
		}
	})
	answerMap.Delete(0)
	for _, q := range formQs {
		if !q.IsRequired {
			continue
		}
		if _, ok := answerMap[q.ID]; !ok {
			return nil, errs.New(errs.BadRequest, fmt.Sprintf("missing required answer for question %d", q.ID), nil)
		}
	}

	var answers []models.AnswerRequest
	// Checks Answers list
	checkedMap := map[int64]bool{}
	for _, a := range req.Answers {
		if _, ok := formQsMap[a.QuestionID]; !ok {
			return nil, errs.New(errs.BadRequest, "invalid question id", nil)
		}
		value := strings.TrimSpace(a.AnswerValue)
		if value == "" {
			// Empty optional answer, should not be saved to the database
			continue
		}
		if formQsMap[a.QuestionID].Type == models.FORM_NUMBER {
			var number models.NumberOption
			if err := json.Unmarshal(*formQsMap[a.QuestionID].Options, &number); err != nil {
				return nil, errs.New(errs.BadRequest, "invalid number options", nil)
			}
			// turn a.AnswerValue into a float64
			answer, err := strconv.ParseFloat(a.AnswerValue, 64)
			if err != nil {
				return nil, errs.New(errs.BadRequest, "invalid number answer", nil)
			}
			if (number.Min != 0 && int(answer) < number.Min) || (number.Max != 0 && int(answer) > number.Max) {
				return nil, errs.New(errs.BadRequest, "answer out of range", nil)
			}
		}
		if checkedMap[a.QuestionID] && formQsMap[a.QuestionID].Type != models.FORM_CHECKBOX {
			return nil, errs.New(errs.BadRequest, "duplicate answer to one question", nil)
		}
		checkedMap[a.QuestionID] = true
		answers = append(answers, a)
	}

	return answers, nil
}

// ====== HELPERS ======

func (s *FormService) getPerQuestionAnalysisData(formID int64) ([]models.QuestionAnalysisDTO, error) {
	rows, err := s.formRepo.GetFormAnalysisData(formID)
	if err != nil {
		return nil, err
	}

	analysisMap := make(map[int64]*models.QuestionAnalysisDTO)
	var orderedKeys []int64

	for _, row := range rows {
		if _, exists := analysisMap[row.QuestionID]; !exists {
			analysisMap[row.QuestionID] = &models.QuestionAnalysisDTO{
				QuestionID:   row.QuestionID,
				QuestionText: row.QuestionText,
				Type:         row.Type,
				Answers:      []models.AnswerCountDTO{},
			}
			orderedKeys = append(orderedKeys, row.QuestionID)
		}

		if row.AnswerValue.Valid && row.Type != models.FORM_TEXT || row.Type == models.FORM_PARAGRAPH {
			analysisMap[row.QuestionID].Answers = append(analysisMap[row.QuestionID].Answers, models.AnswerCountDTO{
				Value: row.AnswerValue.String,
				Count: row.AnswerCount,
			})
		}
	}

	var result []models.QuestionAnalysisDTO
	for _, k := range orderedKeys {
		result = append(result, *analysisMap[k])
	}

	return result, nil
}
