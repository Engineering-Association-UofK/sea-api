package bot

import (
	"fmt"
	"log/slog"
	"sea-api/internal/models"
	"sea-api/internal/repositories"
	"sea-api/internal/services"
)

type BotService struct {
	repo *repositories.BotRepository

	feedbackService *services.FeedbackService
}

func NewBotService(repo *repositories.BotRepository, feedbackService *services.FeedbackService) *BotService {
	return &BotService{
		repo:            repo,
		feedbackService: feedbackService,
	}
}

func (s *BotService) HandleSession(req models.BotRequest, claims *models.ManagedClaims) (*models.BotResponse, error) {
	// Check for existing session
	state, err := s.repo.GetSession(req.SessionID)
	if err != nil {
		slog.Debug("Session not found")
		return nil, err
	}

	var userID *int64
	var currentNodeID int64

	var nextNode *models.NodeRow
	var nodeResponse *models.BotResponse

	if state != nil {
		// Update & check variables
		currentNodeID = state.CurrentNodeID
		if claims != nil {
			userID = &claims.UserID
		} else {
			userID = state.UserID
		}
		slog.Debug(fmt.Sprintf("State: %v", state))

		// Session exists, determine next node based on keyword and language
		nextNode, err = s.repo.GetNextNodeRow(currentNodeID, req.Language, req.Keyword)
		if err != nil {
			slog.Debug("Next node not found")
			return nil, err
		}

		slog.Debug(fmt.Sprintf("Next node: %v", nextNode))
		if nextNode.Type == models.NodeAction {
			nodeResponse, err = s.handleActions(req, nextNode, state)
			if err != nil {
				slog.Debug("Action failed")
				return nil, err
			}
		}
	} else {
		// No session, start from the beginning
		slog.Debug("No session found, starting from the beginning")
		nextNode, err = s.repo.GetStartNode(&req.Language)
		if err != nil {
			slog.Debug("Start node not found")
			return nil, err
		}
	}

	if nodeResponse == nil {
		nodeResponse, err = s.getNodeView(nextNode, req.Language)
		if err != nil {
			slog.Debug("Failed to get node view")
			return nil, err
		}
	}
	slog.Debug(fmt.Sprintf("Node response: %v", nodeResponse))

	// Update session
	if err := s.repo.UpsertSession(req.SessionID, nextNode.ID, userID); err != nil {
		return nil, err
	}

	return nodeResponse, nil
}

func (s *BotService) GoBackView(req models.BotRequest) (*models.BotResponse, error) {
	state, err := s.repo.GetSession(req.SessionID)
	if err != nil {
		slog.Debug("Session not found")
		return nil, err
	}
	if state == nil {
		slog.Debug("No session found, starting from the beginning")
		return s.HandleSession(req, nil)
	}

	parentNode, err := s.repo.GetParentOrStartNodeRow(state.CurrentNodeID, req.Language, req.Keyword)
	if err != nil {
		slog.Debug("Parent node not found")
		return nil, err
	}

	if err := s.repo.UpsertSession(req.SessionID, parentNode.ID, state.UserID); err != nil {
		slog.Debug("Failed to update session")
		return nil, err
	}

	return s.getNodeView(parentNode, req.Language)
}

func (s *BotService) getNodeView(node *models.NodeRow, lang models.Language) (*models.BotResponse, error) {
	// Get Options and translation
	edges, err := s.repo.GetEdgesForNode(node.ID, lang)
	if err != nil {
		slog.Debug("Failed to get edges")
		return nil, err
	}

	// Map to DTO
	response := &models.BotResponse{
		NodeType: node.Type,
		Content:  node.Content,
		Options:  make([]models.BotOptionView, 0, len(edges)),
	}

	for _, edge := range edges {
		response.Options = append(response.Options, models.BotOptionView{
			Keyword: edge.Keyword,
			Label:   edge.Label,
		})
	}

	return response, nil
}

func (s *BotService) ResetDefault() error {
	slog.Debug("Resetting default")
	return s.repo.ResetDefault()
}
