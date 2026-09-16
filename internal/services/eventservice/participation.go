package eventservice

import (
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/models/eventmodels"
	"sea-api/internal/utils/valid"
)

func (s *EventService) GetParticipants(eventID int64, req *eventmodels.EventParticipantListRequest) (*eventmodels.EventParticipantListResponse, error) {
	count, err := s.repo.CountApplications(eventID)
	if err != nil {
		return nil, err
	}

	pages := valid.Limit(&req.ListRequest, count)

	participants, err := s.repo.GetParticipationByEventID(eventID, req.Limit, req.Page)
	if err != nil {
		return nil, err
	}
	return &eventmodels.EventParticipantListResponse{
		List: participants,
		ListResponse: models.ListResponse{
			TotalPages:  pages,
			CurrentPage: req.Page,
			Count:       count,
		},
	}, nil
}

func (s *EventService) RemoveParticipant(eventID, partID int64) error {
	part, err := s.repo.GetParticipation(partID)
	if err != nil {
		return err
	}

	if part.ID != partID {
		return errs.New(errs.Forbidden, "Participant does not belong to event", nil)
	}

	return s.repo.DeleteParticipation(partID)
}
