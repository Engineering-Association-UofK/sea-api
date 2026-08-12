package forms

import (
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/utils"
)

func (s *FormService) UpdateResponseStatus(req *models.UpdateResponseStatusRequest) error {
	if _, err := s.formRepo.GetResponseByID(req.ID); err != nil {
		return errs.New(errs.NotFound, "response not found", nil)
	}
	return s.formRepo.UpdateResponseStatus(req.ID, req.Status)
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

	var dtos = []models.FormResponseDTO{}
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

	var dtos = []models.FormResponseDTO{}
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

func (s *FormService) DeleteResponse(id int64) error {
	if _, err := s.formRepo.GetResponseByID(id); err != nil {
		return errs.New(errs.NotFound, "response not found", nil)
	}
	return s.formRepo.DeleteResponse(id)
}
