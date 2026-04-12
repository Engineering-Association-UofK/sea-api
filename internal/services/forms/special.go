package forms

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/utils"
	"strconv"
	"strings"
	"time"
)

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
			StartDate:            rows[0].StartDate,
			EndDate:              rows[0].EndDate,
			HeaderImageID:        rows[0].HeaderImageID.Int64,
		},
	}

	pageMap := make(map[int64]models.UpdatePageRequest)
	questions := []models.UpdateQuestionRequest{}

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

	pages := []models.UpdatePageRequest{}
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

	var answerModels = []models.FormAnswerModel{}
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

	var answers = []models.AnswerRequest{}
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
