package eventservice

import (
	"context"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/models/eventmodels"
	"sea-api/internal/utils/valid"
)

func (s *EventService) GetParticipants(ctx context.Context, eventID int64, req *eventmodels.ParticipantListRequest) (*eventmodels.ParticipantListResponse, error) {
	count, err := s.repo.CountParticipationByEventID(eventID)
	if err != nil {
		return nil, err
	}

	pages := valid.Limit(&req.ListRequest, count)

	participants, err := s.repo.GetParticipantViews(eventID, req.Limit, req.Page)
	if err != nil {
		return nil, err
	}

	list := make([]eventmodels.ParticipantResponse, 0, len(participants))
	for _, participant := range participants {
		photoURL := ""
		if participant.PhotoKey != nil {
			photoURL, err = s.s3.GenerateDownloadUrlByKey(ctx, *participant.PhotoKey)
			if err != nil {
				return nil, err
			}
		}

		list = append(list, eventmodels.ParticipantResponse{
			ID:       participant.ID,
			EventID:  participant.EventID,
			UserID:   participant.UserID,
			Username: participant.Username,
			PhotoURL: photoURL,
			JoinedAt: participant.JoinedAt,
		})
	}

	return &eventmodels.ParticipantListResponse{
		List: list,
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
