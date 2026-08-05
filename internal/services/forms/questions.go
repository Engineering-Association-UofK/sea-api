package forms

import (
	"encoding/json"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"strings"
)

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

func (s *FormService) DeleteQuestion(id int64) error {
	if _, err := s.formRepo.GetQuestionByID(id); err != nil {
		return errs.New(errs.NotFound, "question not found", nil)
	}
	return s.formRepo.DeleteQuestion(id)
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
