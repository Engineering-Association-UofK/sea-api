package bot

import (
	"fmt"
	"log/slog"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"strings"
)

var actionsMap = map[models.BotActionType]func(req models.BotRequest, action *models.BotAction, s *BotService, userID *int64) (string, error){
	models.BotActionFeedback: func(req models.BotRequest, action *models.BotAction, s *BotService, userID *int64) (string, error) {
		if models.AllowedFeedbackTypes[models.FeedbackType(action.ActionText)] {
			_, err := s.feedbackService.Create(req.Input, userID, models.FeedbackTechnical)
			return "", err
		}
		return "", errs.New(errs.InternalServerError, "invalid feedback type", nil)
	},
	models.BotActionRedirect: func(req models.BotRequest, action *models.BotAction, s *BotService, userID *int64) (string, error) {
		return action.ActionText, nil
	},
}

func (s *BotService) handleActions(req models.BotRequest, nextNode *models.NodeRow, state *models.UserState) (*models.BotResponse, error) {
	actionModel, err := s.repo.GetAction(nextNode.ID)
	if err != nil {
		return nil, err
	}
	if action, ok := actionsMap[actionModel.ActionType]; ok {

		data, err := action(req, actionModel, s, state.UserID)
		if err != nil {
			return nil, err
		}

		metadata := getMetadata(data)

		slog.Debug(fmt.Sprintf("%v", metadata))
		view, err := s.getNodeView(nextNode, req.Language)
		if err != nil {
			return nil, err
		}
		view.Metadata = metadata

		return view, nil
	}
	slog.Error("action not found in actionsMap", "keyword", req.Keyword)
	return nil, errs.New(errs.Conflict, "invalid action", nil)
}

func getMetadata(data string) interface{} {
	if data == "" {
		return nil
	}
	if strings.HasPrefix(data, "http") {
		return models.BotRedirect{
			ActionText: data,
			IsInternal: false,
		}
	}
	return models.BotRedirect{
		ActionText: data,
		IsInternal: true,
	}
}
