package event

import (
	"sea-api/internal/errs"
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

func (s *EventService) UpdateParticipant(req models.ParticipantUpdateRequest) error {
	participant, err := s.EventRepo.GetParticipantByID(req.ID)
	if err != nil {
		return errs.New(errs.NotFound, "Participant not found", nil)
	}

	participant.Status = req.Status
	participant.Completed = req.Completed

	// Update participant basic info
	if err := s.EventRepo.UpdateParticipant(participant); err != nil {
		return err
	}

	// Update grades/scores if provided
	if len(req.Grades) > 0 {
		scores := make([]models.ComponentScoreModel, len(req.Grades))
		for i, g := range req.Grades {
			scores[i] = models.ComponentScoreModel{
				ParticipantID: req.ID,
				ComponentID:   g.ComponentID,
				Score:         g.Score,
			}
		}

		// Get existing scores to determine if we update or create
		existingScores, err := s.EventRepo.GetScoresByParticipantID(req.ID)
		if err != nil {
			return err
		}

		scoreMap := make(map[int64]int64)
		for _, es := range existingScores {
			scoreMap[es.ComponentID] = es.ID
		}

		var toUpdate []models.ComponentScoreModel
		var toCreate []models.ComponentScoreModel

		for _, ns := range scores {
			if id, exists := scoreMap[ns.ComponentID]; exists {
				ns.ID = id
				toUpdate = append(toUpdate, ns)
			} else {
				toCreate = append(toCreate, ns)
			}
		}

		if len(toUpdate) > 0 {
			if err := s.EventRepo.MassUpdateScore(toUpdate); err != nil {
				return err
			}
		}
		if len(toCreate) > 0 {
			if err := s.EventRepo.MassCreateScore(toCreate, nil); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *EventService) BatchUpdateParticipant(req []models.ParticipantUpdateRequest) error {
	for _, r := range req {
		if err := s.UpdateParticipant(r); err != nil {
			return err
		}
	}

	return nil
}

func (s *EventService) Complete(userID, eventID int64) error {
	participant, err := s.EventRepo.GetParticipantByEventAndUserIDs(eventID, userID)
	if err != nil {
		return errs.New(errs.NotFound, "Participant not found", nil)
	}

	participant.Completed = true
	return s.EventRepo.UpdateParticipant(participant)
}
