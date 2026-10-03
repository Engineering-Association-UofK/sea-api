package eventservice

import (
	"context"
	"sea-api/internal/models"
	"sea-api/internal/models/eventmodels"
	"sea-api/internal/utils/valid"
)

func (s *EventService) GetEventViewList(ctx context.Context, req *eventmodels.EventListRequest) (*eventmodels.EventViewListResponse, error) {
	count, err := s.repo.GetCount(req)
	if err != nil {
		return nil, err
	}

	totalPages := valid.Limit(&req.ListRequest, count)

	events, err := s.repo.GetList(req)
	if err != nil {
		return nil, err
	}

	eventIDs := make([]int64, 0, len(*events))
	for _, event := range *events {
		eventIDs = append(eventIDs, event.ID)
	}

	coords, err := s.repo.GetCoordsByEventIDs(eventIDs)
	if err != nil {
		return nil, err
	}
	coordsByEventID := make(map[int64][]eventmodels.CoordResponse, len(eventIDs))
	for _, coord := range coords {
		coordsByEventID[coord.EventID] = append(coordsByEventID[coord.EventID], eventmodels.CoordResponse{
			Name: coord.Name,
			Role: coord.Role,
		})
	}

	var list = []eventmodels.EventViewResponse{}
	for _, e := range *events {
		url, err := s.s3.GenerateDownloadUrlByKey(ctx, e.BackgroundKey)
		if err != nil {
			return nil, err
		}

		list = append(list, eventmodels.EventViewResponse{
			ID:              e.ID,
			Name:            e.Name,
			Description:     e.Description,
			BackgroundURL:   url,
			RequireApplying: e.RequireApplying,
			Belonging:       e.Belonging,
			StartDate:       e.StartDate,
			EndDate:         e.EndDate,
			Coords:          coordsByEventID[e.ID],
		})
	}

	return &eventmodels.EventViewListResponse{
		List: list,
		ListResponse: models.ListResponse{
			TotalPages:  totalPages,
			CurrentPage: req.Page,
			Count:       count,
		},
	}, nil
}
