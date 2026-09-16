package eventservice

import (
	"sea-api/internal/errs"
	"sea-api/internal/models/eventmodels"
)

func (s *EventService) AddCoordinators(eventID int64, list *eventmodels.AddCoordsRequest) error {
	models := []eventmodels.EventCoord{}
	for _, e := range list.List {
		models = append(models, eventmodels.EventCoord{
			EventID: eventID,
			Name:    e.Name,
			Role:    e.Role,
		})
	}
	if err := s.repo.CreateBatchCoords(models); err != nil {
		return err
	}
	return nil
}

func (s *EventService) GetCoordinators(eventID int64) ([]eventmodels.EventCoord, error) {
	coords, err := s.repo.GetCoordsByEventID(eventID)
	if err != nil {
		return nil, err
	}
	return coords, nil
}

func (s *EventService) UpdateCoordinator(eventId, id int64, coord *eventmodels.CoordRequest) error {
	if err := s.repo.UpdateCoord(&eventmodels.EventCoord{
		ID:      id,
		EventID: eventId,
		Name:    coord.Name,
		Role:    coord.Role,
	}); err != nil {
		return err
	}
	return nil
}

func (s *EventService) DeleteCoordinator(coordID, eventID int64) error {
	coord, err := s.repo.GetCoordsByEventID(eventID)
	if err != nil {
		return err
	}

	var found bool
	for _, c := range coord {
		if c.ID == coordID {
			found = true
		}
	}

	if found {
		if err := s.repo.DeleteCoord(coordID); err != nil {
			return err
		}
	} else {
		return errs.New(errs.Forbidden, "Coord does not belong to this event", nil)
	}
	return nil
}

func (s *EventService) DeleteCoordinatorsByEventID(eventID int64) error {
	if err := s.repo.DeleteCoordsByEventID(eventID); err != nil {
		return err
	}
	return nil
}
