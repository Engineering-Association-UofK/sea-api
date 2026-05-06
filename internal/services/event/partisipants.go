package event

import (
	"sea-api/internal/models"
	"sea-api/internal/utils/valid"
)

func (s *EventService) GetEventParticipants(userID int64, req models.ListRequest) ([]models.EventParticipantsResponse, error) {
	total, err := s.EventRepo.GetTotalEvents()
	if err != nil {
		return nil, err
	}

	valid.Limit(&req, total)

	participants, err := s.EventRepo.GetEventParticipantsList(req, userID)
	if err != nil {
		return nil, err
	}

	return []models.EventParticipantsResponse{
		{
			EventID:      userID,
			Participants: participants,
			Page:         req.Page,
			Current:      req.Page,
		},
	}, nil
}
