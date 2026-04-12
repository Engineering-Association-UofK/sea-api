package forms

import (
	"encoding/json"
	"fmt"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"strconv"
)

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

func (s *FormService) getPerQuestionAnalysisData(formID int64) ([]models.QuestionAnalysisDTO, error) {
	rows, err := s.formRepo.GetFormAnalysisData(formID)
	if err != nil {
		return nil, err
	}

	analysisMap := make(map[int64]*models.QuestionAnalysisDTO)
	var orderedKeys = []int64{}

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

		if row.AnswerValue.Valid {
			if row.Type == models.FORM_NUMBER {
				var number models.NumberOption
				if err := json.Unmarshal(row.Options, &number); err != nil {
					return nil, errs.New(errs.BadRequest, "invalid options", nil)
				}

				val, err := strconv.ParseFloat(row.AnswerValue.String, 64)
				if err != nil {
					continue
				}

				binLabel := getRange(number.Min, number.Max, val)

				found := false
				for i, ans := range analysisMap[row.QuestionID].Answers {
					if ans.Value == binLabel {
						analysisMap[row.QuestionID].Answers[i].Count += row.AnswerCount
						found = true
						break
					}
				}

				if !found {
					analysisMap[row.QuestionID].Answers = append(analysisMap[row.QuestionID].Answers, models.AnswerCountDTO{
						Value: binLabel,
						Count: row.AnswerCount,
					})
				}

			} else if row.Type == models.FORM_CHECKBOX || row.Type == models.FORM_DROPDOWN || row.Type == models.FORM_RADIO {
				analysisMap[row.QuestionID].Answers = append(analysisMap[row.QuestionID].Answers, models.AnswerCountDTO{
					Value: row.AnswerValue.String,
					Count: row.AnswerCount,
				})
			}
		}
	}

	var result = []models.QuestionAnalysisDTO{}
	for _, k := range orderedKeys {
		result = append(result, *analysisMap[k])
	}

	return result, nil
}

func (s *FormService) GetEntireFormDetails(formID int64) (*models.FormDerailedResponse, error) {
	form, err := s.formRepo.GetFormByID(formID)
	if err != nil {
		return nil, err
	}

	rows, err := s.formRepo.GetFormDetailedResponses(formID)
	if err != nil {
		return nil, err
	}

	var responses = []models.FormDetailedResponseRow{}

	if len(rows) == 0 {
		return &models.FormDerailedResponse{
			ID:                   form.ID,
			Title:                form.Title,
			Description:          form.Description,
			StartDate:            form.StartDate,
			EndDate:              form.EndDate,
			AllowMultipleEntries: form.AllowMultipleEntries,
			CreatedAt:            form.CreatedAt,
			Responses:            responses,
		}, nil
	}

	var currentRow *models.FormDetailedResponseRow
	var currentResponseID int64 = -1

	for _, row := range rows {
		if row.ResponseID != currentResponseID {
			if currentRow != nil {
				responses = append(responses, *currentRow)
			}

			indexStr := fmt.Sprintf("%d", row.UserIndex)

			currentRow = &models.FormDetailedResponseRow{
				Index:     indexStr,
				NameAr:    row.NameAr,
				NameEn:    row.NameEn,
				Email:     row.Email,
				Questions: []models.ResponseQuestionDetails{},
			}
			currentResponseID = row.ResponseID
		}

		currentRow.Questions = append(currentRow.Questions, models.ResponseQuestionDetails{
			QuestionText: row.QuestionText,
			Type:         row.Type,
			AnswerValue:  row.AnswerValue,
		})
	}

	if currentRow != nil {
		responses = append(responses, *currentRow)
	}

	return &models.FormDerailedResponse{
		ID:                   form.ID,
		Title:                form.Title,
		Description:          form.Description,
		StartDate:            form.StartDate,
		EndDate:              form.EndDate,
		AllowMultipleEntries: form.AllowMultipleEntries,
		CreatedAt:            form.CreatedAt,
		Responses:            responses,
	}, nil
}

// Helpers

func getRange(oMin, oMax int, val float64) string {
	const numBins = 10
	min := float64(oMin)
	max := float64(oMax)

	if max <= min {
		max = min + 10
	}

	binSize := (max - min) / float64(numBins)

	binIndex := int((val - min) / binSize)

	if binIndex >= numBins {
		binIndex = numBins - 1
	} else if binIndex < 0 {
		binIndex = 0
	}

	binStart := min + (float64(binIndex) * binSize)
	binEnd := binStart + binSize
	return fmt.Sprintf("%.1f - %.1f", binStart, binEnd)
}
